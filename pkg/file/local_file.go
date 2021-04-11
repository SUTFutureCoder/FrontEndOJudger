package file

import (
	"io/ioutil"
)

type localFile struct {
	fileBase
}

func (l *localFile) Put(data []byte, dst string) (string, error) {
	return dst, ioutil.WriteFile(dst, data, 0777)
}

func (l *localFile) Get(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func (l *localFile) Delete(filePath string) error {
	panic("implement me")
}

func (l *localFile) Unzip(filePath string, dst string) error {
	return UnzipLocal(filePath, dst)
}
