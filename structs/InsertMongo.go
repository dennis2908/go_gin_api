package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

type InsertCustomer struct {
	IdCustomer string `bson:"title"`
	Operation  string `bson:"content"`
}

type PostInsertCustomer struct {
	Id         primitive.ObjectID `bson:"_id"`
	IdKonsumen string             `bson:"title"`
	Operation  string             `bson:"content"`
}
