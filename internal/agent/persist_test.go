package agent

import (
	"os"
	"testing"
)

func TestPersist(t *testing.T) {
	dbPath := "test.db"
	_ = os.Remove(dbPath)
	if err := InitDB(dbPath); err != nil {
		t.Fatalf("init db: %v", err)
	}
	task := &Task{ID: "t1", Status: "pending"}
	if err := SaveTask(task); err != nil {
		t.Fatalf("save: %v", err)
	}
	status, err := GetTaskStatus("t1")
	if err != nil || status != "pending" {
		t.Fatalf("get: %v, status: %s", err, status)
	}
	_ = os.Remove(dbPath)
}
