package Module

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"

	"log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"strings"
	"errors"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Username string        `bson:"username" json:"username"`
	Password string        `bson:"password" json:"password"`
}


type UserService struct {
	Username string        `bson:"username" json:"username"`
	Password string        `bson:"password" json:"password"`
}


const  (
	DB = "hotmoto_db"
	COLLECTION = "user"
)

var deliminator = "||"
var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	fmt.Println("Mongo server connected")
	db = session.DB(DB)

}

func Generate(s string) (string, error) {
	salt := uuid.New().String()
	saltedBytes := []byte(s + salt)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash + deliminator + salt, nil
}

//Compare string to generated hash
func Compare(hash string, s string) error {
	parts := strings.Split(hash, deliminator)
	if len(parts) != 2 {
		return errors.New("Invalid hash, must have 2 parts")
	}

	incoming := []byte(s + parts[1])
	existing := []byte(parts[0])
	return bcrypt.CompareHashAndPassword(existing, incoming)
}


func FindAll() ([]User, error)  {

	var users []User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	if err != nil {

		return nil, err
	}
	log.Println(COLLECTION)
	return users, err
}

func FindById(username string) (User, error)  {
	var user User
	err := db.C(COLLECTION).Find(bson.M{"username":username}).One(&user)
	return user, err
}
func Insert(user User) (error)  {

	err := db.C(COLLECTION).Insert(user)
	return  err
}

func Update(user User) (error)  {

	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return  err
}
func Delete(userID string) ( error)  {

	err := db.C(COLLECTION).RemoveId(userID)
	return err
}
