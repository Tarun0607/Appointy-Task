package main

import (

	// Built-in Golang packages
	"context" // manage multiple requests
	"fmt"
	"encoding/json"     // Println() function
	 
	// get an object type
	"time"
	"net/http"
	// Official 'mongo-go-driver' packages
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFields struct {
		
	Id           string             `json:"id" bson:"id,omitempty"`
	Title        string             `json:"title" bson:"title,omitempty"`
	Participants [2]int             `json:"participants" bson:"participants,omitempty"`
	StartTime    time.Time          `json:"starttime" bson:"starttime,omitempty"`
	EndTime      time.Time          `json:"endtime" bson:"endtime,omitempty"`
	TimeStamp    time.Time          `json:"timestamp" bson:"timestamp,omitempty"`
}

var client *mongo.Client

func addMeeting(response http.ResponseWriter,request *http.Request){
	
	response.Header().Set("content-type", "application/json")
	var oneDoc MongoFields 
	_ = json.NewDecoder(request.Body).Decode(&oneDoc)
	collection := client.Database("MeetingsAPI").Collection("Meeting_Details")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, oneDoc)
	json.NewEncoder(response).Encode(result)

}

func getMeeting(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)

	var oneDoc MongoFields
	collection := client.Database("MeetingsAPI").Collection("Meeting_Details")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, MongoFields{Id :params["id"]}).Decode(&oneDoc)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(oneDoc)
}




func main() {

	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/meetings", addMeeting).Methods("POST")
	router.HandleFunc("/meetings/{id}", getMeeting).Methods("GET")
	http.ListenAndServe(":12345", router)
	
}
