package boltdb

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"musaitHrMgBotGo/pkg/repository"
	"strconv"
)

type UserRepository struct {
	db *bolt.DB
}

func NewUserRepository(db *bolt.DB) *UserRepository {
	return &UserRepository{db}
}
func (u *UserRepository) Save(chatID int64, value string, bucket repository.Bucket) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		fmt.Println(value)
		return b.Put(intToBytes(chatID), []byte(value))
	})
}
func (u *UserRepository) Get(userID int64, bucket repository.Bucket) (string, error) {
	var value string
	err := u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("userID ga tegishli bucket topilmadi %d", userID)
		}
		v := b.Get(intToBytes(userID))
		if v == nil {
			fmt.Errorf("bucket %s dan malumot topilmadi", bucket)
		}
		value = string(v)
		return nil
	})
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", errors.New("value is empty")
	}
	return value, nil
}
func intToBytes(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
}
