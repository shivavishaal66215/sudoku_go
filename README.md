## What is this repo?
This is an app I worked on while trying to learn Golang. Tonnes of optimisation can be done to the code and the front-end could be a lot prettier but that is not the main focus of this project so I will not be working on that.

## Setting up Sudoku_Go locally
1. Move into the `client/web` directory
2. Use `npm install` to install all the node packages
3. Move into the `api` directory
4. Use `go get .` to get all the go modules

## Running Sudoku_Go
1. Use `go run .` inside the `api` directory to start the server
2. Use `npm start` inside the `client/web` directory to start the React App.

_Note: Entry point for the api is `main.go` and the entry point for the client is `index.js`_
