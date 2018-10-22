package Module

import (
	"fmt"
	"github.com/anachronistic/apns"

	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)
const  (
	DEVICETOKEN_CENTER = "deviceTokenCenter"
)
const localhost  = "localhost:8080"
type Notification struct  {
	Title string        `bson:"title" json:"title"`
	Content string        `bson:"content" json:"content"`
	UserName string        `bson:"username" json:"username"`
	DeviceToken string     `bson:"deviceToken" json:"deviceToken"`
}

type DeviceTokenCenter struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	UserID 	string        `bson:"username" json:"username"`
	DeviceToken string        `bson:"device_token" json:"device_token"`
	Total int        `bson:"total" json:"total"`
}

func InsertDeviceToken(requestDeviceToken RegisterDeviceTokenRequest)  {
	var deviceToken DeviceTokenCenter
	deviceToken.ID = bson.NewObjectId()
	deviceToken.DeviceToken = requestDeviceToken.DeviceToken
	deviceToken.UserID = requestDeviceToken.UserID
	err := db.C(DEVICETOKEN_CENTER).Insert(deviceToken)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendAllDeviceNotification(notiRequest Notification)  {
	var devices []DeviceTokenCenter
	err := db.C(DEVICETOKEN_CENTER).Find(bson.M{}).All(&devices)

	if err == nil {
		 return 
	}

	for _, value := range devices {
		notiRequest.DeviceToken = value.DeviceToken
		SendPushToClient(notiRequest)
	}
}

func PushNotificationSingle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var notiRequest Notification

	if err := json.NewDecoder(r.Body).Decode(&notiRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user,err := FindById(notiRequest.UserName)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	notiRequest.DeviceToken = user.DeviceToken
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
