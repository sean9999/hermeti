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
	"testing/fstest"
)

var ErrHermeti = errors.New("hermeti")

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

// RealEnv creates a real Env for a CLI, suitable for main().
// The rootDir you pass in is used as an [os.Root].
// RealEnv("/") will give you the classic insecure behaviour.
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

// TestEnv creates an Env suitable for testing.
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
		return fmt.Errorf("%w: can't spy: %w", ErrHermeti, err)
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
		return nil, fmt.Errorf("%w: cannot capture output", ErrHermeti)
	}
	return buf, nil
}

// PipeIn pipes a stream into InStream (ex: os.StdIn)
func (env *Env) PipeIn(r io.Reader) error {
	if r == nil {
		return fmt.Errorf("%w: nil reader", ErrHermeti)
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

func (env *Env) PipeInFile(filePath string) error {
	fd, err := env.Filesystem.Open(filePath)
	if err != nil {
		return fmt.Errorf("%w: could not pipe in file. %w", ErrHermeti, err)
	}
	return env.PipeIn(fd)
}

// PipeInFiles pipes in files to the environment's stdin (InputStream)
func (env *Env) PipeInFiles(filePaths ...string) error {
	var e error
	for _, filePath := range filePaths {
		err := env.PipeInFile(filePath)
		if err != nil {
			if env != nil {
				e = fmt.Errorf("%w. %w", e, err)
			} else {
				e = err
			}
		}
	}
	return e
}
