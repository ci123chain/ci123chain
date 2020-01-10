package couchdb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRead(t *testing.T)  {
	db, _ := NewGoCouchDB("ci123", "192.168.2.89:30301", nil)
	for i := 0; i < 10000; i++ {
		fmt.Println("processing: ", i)
		rev, err2 := db.GetRev2([]byte("fc//collectedFees"))
		assert.NotEqual(t, rev, "")
		assert.Nil(t, err2)
	}
}