package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckSudokuMatch(a [9][9]int, b [9][9]int) bool{
	for i:=0;i<9;i++{
		for j:=0;j<9;j++{
			if a[i][j] != b[i][j]{
				return false
			}
		}
	}

	return true
}

func GetLatestGame(result []Document)(int,error){
	if len(result) == 0{
		return -1, errors.New("size is 0")
	}
	latest_ts_index := 0
	latest_ts,err := time.Parse(TIME_FORMAT,fmt.Sprint(result[0].SolvesData.Ts))
	if err != nil{
		return -1,err
	}
	for i:=0;i<len(result);i++{
		current_ts_string := fmt.Sprint(result[i].SolvesData.Ts)
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

	result,err := FindManyMongo(bson.M{"username": username}, "solves")
	if err != nil{
		c.IndentedJSON(http.StatusNotFound, "no generated sudokus")
		return
	}

	latest_ts_index, err := GetLatestGame(result)
	if err != nil{
		c.IndentedJSON(http.StatusNotFound, "trouble parsing game history")
		return
	}

	if result[latest_ts_index].SolvesData.Completed{
		c.IndentedJSON(http.StatusNotFound, "no unsolved sudokus")
		return
	}

	c.IndentedJSON(http.StatusOK, result[latest_ts_index].SolvesData.Current)
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

		err = InsertMongo(Document{
			Type: "solves",
			SolvesData: Solves{
				Username: username,
				Ts: time.Now().Format(TIME_FORMAT),
				Current: result["incomplete"],
				Sudoku: result["complete"],
				Difficulty: difficulty,
				Completed: false,
			},
		})

		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,"trouble generating sudoku")
			return
		}

		c.IndentedJSON(http.StatusAccepted, result["incomplete"])
	}
}

func HandleSaveSudoku(c *gin.Context){
	login_status := CheckLogin(c)
	if !login_status{
		c.IndentedJSON(http.StatusForbidden, "not logged in")
		return
	}

	username,err := c.Cookie("Username")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not parse cookie")
		return
	}

	c.Request.ParseForm()
	result,err := FindManyMongo(bson.M{"username": username},"solves")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not fetch history")
		return
	}

	var passedSudoku [9][9]int
	err = json.Unmarshal([]byte(c.Request.Form["sudoku"][0]),&passedSudoku)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not parse payload")
		return
	}

	latest_index,err := GetLatestGame(result)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "trouble fetching history")
		return
	}

	latest := result[latest_index]
	err = UpdateOneMongo(bson.M{"username": username, "ts": latest.SolvesData.Ts}, "solves", Document{
		Type: "solves",
		SolvesData: Solves{
			Current: passedSudoku,
			Sudoku: latest.SolvesData.Sudoku,
			Ts: latest.SolvesData.Ts,
			Username: latest.SolvesData.Username,
			Completed: latest.SolvesData.Completed,
			Difficulty: latest.SolvesData.Difficulty,
		},
	})
	if err != nil{
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "could not save changes")
		return
	}

	c.IndentedJSON(http.StatusOK, "Saved")
}

func HandleSubmitSudoku(c *gin.Context){
	login_status := CheckLogin(c)
	if !login_status{
		c.IndentedJSON(http.StatusForbidden, "not logged in")
		return
	}

	username,err := c.Cookie("Username")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not parse cookie")
		return
	}

	c.Request.ParseForm()
	result,err := FindManyMongo(bson.M{"username": username},"solves")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not fetch history")
		return
	}

	var passedSudoku [9][9]int
	err = json.Unmarshal([]byte(c.Request.Form["sudoku"][0]),&passedSudoku)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not parse payload")
		return
	}

	latest_index,err := GetLatestGame(result)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "trouble fetching history")
		return
	}

	latest := result[latest_index]

	//check if submitted sudoku and the correct sudoku are the same
	match := CheckSudokuMatch(passedSudoku,latest.SolvesData.Sudoku)
	if match{
		//passed sudoku matches with the correct sudoku
		//set this sudoku as completed
		err = UpdateOneMongo(bson.M{"username": username, "ts": latest.SolvesData.Ts}, "solves", Document{
			Type: "solves",
			SolvesData: Solves{
				Current: passedSudoku,
				Sudoku: latest.SolvesData.Sudoku,
				Ts: latest.SolvesData.Ts,
				Username: latest.SolvesData.Username,
				Completed: true,
				Difficulty: latest.SolvesData.Difficulty,
			},
		})

		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, "could not make changes to database")
			return
		}
		
		c.IndentedJSON(http.StatusOK,"successfully submitted")

	}else{
		c.IndentedJSON(http.StatusForbidden,"sudokus don't match")
		return
	}
}

func HandleGetStats(c *gin.Context){
	login_status := CheckLogin(c)
	if !login_status{
		c.IndentedJSON(http.StatusForbidden, "not logged in")
		return
	}

	username,err := c.Cookie("Username")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not parse cookie")
		return
	}

	result,err := FindManyMongo(bson.M{"username": username},"solves")
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, "could not fetch history")
		return
	}

	easy := 0
	medium := 0
	hard := 0

	for i:=0;i<len(result);i++{
		if result[i].SolvesData.Completed{
			if result[i].SolvesData.Difficulty == 0{
				easy += 1
			}else if result[i].SolvesData.Difficulty == 1{
				medium += 1
			}else{
				hard += 1
			}
		}
	}

	c.IndentedJSON(http.StatusOK, map[string]int{
		"easy": easy,
		"medium": medium,
		"hard": hard,
	})
}
