package domain

import "strings"

const Gateway_prefix = "gateway"

func GetShardDomain(gateway string, shardName string) (domain string) {
	domain = strings.Replace(gateway, Gateway_prefix, shardName, 1)
	return
}
