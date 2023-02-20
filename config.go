/*
 * @Author: SpenserCai
 * @Date: 2023-02-20 10:36:15
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-20 10:39:49
 * @Description: file content
 */
package main

import "golang.org/x/sys/windows"

// 定义微信数据结构
type WeChatData struct {
	Version  string
	NickName string
	Account  string
	Mobile   string
	Key      string
}

var PROCESS_ALL_ACCESS = uint32(
	windows.PROCESS_QUERY_INFORMATION |
		windows.PROCESS_VM_READ |
		windows.PROCESS_VM_WRITE |
		windows.PROCESS_VM_OPERATION |
		windows.PROCESS_CREATE_THREAD |
		windows.PROCESS_DUP_HANDLE |
		windows.PROCESS_TERMINATE |
		windows.PROCESS_SUSPEND_RESUME |
		windows.PROCESS_SET_QUOTA |
		windows.PROCESS_SET_INFORMATION |
		windows.PROCESS_QUERY_LIMITED_INFORMATION)
