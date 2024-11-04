package utils

import (
	"crypto/sha256"
	"fmt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"novaro-server/config"
	"novaro-server/model"
	"novaro-server/service"
	"os"
	"path/filepath"
	"time"
)

func UploadFiles(w http.ResponseWriter, r *http.Request) ([]*model.Imgs, error) {
	if r.Method != http.MethodPost {

		return nil, fmt.Errorf("Method not allowed")
	}

	// 解析multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return nil, err
	}

	files := r.MultipartForm.File["images"]
	//sourceId := r.FormValue("sourceId")
	//if sourceId == "" {
	//	return nil, fmt.Errorf("sourceId is empty")
	//}
	open, err := files[0].Open()
	if err != nil {
		return nil, err
	}
	fileHash, err := HashMultipartFile(open)
	if err != nil {
		return nil, err
	}
	sourceId := fmt.Sprintf("%x", fileHash)
	fmt.Println(sourceId)

	var imglist []*model.Imgs
	var errlist []error

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// 处理每个文件
		imgs, err := handleFile(file, fileHeader, sourceId)

		if err != nil {
			errlist = append(errlist, err)
			continue
		}

		imglist = append(imglist, imgs)
	}

	logger.Infof("upload files fail:total:%d", len(errlist))

	return imglist, nil

}
func handleFile(file multipart.File, fileHeader *multipart.FileHeader, sourceId string) (*model.Imgs, error) {
	// 创建唯一的文件名
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))

	// 获取当前日期
	now := time.Now()
	curDate := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())

	// 生成路径
	uploadDir := config.Get().Client.UploadPath + "/" + curDate

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		return nil, err
	}

	filepath := filepath.Join(uploadDir, filename)

	// 存入db
	uploadFile, err2 := service.NewImgsService().UploadFile(filepath, sourceId)
	if err2 != nil {
		return nil, err2
	}

	// 创建目标文件
	dst, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}
	return uploadFile, nil
}

func HashMultipartFile(file multipart.File) ([]byte, error) {
	h := sha256.New()
	_, err := io.Copy(h, file)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
