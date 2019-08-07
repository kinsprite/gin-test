package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI = "mongodb://my-root:my-root-pw@mongo:27017"

type numberValue struct {
	Found bool
	Name  string
	Value float64
}

func mongoDemo() numberValue {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Println("Error   Invalid mongo URI, error: ", err)
	}

	ctx, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConnect()
	err = client.Connect(ctx)

	if err != nil {
		log.Println("Error   Can't connect to mongo server, error: ", err)
	}

	ctx, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Println("Error   Can't ping to mongo server, error: ", err)
	} else {
		log.Println("INFO   ping to mongo server OK")
	}

	collection := client.Database("testing").Collection("numbers")

	if OK, value := findPI(collection); OK {
		return numberValue{true, "pi", value}
	}

	return insertPI(collection)
}

func findPI(collection *mongo.Collection) (bool, float64) {
	var result struct {
		Value float64
	}

	filter := bson.M{"name": "pi"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		log.Println("INFO   Can't find 'pi', error: ", err)
		return false, 0
	}

	log.Println("INFO   Found 'pi': ", result.Value)
	return true, result.Value
}

func insertPI(collection *mongo.Collection) numberValue {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value := 3.14159
	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": value})

	if err != nil {
		log.Println("ERROR   Insert 'pi', error: ", err)
		return numberValue{false, "", 0}
	}

	id := res.InsertedID
	log.Println("INFO   Inserted 'pi' ID: ", id)

	return numberValue{false, "pi", value}
}
