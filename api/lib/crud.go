package lib

import (
	"bytes"
	"encoding/json"

	"github.com/boltdb/bolt"
)

// @Author KHighness
// @Update 2022-11-01

type User struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Gender    uint8  `json:"gender"`
	Age       uint8  `json:"age"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

type UserDao struct {
	DB *bolt.DB
}

func NewUserDao(db *bolt.DB) *UserDao {
	return &UserDao{db}
}

func (d *UserDao) genCreateUserFunc(user *User) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		user.ID = id

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put(uint64ToBytes(user.ID), buf)
	}
}

func (d *UserDao) CreateUser(user *User) error {
	return d.DB.Batch(d.genCreateUserFunc(user))
}

func (d *UserDao) CreateUserInBatch(user *User) error {
	return d.DB.Batch(d.genCreateUserFunc(user))
}

func (d *UserDao) GetUserByID(id uint64) (user *User, err error) {
	user = &User{}
	err = d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))

		buf := b.Get(uint64ToBytes(id))
		bufCopy := make([]byte, len(buf))
		copy(bufCopy, buf)

		err := json.Unmarshal(bufCopy, user)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (d *UserDao) PutUser(user *User) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put(uint64ToBytes(user.ID), buf)
	})
}

func (d *UserDao) DeleteUerByID(id uint64) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))
		return b.Delete(uint64ToBytes(id))
	})
}

type Event struct {
	Time   int64  `json:"time"`
	Name   string `json:"name"`
	Type   uint8  `json:"type"`
	Cancel bool   `json:"cancel"`
}

type EventDao struct {
	DB *bolt.DB
}

func NewEventDao(db *bolt.DB) *EventDao {
	return &EventDao{DB: db}
}

func (d *EventDao) genCreateEventFunc(event *Event) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketEvents))

		buf, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return b.Put(uint64ToBytes(uint64(event.Time)), buf)
	}
}

func (d *EventDao) CreateEventInBatch(event *Event) error {
	return d.DB.Batch(d.genCreateEventFunc(event))
}

func (d *EventDao) CreateEvent(event *Event) error {
	return d.DB.Update(d.genCreateEventFunc(event))
}

func (d *EventDao) getEventBetween(start, end int64) (events []*Event, err error) {
	err = d.DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BucketEvents)).Cursor()
		min := uint64ToBytes(uint64(start))
		max := uint64ToBytes(uint64(end))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			var event Event
			if err := json.Unmarshal(v, &event); err != nil {
				return err
			}
			events = append(events, &event)
		}
		return nil
	})
	return
}

type BucketDao struct {
	DB *bolt.DB
}

func NewBucketDao(db *bolt.DB) *BucketDao {
	return &BucketDao{db}
}

func (d *BucketDao) CreateBucket(name []byte) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}

func (d *BucketDao) DeleteBucket(name []byte) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(name)
	})
}
