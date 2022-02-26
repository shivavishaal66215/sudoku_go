package main

const(
	TIME_FORMAT = "2006-01-02T15:04:05.000Z"
)

type Logins struct{
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type Sessions struct{
	Username string `bson:"username"`
	AuthToken string `bson:"auth_token"`
	Ts string `bson:"ts"`
}

type Solves struct{
	Current [9][9]int `bson:"current"`
	Sudoku [9][9]int `bson:"sudoku"`
	Completed bool `bson:"completed"`
	Difficulty int `bson:"difficulty"`
	Ts string `bson:"ts"`
	Username string `bson:"username"`
}

type Document struct{
	Type string
	LoginData Logins
	SessionData Sessions
	SolvesData Solves
}


