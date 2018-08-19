# Friend-Management

## About this task
This is first golang application created by author `Hanstevin Tandian`. It includes feature to connect friend, list friends, get common friends, subscribe to an update, block friend / updates, and get update recipients list.

## Libraries used
 - github.com/gorilla/mux
 To handle routing
 - github.com/boltdb/bolt
 To handle data in memory

## About application
Upon start, there will be a seeding process where 10 dummy data is pumped into the bolt bucket to be used by services.

## Running application

```sh
./friend-management
```

## How to use

Available users will be:
- User1@hans.com
- User2@hans.com
- User3@hans.com
- User4@hans.com
- User5@hans.com
- User6@hans.com
- User7@hans.com

Access services using Postman / curl

with hostname : http://friend-management-hantan.c9users.io:8080/friend/

and request payload:
- /connect
```json
{
  "friends":
    [
      "User1@hans.com",
      "User3@hans.com"
    ]
}
```
- /list
```json
{
  "email": "User1@hans.com"
}
```
- /common
```json
{
  "friends": [
     "User2@hans.com",
      "User3@hans.com"
      ]
}
```
- /subscribe
```json
{
  "requestor": "User2@hans.com",
  "target": "User3@hans.com"
}
```
- /block
```json
{
  "requestor": "User3@hans.com",
  "target": "User1@hans.com"
}
```
- /send
```json
{
  "sender": "User1@hans.com",
  "text": "Hi! User4@hans.com"
}
```