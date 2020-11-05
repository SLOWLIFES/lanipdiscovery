package lanipdiscovery

import (
	"fmt"
	"log"
	"net"
	"time"
)

type LanIpRegistrants struct {
	Addr  string `json:"addr"`
	Group string `json:"group"`
	Name  string `json:"name"`
}

func (l *LanIpRegistrants) RegisteredSync() {
	l.Registered()
	select {}
}

func (l *LanIpRegistrants) Registered() {
	go func() {
		defer func() {
			recover()
			time.Sleep(time.Second * 10)
			l.Registered()
		}()

		for {
			ips := getLocalIP()
			for e := range ips {

				httpdo := HttpDo{
					Url: fmt.Sprintf("%s/registered", l.Addr),
					Data: map[string]string{
						"name":   l.Name,
						"group":  l.Group,
						"lan_ip": ips[e],
					},
				}
				_, err := httpdo.Post()
				if err != nil {
					log.Println(err)
				}
				//log.Println("data:", l.Group, string(data))

			}
			time.Sleep(time.Minute / 2)
		}
	}()
}

func getLocalIP() (ips []string) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
		err     error
	)

	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String()) // 192.168.1.1
				return
			}
		}
	}

	return
}
