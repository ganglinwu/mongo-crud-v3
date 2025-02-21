package controllers

import (
	"context"
	"log"

	"github.com/ganglinwu/mongo-crud-v3/config"
	"github.com/ganglinwu/mongo-crud-v3/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var collection = config.Collection

func GetAllEmployees() ([]models.Employee, error) {
	// search filter
	filter := bson.D{bson.E{}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("error in collection.Find")
		return nil, err
	}

	employees := []models.Employee{}

	err = cursor.All(context.TODO(), employees)
	if err != nil {
		log.Println("error in cursor.All")
		return nil, err
	}
	return employees, nil
}
