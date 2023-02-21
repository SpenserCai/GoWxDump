/*
 * @Author: SpenserCai
 * @Date: 2023-02-20 16:23:37
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-21 12:16:13
 * @Description: file content
 */
package main

import (
	"fmt"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

var suggestions = []prompt.Suggest{
	{Text: "show_info", Description: "获取微信基础信息"},
	{Text: "decrypt", Description: "解密数据"},
	{Text: "friends_list", Description: "获取好友列表"},
	{Text: "help", Description: "帮助"},
	{Text: "exit", Description: "退出程序"},
}

func completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func executor(cmd string) {
	cmd = strings.TrimSpace(cmd)
	blocks := strings.Split(cmd, " ")
	if len(blocks) == 0 {
		return
	}
	switch blocks[0] {
	case "show_info":
		ShowInfoCmd()
	case "decrypt":
		DecryptCmd()
	case "friends_list":
		FriendsListCmd()
	case "help":
		// 显示命令的帮助信息
		for _, v := range suggestions {
			fmt.Println(v.Text, v.Description)
		}
	case "exit":
		// 退出命令模式并关闭程序
		return
	default:
		fmt.Println("Unknown command")
	}
}

func RunCommand() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>"),
		prompt.OptionTitle("GoWxDump"),
	)
	p.Run()
}
