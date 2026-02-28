package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-file/common"
	"go-file/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetChatPage renders the chat room page
func GetChatPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"message":  "",
		"option":   common.OptionMap,
		"username": username,
	})
}

// GetChatMessages returns recent chat messages
func GetChatMessages(c *gin.Context) {
	messages, err := model.GetRecentChatMessages(100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取消息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    messages,
	})
}

// SendChatMessage sends a text message
type SendMessageRequest struct {
	Content string `json:"content"`
}

func SendChatMessage(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "请先登录",
		})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求",
		})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "消息内容不能为空",
		})
		return
	}

	message := &model.ChatMessage{
		Username: username,
		Content:  req.Content,
		Type:     "text",
	}

	if err := model.InsertChatMessage(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "发送消息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    message,
	})
}

// UploadChatFile handles file uploads for chat
// Files are saved to {ExplorerRootPath}/chat/{subpath}/
func UploadChatFile(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "请先登录",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "获取文件失败",
		})
		return
	}
	defer file.Close()

	// Determine file type
	ext := filepath.Ext(header.Filename)
	fileType := "file"
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" {
		fileType = "image"
	}

	// Get subpath from form (relative path within chat folder)
	subpath := c.PostForm("path")
	subpath = strings.TrimSpace(subpath)

	// Build save directory: {ExplorerRootPath}/chat/{subpath}
	saveDir := filepath.Join(common.ExplorerRootPath, "chat")
	if subpath != "" {
		saveDir = filepath.Join(saveDir, subpath)
	}

	// Security check: ensure the path is under ExplorerRootPath
	absSaveDir, _ := filepath.Abs(saveDir)
	if !strings.HasPrefix(absSaveDir, common.ExplorerRootPath) {
		saveDir = filepath.Join(common.ExplorerRootPath, "chat")
		absSaveDir, _ = filepath.Abs(saveDir)
	}

	// Generate unique filename
	id := uuid.New().String()
	filename := fmt.Sprintf("%s%s", id, ext)
	uploadFilePath := filepath.Join(saveDir, filename)

	// Create directory if not exists
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建目录失败",
		})
		return
	}

	// Save file
	if err := c.SaveUploadedFile(header, uploadFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存文件失败",
		})
		return
	}

	// Build file URL for access via explorer
	// URL format: /explorer?path=/chat/{subpath}/{filename}
	var fileUrl string
	if subpath != "" {
		fileUrl = fmt.Sprintf("/explorer?path=/chat/%s/%s", subpath, filename)
	} else {
		fileUrl = fmt.Sprintf("/explorer?path=/chat/%s", filename)
	}

	message := &model.ChatMessage{
		Username: username,
		Type:     fileType,
		FileUrl:  fileUrl,
		FileName: header.Filename,
	}

	if err := model.InsertChatMessage(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "发送消息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    message,
	})
}

// Long polling for new messages
func PollChatMessages(c *gin.Context) {
	lastTimestamp := int64(0)
	if ts := c.Query("last"); ts != "" {
		fmt.Sscanf(ts, "%d", &lastTimestamp)
	}

	// Poll for up to 30 seconds
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    []interface{}{},
			})
			return
		case <-ticker.C:
			var messages []*model.ChatMessage
			model.DB.Where("timestamp > ?", lastTimestamp).Order("timestamp asc").Find(&messages)
			if len(messages) > 0 {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data":    messages,
				})
				return
			}
		}
	}
}

// ClearChatMessages clears all chat messages from database (admin only)
func ClearChatMessages(c *gin.Context) {
	err := model.ClearAllChatMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "清空聊天记录失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "聊天记录已清空",
	})
}