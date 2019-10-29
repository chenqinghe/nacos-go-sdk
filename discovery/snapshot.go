package discovery

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Snapshoter interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type Snapshot struct {
	baseDir string
}

func NewSnapshot(dir string) (*Snapshot, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, errors.New("parameter dir must be a dir path")
	}

	return &Snapshot{baseDir: dir}, nil
}

func (s *Snapshot) Set(key string, value []byte) error {
	fd, err := os.OpenFile(filepath.Join(s.baseDir, key), os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = fd.Write(value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Snapshot) Get(key string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(s.baseDir, key))
}

func (s *Snapshot) Delete(key string) error {
	os.Remove(filepath.Join(s.baseDir, key))
	return nil
}
