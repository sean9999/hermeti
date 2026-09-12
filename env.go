package hermeti

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	"github.com/sean9999/pear"
)

type Filesystem interface {
	fs.StatFS
	fs.ReadDirFS
	fs.ReadFileFS
	fs.ReadLinkFS
}

// Env is a computing environment.
type Env struct {
	InStream    io.Reader
	queue       []byte
	OutStream   io.Writer
	ErrStream   io.Writer
	Filesystem  Filesystem
	Randomness  io.Reader
	Args        []string
	Vars        map[string]string
	Exit        func(int)
	Chdir       func(string) error
	UserHomeDir func() (string, error)
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
func RealEnv(rootDir string) Env {

	root, err := os.OpenRoot(rootDir)
	if err != nil {
		panic(err)
	}

	e := Env{
		InStream:    os.Stdin,
		OutStream:   os.Stdout,
		ErrStream:   os.Stderr,
		Filesystem:  root.FS().(Filesystem),
		Randomness:  rand.Reader,
		Args:        os.Args,
		Vars:        stringsToMap(os.Environ()),
		Exit:        os.Exit,
		Chdir:       os.Chdir,
		UserHomeDir: os.UserHomeDir,
	}
	return e
}

// TestEnv creates an Env suitable for testing
func TestEnv() Env {
	env := Env{
		InStream:   new(bytes.Buffer),
		OutStream:  new(bytes.Buffer),
		ErrStream:  new(bytes.Buffer),
		Filesystem: make(fstest.MapFS),
		Args:       []string{},
		Vars:       map[string]string{},
		Exit:       func(_ int) {},
	}
	return env
}

func (env *Env) Spy(ch chan string) error {
	buf, err := env.CaptureOutput()
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(buf)
	go func() {
		for sc.Scan() {
			ch <- sc.Text()
		}
	}()
	return nil
}

func (env *Env) CaptureOutput() (*bytes.Buffer, error) {
	buf, ok := env.OutStream.(*bytes.Buffer)
	if !ok {
		return nil, errors.New("cannot capture output")
	}
	return buf, nil
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

// PipeInFiles pipes in files to the environment's stdin (InputStream)
func (env *Env) PipeInFiles(fPaths ...string) error {
	var e error
	for _, fPath := range fPaths {
		err := env.PipeInFile(fPath)
		if err != nil {
			if env != nil {
				e = fmt.Errorf("%w. %w", env, err)
			} else {
				e = err
			}
		}
	}
	return e
}
