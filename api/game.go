package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetLatestGame(result []bson.M)(int,error){
	latest_ts_index := 0
	latest_ts,err := time.Parse(TIME_FORMAT,fmt.Sprint(result[0]["ts"]))
	if err != nil{
		return -1,err
	}
	for i:=0;i<len(result);i++{
		current_ts_string := fmt.Sprint(result[i]["ts"])
		current_ts, err := time.Parse(TIME_FORMAT,current_ts_string)
		if err != nil{
			return -1,err
		}

		if current_ts.After(latest_ts){
			latest_ts_index = i
			latest_ts = current_ts
		}
	}

	return latest_ts_index,nil
}

func HandleCheckUnsolved(c *gin.Context){
	login_status := CheckLogin(c)
	if !login_status{
		c.IndentedJSON(http.StatusForbidden,"not logged in")
		return
	}
	username, err := c.Cookie("Username")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,"unable to parse cookie")
	}

	result := FindMany(bson.M{"username": username}, "solves")
	if result == nil{
		//user has not generated any sudokus
		c.IndentedJSON(http.StatusNotFound, "no generated sudokus")
		return
	}
	
	latest_ts_index, err := GetLatestGame(result)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "trouble parsing game history")
		return
	}

	if result[latest_ts_index]["completed"] == true{
		c.IndentedJSON(http.StatusNotFound, "no unsolved sudokus")
		return
	}

	c.IndentedJSON(http.StatusOK, result[latest_ts_index]["current"])
}

func HandleGenSudoku(c *gin.Context){
	login_status := CheckLogin(c)
	if !login_status{
		c.IndentedJSON(http.StatusForbidden,"not logged in")
		return
	}

	username,err := c.Cookie("Username")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not parse cookie")
	}
	
	c.Request.ParseForm()
	difficulty,err := strconv.Atoi(strings.Join(c.Request.Form["difficulty"],""))
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "Try again later!")
		return
	}
	if difficulty != 0 && difficulty != 1 && difficulty != 2{
		c.IndentedJSON(http.StatusBadRequest,"Difficulty must be 0,1 or 2")
	}else{
		result := GenerateSudoku(difficulty)

		err = InsertOne(bson.D{
			{Key: "username", Value: username},
			{Key: "ts", Value: time.Now().Format(TIME_FORMAT)},
			{Key: "sudoku", Value: result["complete"]},
			{Key: "difficulty", Value: difficulty},
			{Key: "completed", Value: false},
			{Key: "current", Value: result["incomplete"]},
		},"solves")

		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,"trouble generating sudoku")
			return
		}

		c.IndentedJSON(http.StatusAccepted, result["incomplete"])
	}
}
