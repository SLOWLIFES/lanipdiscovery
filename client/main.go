package main

import (
	"flag"
	"github.com/SLOWLIFES/lanipdiscovery"
)

func main() {
	var addr string
	var group string
	var name string
	flag.StringVar(&addr, "addr", "", "服务端域名,例如:http://a.com")
	flag.StringVar(&group, "group", "", "分组")
	flag.StringVar(&name, "name", "", "名")
	flag.Parse()

	l := lanipdiscovery.LanIpRegistrants{
		Addr:  addr,
		Group: group,
		Name:  name,
	}
	l.RegisteredSync()

}
