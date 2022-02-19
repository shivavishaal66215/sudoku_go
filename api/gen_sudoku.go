package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func randomItemFromMap(m map[int]int) int{
	allowed := make([]int,0)

	for i:=1;i<=9;i++{
		if m[i] == 0{
			allowed = append(allowed, i)
		}
	}
	r := rand.Int() % len(allowed)
	return allowed[r]
}

func generateDiagonalBlock() [3][3]int{
	arr := [3][3]int{}
	i,j := 0,0 // i->row, j->col
	visited := make(map[int]int)
	r := 0

	for i = 0;i < 3;{
		if j >= 3{
			j = 0
			i += 1
			continue
		}
		r = randomItemFromMap(visited)
		arr[i][j] = r
		visited[r] = 1
		j += 1
	}
	return arr
}


func fillMainDiagonal(arr *[9][9]int){
	blocks := [3][3][3]int{}

	for i:=0;i<3;i++{
		blocks[i] = generateDiagonalBlock()
	}

	for k:=0;k<9;k+=3{
		for i:=k;i<k+3;i++{
			for j:=k;j<k+3;j++{
				arr[i][j] = blocks[k/3][i-k][j-k]
			}
		}
	}
}

func isEmpty(arr [9][9]int) bool{
	for i:=0;i<9;i++{
		for j:=0;j<9;j++{
			if arr[i][j] == 0{
				return true
			}
		}
	}
	
	return false
}

func checkSafe(arr *[9][9]int, row int, col int, n int) bool{
	i := 0
	j := 0

	//row check	
	for j=0;j<9;j++{
		if arr[row][j] == n{
			return false
		}
	}

	//col check
	for i=0;i<9;i++{
		if arr[i][col] == n{
			return false
		}
	}

	//block check
	block_i := 3 * (row/3)
	block_j := 3 * (col/3)

	for i:=block_i;i<block_i+3;i++{
		for j:=block_j;j<block_j+3;j++{
			if arr[i][j] == n{
				return false
			}
		}
	}

	return true
}

func fillRest(arr[9][9]int, row int, col int, result *[9][9]int){
	if !isEmpty(*result) {
		return
	}

	if col >= 9{
		col = 0
		row += 1
	}

	if row >= 9{
		for i:=0;i<9;i++{
			for j:=0;j<9;j++{
				result[i][j] = arr[i][j]
			}
		}
		return
	}

	if arr[row][col] != 0{
		fillRest(arr,row,col+1,result)
	}

	for i:=1;i<=9;i++{
		if checkSafe(&arr, row, col, i){
			arr[row][col] = i
			fillRest(arr,row,col+1,result)
			arr[row][col] = 0
		}
	}
}

func unSolve(arr[9][9]int, difficulty int) [9][9]int{
	result := [9][9]int{}
	difficulty_map := map[int]float64{
		0 : 0.8,
		1 : 0.6,
		2 : 0.4,
	}
	cur := difficulty_map[difficulty]
	for i:=0;i<9;i++{
		for j:=0;j<9;j++{
			r := rand.Float64()
			if r > cur{
				result[i][j] = 0
			}else{
				result[i][j] = arr[i][j]
			}
		}
	}
	return result
}

func generateSudoku(difficulty int) (map[string][9][9]int){
	arr := [9][9]int{}
	result := [9][9]int{}

	fillMainDiagonal(&arr)
	fillRest(arr, 0, 0, &result)

	message := map[string][9][9]int{}
	message["complete"] = result
	message["incomplete"] = unSolve(result,difficulty)

	return message
}

func HandleGenSudoku(c *gin.Context){
	login_status := CheckLogin(c)
	if !login_status{
		c.IndentedJSON(http.StatusForbidden,"not logged in")
		return
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
		c.IndentedJSON(http.StatusAccepted, generateSudoku(difficulty))
	}
}
