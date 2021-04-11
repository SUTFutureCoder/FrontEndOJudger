package file

import (
	"FrontEndOJudger/pkg/setting"
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Ifile interface {
	Put(data []byte, dst string) (string, error)
	Get(filePath string) ([]byte, error)
	Delete(filePath string) error
	Unzip(filePath string, dst string) error
}

type fileBase struct {
	UserId uint64
}

const LOCALFILE = "LOCALFILE"
const UCLOUD = "UCLOUD"

var fileTools = map[string]Ifile{
	LOCALFILE: &localFile{},
	UCLOUD: &ucloud{},
}

func GetFileTool(toolType string) (Ifile, error) {
	if _, ok := fileTools[toolType]; !ok {
		return nil, errors.New("file tool type not exist")
	}
	return fileTools[toolType], nil
}

func GetDest(submitId uint64) (string, error){
	testChamberDirName := fmt.Sprintf("%s/%d/", setting.JudgerSetting.TestChamberBaseDir, submitId)
	testChamberFileName := fmt.Sprintf("%sindex.html", testChamberDirName)
	// 检查是否存在
	_, err := os.Stat(testChamberFileName)
	if err == nil || os.IsExist(err) {
		return testChamberDirName, nil
	}
	err = os.MkdirAll(testChamberDirName, 0777)
	if err != nil {
		log.Printf("Can not mkdir [%s] error [%s]", testChamberDirName, err)
	}
	return testChamberDirName, err
}

func PutLocal(data []byte, dst string) error {
	return ioutil.WriteFile(dst, data, 0777)
}

func UnzipLocal(filePath string, dst string) error {
	zipReader, _ := zip.OpenReader(filepath.Join(setting.FileSetting.FileBaseDir, filePath))
	if zipReader == nil {
		return errors.New("zip file not exist")
	}

	unZipFunc := func(file *zip.File) error {

		decodeName := file.Name
		if file.Flags == 0{
			//如果标致位是0  则是默认的本地编码   默认为gbk
			i:= bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content,_:= ioutil.ReadAll(decoder)
			decodeName = string(content)
		}

		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()
		extractedFilePath := filepath.Join(
			dst,
			decodeName,
		)
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(extractedFilePath, file.Mode())
			if err != nil {
				log.Println(err)
			}
		} else {
			d1 := filepath.Dir(extractedFilePath)
			_, err := os.Stat(d1)
			if !os.IsExist(err) {
				os.MkdirAll(d1, file.Mode())
			}

			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, file := range zipReader.Reader.File {
		err := unZipFunc(file)
		if err != nil {
			return err
		}
	}
	return nil
}