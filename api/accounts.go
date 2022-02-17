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
	username := strings.Join(c.Request.Form[USERNAME],"")
	password := strings.Join(c.Request.Form[PASSWORD],"")

	if username == "" || password == ""{
		//this means user already exists
		c.IndentedJSON(http.StatusForbidden,CREDENTIALS_MISSING)
		return
	}
	
	//checking if user already exists
	result := FindOne(bson.M{USERNAME: username},LOGIN)

	if result != nil{
		//this means user already exists
		c.IndentedJSON(http.StatusForbidden,USER_EXISTS)
		return
	}

	//generating hash for password
	hasher := sha1.New()
    hasher.Write([]byte(password))
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	//making new user
	data := bson.D{
		{Key: USERNAME, Value: username},
		{Key: PASSWORD, Value: sha},
	}

	err := InsertOne(data,LOGIN)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"")
		return
	}
	c.IndentedJSON(http.StatusOK,"")
}

func HandleLogin(c *gin.Context){
	c.Request.ParseForm()

	//fetching data from POST body
	username := strings.Join(c.Request.Form[USERNAME],"")
	password := strings.Join(c.Request.Form[PASSWORD],"")

	if username == "" || password == ""{
		c.IndentedJSON(http.StatusForbidden,CREDENTIALS_MISSING)
	}

	c.IndentedJSON(http.StatusOK, "")
}