package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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
		c.IndentedJSON(http.StatusNoContent, "no unsolved sudokus")
		return
	}


	latest_ts_index := 0
	latest_ts := time.Now()
	for i:=0;i<len(result);i++{
		current_ts_string := fmt.Sprint(result[i]["ts"])
		current_ts, err := time.Parse(time.Now().String(),current_ts_string)

		if err != nil{
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError,"could not fetch data")
			return
		}

		if current_ts.After(latest_ts){
			latest_ts_index = i
			latest_ts = current_ts
		}
	}

	fmt.Println(latest_ts_index)

	c.IndentedJSON(http.StatusOK,"")
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
		data,err := json.Marshal(result["complete"])
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, "could not marshal array")
			return
		}

		err = InsertOne(bson.D{
			{Key: "username", Value: username},
			{Key: "ts", Value: time.Now().String()},
			{Key: "sudoku", Value: data},
			{Key: "difficulty", Value: difficulty},
			{Key: "completed", Value: false},
		},"solves")

		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,"trouble generating sudoku")
			return
		}

		c.IndentedJSON(http.StatusAccepted, result["incomplete"])
	}
}
