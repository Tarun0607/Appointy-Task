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
	"go.mongodb.org/mongo-driver/bson"
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



type Participant struct {
	Id           int       `json:"id"`
	Name        string    `json:"name"`
	Email		string 	  `json:"email"`
	Rsvp		string    `json:"rsvp"`
	
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
	response.Header().Set("content-type", "application/json")
	
	params := mux.Vars(request)
	var meetings []MongoFields	
	collection := client.Database("MeetingsAPI").Collection("Meeting_Details")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting MongoFields
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	//json.NewEncoder(response).Encode(meetings)

	for _,item:=range meetings{
		if(item.Id==params["id"]){
			json.NewEncoder(response).Encode(item)
			return
		}
	}
	json.NewEncoder(response).Encode(&MongoFields{})
}



func getMeetingTime(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meetings []MongoFields
	params := mux.Vars(request)
	
	collection := client.Database("MeetingsAPI").Collection("Meeting_Details")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting MongoFields
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	//json.NewEncoder(response).Encode(meetings)
	ts, err := time.Parse(time.RFC3339, params["id1"])
	te, err := time.Parse(time.RFC3339, params["id2"])

	for _,item:=range meetings{
		if(item.StartTime.Before(ts) && item.EndTime.After(te)){
			json.NewEncoder(response).Encode(item)
			return		
		}
	}
	json.NewEncoder(response).Encode(&MongoFields{})
}



func getMeetingsBasedOnParticipant(response http.ResponseWriter, request *http.Request) {
	
	response.Header().Set("content-type", "application/json")
	var meetings []MongoFields
	var participant_info []Participant	
	params := mux.Vars(request)
	
	
	collection := client.Database("MeetingsAPI").Collection("Meeting_Details")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting MongoFields
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}


	collection = client.Database("MeetingsAPI").Collection("Participants_Details")
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err = collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var each_participant Participant
		cursor.Decode(&each_participant)
		participant_info = append(participant_info, each_participant)
	}
	if err = cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	var pid int
	for _,item:=range participant_info{
		if(item.Email == params["id3"] ){
			pid = item.Id
		}
	}
	fmt.Println(pid)
	for _,item:=range meetings{
		for _,each:= range item.Participants{
			if(each==pid){
				json.NewEncoder(response).Encode(item)
				
			}
		}
	}

	//json.NewEncoder(response).Encode(&Participant{})
	
	
}




func main() {

	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/meetings", addMeeting).Methods("POST")
	router.HandleFunc("/meetings/{id}", getMeeting).Methods("GET")
	router.HandleFunc("/meetings?start={id1}&end={id2}",getMeetingTime).Methods("GET") 
	router.HandleFunc("/participants/{id3}", getMeetingsBasedOnParticipant).Methods("GET")
	http.ListenAndServe(":12345", router)
	
}
