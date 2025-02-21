package config

import (
	"log"
	"os"

	"github.com/ganglinwu/mongo-crud-v3/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Collection *mongo.Collection
	Client     *mongo.Client
)

var employee models.Employee // temp, delete later

func ConnectDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	uri := os.Getenv("CONNECTION_STRING")
	database := os.Getenv("DATABASE_NAME")
	collection := os.Getenv("COLLECTION_NAME")

	bsonOpts := &options.BSONOptions{
		ObjectIDAsHexString: true,
	}

	opts := options.Client().ApplyURI(uri).SetBSONOptions(bsonOpts)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Println("error connecting to mongoDB")
		return err
	}

	Client = client
	Collection = Client.Database(database).Collection(collection)
	return nil
}

func GetClientPointer() *mongo.Client {
	return Client
}

func GetCollectionPointer() *mongo.Collection {
	return Collection
}
