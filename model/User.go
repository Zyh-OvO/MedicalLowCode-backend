package model

import (
	"MedicalLowCode-backend/util"
	"database/sql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type User struct {
	UserId       int `gorm:"primaryKey;autoIncrement"`
	UserName     string
	UserPassword string
	UserEmail    string
}

func (u User) TableName() string {
	return "user"
}

type RegisterCode struct {
	RegisterId int
	UserEmail  string
	Code       string
	ExpireTime time.Time
}

func (r RegisterCode) TableName() string {
	return "register_code"
}

type ResetCode struct {
	ResetId    int
	UserEmail  string
	Code       string
	ExpireTime time.Time
}

func (r ResetCode) TableName() string {
	return "reset_code"
}

func AddRegisterCode(userEmail string, code string) {
	registerCode := RegisterCode{}
	result := DB.Where("user_email = ?", userEmail).First(&registerCode)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			registerCode.UserEmail = userEmail
			registerCode.Code = code
			registerCode.ExpireTime = time.Now().Add(time.Minute * 10)
			DB.Select("user_email", "code", "expire_time").Create(&registerCode)
		} else {
			panic(result.Error)
		}
	} else {
		registerCode.Code = code
		registerCode.ExpireTime = time.Now().Add(time.Minute * 10)
		DB.Where("register_id = ?", registerCode.RegisterId).Select("code", "expire_time").Updates(&registerCode)
	}
}

func QueryRegisterCode(userEmail string) string {
	registerCode := RegisterCode{}
	result := DB.Where("user_email = ?", userEmail).First(&registerCode)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ""
		} else {
			panic(result.Error)
		}
	}
	if time.Now().After(registerCode.ExpireTime) {
		DB.Where("register_id = ?", registerCode.RegisterId).Delete(&RegisterCode{})
		return ""
	}
	return registerCode.Code
}

func RegisterNewUser(userEmail string, userName string, userPassword string) {
	user := User{
		UserName:     userName,
		UserPassword: util.BcryptPassword(userPassword),
		UserEmail:    userEmail,
	}
	result := DB.Select("user_name", "user_password", "user_email").Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	//创建用户根目录
	dir := Directory{
		UserId:  user.UserId,
		DirName: "user_" + strconv.Itoa(user.UserId),
		IsRoot:  true,
	}
	err := DB.Transaction(func(tx *gorm.DB) error {
		//sql1
		if err := tx.Create(&dir).Error; err != nil {
			return err
		}
		//sql2
		insertPathSql := "insert into directory_path (ancestor_id, descendant_id, depth) values (@newDirId, @newDirId, 0)"
		if err := tx.Exec(insertPathSql, sql.Named("newDirId", dir.DirId)).Error; err != nil {
			return err
		}
		//检查目录是否存在
		targetDirPath := filepath.Join(UserFileRootDirPath, dir.DirName)
		if _, err := os.Stat(targetDirPath); os.IsNotExist(err) {
			// 使用 MkdirAll 函数递归创建目录
			if err := os.MkdirAll(targetDirPath, os.ModePerm); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			return os.ErrExist
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func QueryUser(userEmail string) *User {
	user := User{}
	result := DB.Where("user_email = ?", userEmail).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(result.Error)
		}
	}
	return &user
}

func AddResetCode(userEmail string, code string) {
	resetCode := ResetCode{}
	result := DB.Where("user_email = ?", userEmail).First(&resetCode)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			resetCode.UserEmail = userEmail
			resetCode.Code = code
			resetCode.ExpireTime = time.Now().Add(time.Minute * 10)
			DB.Select("user_email", "code", "expire_time").Create(&resetCode)
		} else {
			panic(result.Error)
		}
	} else {
		resetCode.Code = code
		resetCode.ExpireTime = time.Now().Add(time.Minute * 10)
		DB.Where("reset_id = ?", resetCode.ResetId).Select("code", "expire_time").Updates(&resetCode)
	}
}

func QueryResetCode(userEmail string) string {
	resetCode := ResetCode{}
	result := DB.Where("user_email = ?", userEmail).First(&resetCode)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ""
		} else {
			panic(result.Error)
		}
	}
	if time.Now().After(resetCode.ExpireTime) {
		DB.Where("reset_id = ?", resetCode.ResetId).Delete(&ResetCode{})
		return ""
	}
	return resetCode.Code
}

func ResetPassword(userEmail string, newPassword string) {
	user := User{}
	result := DB.Where("user_email = ?", userEmail).First(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	user.UserPassword = util.BcryptPassword(newPassword)
	DB.Where("user_id = ?", user.UserId).Select("user_password").Updates(&user)
}
