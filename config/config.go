package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ganglinwu/mongo-crud-v3/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoInstance struct {
	collection *mongo.Collection
	client     *mongo.Client
}

var mgi mongoInstance

var employee models.Employee // temp, delete later

func ConnectDB() error {
	// load .env using godotenv package
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get env variables
	uri := os.Getenv("CONNECTION_STRING")
	database := os.Getenv("DATABASE_NAME")
	collection := os.Getenv("COLLECTION_NAME")

	// connection option to convert bson ObjectID as hex string
	bsonOpts := &options.BSONOptions{
		ObjectIDAsHexString: true,
	}

	// load bsonopts into mongo connection opts
	opts := options.Client().ApplyURI(uri).SetBSONOptions(bsonOpts)

	// connect to mongoDB
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Println("error connecting to mongoDB")
		return err
	}

	// check if collection already exists
	filter := bson.D{bson.E{Key: "name", Value: collection}}
	collectionNames, err := client.Database(database).ListCollectionNames(context.TODO(), filter)
	if err != nil {
		log.Println("error listing collection names")
		return err
	}

	// we expect collectionNames to have length 0 or 1
	// mongoDB atlas does not allow duplicate collection names
	if len(collectionNames) != 1 {
		fmt.Println("Collection does not exist. Proceeding to create new collection")
		validator := bson.M{
			"$jsonSchema": models.JsonSchema,
		}
		opts := options.CreateCollection().SetValidator(validator)
		client.Database(database).CreateCollection(context.TODO(), collection, opts)
	}

	// populate Client and Collection variables
	mgi.client = client
	mgi.collection = client.Database(database).Collection(collection)
	return nil
}

// exported function that fetches mongo Client pointer from config package
func GetClientPointer() *mongo.Client {
	return mgi.client
}

// exported function that fetches mongo Collection pointer from config package
func GetCollectionPointer() *mongo.Collection {
	return mgi.collection
}
