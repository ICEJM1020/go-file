package model

import (
	"strings"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username" gorm:"unique;"`
	Password    string `json:"password" gorm:"not null;"`
	DisplayName string `json:"displayName"`
	Role        int    `json:"role" gorm:"type:int;default:1"`   // admin, common
	Status      int    `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Token       string `json:"token"`
}

func (user *User) Insert() error {
	err := DB.Create(user).Error
	return err
}

func (user *User) Update() error {
	err := DB.Model(user).Updates(user).Error
	return err
}

func (user *User) Delete() error {
	err := DB.Delete(user).Error
	return err
}

func (user *User) ValidateAndFill() {
	// When querying with struct, GORM will only query with non-zero fields,
	// that means if your field’s value is 0, '', false or other zero values,
	// it won’t be used to build query conditions
	DB.Where(&user).First(&user)
}

func ValidateUserToken(token string) (user *User) {
	if token == "" {
		return nil
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	user = &User{}
	if DB.Where("token = ?", token).First(user).RowsAffected == 1 {
		return user
	}
	return nil
}

func ValidateUserPassword(username string, password string) (user *User) {
	if password == "" || username == "" {
		return nil
	}
	// password = strings.Replace(password, "Bearer ", "", 1)
	user = &User{}
	if DB.Where("username = ?", username).First(user).RowsAffected == 1 {
		if user.Password == password {
			return user
		}
		return nil
	}
	return nil
}
