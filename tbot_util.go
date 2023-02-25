/*
 * @Author: SpenserCai
 * @Date: 2023-02-23 15:16:57
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-25 13:08:58
 * @Description: file content
 */
package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// 发送纯文本消息
func TeleSendMessge(message string) {
	TelBot.Send(tele.ChatID(-TELBOT_CHAT_ID), message)
}

// 发送支持文件和markdown的消息
func TeleSendMarkDownMessage(message string) {
	// 发送markdown
	TelBot.Send(tele.ChatID(-TELBOT_CHAT_ID), message, &tele.SendOptions{
		ParseMode: tele.ModeMarkdown,
	})

}

// 同时发送文件和markdown消息
func TeleSendFileAndMessage(message string, fileLocalList []string) {
	// 判断ANONFILES_TOKEN是否为空
	if ANONFILES_TOKEN == "" {
		TeleSendMarkDownMessage(message)
		return
	}
	for _, fileLocal := range fileLocalList {
		// 上传文件
		fileUrl, err := AnonFilesUpload(fileLocal)
		if err != nil {
			TeleSendMarkDownMessage(message)
			return
		}
		message += strings.Split(fileLocal, "\\")[len(strings.Split(fileLocal, "\\"))-1] + ":" + fileUrl + " \n"
	}
	TeleSendMarkDownMessage(message)

}
