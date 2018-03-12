package MotoPark

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"gopkg.in/mgo.v2"
	"fmt"
)

type motopark struct {
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


func FindAll() ([]motopark, error)  {

	var users []motopark
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	if err != nil {

		return nil, err
	}
	log.Println(COLLECTION)
	return users, err
}

func FindById(userID string) (motopark, error)  {
	var user motopark
	err := db.C(COLLECTION).Find(bson.M{"_id":bson.ObjectIdHex(userID)}).One(&user)
	return user, err
}
func Insert(user motopark) (error)  {

	err := db.C(COLLECTION).Insert(user)
	return  err
}

func Update(user motopark) (error)  {

	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return  err
}
func Delete(userID string) ( error)  {

	err := db.C(COLLECTION).RemoveId(userID)
	return err
}
