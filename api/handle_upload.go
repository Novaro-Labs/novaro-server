package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"novaro-server/model"
	"novaro-server/service"
	"novaro-server/utils"
	"os"
	"path/filepath"
	"strings"
)

type UploadApi struct {
}

func NewUploadApi() *UploadApi {
	return &UploadApi{}
}

func (api *UploadApi) UploadFile(c *gin.Context) {
	files, err := utils.UploadFiles(c.Writer, c.Request)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": files,
		"msg":  "success",
	})
}

func (api *UploadApi) LoadSql(c *gin.Context) {
	sqlDir := "./db"
	var err error
	// 遍历目录中的所有.sql文件
	err = filepath.Walk(sqlDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
			// 读取SQL文件内容
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil
			}

			// 执行SQL语句
			err = model.GetDB().Exec(string(content)).Error
			if err != nil {
				log.Printf("Error executing SQL from file %s: %v", path, err)
			} else {
				log.Printf("Successfully executed SQL from file: %s", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Error walking through directory:", err)
	}
}

// TokenImg godoc
// @Summary get image
// @Description get the image based on the tokenId
// @Tags upload
// @Accept json
// @Produce json
// @Param sourceId query string true "sourceId"
// @Success 200
// @Failure 400
// @Router /v1/api/upload/getTokenImg [get]
func (api *UploadApi) TokenImg(c *gin.Context) {
	path, err := service.NewImgsService().GetBySourceId(c.Query("sourceId"))
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": path,
	})
}
