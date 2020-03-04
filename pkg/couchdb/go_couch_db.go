package couchdb

import (
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

const DBAuthUser = "DBAuthUser"
const DBAuthPwd = "DBAuthPwd"
const DBAuth = "DBAuth"
const DBName = "DBName"
type GoCouchDB struct {
	db 		*Database
}

type KVRead struct {
	_id 	string
	_rev	string
	Value 	string
}

type KVWrite struct {
	Value 	string
}

func NewGoCouchDB(name, address string, auth Auth) (*GoCouchDB, error) {
	conn, err := NewConnection(address, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	err = conn.CreateDB(name, auth)
	if err != nil {
		switch t := err.(type) {
		case *Error:
			if t.ErrorCode == "file_exists" {
				// do nothing
			}
		default:
			return nil, err
		}
	}
	db := conn.SelectDB(name, auth)
	return &GoCouchDB{db:db}, nil
}

// Implements DB.
func (cdb *GoCouchDB) Get(key []byte) []byte {
	retry := 0
	for {
		key = nonNilBytes(key)
		var doc KVRead

		_, err := cdb.db.Read(hex.EncodeToString(key), &doc,nil)

		if err != nil {
			switch t := err.(type) {
			case *Error:
				if t.ErrorCode == "not_found" {
					return nil
				}
			default:
				panic(err)
			}
		}
		res, err := hex.DecodeString(doc.Value)
		if err != nil {
			panic(err)
		}
		if len(res) == 0 {
			if err != nil {
				fmt.Println("***********Retry***********", retry)
				fmt.Println("Method: Get, ","key: ", string(key), " id:", hex.EncodeToString(key))
				fmt.Println(err)
			}
		} else {
			return res
		}
		retry++
	}

}

// Implements DB.
func (cdb *GoCouchDB) Has(key []byte) bool {
	return cdb.Get(key) != nil
}

// Implements DB.
func (cdb *GoCouchDB) Set(key []byte, value []byte) {
	retry := 0
	for {
		key = nonNilBytes(key)
		value = nonNilBytes(value)
		id := hex.EncodeToString(key)
		rev := cdb.GetRev(key)

		var newDoc KVWrite
		newDoc = KVWrite{
			Value:	hex.EncodeToString(value),
		}
		// save newDoc
		rev, err := cdb.db.Save(newDoc, id, rev)
		if err != nil {
			fmt.Println("***************Retry******************", retry)
			fmt.Println("Method: Set, ","key: ", string(key), ", id:", hex.EncodeToString(key))
			fmt.Println(err)
		} else {
			return
		}
		retry++
	}
}

// Implements DB.
func (cdb *GoCouchDB) SetSync(key []byte, value []byte) {
	cdb.Set(key, value)
}


// Implements DB.
func (cdb *GoCouchDB) Delete(key []byte) {
	retry := 0
	for {
		key = nonNilBytes(key)
		id := hex.EncodeToString(key)
		// read oldDoc & now rev
		rev := cdb.GetRev(key)
		if rev == "" {
			return
		}
		rev, err := cdb.db.Delete(id, rev)
		if err != nil {
			fmt.Println("***************Retry******************", retry)
			fmt.Println("Method: Delete, ","key: ", string(key), ", id:", hex.EncodeToString(key))
			fmt.Println(err)
			cdb.Delete(key)
		} else {
			return
		}
		retry++
	}
}

// Implements DB.
func (cdb *GoCouchDB) DeleteSync(key []byte) {
	cdb.Delete(key)
}

// Implements DB.
func (cdb *GoCouchDB) Close() {

}

// Implements DB.
func (cdb *GoCouchDB) Print() {

}

// Implements DB.
func (cdb *GoCouchDB) Stats() map[string]string {
	return nil
}

// Implements DB.
func (cdb *GoCouchDB) NewBatch() dbm.Batch {
	batch := cdb.db.NewBulkDocument()
	return &goCouchDBBatch{cdb,batch}
}

type CouchIterator struct {
	cdb 		*GoCouchDB
	results 	*RangeQueryResponse
	cursor		int
	start		[]byte
	end			[]byte
	isReverse	bool
	valid       bool
}

// Implements DB.
func (cdb *GoCouchDB) Iterator(start, end []byte) dbm.Iterator {

	return cdb.newCouchIterator(start, end, false)
}

// Implements DB.
func (cdb *GoCouchDB) ReverseIterator(start, end []byte) dbm.Iterator {

	return cdb.newCouchIterator(start, end, true)
}

func (cdb *GoCouchDB) newCouchIterator(start, end []byte, isReverse bool) dbm.Iterator{
	var results *RangeQueryResponse
	results, err := cdb.db.ReadRange(hex.EncodeToString(start), hex.EncodeToString(end))
	if err != nil {
		panic(err)
	}
	return &CouchIterator{cdb,results,0,start,end,isReverse,true}
}

// Implements Iterator.
func (iter *CouchIterator) Domain() (start []byte, end []byte){
	return iter.start, iter.end
}

// Implements Iterator.
func (iter *CouchIterator) Valid() bool{
	if len(iter.results.Rows) > 0 && iter.valid {
		return true
	}
	return false
}

// Implements Iterator.
func (iter *CouchIterator) Next() {
	iter.assertValid()
	if iter.cursor < len(iter.results.Rows) - 1{
		iter.cursor++
	}else{
		iter.valid = false
	}
}

// Implements Iterator.
func (iter *CouchIterator) Key() (key []byte){
	iter.assertValid()
	if iter.isReverse {
		key, err := hex.DecodeString(iter.results.Rows[len(iter.results.Rows)-iter.cursor-1].Key)
		if err != nil {
			panic(err)
		}
		return key
	}
	key, err := hex.DecodeString(iter.results.Rows[iter.cursor].Key)
	if err != nil {
		panic(err)
	}
	return key
}

// Implements Iterator.
func (iter *CouchIterator) Value() (value []byte){
	iter.assertValid()
	if iter.isReverse {
		value, err := hex.DecodeString(iter.results.Rows[len(iter.results.Rows)-iter.cursor-1].Doc.Value)
		if err != nil {
			panic(err)
		}
		return value
	}
	value, err := hex.DecodeString(iter.results.Rows[iter.cursor].Doc.Value)
	if err != nil {
		panic(err)
	}
	return value
}

// Implements Iterator.
func (iter *CouchIterator) Close() {
	iter = nil
}

func (iter *CouchIterator) assertValid() {
	if !iter.Valid() {
		panic("goCouchDBIterator is invalid")
	}
}

type goCouchDBBatch struct{
	cdb *GoCouchDB
	batch *BulkDocument
}

// Implements Batch.
func (mBatch *goCouchDBBatch) Set(key, value []byte) {
	var newDoc KVWrite
	value = nonNilBytes(value)
id := hex.EncodeToString(key)
	rev := mBatch.cdb.GetRev(key)
	newDoc = KVWrite{
		Value:	hex.EncodeToString(value),
	}
	err := mBatch.batch.Save(newDoc, id, rev)
	if err != nil {
		panic(err)
	}
}

// Implements Batch.
func (mBatch *goCouchDBBatch) Delete(key []byte) {
	id := hex.EncodeToString(key)
	rev := mBatch.cdb.GetRev(key)
	if rev == "" {
		return
	}
	err := mBatch.batch.Delete(id, rev)
	if err != nil {
		panic(err)
	}
}

// Implements Batch.
func (mBatch *goCouchDBBatch) Write() {
	if mBatch.batch.docs == nil {
		return
	}
	_, err := mBatch.batch.Commit()
	if err != nil {
		panic(err)
	}
}

// Implements Batch.
func (mBatch *goCouchDBBatch) WriteSync() {
	mBatch.Write()
}

// Implements Batch.
func (mBatch *goCouchDBBatch) Close() {

}

func nonNilBytes(bz []byte) []byte {
	if bz == nil {
		return []byte{}
	}
	return bz
}


func (cdb *GoCouchDB) GetRev2(key []byte) (string, error) {
	id := hex.EncodeToString(key)
	// read oldDoc & now rev
	rev, err := cdb.db.Read(id, nil, nil)
	if err != nil {
		switch t := err.(type) {
		case *Error:
			if t.ErrorCode == "not_found" {
				return "", err
			}
		default:
			fmt.Println("***********")
			fmt.Println(err)
			fmt.Println("***********")
			panic(err)
		}
	}
	return rev, nil
}


func (cdb *GoCouchDB) GetRev(key []byte) string {
	id := hex.EncodeToString(key)
	// read oldDoc & now rev
	rev, err := cdb.db.Read(id, nil, nil)
	if err != nil {
		switch t := err.(type) {
		case *Error:
			if t.ErrorCode == "not_found" {
				return ""
			}
		default:
			fmt.Println("***********")
			fmt.Println(err)
			fmt.Println("***********")
			panic(err)
		}
	}
	return rev
}

func (cdb *GoCouchDB) SetRev(key, value []byte, rev string) (string, error) {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	id := hex.EncodeToString(key)
	var newDoc KVWrite
	newDoc = KVWrite{
		Value:	hex.EncodeToString(value),
	}
	// save newDoc
	rev, err := cdb.db.Save(newDoc, id, rev)
	if err != nil {
		return "",err
	}
	return rev, nil
}

// Implements DB.
func (cdb *GoCouchDB) GetRevAndValue(key []byte) (string, []byte) {
	key = nonNilBytes(key)
	var doc KVRead
	rev, err := cdb.db.Read(hex.EncodeToString(key), &doc,nil)
	if err != nil {
		switch t := err.(type) {
		case *Error:
			if t.ErrorCode == "not_found" {
				return "", nil
			}
		default:
			fmt.Println("***********")
			fmt.Println(err)
			fmt.Println("***********")
			panic(err)
		}
	}
	res, err := hex.DecodeString(doc.Value)
	if err != nil {
		panic(err)
	}
	if len(res) == 0 {
		return "", nil
	}
	return rev, res
}

func (cdb *GoCouchDB) ResetDB() error {
	var auth *BasicAuth
	name := viper.GetString(DBName)
	user := viper.GetString(DBAuthUser)
	pwd := viper.GetString(DBAuthPwd)
	if user == "" && pwd == "" {
		err := cdb.db.connection.DeleteDB(name, nil)
		if err != nil {
			return err
		}
		err = cdb.db.connection.CreateDB(name, nil)
		if err != nil {
			return err
		}
	} else {
		auth = &BasicAuth{user,pwd}
		err := cdb.db.connection.DeleteDB(name, auth)
		if err != nil {
			return err
		}
		err = cdb.db.connection.CreateDB(name, auth)
		if err != nil {
			return err
		}
	}


	return nil
}