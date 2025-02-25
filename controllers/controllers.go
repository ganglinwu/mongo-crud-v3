package controllers

import (
	"context"
	"log"

	"github.com/ganglinwu/mongo-crud-v3/config"
	"github.com/ganglinwu/mongo-crud-v3/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func FindAllEmployees() ([]models.Employee, error) {
	// search filter
	filter := bson.D{bson.E{}}

	collection := config.GetCollectionPointer()
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("error in collection.Find")
		log.Println("msg: ", err.Error())
		return nil, err
	}

	employees := []models.Employee{}

	err = cursor.All(context.TODO(), &employees)
	if err != nil {
		log.Println("error in cursor.All")
		return nil, err
	}
	return employees, nil
}

func FindEmployeeByID(ID string) (interface{}, error) {
	// convert ID string to ObjectID
	objectID, err := bson.ObjectIDFromHex(ID)
	if err != nil {
		log.Println("error while converting ID to ObjectID")
		log.Println("errMsg: ", err.Error())
		return nil, err
	}

	// search filter
	filter := bson.D{{Key: "_id", Value: objectID}}

	var employee models.Employee

	collection := config.GetCollectionPointer()
	err = collection.FindOne(context.TODO(), filter).Decode(&employee)
	if err != nil {
		log.Println("error in collection.FindOne")
		log.Println("errMsg: ", err.Error())
		return nil, err
	}
	return employee, nil
}

func CreateEmployee(emp models.Employee) (interface{}, error) {
	collection := config.GetCollectionPointer()

	result, err := collection.InsertOne(context.TODO(), emp)
	if err != nil {
		log.Println("error in collection.InsertOne")
		log.Println("msg: ", err.Error())
		return nil, err
	}

	return result, nil
}

func DeleteEmployeeByID(ID string) (interface{}, error) {
	collection := config.GetCollectionPointer()

	objectID, err := bson.ObjectIDFromHex(ID)
	if err != nil {
		log.Println("error while converting ID to ObjectID")
		log.Println("errMsg: ", err.Error())
		return nil, err
	}

	result := collection.FindOneAndDelete(context.TODO(), bson.D{bson.E{Key: "_id", Value: objectID}})

	// NOTE: FindOneAndDelete does not return error.
	// we must check for .Err under result!
	if result.Err() != nil {
		return nil, result.Err()
	} else {
		return result, nil
	}
}

func UpdateEmployeeByID(ID string, update bson.D) (interface{}, error) {
	collection := config.GetCollectionPointer()

	objectID, err := bson.ObjectIDFromHex(ID)
	if err != nil {
		log.Println("error while converting ID to ObjectID")
		log.Println("errMsg: ", err.Error())
		return nil, err
	}

	result := collection.FindOneAndUpdate(context.TODO(), bson.D{bson.E{Key: "_id", Value: objectID}}, update)

	// NOTE: FindOneAndUpdate does not return error.
	// we must check for .Err under result!
	if result.Err() != nil {
		return nil, result.Err()
	} else {
		return result, nil
	}
}
