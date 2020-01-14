package screws

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

//IFiling 文件处理器接口
type IFiling interface {
	SuffixOfFile(fileHeader *multipart.FileHeader) string
	NewFilePath(dir string) string
	CheckUploadFile(fileHeader *multipart.FileHeader, fileType string) error
	SaveUploadFile(uniqueNumber uint, filePath, savePath /*root*/ string, fileHeaders ...*multipart.FileHeader) ([]string, error)
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
	return "." + s[len(s)-1]
}

//NewFilePath 生成文件目录
func (f *filing) NewFilePath(dir string) string {
	return dir + "/" + time.Now().Format("2006/01") + "/"
}

//CheckUploadFile 检查上传文件：doc/img
func (f *filing) CheckUploadFile(fileHeader *multipart.FileHeader, fileType string) error {
	switch fileType {
	case "doc":
		if fileHeader.Size > 6000000 {
			return errors.New(fileHeader.Filename + "-File size is too large: >6M")
		}
		if !strings.HasSuffix(fileHeader.Filename, ".jpg") && !strings.HasSuffix(fileHeader.Filename, ".jpeg") && !strings.HasSuffix(fileHeader.Filename, ".png") && !strings.HasSuffix(fileHeader.Filename, ".pdf") {
			return errors.New(fileHeader.Filename + "-File format error: jpg/jpeg/png/pdf")
		}
	case "img":
		if fileHeader.Size > 4000000 {
			return errors.New(fileHeader.Filename + "-File size is too large: >4M")
		}
		if !strings.HasSuffix(fileHeader.Filename, ".jpg") && !strings.HasSuffix(fileHeader.Filename, ".jpeg") && !strings.HasSuffix(fileHeader.Filename, ".png") {
			return errors.New(fileHeader.Filename + "-File format error: jpg/jpeg/png")
		}
	default:
		return errors.New("fileType not exist")
	}
	return nil
}

//SaveUploadFile 保存上传文件
func (f *filing) SaveUploadFile(uniqueNumber uint, filePath, savePath /*root*/ string, fileHeaders ...*multipart.FileHeader) ([]string, error) {
	var fileNames []string
	for _, fileHeader := range fileHeaders {
		fileHash, err := NewTinyTools().SHA256OfFile(fileHeader)
		if err != nil {
			return nil, err
		}
		if err := os.MkdirAll(savePath, 0777); err != nil {
			return nil, err
		}
		var newFileName string
		if uniqueNumber != 0 {
			uniqueNumber++
			newFileName = NewTinyTools().SHA256OfString(strconv.Itoa(int(uniqueNumber))) + f.SuffixOfFile(fileHeader)
		}
		newFileName = fileHash + f.SuffixOfFile(fileHeader)

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
