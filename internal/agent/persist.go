package agent

import (
	"sync"

	"go.etcd.io/bbolt"
)

var (
	db     *bbolt.DB
	dbOnce sync.Once
)

const taskBucket = "tasks"

func InitDB(path string) error {
	var err error
	dbOnce.Do(func() {
		db, err = bbolt.Open(path, 0600, nil)
		if err == nil {
			err = db.Update(func(tx *bbolt.Tx) error {
				_, e := tx.CreateBucketIfNotExists([]byte(taskBucket))
				return e
			})
		}
	})
	return err
}

func SaveTask(task *Task) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		return b.Put([]byte(task.ID), []byte(task.Status))
	})
}

func GetTaskStatus(id string) (string, error) {
	var status string
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		v := b.Get([]byte(id))
		if v != nil {
			status = string(v)
		}
		return nil
	})
	return status, err
}
