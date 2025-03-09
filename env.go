package hermeti

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// Env is a computing environment.
type Env struct {
	InStream   io.Reader
	OutStream  io.Writer
	ErrStream  io.Writer
	Filesystem afero.Fs
	Randomness io.Reader
	Args       []string
	Vars       map[string]string
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
