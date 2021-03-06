package MotoPark

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MotoPark struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Location struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"location"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	Total           int    `json:"total"`
	AvailableSlot   int    `json:"AvailableSlot"`
	OpenTime        int    `json:"openTime"`
	CloseTime       int    `json:"closeTime"`
	Status          int    `json:"status"`
	Fullname        string `json:"fullname"`
	Cost            int    `json:"cost"`
	NumberHours     int    `json:"numberHours"`
	Username        string `json:"username"`
	UserID          string `json:"userID"`
	Email           string `json:"email"`
	ImageURL        string `json:"imageUrl"`
	DescriptionPark string `json:"description_park"`
	Type            int    `json:"type"`
}

type PriceParkRequest struct {
	Cost   int    `json:"cost"`
	IdPark string `json:"id_park"`
}

type SlotParkRequest struct {
	Slot   int    `json:"slot"`
	IdPark string `json:"id_park"`
}

type SlotPark struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	AvailableSlot int           `json:"available_slot"`
	ParkID        string        `json:"park_id"`
	MaximumSlot   int           `json:"maximum_slot"`
}

type PositionParking struct {
	Type        string     `json:"type"`
	coordinates [2]float64 `json:"coordinates"`
}

type FindingNearLocation struct {
	Position [2]float64 `bson:"position" json:"position"`
	Scope    int        `bson:"scope" json:"scope"`
}

const (
	DB         = "hotmoto_db"
	COLLECTION = "motopark"
)

var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	fmt.Println("Mongo server connected")
	db = session.DB(DB)

}

func checkDBLostConnect() {
	if db == nil {
		session, err := mgo.Dial("localhost")
		if err != nil {
			fmt.Println("Failed to establish connection to Mongo server:", err)
		}
		fmt.Println("Mongo server connected")
		db = session.DB(DB)
	}

}

func FindAll() ([]MotoPark, error) {
	checkDBLostConnect()
	var users []MotoPark
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	if err != nil {

		return nil, err
	}
	log.Println(COLLECTION)
	return users, err
}

func FindById(userID string) (MotoPark, error) {
	var park MotoPark
	err := db.C(COLLECTION).Find(bson.M{"_id": bson.ObjectIdHex(userID)}).One(&park)
	return park, err
}

func FindParksByUserID(userID string) ([]MotoPark, error) {
	var parks []MotoPark
	err := db.C(COLLECTION).Find(bson.M{"username": userID}).All(&parks)
	return parks, err
}

func Insert(park MotoPark) error {

	err := db.C(COLLECTION).Insert(park)
	return err
}

func Update(park MotoPark) error {
	err := db.C(COLLECTION).UpdateId(park.ID, &park)
	return err
}

func UpdateCost(parkId string, cost int) (MotoPark, error) {
	var park MotoPark
	err := db.C(COLLECTION).Find(bson.M{"_id": bson.ObjectIdHex(parkId)}).One(&park)
	park.Cost = cost
	return park, err
}

func UpdateAvailableSlot(parkId string, slot int) (MotoPark, error) {
	var park MotoPark
	err := db.C(COLLECTION).Find(bson.M{"_id": bson.ObjectIdHex(parkId)}).One(&park)
	if slot >= park.Total {
		park.AvailableSlot = 0
	} else {
		park.AvailableSlot = park.Total - slot
	}
	return park, err
}

func Delete(parkingId string) error {

	err := db.C(COLLECTION).RemoveId(bson.ObjectIdHex(parkingId))
	return err
}

func FindNearLocationParking(findingLocation FindingNearLocation) ([]MotoPark, error) {
	var parks []MotoPark

	collection := db.C(COLLECTION)
	long := findingLocation.Position[0]
	lat := findingLocation.Position[1]

	scope := findingLocation.Scope
	log.Print(findingLocation.Position)
	err := collection.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": scope,
			},
		},
	}).All(&parks)

	return parks, err
}
