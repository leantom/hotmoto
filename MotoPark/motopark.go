package MotoPark

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"gopkg.in/mgo.v2"
	"fmt"
)

type MotoPark struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Position  [2]float64 `bson:"position" json:"position"`
	Name  		string        `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Phone string `bson:"phone" json:"phone"`
	Cost string `bson:"cost" json:"cost"`
	Total int `bson:"total" json:"total"`
}

const  (
	DB = "hotmoto_db"
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


func FindAll() ([]MotoPark, error)  {

	var users []MotoPark
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	if err != nil {

		return nil, err
	}
	log.Println(COLLECTION)
	return users, err
}

func FindById(userID string) (MotoPark, error)  {
	var user MotoPark
	err := db.C(COLLECTION).Find(bson.M{"_id":bson.ObjectIdHex(userID)}).One(&user)
	return user, err
}
func Insert(user MotoPark) (error)  {

	err := db.C(COLLECTION).Insert(user)
	return  err
}

func Update(user MotoPark) (error)  {

	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return  err
}
func Delete(userID string) ( error)  {

	err := db.C(COLLECTION).RemoveId(userID)
	return err
}

func  FindNearLocationParking() ([]MotoPark, error) {
	var parks []MotoPark

	// search criteria
	long := -73.8601152
	lat := 	40.7311739

	scope := 3000 // max distance in metres

	collection := db.C(COLLECTION)

	err := collection.Find(bson.M{
		"location" : bson.M{
			"$nearSphere":bson.M{
				"$geometry": bson.M{
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": scope,
			},
		},
	}).Limit(10).All(&parks)


	return parks, err
}
