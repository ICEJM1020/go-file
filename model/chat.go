package model

import (
	"time"
)

// ChatMessage represents a chat message
type ChatMessage struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"index"`
	Content   string `json:"content" gorm:"type:text"`
	Type      string `json:"type"`      // "text", "image", "file"
	FileUrl   string `json:"file_url"`  // URL for image/file if Type is not "text"
	FileName  string `json:"file_name"` // Original filename for files
	Timestamp int64  `json:"timestamp"`
}

func InitChatTable() error {
	return DB.AutoMigrate(&ChatMessage{}).Error
}

func InsertChatMessage(message *ChatMessage) error {
	message.Timestamp = time.Now().Unix()
	return DB.Create(message).Error
}

func GetRecentChatMessages(limit int) ([]*ChatMessage, error) {
	var messages []*ChatMessage
	err := DB.Order("timestamp desc").Limit(limit).Find(&messages).Error
	// Reverse to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, err
}

func DeleteOldChatMessages(keepCount int) error {
	var count int64
	DB.Model(&ChatMessage{}).Count(&count)
	if count > int64(keepCount) {
		var messages []*ChatMessage
		DB.Order("timestamp desc").Limit(int(count) - keepCount + 1).Find(&messages)
		if len(messages) > 0 {
			oldest := messages[len(messages)-1]
			return DB.Where("timestamp < ?", oldest.Timestamp).Delete(&ChatMessage{}).Error
		}
	}
	return nil
}

// ClearAllChatMessages deletes all chat messages from database
// Note: This does NOT delete the actual files stored on disk
func ClearAllChatMessages() error {
	return DB.Where("1 = 1").Delete(&ChatMessage{}).Error
}