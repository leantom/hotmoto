package MotoPark

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MotoPark struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Position [2]float64    `bson:"position" json:"position"`
	Name     string        `bson:"name" json:"name"`
	Address  string        `bson:"address" json:"address"`
	Phone    string        `bson:"phone" json:"phone"`
	Cost     string        `bson:"cost" json:"cost"`
	Total    int           `bson:"total" json:"total"`
}

type FindingNearLocation struct {
	Position [2]float64    `bson:"position" json:"position"`
	scope float64 `bson:"scope" json:"scope"`
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

func FindAll() ([]MotoPark, error) {

	var users []MotoPark
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	if err != nil {

		return nil, err
	}
	log.Println(COLLECTION)
	return users, err
}

func FindById(userID string) (MotoPark, error) {
	var user MotoPark
	err := db.C(COLLECTION).Find(bson.M{"_id": bson.ObjectIdHex(userID)}).One(&user)
	return user, err
}

func Insert(park MotoPark) error {

	err := db.C(COLLECTION).Insert(park)
	return err
}

func Update(user MotoPark) error {

	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}
func Delete(userID string) error {

	err := db.C(COLLECTION).RemoveId(userID)
	return err
}

func FindNearLocationParking(findingLocation FindingNearLocation) ([]MotoPark, error) {
	var parks []MotoPark

	collection := db.C(COLLECTION)

	log.Print(findingLocation.Position)
	err := collection.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates":  findingLocation.Position,
				},
				"$maxDistance": findingLocation.scope,
			},
		},
	}).Limit(100).All(&parks)

	return parks, err
}
