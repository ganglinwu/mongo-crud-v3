package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ganglinwu/mongo-crud-v3/config"
	"github.com/ganglinwu/mongo-crud-v3/routes"
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
	router := http.NewServeMux()

	// add routes
	router.HandleFunc("GET /employees", routes.GetAllEmployees)
	router.HandleFunc("POST /employees", routes.InsertEmployee)
	router.HandleFunc("GET /employees/{id}", routes.GetEmployeeByID)
	router.HandleFunc("DELETE /employees/{id}", routes.DropEmployeeByID)
	router.HandleFunc("PATCH /employees/{id}", routes.PatchEmployeeByID)

	// initialize server with opts
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// listen and serve
	fmt.Println("Server listening..")
	log.Fatal(server.ListenAndServe())
}
