package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ganglinwu/mongo-crud-v3/controllers"
	"github.com/ganglinwu/mongo-crud-v3/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	employees, err := controllers.FindAllEmployees()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error(), "msg": "error in FindAllEmployees"})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employees)
}

func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	employee, err := controllers.FindEmployeeByID(id)
	// err handling
	if err != nil {
		// send error in text form
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// no such document exists in mongoDB
		if err.Error()[0:19] == "mongo: no documents" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("ID not found"))
			fmt.Println(err.Error())
		} else {

			// probably not user issue
			w.WriteHeader(http.StatusInternalServerError)
			errMsg := fmt.Sprintf("Something went wrong. Error: %s", err.Error())
			w.Write([]byte(errMsg))
		}
		return
	}

	// success and send result as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

func InsertEmployee(w http.ResponseWriter, r *http.Request) {
	var createdEmp models.Employee
	json.NewDecoder(r.Body).Decode(&createdEmp)
	defer r.Body.Close()

	insertedEmp, err := controllers.CreateEmployee(createdEmp)
	// error handling
	// with schema validation error handling is more complex
	// it is now possible for user to send json with fields missing
	// we check by asserting err to be type mongo.WriteException
	if err != nil {
		_, ok := err.(mongo.WriteException)
		if ok {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/html ; charset=utf-8")
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedEmp)
}

func DropEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := controllers.DeleteEmployeeByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "text/html ; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func PatchEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// check if employee to be updated even exists in mongoDB
	currentEmpRecord, err := controllers.FindEmployeeByID(id)
	if err != nil {
		// send error in text form
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// no such document exists in mongoDB
		if err.Error()[0:19] == "mongo: no documents" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("ID not found"))
			fmt.Println(err.Error())
		} else {

			// probably not user issue
			w.WriteHeader(http.StatusInternalServerError)
			errMsg := fmt.Sprintf("Something went wrong. Error: %s", err.Error())
			w.Write([]byte(errMsg))
		}
		return
	}
	var currentEmployee, empToUpdate models.Employee

	// type assert currentEmpRecord to models.Employee
	currentEmployee, ok := currentEmpRecord.(models.Employee)
	if !ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong while fetching the employee details"))
		return
	}

	json.NewDecoder(r.Body).Decode(&empToUpdate)

	updateValue := bson.M{}

	if empToUpdate.Name == nil {
		updateValue["name"] = currentEmployee.Name
	} else {
		updateValue["name"] = empToUpdate.Name
	}
	if empToUpdate.Salary == nil {
		updateValue["salary"] = currentEmployee.Salary
	} else {
		updateValue["salary"] = empToUpdate.Salary
	}
	if empToUpdate.Age == nil {
		updateValue["age"] = currentEmployee.Age
	} else {
		updateValue["age"] = empToUpdate.Age
	}
	if empToUpdate.Position == nil {
		updateValue["position"] = currentEmployee.Position
	} else {
		updateValue["position"] = empToUpdate.Position
	}

	update := bson.D{{
		Key:   "$set",
		Value: updateValue,
	}}

	deleteResult, err := controllers.UpdateEmployeeByID(id, update)
	if err != nil {
		w.Header().Set("Content-Type", "text/html ; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deleteResult)
}
