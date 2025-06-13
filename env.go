package hermeti

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sean9999/pear"
	"github.com/spf13/afero"
)

// Env is a computing environment.
type Env struct {
	InStream   io.Reader
	queue      []byte
	OutStream  io.Writer
	ErrStream  io.Writer
	Filesystem afero.Fs
	Randomness io.Reader
	Args       []string
	Vars       map[string]string
	Exit       func(int)
}

// take strings of the form "foo=bar" and return a map
func stringsToMap(kvs []string) map[string]string {
	m := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		x := strings.Split(kv, "=")
		if len(x) == 2 {
			m[x[0]] = x[1]
		}
		if len(x) == 1 {
			m[kv] = ""
		}
		if len(x) > 2 {
			m[x[0]] = strings.Join(x[1:], "=")
		}
	}
	return m
}

// RealEnv creates a real Env for a CLI, using standard OS resources
func RealEnv() Env {
	e := Env{
		InStream:   os.Stdin,
		OutStream:  os.Stdout,
		ErrStream:  os.Stderr,
		Filesystem: afero.NewOsFs(),
		Randomness: rand.Reader,
		Args:       os.Args,
		Vars:       stringsToMap(os.Environ()),
		Exit:       os.Exit,
	}
	return e
}

// TestEnv creates an Env suitable for testing
func TestEnv() Env {
	env := Env{
		InStream:   new(bytes.Buffer),
		OutStream:  new(bytes.Buffer),
		ErrStream:  new(bytes.Buffer),
		Filesystem: afero.NewMemMapFs(),
		Args:       []string{},
		Vars:       map[string]string{},
		Exit:       func(_ int) {},
	}
	return env
}

// mount a subdirectory into an environment. Useful for testing. Probably dangerous otherwise
func (env *Env) Mount(dirfs fs.ReadDirFS, at string) error {
	if env.Filesystem == nil {
		return errors.New("nil filesystem")
	}
	entries, err := dirfs.ReadDir(".")
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			srcFile, err := dirfs.Open(e.Name())
			if err != nil {
				return err
			}
			destFile, err := env.Filesystem.Create(filepath.Join(at, e.Name()))
			if err != nil {
				return err
			}
			io.Copy(destFile, srcFile)
			srcFile.Close()
		}
	}
	return nil
}

// PipeIn pipes a stream into stdIn
func (env *Env) PipeIn(r io.Reader) error {

	if r == nil {
		return pear.New("nil reader")
	}

	buf := new(bytes.Buffer)
	if env.InStream != nil {
		existingBytes, err := io.ReadAll(env.InStream)
		if err != nil {
			return err
		}
		buf.Write(existingBytes)
	}
	newBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	buf.Write(newBytes)
	env.InStream = buf
	return nil

}

func (env *Env) PipeInFile(fpath string) error {

	fd, err := env.Filesystem.Open(fpath)
	if err != nil {
		return fmt.Errorf("could not pipe in file. %w", err)
	}
	return env.PipeIn(fd)

}

func (env *Env) PipeInFiles(fpaths ...string) error {
	var e error

	for _, fpath := range fpaths {
		err := env.PipeInFile(fpath)
		if err != nil {
			if e != nil {
				e = fmt.Errorf("%w. %w", e, err)
			} else {
				e = err
			}
		}
	}
	return e
}
