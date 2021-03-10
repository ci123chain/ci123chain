package redissource

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	r "github.com/ci123chain/ci123chain/pkg/redis"
	"github.com/ci123chain/ci123chain/sdk/domain"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

const Domain = "DOMAIN"
const SharedKey  = "s/k:order/OrderBook"
const HostPattern  = "[*]+"
func NewRedisSource(host, urlreg string) *RedisDBSourceImp {
	imp := &RedisDBSourceImp{
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

type RedisDBSourceImp struct {
	hostStr string
	urlreg  string
	conn    *r.RedisDB
}

func (s *RedisDBSourceImp) FetchSource() (hostArr []string) {
	logger.Debug("Start fetch from redisDB")
	if s.conn.DB.Client == nil {
		opt, err := getOption(s.hostStr)
		if err != nil {
			logger.Error("Connection Error: ", err)
		}
		conn := r.NewRedisDB(opt)
		s.conn = conn
	}

	bz := s.conn.Get([]byte(SharedKey))
	var shared map[string]interface{}
	err := json.Unmarshal(bz, &shared)
	if err != nil {
		logger.Error("fetch data from redisDB error: ", err)
	}

	orderDict, ok := shared["value"].(map[string]interface{})
	if !ok {
		return
	}
	lists, ok := orderDict["lists"].([]interface{})
	if !ok {
		return
	}

	//gateway.nodekey.suffix
	var host string
	for _, value := range lists {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		shardDomain := item["domain"].(string)
		if shardDomain != "" {
			host = s.getAdjustHost(HostPattern, shardDomain)
		} else {
			name := item["name"].(string)
			selfDomain := viper.GetString(Domain)
			if len(selfDomain) > 0 {
				name = domain.GetShardDomain(selfDomain, name)
			}
			host = s.getAdjustHost(HostPattern, name)
		}
		//if !strings.HasPrefix(name, "http") {
		//	name = "http://" + name + ":80"
		//}
		if len(host) > 0 {
			hostArr = append(hostArr, host)
		}
	}
	logger.Debug("End fetch from redis")
	return
}

func (s *RedisDBSourceImp) GetDBConnection() (db *r.RedisDB, err error) {
	opt, err := getOption(s.hostStr)
	if err != nil {
		return nil, err
	}
	db = r.NewRedisDB(opt)
	err = r.DBIsValid(db)
	return
}

func getOption(statedb string) (*redis.Options, error) {
	// redisdb://admin:password@192.168.2.89:5984
	// redisdb://192.168.2.89:5984
	s := strings.Split(statedb, "://")
	if len(s) < 2 {
		return nil, errors.New("redisdb format error")
	}
	if s[0] != "redisdb" {
		return nil, errors.New("redisdb format error")
	}
	auths := strings.Split(s[1], "@")

	if len(auths) < 2 { // 192.168.2.89:5984 无用户名 密码
		opt := &redis.Options{
			Addr: auths[0],
			DB:   0,
		}
		opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		return opt, nil
	} else { // admin:password@192.168.2.89:5984
		info := auths[0] // admin:password
		userandpass := strings.Split(info, ":")
		if len(userandpass) < 2 {
			opt := &redis.Options{
				Addr: auths[1],
				DB:   0,
			}
			opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			return opt, nil
		} else {
			opt := &redis.Options{
				Addr:               auths[1],
				Username:           userandpass[0],
				Password:           userandpass[1],
				DB:                 0,
			}
			opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			return opt, nil
		}
	}
}

func (s *RedisDBSourceImp) getAdjustHost(pattern, name string) string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	host := reg.ReplaceAllString(s.urlreg, name)
	return host
}