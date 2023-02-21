/*
 * @Author: SpenserCai
 * @Date: 2023-02-21 11:35:50
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-21 14:02:29
 * @Description: file content
 */
package db

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FriendInfo struct {
	NickName string
	Alias    string
	UserName string
	Remark   string
}

// 创建一个数据库操作对象

type WeChatDb struct {
	// 数据库对象
	Db *gorm.DB
}

// 初始化数据库，入参为数据库文件路径
func (w *WeChatDb) InitDb(dbPath string) error {
	// 判断文件是否存在
	_, err := os.Stat(dbPath)
	if err != nil {
		return err
	}
	// 打开数据库
	w.Db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

// 对象被销毁时，关闭数据库连接
func (w *WeChatDb) Close() error {
	db, err := w.Db.DB()
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

// 查询最近聊天的好友
func (w *WeChatDb) GetNearChatFriends(topNumber int) ([]string, error) {
	// 查询ChatInfo表，UserName中不含有@的，LastReadedCreateTime最大的前topNumber条记录，只需要UserName字段
	var userNameList []string
	err := w.Db.Table("ChatInfo").Select("UserName").Where("UserName NOT LIKE ?", "%@%").Order("LastReadedCreateTime DESC").Limit(topNumber).Find(&userNameList).Error
	if err != nil {
		return nil, err
	}
	return userNameList, nil
}

// 通过UserName列表查询好友详细信息
func (w *WeChatDb) GetFriendInfoListWithUserList(user []string) ([]FriendInfo, error) {
	var friendInfoList []FriendInfo
	err := w.Db.Table("Contact").Where("UserName IN ?", user).Find(&friendInfoList).Error
	if err != nil {
		return nil, err
	}
	return friendInfoList, nil
}
