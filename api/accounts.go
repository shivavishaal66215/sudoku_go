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
		//this means credentials are missing
		c.IndentedJSON(http.StatusForbidden,"credentials missing")
		return
	}

	//checking if user already exists
	_,err := FindOneMongo(bson.M{"username": username},"logins")
	if err == nil{
		//user exists or connection to mongo failed. Either way, return user already exists
		c.IndentedJSON(http.StatusForbidden, "user exists")
		return
	}

	//generating hash for password
	hasher := sha1.New()
    hasher.Write([]byte(password))
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	doc := Document{
		Type: "logins",
		LoginData: Logins{
			Username: username,
			Password: sha,
		},
	}

	err = InsertMongo(doc)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"could not create user")
		return
	}
	c.IndentedJSON(http.StatusOK,"user created")
}

func HandleLogin(c *gin.Context){
	c.Request.ParseForm()
	current_time := time.Now().String()

	//fetching data from POST body
	username := strings.Join(c.Request.Form["username"],"")
	password := strings.Join(c.Request.Form["password"],"")

	if username == "" || password == ""{
		c.IndentedJSON(http.StatusForbidden,"credentials missing")
		return
	}

	mongoResult,err := FindOneMongo(bson.M{"username" : username}, "logins")
	if err != nil{
		//user doesnt exist
		c.IndentedJSON(http.StatusNotFound, "user not found")
		return
	}

	stored_hash := mongoResult.LoginData.Password

	//generating hash for current password
	hasher := sha1.New()
    hasher.Write([]byte(password))
    current_hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if stored_hash != current_hash{
		//password entered is invalid
		c.IndentedJSON(http.StatusForbidden,"incorrect username or password")
		return
	}

	//delete all existing sessions to prevent user from submitting on more than 1 device
	DeleteManyMongo(Document{
		Type: "sessions",
		SessionData: Sessions{
			Username: username,
		},
	})

	//generate session token
	hasher = sha1.New()
	hasher.Write([]byte(current_time))
	auth_token := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	doc := Document{
		Type: "sessions",
		SessionData: Sessions{
			Username: username,
			AuthToken: auth_token,
			Ts: time.Now().Format(TIME_FORMAT),
		},
	}
	err = InsertMongo(doc)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"could not create new session")
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

	_,err = FindOneMongo(bson.M{"username" : username,"auth_token":auth_token},"sessions")
	return err == nil
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
		c.IndentedJSON(http.StatusInternalServerError,"auth token missing")
		return
	}
	username, err := c.Cookie("Username")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"username missing")
		return
	}

	_,err = FindOneMongo(bson.M{"username": username, "auth_token": auth_token}, "sessions")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"not logged in")
		return
	}

	err = DeleteManyMongo(Document{
		Type: "sessions",
		SessionData: Sessions{
			Username: username,
		},
	})
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"could not log you out")
		return
	}

	c.IndentedJSON(http.StatusOK,"logged out")
}