package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ganglinwu/mongo-crud-v3/config"
)

func main() {
	// connect to DB
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongoDB")

	// defer closing of DB connection
	client := config.GetClientPointer()
	defer client.Disconnect(context.TODO())

	// initialize router, consider only using "net/http"
	httpHandler := new(http.Handler)

	// add routes

	// listen and serve
	log.Fatal(http.ListenAndServe(":8080", *httpHandler))
}
