package couchdb

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/logger"
	dbm "github.com/tendermint/tm-db"
)

const DBAuthUser = "DBAuthUser"
const DBAuthPwd = "DBAuthPwd"
const DBAuth = "DBAuth"
const DBName = "DBName"
type GoCouchDB struct {
	db 		*Database
	lg      logger.Logger
}

type KVWrite struct {
	Value 	string
}

func NewGoCouchDB(name, address string, auth Auth) (*GoCouchDB, error) {
	lg := logger.GetLogger()
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
			} else {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	db := conn.SelectDB(name, auth)
	return &GoCouchDB{
		db:	db,
		lg:	lg,
	}, nil
}

// Implements DB.
func (cdb *GoCouchDB) Get(key []byte) []byte {
	retry := 0
	for {
		key = nonNilBytes(key)
		var doc KVWrite

		_, err := cdb.db.Read(hex.EncodeToString(key), &doc,nil)
		if err != nil {
			switch t := err.(type) {
			case *Error:
				if t.ErrorCode == "not_found" {
					return nil
				}
			default:
				cdb.lg.Info("***************Retry******************")
				cdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
				cdb.lg.Info(fmt.Sprintf("Method: Get, key: %s, id: %s", string(key), hex.EncodeToString(key)))
				cdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
				retry ++
				continue
			}
		}
		res, _ := hex.DecodeString(doc.Value)
		return res
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
			cdb.lg.Info("***************Retry******************")
			cdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			cdb.lg.Info(fmt.Sprintf("Method: Set, key: %s, id: %s", string(key), hex.EncodeToString(key)))
			cdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
			retry++
		} else {
			return
		}
	}
}

func (cdb *GoCouchDB) SetDoc(id, rev string, doc interface{}) {
	retry := 0
	for {
		// save newDoc
		_, err := cdb.db.Save(doc, id, rev)
		if err != nil {
			cdb.lg.Info("***************Retry******************")
			cdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			cdb.lg.Info(fmt.Sprintf("Method: SetDoc, id: %s, rev: %s", id, rev))
			cdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
			retry++
		} else {
			return
		}
	}
}

// Implements DB.
func (cdb *GoCouchDB) SetSync(key []byte, value []byte) {
	cdb.Set(key, value)
}

// Implements DB.
func (cdb *GoCouchDB) Delete(key []byte) {
	retry := 0
	key = nonNilBytes(key)
	id := hex.EncodeToString(key)
	for {
		// read oldDoc & now rev
		rev := cdb.GetRev(key)
		if rev == "" {
			return
		}
		rev, err := cdb.db.Delete(id, rev)
		if err != nil {
			cdb.lg.Info("***************Retry******************")
			cdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			cdb.lg.Info(fmt.Sprintf("Method: Delete, key: %s, id: %s", string(key), hex.EncodeToString(key)))
			cdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
			retry++
		} else {
			return
		}
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
	retry := 0
	for {
		results, err := cdb.db.ReadRange(hex.EncodeToString(start), hex.EncodeToString(end))
		if err != nil {
			cdb.lg.Info("***************Retry******************")
			cdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			cdb.lg.Info(fmt.Sprintf("Method: ReadRange, start: %s, end: %s", hex.EncodeToString(start), hex.EncodeToString(end)))
			cdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
			retry++
		} else {
			return &CouchIterator{cdb,results,0,start,end,isReverse,true}
		}
	}
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
	//mBatch.cdb.lg.Info(fmt.Sprintf("BatchSet: key:%s GetRev:%s", string(key), rev))
	newDoc = KVWrite{
		Value:	hex.EncodeToString(value),
	}
	err := mBatch.batch.Save(newDoc, id, rev)
	if err != nil {
		mBatch.cdb.lg.Error("BatchSet Error", err)
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
		mBatch.cdb.lg.Error("BatchDelete Error", err)
		panic(err)
	}
}

type logBulk struct {
	RealKey string `json:"real_key"`
	BulkDoc bulkDoc
}

type logBulkResult struct {
	RealKey string `json:"real_key"`
	BulkDoc BulkDocumentResult
}

// Implements Batch.
func (mBatch *goCouchDBBatch) Write() {
	if mBatch.batch.docs == nil || mBatch.batch.closed{
		return
	}

	batchDocsMap := make(map[string]logBulk)
	for _, v := range mBatch.batch.docs {
		rk, _ := hex.DecodeString(v._id)
		batchDocsMap[v._id] = logBulk{
			RealKey: string(rk),
			BulkDoc: v,
		}
	}
	//mBatch.cdb.lg.Info(fmt.Sprintf("Before batchCommit docsMap: %v", batchDocsMap))
	retry := 0
	for {
		resp, err := mBatch.batch.Commit()
		if err != nil {
			//mBatch.cdb.lg.Error(fmt.Sprintf("BatchCommit Error: %s", err.Error()))
			//mBatch.cdb.lg.Info("***************Retry******************")
			//mBatch.cdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			mBatch.batch.closed = false
			retry++
			continue
		} else {
			//mBatch.cdb.lg.Info(fmt.Sprintf("BatchCommit response: %v", resp))
			//var docs []bulkDoc
			respDocsMap := make(map[string]logBulkResult)
			for _, v := range resp {
				rk, _ := hex.DecodeString(v.ID)
				respDocsMap[v.ID] = logBulkResult{
					RealKey: string(rk),
					BulkDoc: v,
				}
				if v.Ok {
					delete(batchDocsMap, v.ID)
				}
			}
			//mBatch.cdb.lg.Info(fmt.Sprintf("BatchCommit response docsMap: %v", respDocsMap))
			if len(batchDocsMap) != 0 {
				for _, v := range batchDocsMap {
					mBatch.cdb.SetDoc(v.BulkDoc._id, v.BulkDoc._rev, v.BulkDoc.doc)
				}
				mBatch.cdb.lg.Info("BatchCommit retry success")
			} else {
				mBatch.cdb.lg.Info("BatchCommit success")
			}
			return
		}
	}
}

// Implements Batch.
func (mBatch *goCouchDBBatch) WriteSync() {
	mBatch.Write()
}

// Implements Batch.
func (mBatch *goCouchDBBatch) Close() {
	mBatch = nil
}

func nonNilBytes(bz []byte) []byte {
	if bz == nil {
		return []byte{}
	}
	return bz
}

func (cdb *GoCouchDB) GetRev(key []byte) string {
	id := hex.EncodeToString(key)
	// read oldDoc & now rev
	retry := 0
	for {
		rev, err := cdb.db.Read(id, nil, nil)
		if err != nil {
			switch t := err.(type) {
			case *Error:
				if t.ErrorCode == "not_found" {
					cdb.lg.Info(fmt.Sprintf("GetRev not_found, key: %s, id: %s", string(key), hex.EncodeToString(key)))
					return ""
				}
			default:
				cdb.lg.Info("***************Retry******************")
				cdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
				cdb.lg.Info(fmt.Sprintf("Method: GetRev, key: %s, id: %s", string(key), hex.EncodeToString(key)))
				cdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
				retry++
			}
		} else {
			return rev
		}
	}
}