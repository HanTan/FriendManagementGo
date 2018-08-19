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