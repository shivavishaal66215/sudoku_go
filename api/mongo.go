package main

import (
	"context"
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

func InsertOne(data bson.D, collection string) error{
	client,ctx,cancel,err := Connect()
	defer cancel()

	if err != nil{
		return err
	}

	_, err = client.Database("sudoku_go").Collection(collection).InsertOne(ctx,data)
	if err != nil{
		return err
	}

	Close(client,ctx,cancel)

	return nil
}

func FindOne(query bson.M, collection string) (bson.M){
	client,ctx,cancel,err := Connect()
	defer cancel()

	if err != nil{
		return nil
	}

	var result bson.M
	err = client.Database("sudoku_go").Collection(collection).FindOne(ctx,query).Decode(&result)

	if err != nil{
		return nil
	}

	Close(client,ctx,cancel)

	return result
}

func UpdateSession(auth_token string, username string) error{
	client,ctx,cancel,err := Connect()
	defer cancel()

	if err != nil{
		return err
	}

	_,err = client.Database("sudoku_go").Collection("login").UpdateOne(
		ctx,
		bson.M{"username":username},
		bson.D{
			{Key: "$set", Value: bson.D{{Key:"auth_token", Value:auth_token}}},
	})

	if err != nil{
		return err
	}

	return nil
}