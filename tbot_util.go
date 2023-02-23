/*
 * @Author: SpenserCai
 * @Date: 2023-02-23 15:16:57
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-23 22:42:51
 * @Description: file content
 */
package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

func InitBot() {
	tempBot, err := tele.NewBot(tele.Settings{
		Token:  TELBOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		Client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(&url.URL{
					Scheme: "socks5",
					Host:   "127.0.0.1:" + strconv.Itoa(LOCAL_PROXY_PORT),
				}),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	TelBot = tempBot
}

func SendMessge(message string) {
	TelBot.Send(tele.ChatID(-TELBOT_CHAT_ID), message)
}
