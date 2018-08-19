package controller

import (
	"encoding/json"
	"friend-management/model"
	"friend-management/repository"
	"friend-management/util"
	"log"
	"net/http"
)

var UserRepo repository.IRepository

func StartWebServer(port string) {

	r := NewRouter()
	http.Handle("/", r)

	log.Println("Starting HTTP service at " + port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Println("An error occured starting HTTP listener at port " + port)
		log.Println("Error: " + err.Error())
	}
}

func ConnectFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	req := &model.FriendConnectionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("Error decoding body: %s", err)
		response := &model.BasicResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	// Read the user struct BoltDB
	userA, errA := UserRepo.GetUser(req.Friends[0])
	userB, errB := UserRepo.GetUser(req.Friends[1])

	// If err, return
	if errA != nil || errB != nil {
		log.Printf("Error QueryA: %s", errA)
		log.Printf("Error QueryB: %s", errB)
		response := &model.BasicResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	//Add Friend
	var bISaFriend bool = util.Contains(userA.Friends, userB.Email)
	if !bISaFriend {
		userA.Friends = append(userA.Friends, userB.Email)
		log.Printf("B added to A Friend's")
	}

	var aISbFriend bool = util.Contains(userB.Friends, userA.Email)

	if !aISbFriend {
		userB.Friends = append(userB.Friends, userA.Email)
		log.Printf("A added to B Friend's")
	}

	errWA := UserRepo.UpdateUser(userA)
	errWB := UserRepo.UpdateUser(userB)

	// If err, return
	if errWA != nil || errWB != nil {
		log.Printf("Error QueryA: %s", errA)
		log.Printf("Error QueryB: %s", errB)
		response := &model.BasicResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &model.BasicResponse{}
	response.Success = true

	json.NewEncoder(w).Encode(response)

}
