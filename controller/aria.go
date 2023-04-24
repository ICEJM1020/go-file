package controller

import (
	"context"
	"fmt"
	"go-file/common"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func ariaDownload(saveURL, savePath string) (err error) {

	ariaCmd := exec.Command("/aria2c", fmt.Sprintf("--dir=%s", savePath), saveURL)
	ariaCmd.Stderr = os.Stderr
	err = ariaCmd.Start()
	return err
}

func ServerDowanload(c *gin.Context) {
	uploadPath := common.UploadPath
	path := c.PostForm("path")
	saveurl := c.PostForm("saveurl")
	if path != "" {
		uploadPath = filepath.Join(common.ExplorerRootPath, path)
		if !strings.HasPrefix(uploadPath, common.ExplorerRootPath) {
			// In this case the given path is not valid, so we reset it to ExplorerRootPath.
			uploadPath = common.ExplorerRootPath
		}

		// Start a go routine to delete explorer' cache
		if common.ExplorerCacheEnabled {
			go func() {
				ctx := context.Background()
				rdb := common.RDB
				key := "cacheExplorer:" + uploadPath
				rdb.Del(ctx, key)
			}()
		}
	}

	err := ariaDownload(saveurl, uploadPath)
	if err != nil {
		message := "Failed to create offline download" + err.Error()
		common.SysError(message)
		c.String(http.StatusInternalServerError, message)
		return
	}
	// uploader := c.GetString("username")
	// if uploader == "" {
	// 	uploader = "匿名用户"
	// }
	// currentTime := time.Now().Format("2006-01-02 15:04:05")
	// if saveToDatabase {
	// 	fileObj := &model.File{
	// 		Uploader: uploader,
	// 		Time:     currentTime,
	// 		Link:     link,
	// 		Filename: filename,
	// 	}
	// 	err = fileObj.Insert()
	// 	if err != nil {
	// 		common.SysError("failed to insert file to database: " + err.Error())
	// 		continue
	// 	}
	// }
}
