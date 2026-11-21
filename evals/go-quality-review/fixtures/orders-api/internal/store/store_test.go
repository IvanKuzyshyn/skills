package store

import "testing"

func TestSave(t *testing.T) {
	err := Save(&Order{ID: "1", Item: "book", Amount: 2})
	if err != nil {
		t.Fatal("failed")
	}
}
