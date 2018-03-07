package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"fmt"
)


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
	session, err := mgo.Dial("ec2-52-55-50-216.compute-1.amazonaws.com")
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	fmt.Println("Mongo server connected")
	db = session.DB("test")
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




var locationParkingDAO = LocationParkingDAO{}


// Fetch Example

func LocationFisrtParking(w http.ResponseWriter, r *http.Request) {

	res, err := locationParkingDAO.FindAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, res)
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func init() {
	Connect()
}

func main() {
	fmt.Print("Hello world")
	r := mux.NewRouter()

	r.HandleFunc("/parkings", LocationFisrtParking).Methods("GET")

}