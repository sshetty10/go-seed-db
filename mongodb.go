package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type MongoDb struct {
	MongoCli *mongo.Client
	MongoDb  *mongo.Database
}

type MTrainer struct {
	ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
	Age  int                `bson:"age" json:"age"`
	City string             `bson:"city" json:"city"`
}

func (d *MongoDb) GetTrainers() ([]*MTrainer, error) {
	result := []*MTrainer{}

	collection := d.MongoDb.Collection("trainers")
	filter := bson.M{}

	ctx := context.TODO()
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var t MTrainer
		err := cur.Decode(&t)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, &t)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
	return result, nil
}

func (d *MongoDb) GetTrainerByName(name string) (*MTrainer, error) {
	var result MTrainer

	collection := d.MongoDb.Collection("trainers")
	filter := bson.M{"name": name}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found documents: %+v\n", result)
	return &result, nil
}

func (d *MongoDb) GetTrainerByID(id string) (*MTrainer, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result MTrainer
	collection := d.MongoDb.Collection("trainers")
	filter := bson.M{"_id": objID}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found a single document: %+v\n", result)
	return &result, nil
}

func (d *MongoDb) AddTrainer(t *MTrainer) error {
	t.ID = primitive.NewObjectID()
	collection := d.MongoDb.Collection("trainers")
	result, err := collection.InsertOne(context.TODO(), t)
	if err != nil {
		return err
	}

	t.ID = result.InsertedID.(primitive.ObjectID)
	fmt.Println("Inserted multiple documents: ", t.ID.Hex())
	return nil
}

func (d *MongoDb) UpdateTrainer(t *MTrainer) error {
	filter := bson.M{"_id": t.ID}
	collection := d.MongoDb.Collection("trainers")

	update := bson.M{
		"$set": bson.M{
			"age":  t.Age,
			"name": t.Name,
			"city": t.City,
		}}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}

func (d *MongoDb) DeleteTrainer(id string) error {
	collection := d.MongoDb.Collection("trainers")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil
}
