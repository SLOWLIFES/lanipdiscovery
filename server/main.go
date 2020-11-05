package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/muesli/cache2go"
	"github.com/wailovet/osmanthuswine"
	"github.com/wailovet/osmanthuswine/src/core"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func main() {
	config := core.GetInstanceConfig()
	config.Host = "0.0.0.0"
	flag.StringVar(&config.Port, "port", "8808", "端口号")
	flag.Parse()
	osmanthuswine.HandleFunc("/get", func(request core.Request, response core.Response) {

		group := request.REQUEST["group"]
		if group == "" {
			group = "lanipdiscovery"
		}

		ip, err := parseIP(request.OriginRequest.RemoteAddr)
		response.CheckErrDisplayByError(err)

		result := map[string]interface{}{}
		getCache(fmt.Sprintf("[%s]-[%s]", group, ip)).Foreach(func(key interface{}, item *cache2go.CacheItem) {
			result[key.(string)] = item.Data()
		})

		response.DisplayByData(result)

	})

	osmanthuswine.HandleFunc("/registered", func(request core.Request, response core.Response) {
		group := request.REQUEST["group"]
		if group == "" {
			group = "lanipdiscovery"
		}

		name := request.REQUEST["name"]
		if name == "" {
			response.DisplayByError("name is empty", 500)
		}

		lanIp := request.REQUEST["lan_ip"]
		if lanIp == "" {
			response.DisplayByError("lan ip is empty", 500)
		}

		ip, err := parseIP(request.OriginRequest.RemoteAddr)
		response.CheckErrDisplayByError(err)

		lip := LanIp{
			Name: name,
			Ip:   lanIp,
		}
		getCache(fmt.Sprintf("[%s]-[%s]", group, ip)).Add(lip.Name, time.Minute, lip.Ip)
		response.DisplayByData(lip)
	})
	osmanthuswine.Run()
}

func parseIP(addr string) (string, error) {
	log.Println("addr:", addr)
	ip := ""
	tmp := strings.Split(addr, ":")
	if len(tmp) > 0 {
		ip = tmp[0]
	}
	wanIp := net.ParseIP(ip)
	if wanIp == nil {
		return "", errors.New("get wan ip error:" + addr)
	}

	return ip, nil
}

var cacheMap sync.Map

func getCache(key string) *cache2go.CacheTable {

	value, ok := cacheMap.Load(key)
	if ok {
		return value.(*cache2go.CacheTable)
	}
	cache := cache2go.Cache(key)
	cacheMap.Store(key, cache)
	return cache
}
