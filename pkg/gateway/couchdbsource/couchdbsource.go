package couchdbsource



import (
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"log"
	"strings"
)


const Name  = "ci123"
const SharedKey  = "order//OrderBook"
//const StateDB = "couchdb://couchdb-service:5984"
const StateDB = "couchdb://192.168.2.89:30301"

func NewCouchSource() *CouchDBSourceImp {
	return &CouchDBSourceImp{}
}

type CouchDBSourceImp struct {
}

func (s *CouchDBSourceImp) FetchSource() (hostArr []string) {
	conn, err := s.GetDBConnection()
	if err != nil {
		log.Println(err)
	}
	bz := conn.Get([]byte(SharedKey))
	var shared map[string]interface{}
	err = json.Unmarshal(bz, &shared)
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
		name = "http://192.168.2.89:1317"
		hostArr = append(hostArr, name)
	}
	log.Println(hostArr)
	return
	//s.ConfigServerPool(hostArr)
}

func (svr *CouchDBSourceImp) GetDBConnection() (db *couchdb.GoCouchDB, err error) {
	s := strings.Split(StateDB, "://")
	if len(s) < 2 {
		return nil, errors.New("statedb format error")
	}
	if s[0] != "couchdb" {
		return nil, errors.New("statedb format error")
	}
	auths := strings.Split(s[1], "@")

	if len(auths) < 2 {
		db, err = couchdb.NewGoCouchDB(Name, auths[0],nil)
	} else {
		info := auths[0]
		userpass := strings.Split(info, ":")
		if len(userpass) < 2 {
			db, err = couchdb.NewGoCouchDB(Name, auths[1],nil)
		}
		auth := &couchdb.BasicAuth{Username: userpass[0], Password: userpass[1]}
		db, err = couchdb.NewGoCouchDB(Name, auths[1], auth)
	}
	return
}