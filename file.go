package screws

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

//IFiling 文件处理器接口
type IFiling interface {
	SuffixOfFile(fileHeader *multipart.FileHeader) string
	DateDir(dir string) string
	CheckUploadFile(requiredSize int64, requiredSuffix /* separate with "/" like ".jpg/.png" */ string, fileHeaders ...*multipart.FileHeader) error
	SaveUploadFile(rootDir, filePath string, fileHeaders ...*multipart.FileHeader) ([]string, error)
	DeleteUploadedFile(rootDir string, filePaths ...string) error
	ReadDirItems(dir string, s *[]string) error
}

//NewFiling 初始化文件处理器
func NewFiling() IFiling {
	return &filing{}
}

//filing 文件处理器
type filing struct {
}

//SuffixOfFile 获取文件后缀
func (f *filing) SuffixOfFile(fileHeader *multipart.FileHeader) string {
	s := strings.Split(fileHeader.Filename, ".")
	if len(s) < 2 {
		return ".unknown"
	}
	return "." + strings.ToLower(s[len(s)-1])
}

//DateDir 日期目录: "/dir/2006/01"
func (f *filing) DateDir(dir string) string {
	return dir + "/" + time.Now().Format("2006/01")
}

//CheckUploadFile 检查上传文件：doc/img
func (f *filing) CheckUploadFile(requiredSize int64, requiredSuffix /* separate with "/" like ".jpg/.png" */ string, fileHeaders ...*multipart.FileHeader) error {
	for _, fileHeader := range fileHeaders {
		if fileHeader.Size > requiredSize {
			return fmt.Errorf("%s is too large: > %d MB", fileHeader.Filename, requiredSize/1000000)
		}
		if !strings.Contains(requiredSuffix, f.SuffixOfFile(fileHeader)) {
			return fmt.Errorf("%s format error: need %s", fileHeader.Filename, requiredSuffix)
		}
	}
	return nil
}

//SaveUploadFile 保存上传文件
func (f *filing) SaveUploadFile(rootDir, filePath string, fileHeaders ...*multipart.FileHeader) ([]string, error) {
	var savePath = rootDir + filePath
	var fileNames []string
	for _, fileHeader := range fileHeaders {
		var newFileName string
		fileHash, err := NewTinyTools().SHA1OfFile(fileHeader)
		if err != nil {
			return nil, err
		}
		newFileName = fileHash + f.SuffixOfFile(fileHeader)

		if err := os.MkdirAll(savePath, 0777); err != nil {
			return nil, err
		}
		src, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		out, err := os.Create(savePath + newFileName)
		if err != nil {
			return nil, err
		}
		defer out.Close()
		_, err = io.Copy(out, src)
		if err != nil {
			return nil, err
		}

		fileNames = append(fileNames, filePath+newFileName)
	}
	return fileNames, nil
}

//DeleteUploadedFile  删除已上传文件
func (f *filing) DeleteUploadedFile(rootDir string, filePaths ...string) error {
	e := ""
	for _, filePath := range filePaths {
		if err := os.Remove(rootDir + filePath); err != nil {
			e = e + err.Error()
			continue
		}
	}
	if e != "" {
		return errors.New(e)
	}
	return nil
}

//ReadDirItems 递归遍历目录项
func (f *filing) ReadDirItems(dir string, s *[]string) error {
	file, err := os.OpenFile(dir, os.O_RDONLY, os.ModeDir)
	if err != nil {
		log.Println(err)
		return err
	}
	fileInfos, err := file.Readdir(-1)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			newDir := dir + "/" + fileInfo.Name()
			if err := f.ReadDirItems(newDir, s); err != nil {
				log.Println(err)
				return err
			}
		} else {
			fileName := fileInfo.Name()
			*s = append(*s, dir+"/"+fileName)
		}
	}
	return nil
}
