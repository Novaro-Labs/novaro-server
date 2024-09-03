package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"novaro-server/config"
	"novaro-server/service"
	"os"
	"path/filepath"
	"time"
)

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	fmt.Println(files)
	sourceId := r.FormValue("sourceId")

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// 处理每个文件
		err = handleFile(file, fileHeader, sourceId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
func handleFile(file multipart.File, fileHeader *multipart.FileHeader, sourceId string) error {
	// 创建唯一的文件名
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))

	// 获取当前日期
	now := time.Now()
	curDate := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())

	// 生成路径
	uploadDir := config.UploadPath + "/" + curDate

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		return err
	}

	filepath := filepath.Join(uploadDir, filename)

	// 存入db
	if err := service.NewImgsService().UploadFile(filepath, sourceId); err != nil {
		return err
	}

	// 创建目标文件
	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil
	}
	return nil
}
