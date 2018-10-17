package Module

import (
	"fmt"
	"github.com/anachronistic/apns"

	"net/http"
	"encoding/json"
)

const localhost  = "localhost:8080"
type Notification struct  {
	Title string        `bson:"title" json:"title"`
	Content string        `bson:"content" json:"content"`
	DeviceToken string     `bson:"deviceToken" json:"deviceToken"`
}

func PushNotificationSingle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var notiRequest Notification

	if err := json.NewDecoder(r.Body).Decode(&notiRequest); err != nil {
		println(err.Error())
		println(r.Body)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := SendPushToClient(notiRequest)

	if result.Error != nil {
		respondWithError(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	respondWithJson(w,200,result)


}


func SendPushToClient(notiRequest Notification) (NotificationResult)   {

	fmt.Println("pushText: ", notiRequest.Title)
	fmt.Println("content: ", notiRequest.Content)
	fmt.Println("pushToken: ", notiRequest.DeviceToken)

	dict := apns.NewAlertDictionary()
	dict.Title = notiRequest.Title
	dict.Body = notiRequest.Content

	payload := apns.NewPayload()
	payload.Badge = 1
	payload.Sound = "bingbong.aiff"
	payload.Alert = dict

	pn := apns.NewPushNotification()
	pn.DeviceToken = notiRequest.DeviceToken
	pn.AddPayload(payload)

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "./config/pushcert.pem", "./config/pushcert.pem")

	resp := client.Send(pn)

	alert, _ := pn.PayloadString()
	fmt.Println("  Alert:", alert)
	fmt.Println("Success:", resp.Success)
	fmt.Println("  Error:", resp.Error)

	result := NotificationResult{Alert:alert,Success:resp.Success,Error:resp.Error}
	return  result
}



type NotificationResult struct  {
	Alert string        `bson:"alert" json:"alert"`
	Success bool        `bson:"success" json:"success"`
	Error error     `bson:"error" json:"error"`
}
