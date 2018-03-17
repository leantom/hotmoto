package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"encoding/json"

	"gopkg.in/mgo.v2/bson"


	"fmt"
	"log"

	//"github.com/jpillora/overseer"

	"./Module"
	"./MotoPark"
)


// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document

// Fetch Example

func LocationFisrtParking(w http.ResponseWriter, r *http.Request) {

	res, err := MotoPark.FindAll()
	if err != nil {
		fmt.Print(res)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, res)
}

func FindAllUser(w http.ResponseWriter, r *http.Request) {
	defer  r.Body.Close()

	users,err := Module.FindAll()

	if err != nil {

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, users)
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	defer  r.Body.Close()

	var users Module.Users
	err := 	json.NewDecoder(r.Body).Decode(&users);
	if err != nil {
		log.Print(&users)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	users.ID = bson.NewObjectId()

	if err := Module.Insert(users); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, users)

}


func UpdateParking(w http.ResponseWriter, r *http.Request) {
	defer  r.Body.Close()

	var parking MotoPark.MotoPark
	err := 	json.NewDecoder(r.Body).Decode(&parking);
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

func deleteParking(w http.ResponseWriter, r *http.Request)  {
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
	//		URL:      "http://localhost",
	//		Interval: 1 * time.Second,
	//	},
	//})



	r := mux.NewRouter()

	r.HandleFunc("/home",LocationFisrtParking).Methods("GET")

	r.HandleFunc("/users",FindAllUser).Methods("GET")

	r.HandleFunc("/users",InsertUser).Methods("POST")

	r.HandleFunc("/users",InsertUser).Methods("POST")

	r.HandleFunc("/parkings",UpdateParking).Methods("POST")

	r.HandleFunc("/parkings", deleteParking).Methods("DELETE")

	r.HandleFunc("/parkings", LocationFisrtParking).Methods("GET")

	server := &http.Server{Addr: ":8080", Handler: r}
	server.SetKeepAlivesEnabled(true)
	server.ListenAndServe()

	//
	//if err := http.ListenAndServe(":8080", r); err != nil {
	//	log.Fatal(err)
	//}
}