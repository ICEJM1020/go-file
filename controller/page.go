package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-file/common"
	"go-file/model"
	"net/http"
	"runtime"
	"time"
)

func GetManagePage(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	var uptime = time.Since(common.StartTime)
	session := sessions.Default(c)
	role := session.Get("role")
	c.HTML(http.StatusOK, "manage.html", gin.H{
		"message":                 "",
		"option":                  common.OptionMap,
		"username":                c.GetString("username"),
		"memory":                  fmt.Sprintf("%d MB", m.Sys/1024/1024),
		"uptime":                  common.Seconds2Time(int(uptime.Seconds())),
		"userNum":                 model.CountTable("users"),
		"fileNum":                 model.CountTable("files"),
		"imageNum":                model.CountTable("images"),
		"FileUploadPermission":    common.FileUploadPermission,
		"FileDownloadPermission":  common.FileDownloadPermission,
		"ImageUploadPermission":   common.ImageUploadPermission,
		"ImageDownloadPermission": common.ImageDownloadPermission,
		"isAdmin":                 role == common.RoleAdminUser,
		"StatEnabled":             common.StatEnabled,
	})
}

func GetLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"message":  "",
		"option":   common.OptionMap,
		"username": c.GetString("username"),
	})
}

func Get404Page(c *gin.Context) {
	c.HTML(http.StatusOK, "404.html", gin.H{
		"message":  "",
		"option":   common.OptionMap,
		"username": c.GetString("username"),
	})
}