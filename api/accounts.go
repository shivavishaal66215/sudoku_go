package main

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleRegister(c *gin.Context){
	c.Request.ParseForm()

	//fetching data from POST body
	username := strings.Join(c.Request.Form["username"],"")
	password := strings.Join(c.Request.Form["username"],"")

	if username == "" || password == ""{
		//this means user already exists
		c.IndentedJSON(http.StatusForbidden,"credentials missing")
		return
	}
	
	//checking if user already exists
	result := FindOne(bson.M{"username": username},"login")

	if result != nil{
		//this means user already exists
		c.IndentedJSON(http.StatusForbidden,"user exists")
		return
	}

	//generating hash for password
	hasher := sha1.New()
    hasher.Write([]byte(password))
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	//making new user
	data := bson.D{
		{Key: "username", Value: username},
		{Key: "password", Value: sha},
	}

	err := InsertOne(data,"login")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"")
		return
	}
	c.IndentedJSON(http.StatusOK,"")
}

func HandleLogin(c *gin.Context){
	c.Request.ParseForm()

	//fetching data from POST body
	username := strings.Join(c.Request.Form["username"],"")
	password := strings.Join(c.Request.Form["password"],"")

	if username == "" || password == ""{
		c.IndentedJSON(http.StatusForbidden,"credentials missing")
		return
	}

	result := FindOne(bson.M{"username" : username}, "login")
	if result == nil{
		//user doesn't exist
		c.IndentedJSON(http.StatusForbidden,"user doesn't exist")
		return
	}

	c.IndentedJSON(http.StatusOK, "")
}