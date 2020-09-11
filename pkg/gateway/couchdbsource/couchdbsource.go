package couchdbsource

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/couchdb"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"regexp"
	"strings"
)


const SharedKey  = "s/k:order/OrderBook"
const HostPattern  = "[*]+"
const DefaultDBName = "ci123"
func NewCouchSource(host, urlreg string) *CouchDBSourceImp {
	imp := &CouchDBSourceImp{
		hostStr: 	host,
		urlreg:  	urlreg,
	}
	conn, err := imp.GetDBConnection()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cann't connect to %s: %s", host, err.Error())))
	}
	imp.conn = conn
	return imp
}

type CouchDBSourceImp struct {
	hostStr string
	urlreg  string
	conn    *couchdb.GoCouchDB
}

func (s *CouchDBSourceImp) FetchSource() (hostArr []string) {
	logger.Debug("Start fetch from couchdb")
	if s.conn == nil {
		conn, err := s.GetDBConnection()
		if err != nil {
			logger.Error("Connection Error: ", err)
		}
		s.conn =conn
	}

	bz := s.conn.Get([]byte(SharedKey))
	var shared map[string]interface{}
	err := json.Unmarshal(bz, &shared)
	if err != nil {
		logger.Error("fetch data from couchdb error: ", err)
	}

	orderDict, ok := shared["value"].(map[string]interface{})
	if !ok {
		return
	}
	lists, ok := orderDict["lists"].([]interface{})
	if !ok {
		return
	}

	for _, value := range lists {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		name := item["name"].(string)

		host := s.getAdjustHost(HostPattern, name)

		//if !strings.HasPrefix(name, "http") {
		//	name = "http://" + name + ":80"
		//}
		if len(host) > 0 {
			hostArr = append(hostArr, host)
		}
	}
	logger.Debug("End fetch from couchdb")
	return
}

func (svr *CouchDBSourceImp) GetDBConnection() (db *couchdb.GoCouchDB, err error) {
	var dbname string
	s := strings.Split(svr.hostStr, "://")
	if len(s) < 2 {
		return nil, errors.New("statedb format error")
	}
	if s[0] != "couchdb" {
		return nil, errors.New("statedb format error")
	}
	auths := strings.Split(s[1], "@")

	if len(auths) < 2 {
		info := auths[0]
		split := strings.Split(info, "/")
		if len(split) < 2 {
			dbname = DefaultDBName
		} else {
			dbname = split[1]
		}
		db, err = couchdb.NewGoCouchDB(dbname, split[0],nil)
	} else {
		info := auths[0] // admin:password
		userandpass := strings.Split(info, ":")
		if len(userandpass) < 2 {
			hostandpath := auths[1]
			split := strings.Split(hostandpath, "/")
			if len(split) < 2 {
				dbname = DefaultDBName
			} else {
				dbname = split[1]
			}
			db, err = couchdb.NewGoCouchDB(dbname, split[0],nil)
		} else {
			auth := &couchdb.BasicAuth{Username: userandpass[0], Password: userandpass[1]}
			hostandpath := auths[1]
			split := strings.Split(hostandpath, "/")
			if len(split) < 2 {
				dbname = DefaultDBName
			} else {
				dbname = split[1]
			}
			db, err = couchdb.NewGoCouchDB(dbname, split[0], auth)
		}
	}
	//if err != nil {
	//	err = errors.New(fmt.Sprintf("cannot connect to couchdb, expect couchdb://xxxxx:5984 or couchdb://user:pass@xxxxx:5984, got %s", svr.hostStr))
	//}
	return
}

func (s *CouchDBSourceImp)getAdjustHost(pattern, name string) string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	host := reg.ReplaceAllString(s.urlreg, name)
	return host
}