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

	var bBlocked bool = util.Contains(userA.Blocked, userB.Email)
	var aBlocked bool = util.Contains(userB.Blocked, userA.Email)

	if aBlocked || bBlocked {
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

func friendList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	req := &model.FriendListRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("Error decoding body: %s", err)
		response := &model.FriendListResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := UserRepo.GetUser(req.Email)

	// If err, return
	if err != nil {
		log.Printf("Error QueryA: %s", req.Email)
		response := &model.FriendListResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &model.FriendListResponse{}
	response.Success = true
	response.Friends = user.Friends
	response.Count = len(user.Friends)

	json.NewEncoder(w).Encode(response)

}

func commonFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	req := &model.CommonFriendRequest{}
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

	Commons := []string{}
	for _, a := range userA.Friends {
		for _, b := range userB.Friends {
			if a == b {
				Commons = append(Commons, a)
			}
		}
	}

	response := &model.FriendListResponse{}
	response.Success = true
	response.Friends = Commons
	response.Count = len(Commons)
	json.NewEncoder(w).Encode(response)
}

func subscribeFriend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	req := &model.SubscriptionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("Error decoding body: %s", err)
		response := &model.BasicResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	// Read the user struct BoltDB
	requestor, errA := UserRepo.GetUser(req.Requestor)

	// If err, return
	if errA != nil {
		log.Printf("Error QueryA: %s", requestor)
		response := &model.BasicResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	//Subscribe Friend
	var bIsSubscriber bool = false
	for _, u := range requestor.Subscription {
		if u == req.Target {
			bIsSubscriber = true
		}
	}

	if !bIsSubscriber {
		requestor.Subscription = append(requestor.Subscription, req.Target)
		log.Printf("B added to A Subscription's")
	}

	UserRepo.UpdateUser(requestor)

	response := &model.BasicResponse{}
	response.Success = true
	json.NewEncoder(w).Encode(response)
}

func blockFriend(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	req := &model.SubscriptionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("Error decoding body: %s", err)
		response := &model.BasicResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	// Read the user struct BoltDB
	requestor, errA := UserRepo.GetUser(req.Requestor)

	// If err, return
	if errA != nil {
		log.Printf("Error QueryA: %s", requestor)
		response := &model.BasicResponse{}
		response.Success = false
		json.NewEncoder(w).Encode(response)
		return
	}

	//Block Friend
	var bIsBlocked bool = false
	for _, u := range requestor.Blocked {
		if u == req.Target {
			bIsBlocked = true
		}
	}

	if !bIsBlocked {
		requestor.Blocked = append(requestor.Blocked, req.Target)
		log.Printf("B added to A Blocked's")
	}

	UserRepo.UpdateUser(requestor)

	response := &model.BasicResponse{}
	response.Success = true
	json.NewEncoder(w).Encode(response)

}
