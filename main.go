package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"fmt"
	"log"

	"github.com/jpillora/overseer"

	"./Module"
	"./MotoPark"

	"time"

	"github.com/jpillora/overseer/fetcher"
	"io/ioutil"

)

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document

// Fetch Example


func LocationFisrtParking(w http.ResponseWriter, r *http.Request) {

	res, err := MotoPark.FindAll()
	if err != nil {

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Print(res)
	respondWithJson(w, http.StatusOK, res)
}
func FindingParkingWithCurrentLocation(w http.ResponseWriter, r *http.Request) {
	var findingNear MotoPark.FindingNearLocation
	err := json.NewDecoder(r.Body).Decode(&findingNear)
	if err != nil {

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	 result,err := MotoPark.FindNearLocationParking(findingNear)

	 if err != nil {
		 respondWithError(w, http.StatusInternalServerError, err.Error())
		 return
	 }

	if len(result) == 0 {
		respondWithJson(w, http.StatusCreated, "Khong co du lieu")
		return
	} else {
		for i := 0; i< len(result) ;i++  {
			result[i].Total = len(result)
		}
	}
	respondWithJson(w, http.StatusCreated, result)
}

func InsertParking(w http.ResponseWriter, r *http.Request) {

	var parking MotoPark.MotoPark
	body, errRead :=  ioutil.ReadAll(r.Body)
	log.Print(r.Body)


	if errRead != nil {
		respondWithError(w, http.StatusBadRequest, errRead.Error())
		return
	}

	err := json.Unmarshal(body,&parking)

	if err != nil {
		log.Print(&parking)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Print(&parking)
	parking.ID = bson.NewObjectId()
	if err := MotoPark.Insert(parking); err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, parking)
}
//login
func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var users Module.UserService
	log.Print(r.Body)
	err := json.NewDecoder(r.Body).Decode(&users)

	if err != nil {
		log.Print(&users)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	usercurrent, errUser := Module.FindById(users.Username)

	if errUser != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashPwd,errHashPwd := Module.Generate(users.Password)
	log.Print(hashPwd)
	if errHashPwd != nil {
		respondWithError(w, http.StatusInternalServerError, "Lỗi hệ thống vui lòng thử lại sau")
		return
	}

	errCompare := Module.Compare(hashPwd,usercurrent.Password)
	if  errCompare != nil {
		respondWithError(w, http.StatusInternalServerError, "Password không đúng")
		return
	}

	respondWithJson(w, http.StatusCreated, usercurrent)
}

func FindAllUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	users, err := Module.FindAll()

	if err != nil {

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, users)
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var users Module.User

	log.Print(json.NewDecoder(r.Body))

	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		log.Print(&users)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(users.Password) == 0 || len(users.Username) == 0 {
		respondWithError(w, http.StatusInternalServerError, "Không được để rỗng Username hoặc Password")
		return
	}

	_,errExist := Module.FindById(users.Username)
	if errExist == nil{
		respondWithError(w, http.StatusInternalServerError, "Username đã tồn tại")
		return
	}

	users.ID = bson.NewObjectId()

	if err := Module.Insert(users); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, users)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var users Module.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		log.Print(&users)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Module.Update(users); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, users)

}

func UpdateParking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var parking MotoPark.MotoPark
	err := json.NewDecoder(r.Body).Decode(&parking)
	if err != nil {
		log.Print(&parking)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := MotoPark.Update(parking); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, parking)

}

func deleteParking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var parking MotoPark.MotoPark

	if err := json.NewDecoder(r.Body).Decode(&parking); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	currentParking, err := MotoPark.FindById(parking.ID.Hex())

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Not Found")
		return
	}

	if err := MotoPark.Delete(currentParking.ID.Hex()); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, "Sucessfully delete")

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

//prog(state) runs in a child process
func prog(state overseer.State) {
	log.Printf("app (%s) listening...", state.ID)
	r := mux.NewRouter()

	r.HandleFunc("/home", LocationFisrtParking).Methods("GET")

	r.HandleFunc("/users", FindAllUser).Methods("GET")

	r.HandleFunc("/users/login", Login).Methods("POST")

	r.HandleFunc("/users/register", InsertUser).Methods("POST")

	r.HandleFunc("/users", UpdateUser).Methods("PUT")

	r.HandleFunc("/parkings", InsertParking).Methods("POST")

	r.HandleFunc("/parkings", UpdateParking).Methods("PUT")

	r.HandleFunc("/parkings", deleteParking).Methods("DELETE")

	r.HandleFunc("/parkings", LocationFisrtParking).Methods("GET")

	r.HandleFunc("/parkings/getNearCurrents", FindingParkingWithCurrentLocation).Methods("POST")

	r.HandleFunc("/uploads", Module.UploadFiles).Methods("POST")

	http.Serve(state.Listener, r)

}

func main() {

	overseer.Run(overseer.Config{
		Program: prog,
		Address: ":8080",
		Fetcher: &fetcher.HTTP{
			URL:      "http://localhost:4000",
			Interval: 1 * time.Second,
		},
	})

}
