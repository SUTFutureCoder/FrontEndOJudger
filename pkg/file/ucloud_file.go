package file

type ucloud struct {
	fileBase
}

func (u *ucloud) Put(data []byte, dst string) (string, error) {
	panic("implement me")
}

func (u *ucloud) Get(filePath string) ([]byte, error) {
	panic("implement me")
}

func (u *ucloud) Delete(filePath string) error {
	panic("implement me")
}

func (u *ucloud) Unzip(filePath string, dst string) error {
	panic("implement me")
}
