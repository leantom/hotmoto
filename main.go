package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"log"

	"github.com/jpillora/overseer"

	"./Module"
	"./MotoPark"

	"io/ioutil"

	"time"

	"github.com/jpillora/overseer/fetcher"
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
	if len(res) == 0 {
		respondWithJson(w, http.StatusCreated, "Khong co du lieu")
		return
	}
	respondWithJson(w, http.StatusOK, res)
}

func FindingParkingWithCurrentLocation(w http.ResponseWriter, r *http.Request) {
	var findingNear MotoPark.FindingNearLocation
	err := json.NewDecoder(r.Body).Decode(&findingNear)
	if err != nil {

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := MotoPark.FindNearLocationParking(findingNear)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(result) == 0 {
		respondWithJson(w, http.StatusCreated, "Khong co du lieu")
		return
	}
	respondWithJson(w, http.StatusCreated, result)
}

func InsertParking(w http.ResponseWriter, r *http.Request) {

	var parking MotoPark.MotoPark
	body, errRead := ioutil.ReadAll(r.Body)

	if errRead != nil {
		respondWithError(w, http.StatusBadRequest, errRead.Error())
		return
	}

	err := json.Unmarshal(body, &parking)

	if err != nil {
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

	err := json.NewDecoder(r.Body).Decode(&users)

	if err != nil {
		log.Print(&users)
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	usercurrent, errUser := Module.FindById(users.Username)
	println(usercurrent.Username)
	if errUser != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashPwd, errHashPwd := Module.Generate(users.Password)
	log.Print(hashPwd)
	if errHashPwd != nil {
		respondWithError(w, http.StatusInternalServerError, "Lỗi hệ thống vui lòng thử lại sau")
		return
	}

	errCompare := Module.Compare(hashPwd, usercurrent.Password)

	if errCompare != nil {
		respondWithError(w, http.StatusInternalServerError, "Password không đúng")
		return
	}

	respondWithJson(w, http.StatusCreated, usercurrent)
}

func FindAllUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	users, err := Module.FindAll()
	println("user finding all")
	if err != nil {

		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, users)
}

type UserRequest struct {
	Username string
}

func FindParksByUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userRequest UserRequest
	err := decoder.Decode(&userRequest)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, err.Error())
		panic(err)
	}
	result, err := MotoPark.FindParksByUserID(userRequest.Username)

	if len(result) == 0 {
		respondWithJson(w, http.StatusCreated, "Khong co du lieu")
		return
	}

	if err != nil {
		respondWithJson(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, result)
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

	_, errExist := Module.FindById(users.Username)
	if errExist == nil {
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

func UpdatePriceParking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var parkingRequest MotoPark.PriceParkRequest

	err := json.NewDecoder(r.Body).Decode(&parkingRequest)
	if err != nil {
		log.Print(&parkingRequest)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	park, err := MotoPark.UpdateCost(parkingRequest.IdPark, parkingRequest.Cost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, park)

}

func UpdateSlotParking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var parkingSlotRequest MotoPark.SlotParkRequest

	err := json.NewDecoder(r.Body).Decode(&parkingSlotRequest)
	if err != nil {
		log.Print(&parkingSlotRequest)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	park, err := MotoPark.UpdateAvailableSlot(parkingSlotRequest.IdPark, parkingSlotRequest.Slot)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, park)

}

func deleteParking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var parking Module.UserFindIDRequest

	if err := json.NewDecoder(r.Body).Decode(&parking); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := MotoPark.Delete(parking.ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, 200, "Sucessfully delete")

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

	r.HandleFunc("/api/home", LocationFisrtParking).Methods("GET")

	r.HandleFunc("/api/users", FindAllUser).Methods("GET")

	r.HandleFunc("/api/users/login", Login).Methods("POST")

	r.HandleFunc("/api/users/register", InsertUser).Methods("POST")

	r.HandleFunc("/api/users", UpdateUser).Methods("PUT")

	r.HandleFunc("/api/parkings", InsertParking).Methods("POST")

	r.HandleFunc("/api/parkings/updateCost", UpdatePriceParking).Methods("POST")

	r.HandleFunc("/api/parkings/updateSlot", UpdateSlotParking).Methods("POST")

	r.HandleFunc("/api/parkings", UpdateParking).Methods("PUT")

	r.HandleFunc("/api/parkings", deleteParking).Methods("DELETE")

	r.HandleFunc("/api/parkings", LocationFisrtParking).Methods("GET")

	r.HandleFunc("/api/parkings/users", FindParksByUser).Methods("POST")

	r.HandleFunc("/api/parkings/getNearCurrents", FindingParkingWithCurrentLocation).Methods("POST")

	r.HandleFunc("/api/uploads", Module.UploadFiles).Methods("POST")

	// push notification
	r.HandleFunc("/api/pushNotificationSingle", Module.PushNotificationSingle).Methods("POST")

	r.HandleFunc("/api/pushNotificationBookingPark", Module.PushNotificationBookParking).Methods("POST")

	r.HandleFunc("/api/deleteDeviceToken", Module.DeleteDeviceTokenByUserID).Methods("POST")

	r.HandleFunc("/api/registerDeviceToken", Module.RegisterDeviceTokenByUserID).Methods("POST")

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
