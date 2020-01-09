package couchdbsource

import (
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"log"
	"strings"
)


const SharedKey  = "order//OrderBook"

func NewCouchSource(dbname, host string) *CouchDBSourceImp {
	imp := &CouchDBSourceImp{
		dbname: dbname,
		hostStr: host,
	}
	conn, err := imp.GetDBConnection()
	if err != nil {
		panic(err)
	}
	imp.conn = conn
	return imp
}

type CouchDBSourceImp struct {
	dbname  string
	hostStr string
	conn    *couchdb.GoCouchDB
}

func (s *CouchDBSourceImp) FetchSource() (hostArr []string) {

	if s.conn == nil {
		conn, err := s.GetDBConnection()
		if err != nil {
			log.Println(err)
		}
		s.conn =conn
	}

	bz := s.conn.Get([]byte(SharedKey))
	var shared map[string]interface{}
	err := json.Unmarshal(bz, &shared)
	if err != nil {
		log.Println(err)
	}

	orderDict, ok := shared["value"].(map[string]interface{})
	if !ok {
		return
	}
	lists := orderDict["lists"].([]interface{})
	if !ok {
		return
	}

	for _, value := range lists {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		name := item["name"].(string)
		if !strings.HasPrefix(name, "http") {
			name = "http://" + name + ":80"
		}
		name = "http://192.168.2.89:1317"
		hostArr = append(hostArr, name)
	}
	return
}

func (svr *CouchDBSourceImp) GetDBConnection() (db *couchdb.GoCouchDB, err error) {
	s := strings.Split(svr.hostStr, "://")
	if len(s) < 2 {
		return nil, errors.New("statedb format error")
	}
	if s[0] != "couchdb" {
		return nil, errors.New("statedb format error")
	}
	auths := strings.Split(s[1], "@")

	if len(auths) < 2 {
		db, err = couchdb.NewGoCouchDB(svr.dbname, auths[0],nil)
	} else {
		info := auths[0]
		userpass := strings.Split(info, ":")
		if len(userpass) < 2 {
			db, err = couchdb.NewGoCouchDB(svr.dbname, auths[1],nil)
		}
		auth := &couchdb.BasicAuth{Username: userpass[0], Password: userpass[1]}
		db, err = couchdb.NewGoCouchDB(svr.dbname, auths[1], auth)
	}
	//if err != nil {
	//	err = errors.New(fmt.Sprintf("cannot connect to couchdb, expect couchdb://xxxxx:5984 or couchdb://user:pass@xxxxx:5984, got %s", svr.hostStr))
	//}
	return
}