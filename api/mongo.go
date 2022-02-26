package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc){
	defer cancel()
	defer func (){
		if err := client.Disconnect(ctx); err != nil{
			panic(err)
		}
	}()
}

func Connect()(*mongo.Client,context.Context,context.CancelFunc, error){
	ctx, cancel := context.WithTimeout(context.Background(),
                                       30 * time.Second)
     
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    return client, ctx, cancel, err
}

//will throw error if no document is found
func FindOneMongo(query bson.M, collection string) (Document,error){
	client,ctx,cancel,err := Connect()
	defer Close(client,ctx,cancel)
	if err != nil{
		return Document{},err
	}

	mongoResult := client.Database("sudoku_go").Collection(collection).FindOne(ctx, query)
	err = mongoResult.Err()
	if err != nil{
		return Document{},err
	}

	var doc Document
	if collection == "logins"{
		var login Logins
		mongoResult.Decode(&login)
		doc = Document{
			Type: "logins",
			LoginData: login,
		}
	}else if collection == "sessions"{
		var session Sessions
		mongoResult.Decode(&session)
		doc = Document{
			Type: "sessions",
			SessionData: session,
		}
	}else if collection == "solves"{
		var solve Solves
		mongoResult.Decode(&solve)
		doc = Document{
			Type: "solves",
			SolvesData: solve,
		}
	}else{
		return Document{},errors.New("invalid collection name")
	}
	return doc,nil
}

func FindManyMongo(query bson.M, collection string) ([]Document,error){
	client,ctx,cancel,err := Connect()
	defer Close(client,ctx,cancel)
	if err != nil{
		return []Document{},err
	}

	cursor, err := client.Database("sudoku_go").Collection(collection).Find(ctx,query)
	if err != nil{
		return []Document{},err
	}

	if collection == "logins"{
		//findmany function is not required for this right now
		return []Document{}, nil
	}else if collection == "sessions"{
		//findmany function is not required for this right now
		return []Document{}, nil
	}else if collection == "solves"{
		var solves []Solves
		err = cursor.All(ctx, &solves)
		if err != nil{
			return []Document{},err
		}

		var doc []Document
		n := len(solves)
		for i:=0;i<n;i++{
			doc = append(doc, Document{
				Type: "solves",
				SolvesData: solves[i],
			})
		}

		return doc,nil
	}else{
		return []Document{},errors.New("invalid collection name")
	}
}

func InsertMongo(doc Document) error{
	client,ctx,cancel,err := Connect()
	defer Close(client,ctx,cancel)
	if err != nil{
		return err
	}

	if doc.Type == "logins"{
		_,err = client.Database("sudoku_go").Collection("logins").InsertOne(ctx,doc.LoginData)
	}else if doc.Type == "sessions"{
		_,err = client.Database("sudoku_go").Collection("sessions").InsertOne(ctx,doc.SessionData)
	}else if doc.Type == "solves"{
		_,err = client.Database("sudoku_go").Collection("solves").InsertOne(ctx,doc.SolvesData)
	}else{
		return errors.New("invalid type for document")
	}

	if err != nil{
		return err
	}
	return nil
}

func DeleteOneMongo(doc Document) error{
	client,ctx,cancel,err := Connect()
	defer Close(client,ctx,cancel)
	if err != nil{
		return err
	}

	if doc.Type == "logins"{
		//deleting this type of document is not required now
		return nil
	}else if doc.Type == "solves"{
		//deleting this type of document is not required now
		return nil
	}else if doc.Type == "sessions"{
		_,err = client.Database("sudoku_go").Collection("sessions").DeleteOne(ctx, bson.M{"username": doc.SessionData.Username, "auth_token": doc.SessionData.AuthToken})
		if err != nil{
			return err
		}
	}else{
		return errors.New("document type not supported")
	}

	return nil
}

func DeleteManyMongo(doc Document) error{
	client,ctx,cancel,err := Connect()
	defer Close(client,ctx,cancel)
	if err != nil{
		return err
	}

	if doc.Type == "logins"{
		//deleting this type of document is not required now
		return nil
	}else if doc.Type == "solves"{
		//deleting this type of document is not required now
		return nil
	}else if doc.Type == "sessions"{
		_,err = client.Database("sudoku_go").Collection("sessions").DeleteMany(ctx, bson.M{"username": doc.SessionData.Username})
		if err != nil{
			return err
		}
	}else{
		return errors.New("document type not supported")
	}

	return nil
}

func UpdateOneMongo(query bson.M, collection string, newData Document) error{
	client,ctx,cancel,err := Connect()
	defer Close(client,ctx,cancel)
	if err != nil{
		return err
	}

	if collection == "logins"{
		//update for this type is not required right now
		return nil
	}else if collection == "sessions"{
		//update for this type is not required right now
		return nil
	}else if collection == "solves"{
		_,err := client.Database("sudoku_go").Collection("solves").UpdateOne(
			ctx,query,bson.D{{Key: "$set", Value: newData.SolvesData}},
		)
		if err != nil{
			return err
		}
	}

	return nil
}