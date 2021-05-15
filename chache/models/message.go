package models


import (
	"fmt"
	"github.com/jinzhu/gorm"
	u "go-contacts/utils"
)

type Message struct {
	gorm.Model
	Text string   `json:"text"`
	UserId int	  `json:"userid"`
	FileId int    `json:"fileid"`
}

func (message *Message) Validate() (map[string]interface{}, bool) {

	if message.Text == "" {
		return u.Message(false, "Text should be on the payload") , false
	}

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (message *Message) Create() (map[string]interface{}) {

	if resp, ok := message.Validate(); !ok {
		return resp
	}

	GetDB().Create(message)

	resp := u.Message(true, "success")
	resp["message"] = message
	return resp
}

func GetMessages(userid int) []Message {
	var	message []Message
	err := GetDB().Table("messages").Where("user_id = ?", userid).Find(&message).Error
	if err != nil {
		return nil
	}
	return message
}

func GetMessage(id int) *Message {

	var message *Message
	err := GetDB().Table("files").Where("id = ?", id).First(message).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return message
}
