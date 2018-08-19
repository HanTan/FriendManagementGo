package model

type User struct {
	Email        string   `json:"email"`
	Friends      []string `json:"friends"`
	Subscription []string `json:"subscription"`
	Blocked      []string `json:"blocked"`
}

type BasicResponse struct {
	Success bool `json:"success"`
}

type FriendConnectionRequest struct {
	Friends []string `json:"friends"`
}
type FriendListRequest struct {
	Email string `json:"email"`
}

type FriendListResponse struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

type CommonFriendRequest struct {
	Friends []string `json:"friends"`
}

type SubscriptionRequest struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}
