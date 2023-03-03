<!--
 * @Author: SpenserCai
 * @Date: 2023-02-17 18:04:27
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-03-03 16:34:02
 * @Description: file content
-->
# GoWxDump
## 特别说明
GoWxDump是<a href="https://github.com/AdminTest0/SharpWxDump">SharpWxDump</a>的Go语言版本。
同时添加了一些新的功能。
## 使用方法
### 1.下载
```
git clone https://github.com/SpenserCai/GoWxDump.git
```
### 2.编译
需要安装mingw-w32
```
build.bat
```
### 3.使用
```
GoWxDump.exe
```
## GoWxDump原创功能
### 1.支持获取数据目录
### 2.支持自动解密
由AdminTest0发布的<a href="https://mp.weixin.qq.com/s/4DbXOS5jDjJzM2PN0Mp2JA">解密脚本</a>翻译成Go语言而来，支持自动解密。
### 3.支持交互式命令
```bash
show_info 获取微信基础信息
decrypt 解密数据
friends_list 获取好友列表 （目前支持：获取最近十个聊天的好友信息，需要解密后才能获取）
```
### 4.非交互式命令
```bash
GoWxDump.exe -spy
```
## 免责声明
本项目仅允许在授权情况下对数据库进行备份，严禁用于非法目的，否则自行承担所有相关责任。使用该工具则代表默认同意该条款;

请勿利用本项目的相关技术从事非法测试，如因此产生的一切不良后果与项目作者无关。