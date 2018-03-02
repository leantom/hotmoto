package Control

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	"hotMoto/Model"
)

type LocationParkingDAO struct {
	Server   string
	Database string
}

var db *mgo.Database


const (
	COLLECTION = "LocationParking"
)

// Establish a connection to database

func (m *LocationParkingDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("connect database")
	db = session.DB(m.Database)
}

// Find list of movies
func (m *LocationParkingDAO) FindAll() ([]LocationParkingDAO, error) {
	var movies []LocationParkingDAO
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	return movies, err
}

func (m *LocationParkingDAO) Insert(res *Model.LocationParking)  error {
	err := db.C(COLLECTION).Insert(res)
	return err
}

func (m *LocationParkingDAO) Update(res *Model.LocationParking) error {
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


