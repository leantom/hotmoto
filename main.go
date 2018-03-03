package hotmoto

import (
	"net/http"
	"github.com/gorilla/mux"

	"encoding/json"
	"hotmoto/Control"

)

var locationParkingDAO = Control.LocationParkingDAO{}


// Fetch Example

func LocationFisrtParking(w http.ResponseWriter, r *http.Request) {

	res, err := locationParkingDAO.FindAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, res)
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


func init() {
	Control.Connect()
}

func main() {

	r := mux.NewRouter()
	//LocationFisrtParking
	r.HandleFunc("/parkings", LocationFisrtParking).Methods("GET")

}