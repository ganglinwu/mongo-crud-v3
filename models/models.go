package models

import "go.mongodb.org/mongo-driver/v2/bson"

/*
type Employee struct {
	Id       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty" bson:"name,omitempty"`
	Age      int           `json:"age,string" bson:"age"`
	Salary   float32       `json:"salary,string" bson:"salary"`
	Position string        `json:"position,omitempty" bson:"position,omitempty"`
}
*/

// pointer version of Employee
// helps to check if field is empty
// by checking if we have a nil pointer
type Employee struct {
	Id       *bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     *string        `json:"name,omitempty" bson:"name,omitempty"`
	Age      *int           `json:"age,string" bson:"age"`
	Salary   *float32       `json:"salary,string" bson:"salary"`
	Position *string        `json:"position,omitempty" bson:"position,omitempty"`
}

// to build validator (option when creating new mongo collection)
var JsonSchema = bson.M{
	"bsonType": "object",
	"title":    "Employee object validation",
	"required": [2]string{"name", "position"},
	"properties": bson.M{
		"name": bson.M{
			"bsonType":    "string",
			"description": "name field must be type: string and non-empty",
		},
		"age": bson.M{
			"bsonType":    "int",
			"minimum":     0,
			"maximum":     120,
			"description": "age field must be an integer that is logically sound",
		},
		"salary": bson.M{
			"bsonType":    "double",
			"minimum":     0,
			"description": "salary field must be type: floating point and non-negative",
		},
		"position": bson.M{
			"bsonType":    "string",
			"description": "description field must be type: string and non-empty",
		},
	},
}
