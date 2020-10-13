package couchdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRead2(t *testing.T)  {
	db, _ := NewGoCouchDB("ci123", "192.168.2.89:30301", nil)
	key := []byte("fc//collectedFeeserwrw")
	assert := assert.New(t)
	for i := 0; i < 1000; i++ {
		//fmt.Println("processing: ", i)
		//_, err2 := db.GetRev2(key)
		//if err2 != nil {
		//	panic(err2)
		//}
		bz := db.Get(key)
		assert.NotEmpty(bz)
		//db.Set(key, []byte("dfasfdsagagfe2323wf"))
	}
}

