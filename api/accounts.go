package main

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

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
	//at this point, all necessary checks have passed
	current_time := time.Now().String()

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

	stored_hash := result["password"]

	//generating hash for current password
	hasher := sha1.New()
    hasher.Write([]byte(password))
    current_hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if stored_hash != current_hash{
		//password entered is invalid
		c.IndentedJSON(http.StatusForbidden,"incorrect username or password")
		return
	}

	//generate session token
	hasher = sha1.New()
	hasher.Write([]byte(current_time))
	auth_token := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	err := UpdateSession(auth_token, username)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"")
		return
	}
	
	c.SetCookie("AuthToken", auth_token, 0, "/", "localhost",true,true)
	c.SetCookie("Username", username, 0, "/", "localhost",true,true)

	c.IndentedJSON(http.StatusOK, "")
}

func CheckLogin(c *gin.Context) bool{
	auth_token,err := c.Cookie("AuthToken")
	if err != nil{
		return false
	}

	username,err := c.Cookie("Username")
	if err != nil{
		return false
	}

	result := FindOne(bson.M{"username" : username},"login")
	if result == nil{
		//user is not logged in
		return false
	}

	return result["auth_token"] == auth_token
}

func HandleCheckLogin(c *gin.Context){
	login_status := CheckLogin(c)
	if login_status{
		c.IndentedJSON(http.StatusOK,"logged in")
	}else{
		c.IndentedJSON(http.StatusForbidden,"not logged in")
	}
}

func HandleLogout(c *gin.Context){
	auth_token, err := c.Cookie("AuthToken")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"could not log you out")
		return
	}
	username, err := c.Cookie("Username")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"could not log you out")
		return
	}

	result := FindOne(bson.M{"username": username}, "login")
	if result["auth_token"] != auth_token{
		c.IndentedJSON(http.StatusInternalServerError,"could not log you out")
		return
	}

	//replace with fake auth_token
	current_time := time.Now().String()
	hasher := sha1.New()
	hasher.Write([]byte(current_time))
	auth_token = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	err = UpdateSession(auth_token,username)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"could not log you out")
		return
	}

	c.IndentedJSON(http.StatusOK,"logged out")
}