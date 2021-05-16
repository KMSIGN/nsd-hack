package controllers

import (
	"bufio"
	"encoding/json"
	"github.com/KMSIGN/nsd-hack/MainServer/models"
	u "github.com/KMSIGN/nsd-hack/MainServer/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func CreateFile(c *gin.Context) {
	user := c.MustGet("user") //Grab the id of the user that send the request
	file := &models.File{}

	err := json.NewDecoder(c.Request.Body).Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message":"Error while decoding request body"})
		return
	}

	var serverAddr string


	serverAddr, err = GetFileServer()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{"message":"No fileserver found"})
		c.Abort()
		return
	}

	file.ServerAddr = serverAddr
	log.Println(serverAddr)

	data := url.Values{
		"hash":        {file.Hash},
		"hashes": 	   {file.Hashes},
		"LastPartSize":{file.LastPartSize},
	}
	resp, err := http.PostForm(serverAddr+"/upload", data)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{"message":"Fileserver response error"})
		c.Abort()
		return
	}

	reader := bufio.NewReader(resp.Body)
	res, err := reader.ReadByte()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{"message":"Fileserver bad response"})
		c.Abort()
		return
	}

	file.UserId = int(user.(uint))
	response := file.Create()
	response["fileserverTcpPort"] = string(res)
	c.JSON(http.StatusOK, response)
}


type ServersList struct {
	Servers []struct {
		Addr string `json:"addr"`
	} `json:"servers"`
}

func GetFileServer() (string, error) {
	// some difficult analyses and logic
	fl, err := ioutil.ReadFile("fileservers.json")
	if err != nil {	return "", err }
	var temp ServersList
	err = json.Unmarshal(fl, &temp)
	if err != nil { return "", err}
	return temp.Servers[0].Addr, nil
}

func GetFilesFor(c *gin.Context) {
	id := c.GetInt("user")
	data := models.GetFiles(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	c.JSON(http.StatusOK, resp)
}

type DownloadTemp struct {
	Id int
	Addr string
}

func GetFile(c *gin.Context) {
	//user := c.GetInt("user")
	inf := &DownloadTemp{}

	err := json.NewDecoder(c.Request.Body).Decode(inf)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message":"Error while decoding request body"})
		return
	}
	file := models.GetFile(inf.Id)

	addr := file.ServerAddr

	data := url.Values{
		"addr": {inf.Addr},
		"hash": {file.Hash},
	}
	resp, err := http.PostForm(addr+"/download", data)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{"message":"No fileserver found"})
		c.Abort()
		return
	}

	reader := bufio.NewReader(resp.Body)
	res, err := reader.ReadByte()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{"message":"No fileserver found"})
		c.Abort()
		return
	}

	response := u.Message(true, string(res))
	c.JSON(http.StatusOK, response)
}
