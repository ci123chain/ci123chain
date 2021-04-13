package redissource

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	r "github.com/ci123chain/ci123chain/pkg/redis"
	"github.com/go-redis/redis/v8"
	"regexp"
	"strings"
)

const (
	HostPattern  = "[*]+"
	FlagNodeList   = "node-list"
)

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

	bz, err := s.conn.Get([]byte(FlagNodeList))
	if err != nil {
		logger.Error("db connection get failed: ", err)
	}
	var node_list []string
	err = json.Unmarshal(bz, &node_list)
	if err != nil {
		logger.Error("fetch data from redisDB error: ", err)
	}

	var host string
	for _, value := range node_list {
		host = s.getAdjustHost(HostPattern, value)
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
	// redisdb://admin:password@192.168.2.89:11001@tls
	// redisdb://192.168.2.89:11001@tls
	s := strings.Split(statedb, "://")
	if len(s) < 2 {
		return nil, errors.New("redisdb format error")
	}
	if s[0] != "redisdb" {
		return nil, errors.New("redisdb format error")
	}
	auths := strings.Split(s[1], "@")

	if len(auths) < 2 { // 192.168.2.89:11001#tls 无用户名 密码
		pd := strings.Split(auths[0], "#")
		opt := &redis.Options{
			Addr: pd[0],
			DB:   0,
		}
		if len(pd) == 1 {
			return opt, nil
		}
		if len(pd) == 2 {
			if pd[1] == "tls" {
				opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
				return opt, nil
			}else {
				return nil, errors.New(fmt.Sprintf("unexpected tls setting: %s", statedb))
			}
		}else {
			return nil, errors.New(fmt.Sprintf("unexpected setting of db tls: %v", statedb))
		}
	} else { // admin:password@192.168.2.89:5984#tls
		info := auths[0] // admin:password
		userandpass := strings.Split(info, ":")
		if len(userandpass) < 2 {
			return nil, errors.New(fmt.Sprintf("unexpected setting of username and password %s", statedb))
		} else {
			pd := strings.Split(auths[1], "#")
			opt := &redis.Options{
				Addr:               pd[0],
				Username:           userandpass[0],
				Password:           userandpass[1],
				DB:                 0,
			}
			if len(pd) == 1 {
				return opt, nil
			}
			if len(pd) == 2 {
				if pd[1] == "tls" {
					opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
					return opt, nil
				}else {
					return nil, errors.New(fmt.Sprintf("unexpected tls setting: %s", statedb))
				}
			}else {
				return nil, errors.New(fmt.Sprintf("unexpected setting of db tls: %v", statedb))
			}
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