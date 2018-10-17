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
	"net/http"
	"encoding/json"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Username string        `bson:"username" json:"username"`
	Password string        `bson:"password" json:"password"`
	DeviceToken string        `bson:"devicetoken" json:"devicetoken"`
}


type UserService struct {
	Username string        `bson:"username" json:"username"`
	Password string        `bson:"password" json:"password"`
	DeviceToken string        `bson:"devicetoken" json:"devicetoken"`
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
	fmt.Print(user.ID)
	return  err
}

func Delete(userID string) ( error)  {

	err := db.C(COLLECTION).RemoveId(userID)
	return err
}

func RegisterDeviceToken(userID string, deviceToken string) (error) {
	var user User

	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(userID)).One(&user)
	println(err.Error())
	if err == nil {
		user.DeviceToken = deviceToken
	}
	err = Update(user)
	return  err
}

func DeleteDeviceToken(userID string) (error) {
	var user User

	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(userID)).One(&user)
	println(user.ID)
	if err == nil {
		user.DeviceToken = ""
	}
	err = Update(user)
	return  err
}

type UserIDRequest struct {
	UserID string `bson:"userID" json:"userID"`
}

type RegisterDeviceTokenRequest struct {
	UserID string `bson:"userID" json:"userID"`
	DeviceToken string `bson:"deviceToken" json:"deviceToken"`
}

func DeleteDeviceTokenByUserID(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userRequest UserIDRequest
	err := decoder.Decode(&userRequest)

	if err != nil {
		respondWithJson(w, http.StatusBadRequest, err.Error())
		panic(err)
	}

	err = DeleteDeviceToken(userRequest.UserID)
	if err != nil {
		respondWithJson(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJson(w, 200, "Xoá thành công")
}

func RegisterDeviceTokenByUserID(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var registerDeviceToken RegisterDeviceTokenRequest
	err := decoder.Decode(&registerDeviceToken)

	if err != nil {
		respondWithJson(w, http.StatusBadRequest, err.Error())
		panic(err)
	}

	err = RegisterDeviceToken(registerDeviceToken.UserID,registerDeviceToken.DeviceToken)
	if err != nil {
		respondWithJson(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJson(w, 200, "Đăng kí thành công")
}
