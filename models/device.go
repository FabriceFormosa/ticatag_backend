package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Adress        string             `json:"adress" bson:"adress"`
	Latitude      string             `json:"latitude" bson:"latitude"`
	Longitude     string             `json:"longitude" bson:"longitude"`
	Adresspostale string             `json:"addresspostale" bson:"adresspostale"`
	CreatedAt     int64              `bson:"created_at" json:"created_at"`
}
