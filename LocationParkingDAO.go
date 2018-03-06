package main

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type LocationParkingDAO struct {
	Server   string
	Database string
}

var db *mgo.Database


const (
	COLLECTION = "LocationParking"
	DB = "location_db"
	MongoDBHosts = "ds251598.mlab.com"
)

// Establish a connection to database

func Connect() {

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: DB,
		Username: "admin",
		Password: "Quang12345@",
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	log.Print("connect database")
	db = mongoSession.DB(mongoDBDialInfo.Database)
}

// Find list of movies
func (m *LocationParkingDAO) FindAll() ([]LocationParkingDAO, error) {
	var movies []LocationParkingDAO
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	return movies, err
}

func (m *LocationParkingDAO) Insert(res *LocationParking)  error {
	err := db.C(COLLECTION).Insert(res)
	return err
}

func (m *LocationParkingDAO) Update(res *LocationParking) error {
	err := db.C(COLLECTION).UpdateId(res.ID, &res)
	return err
}


func (m *LocationParkingDAO) FindNearLocationParking() ([]LocationParkingDAO, error) {
	var restaurants []LocationParkingDAO

	// search criteria
	long := -73.8601152
	lat := 	40.7311739

	scope := 3000 // max distance in metres

	collection := db.C(COLLECTION)

	err := collection.Find(bson.M{
		"location" : bson.M{
			"$nearSphere":bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": scope,
			},
		},
	}).Limit(10).All(&restaurants)


	return restaurants, err
}


