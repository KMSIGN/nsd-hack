package models

import (
	"fmt"
	u "github.com/KMSIGN/nsd-hack/MainServer/utils"
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Name   string `json:"name"`
	Status string `json:"status"`
	Hash string   `json:"hash"`
	Hashes string `json:"hashes" ;sql:"-"`
	PublicHashes string `json:"public_hashes"`
	UserId int	  `json:"userid"`
	ServerAddr string `json:"server_addr"`
	Cached string `json:"cached" ;sql:"-"`
}

func (file *File) Validate() (map[string]interface{}, bool) {

	if file.Name == "" {
		return u.Message(false, "Name should be on the payload") , false
	}

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (file *File) Create() map[string]interface{} {

	if resp, ok := file.Validate(); !ok {
		return resp
	}

	GetDB().Create(file)

	resp := u.Message(true, "success")
	resp["file"] = file
	return resp
}

func GetFiles(userid int) []File {
	var	file []File
	err := GetDB().Table("files").Where("user_id = ?", userid).Find(&file).Error
	if err != nil {
		return nil
	}
	return file
}

func GetFile(id int) *File {
	var file *File
	err := GetDB().Table("files").Where("id = ?", id).First(file).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return file
}
