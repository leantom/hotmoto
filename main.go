package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	//"github.com/jpillora/overseer"
	//"github.com/jpillora/overseer/fetcher"

	"fmt"
	"log"
	//"time"
)


// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document

type location struct {

	Position  []int `bson:"position" json:"position"`
	Name string `bson:"name" json:"name"`
	Type_Location	string `bson:"type" json:"type"`

}

type LocationParking struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Position  []int `bson:"position" json:"position"`
	Name  		string        `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Phone string `bson:"phone" json:"phone"`
	Cost string `bson:"cost" json:"cost"`
	Total int `bson:"total" json:"total"`

}


type LocationParkingDAO struct {
	Server   string
	Database string
}

var db *mgo.Database


const (
	COLLECTION = "motopark"
	DB = "hotmoto_db"
)

// Establish a connection to database

func Connect() {

	// We need this object to establish a session to our MongoDB.
	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	fmt.Println("Mongo server connected")
	db = session.DB(DB)

}

// Find list of movies
func (m *LocationParkingDAO) FindAll() ([]LocationParking, error) {
	var movies []LocationParking
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
		fmt.Print(res)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, res)
}


func UpdateParking(w http.ResponseWriter, r *http.Request) {
	defer  r.Body.Close()

	var parking LocationParking
	err := 	json.NewDecoder(r.Body).Decode(&parking);
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	parking.ID = bson.NewObjectId()

	if err := locationParkingDAO.Insert(&parking); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, parking)

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

func findALL() ([]LocationParking, error) {
	var results []LocationParking

	err := db.C(COLLECTION).Find(bson.M{}).All(&results)

	if err != nil {
		// TODO: Do something about the error
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Results All: ", results)
	}
	return results, err
}

//func prog(state overseer.State) {
//	log.Printf("app (%s) listening...", state.ID)
//	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "app (%s) says hello\n", state.ID)
//	}))
//	http.Serve(state.Listener, nil)
//}

func main() {

	//overseer.Run(overseer.Config{
	//	Program: prog,
	//	Address: ":3000",
	//	Fetcher: &fetcher.HTTP{
	//		URL:      "http://localhost/",
	//		Interval: 1 * time.Second,
	//	},
	//})


	names, err := db.CollectionNames()
	if err != nil {
		// Handle error
		log.Printf("Failed to get coll names: %v", err)
		return
	}
	log.Println(" get coll names:", names)

	r := mux.NewRouter()

	r.HandleFunc("/home",LocationFisrtParking).Methods("GET")
	//UpdateParking
	r.HandleFunc("/parkings",UpdateParking).Methods("POST")
	r.HandleFunc("/parkings", LocationFisrtParking).Methods("GET")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}