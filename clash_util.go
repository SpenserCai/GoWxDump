/*
 * @Author: SpenserCai
 * @Date: 2023-02-23 17:29:55
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-23 22:50:42
 * @Description: file content
 */
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Dreamacro/clash/adapter/outbound"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/listener/socks"
)

func RunClashClient() {
	// 将CLASH用:分割 格式type:server:port:cipher:password:udp
	clashConfig := strings.Split(CLASH_CONN_STR, ":")
	clashType := clashConfig[0]
	clashServer := clashConfig[1]
	clashPort, err := strconv.Atoi(clashConfig[2])
	clashCipher := clashConfig[3]
	clashPassword := clashConfig[4]
	clashUDP, err := strconv.ParseBool(clashConfig[5])
	if err != nil {
		fmt.Println("clashFormat error: ", err)
		return
	}
	in := make(chan constant.ConnContext, 100)
	defer close(in)
	var l *socks.Listener
	// 从LOCAL_PROXY_PORT开始尝试如果报错则+1
	for {
		tmpl, err := socks.New("127.0.0.1:"+strconv.Itoa(LOCAL_PROXY_PORT), in)
		if err == nil {
			l = tmpl
			break
		}
		// 等待100ms
		time.Sleep(100 * time.Millisecond)
		LOCAL_PROXY_PORT++
	}
	defer l.Close()

	println("listen at:", l.Address())

	proxy, err := outbound.NewShadowSocks(
		outbound.ShadowSocksOption{
			Name:     clashType,
			Server:   clashServer,
			Port:     clashPort,
			Cipher:   clashCipher,
			Password: clashPassword,
			UDP:      clashUDP,
		},
	)
	if err != nil {
		panic(err)
	}

	for c := range in {
		conn := c
		metadata := conn.Metadata()
		fmt.Printf("request incoming from %s to %s\n", metadata.SourceAddress(), metadata.RemoteAddress())
		go func() {
			remote, err := proxy.DialContext(context.Background(), metadata)
			if err != nil {
				fmt.Printf("dial error: %s\n", err.Error())
				return
			}
			relay(remote, conn.Conn())
		}()
	}
}

func relay(l, r net.Conn) {
	go io.Copy(l, r)
	io.Copy(r, l)
}
