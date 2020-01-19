package couchdbsource

import (
	"fmt"
	"regexp"
	"testing"
)

func TestUrlReg(t *testing.T)  {
	urlreg := "******"
	if ok, err :=  regexp.MatchString("[*]+", urlreg); !ok {
		panic(err)
	}


	reg, err := regexp.Compile(HostPattern)
	if err != nil {
		panic(err)
	}
	host := reg.ReplaceAllString("afdaf:*******", "hahaha")
	fmt.Println(host)
}