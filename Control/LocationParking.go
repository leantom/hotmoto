package Control

import "gopkg.in/mgo.v2/bson"

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document

type location struct {
	Coordinates []int `bson:"coordinates" json:"coordinates"`
	Type_Location	string `bson:"type" json:"type"`
}

type LocationParking struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Location     location  `bson:"location" json:"location"`
	Name  		string        `bson:"name" json:"name"`
}

