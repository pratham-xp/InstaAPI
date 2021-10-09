package main

import (
	"InstaClone/helper"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"pass,omitempty" bson:"pass,omitempty"`
}

type Post struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption         string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL        string             `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	PostedTimestamp string             `json:"postedTimestamp,omitempty" bson:"postedTimestamp,omitempty"`
}

var UserCollection = helper.ConnectDB("users")
var PostCollection = helper.ConnectDB("posts")

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)

	result, err := UserCollection.InsertOne(context.TODO(), user)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func getUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var user User
	// we get params with mux.
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	err := UserCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var post Post

	_ = json.NewDecoder(r.Body).Decode(&post)

	result, err := PostCollection.InsertOne(context.TODO(), post)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func getPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var post Post

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	err := PostCollection.FindOne(context.TODO(), filter).Decode(&post)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []Post

	cur, err := PostCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var post Post

		err := cur.Decode(&post)

		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(posts)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/InstaClone/users", createUser).Methods("POST")
	r.HandleFunc("/InstaClone/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/InstaClone/posts", createPost).Methods("POST")
	r.HandleFunc("/InstaClone/posts/{id}", getPost).Methods("GET")
	r.HandleFunc("/InstaClone/posts/users/{id}", getPosts).Methods("GET")

	log.Fatal(http.ListenAndServe(":12345", r))
}
