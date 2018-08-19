package repository

import (
	"encoding/json"
	"fmt"
	"friend-management/model"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

type IRepository interface {
	OpenBoltDb()
	GetUser(email string) (model.User, error)
	GetAllUser() ([]model.User, error)
	UpdateUser(user model.User) error
	Seed()
}

type Repository struct {
	boltDB *bolt.DB
}

func (bc *Repository) OpenBoltDb() {
	var err error
	bc.boltDB, err = bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Start seeding users
func (bc *Repository) Seed() {
	bc.initializeBucket()
	bc.seedUsers()
}

func (bc *Repository) initializeBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("UserBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

// Seed (n) make-believe account objects into the AcountBucket bucket.
func (bc *Repository) seedUsers() {

	total := 10
	for i := 0; i < total; i++ {

		// Generate a key 10000 or larger
		Id := strconv.Itoa(i)
		Key := "User" + Id + "@hans.com"

		// Create an instance of our Account struct
		usr := model.User{Email: Key, Friends: []string{}, Subscription: []string{}, Blocked: []string{}}

		// Serialize the struct to JSON
		jsonBytes, _ := json.Marshal(usr)

		// Write the data to the AccountBucket
		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("UserBucket"))
			err := b.Put([]byte(Key), jsonBytes)
			return err
		})
	}
	fmt.Printf("Seeded %v fake users...\n", total)
}

func (bc *Repository) GetUser(email string) (model.User, error) {
	// Allocate an empty User instance we'll let json.Unmarhal populate for us in a bit.
	user := model.User{}

	// Read an object from the bucket using boltDB.View
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		// Read the bucket from the DB
		b := tx.Bucket([]byte("UserBucket"))

		// Read the value identified by our email supplied as []byte
		userBytes := b.Get([]byte(email))
		if userBytes == nil {
			return fmt.Errorf("No user found for " + email)
		}
		// Unmarshal the returned bytes into the user struct we created at
		// the top of the function
		json.Unmarshal(userBytes, &user)

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})
	// If there were an error, return the error
	if err != nil {
		return model.User{}, err
	}
	// Return the User struct and nil as error.
	return user, nil
}

func (bc *Repository) GetAllUser() ([]model.User, error) {
	users := []model.User{}

	err := bc.boltDB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("UserBucket"))
		b.ForEach(func(k, v []byte) error {
			user := model.User{}
			json.Unmarshal(v, &user)
			users = append(users, user)
			return nil
		})
		return nil
	})

	if err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (bc *Repository) UpdateUser(user model.User) error {

	// Serialize the struct to JSON
	jsonBytes, _ := json.Marshal(user)

	// Write the data to the AccountBucket
	err := bc.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("UserBucket"))
		err := b.Put([]byte(user.Email), jsonBytes)
		return err
	})

	return err
}
