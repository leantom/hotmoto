package Module

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"

	"log"
	//"fmt"

	"fmt"
)

type Users struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name  		string        `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Phone string `bson:"phone" json:"phone"`

}

const  (
	DB = "hotmoto_db"
	COLLECTION = "users"
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


func FindAll() ([]Users, error)  {

	var users []Users
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	if err != nil {

		return nil, err
	}
	log.Println(COLLECTION)
	return users, err
}

func FindById(userID string) (Users, error)  {
	var user Users
	err := db.C(COLLECTION).Find(bson.M{"_id":bson.ObjectIdHex(userID)}).One(&user)
	return user, err
}
func Insert(user Users) (error)  {

	err := db.C(COLLECTION).Insert(user)
	return  err
}

func Update(user Users) (error)  {

	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return  err
}
func Delete(userID string) ( error)  {

	err := db.C(COLLECTION).RemoveId(userID)
	return err
}
