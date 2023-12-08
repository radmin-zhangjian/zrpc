package dao

import (
	"log"
	"zrpc/example/model"
	"zrpc/example/utils"
)

// 文档：https://gorm.io/zh_CN/docs/query.html

// MessageCreate Create无论什么情况都执行插入
func MessageCreate(message *model.Message) {
	db := utils.GetDB()
	if result := db.Model(&message).Create(&message); result.Error != nil {
		log.Printf("MessageCreate: %v", result.Error)
	}
	return
}

// MessageSave Save需要插入的数据存在则不进行插入
func MessageSave(message *model.Message) {
	db := utils.GetDB()
	if result := db.Model(&message).Create(&message); result.Error != nil {
		log.Printf("MessageSave: %v", result.Error)
	}
	return
}
