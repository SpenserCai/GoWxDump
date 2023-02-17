package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// 一个map，key是版本号，value一个list，list里面偏移的10进制
var OffSetMap = map[string][]int{
	"3.2.1.15": {
		328121948,
		328122328,
		328123056,
		328121976,
		328123020,
	},
	"3.3.0.115": {
		31323364,
		31323744,
		31324472,
		31323392,
		31324436,
	},
	"3.3.0.84": {
		31315212,
		31315592,
		31316320,
		31315240,
		31316284,
	},
	"3.3.0.93": {
		31323364,
		31323744,
		31324472,
		31323392,
		31324436,
	},
	"3.3.5.34": {
		30603028,
		30603408,
		30604120,
		30603056,
		30604100,
	},
	"3.3.5.42": {
		30603012,
		30603392,
		30604120,
		30603040,
		30604084,
	},
	"3.3.5.46": {
		30578372,
		30578752,
		30579480,
		30578400,
		30579444,
	},
	"3.4.0.37": {
		31608116,
		31608496,
		31609224,
		31608144,
		31609188,
	},
	"3.4.0.38": {
		31604044,
		31604424,
		31605152,
		31604072,
		31605116,
	},
	"3.4.0.50": {
		31688500,
		31688880,
		31689608,
		31688528,
		31689572,
	},
	"3.4.0.54": {
		31700852,
		31701248,
		31700920,
		31700880,
		31701924,
	},
	"3.4.5.27": {
		32133788,
		32134168,
		32134896,
		32133816,
		32134860,
	},
	"3.4.5.45": {
		32147012,
		32147392,
		32147064,
		32147040,
		32148084,
	},
	"3.5.0.20": {
		35494484,
		35494864,
		35494536,
		35494512,
		35495556,
	},
	"3.5.0.29": {
		35507980,
		35508360,
		35508032,
		35508008,
		35509052,
	},
	"3.5.0.33": {
		35512140,
		35512520,
		35512192,
		35512168,
		35513212,
	},
	"3.5.0.39": {
		35516236,
		35516616,
		35516288,
		35516264,
		35517308,
	},
	"3.5.0.42": {
		35512140,
		35512520,
		35512192,
		35512168,
		35513212,
	},
	"3.5.0.44": {
		35510836,
		35511216,
		35510896,
		35510864,
		35511908,
	},
	"3.5.0.46": {
		35506740,
		35507120,
		35506800,
		35506768,
		35507812,
	},
	"3.6.0.18": {
		35842996,
		35843376,
		35843048,
		35843024,
		35844068,
	},
	"3.6.5.7": {
		35864356,
		35864736,
		35864408,
		35864384,
		35865428,
	},
	"3.6.5.16": {
		35909428,
		35909808,
		35909480,
		35909456,
		35910500,
	},
	"3.7.0.26": {
		37105908,
		37106288,
		37105960,
		37105936,
		37106980,
	},
	"3.7.0.29": {
		37105908,
		37106288,
		37105960,
		37105936,
		37106980,
	},
	"3.7.0.30": {
		37118196,
		37118576,
		37118248,
		37118224,
		37119268,
	},
	"3.7.5.11": {
		37883280,
		37884088,
		37883136,
		37883008,
		37884052,
	},
	"3.7.5.23": {
		37895736,
		37896544,
		37895592,
		37883008,
		37896508,
	},
	"3.7.5.27": {
		37895736,
		37896544,
		37895592,
		37895464,
		37896508,
	},
	"3.7.5.31": {
		37903928,
		37904736,
		37903784,
		37903656,
		37904700,
	},
	"3.7.6.24": {
		38978840,
		38979648,
		38978696,
		38978604,
		38979612,
	},
	"3.7.6.29": {
		38986376,
		38987184,
		38986232,
		38986104,
		38987148,
	},
	"3.7.6.44": {
		39016520,
		39017328,
		39016376,
		38986104,
		39017292,
	},
	"3.8.0.31": {
		46064088,
		46064912,
		46063944,
		38986104,
		46064876,
	},
	"3.8.0.33": {
		46059992,
		46060816,
		46059848,
		38986104,
		46060780,
	},
	"3.8.0.41": {
		46064024,
		46064848,
		46063880,
		38986104,
		46064812,
	},
	"3.8.1.26": {
		46409448,
		46410272,
		46409304,
		38986104,
		46410236,
	},
	"3.9.0.28": {
		48418376,
		48419280,
		48418232,
		38986104,
		48419244,
	},
}

// 获取微信进程对象，包含进程ID、进程句柄和Module列表
func GetWeChatProcess() (windows.ProcessEntry32, error) {
	var process windows.ProcessEntry32
	process.Size = uint32(unsafe.Sizeof(process))
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return process, err
	}
	defer windows.CloseHandle(snapshot)
	for {
		err = windows.Process32Next(snapshot, &process)
		if err != nil {
			return process, err
		}
		if windows.UTF16ToString(process.ExeFile[:]) == "WeChat.exe" {
			return process, nil
		}
	}
}

// 获取微信进程的WeChatWin.dll模块对象，包含模块基址、模块大小和模块路径()
func GetWeChatWinModule(process windows.ProcessEntry32) (windows.ModuleEntry32, error) {
	var module windows.ModuleEntry32
	module.Size = uint32(unsafe.Sizeof(module))
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, process.ProcessID)
	if err != nil {
		return module, err
	}
	defer windows.CloseHandle(snapshot)
	for {
		err = windows.Module32Next(snapshot, &module)
		if err != nil {
			return module, err
		}
		if windows.UTF16ToString(module.Module[:]) == "WeChatWin.dll" {
			return module, nil
		}
	}
}

// 通过模块获取版本号 c#代码为：string FileVersion = processModule.FileVersionInfo.FileVersion;转成go代码如下
func GetVersion(module windows.ModuleEntry32) (string, error) {
	image, imgErr := windows.LoadLibraryEx(windows.UTF16ToString(module.ExePath[:]), 0, windows.LOAD_LIBRARY_AS_DATAFILE)
	if imgErr != nil {
		return "", fmt.Errorf("LoadLibraryEx error: %v", imgErr)
	}
	resInfo, infoErr := windows.FindResource(image, windows.ResourceID(1), windows.RT_VERSION)
	if infoErr != nil {
		return "", fmt.Errorf("FindResource error: %v", infoErr)
	}
	resData, dataErr := windows.LoadResourceData(image, resInfo)
	if dataErr != nil {
		return "", fmt.Errorf("LoadResourceData error: %v", dataErr)
	}
	var info *windows.VS_FIXEDFILEINFO
	size := uint32(unsafe.Sizeof(*info))
	err := windows.VerQueryValue(unsafe.Pointer(&resData[0]), `\`, unsafe.Pointer(&info), &size)
	if err != nil {
		return "", fmt.Errorf("VerQueryValue error: %v", err)
	}
	// 从低位到高位，分别为主版本号、次版本号、修订号、编译号
	version := fmt.Sprintf("%d.%d.%d.%d", info.FileVersionMS>>16, info.FileVersionMS&0xffff, info.FileVersionLS>>16, info.FileVersionLS&0xffff)
	return version, nil
}

// 获取微信数据：入参为微信进程句柄，偏移地址，返回值为昵称和错误信息
func GetWeChatData(process windows.Handle, offset uintptr, nSize int) (string, error) {
	var buffer = make([]byte, nSize)
	err := windows.ReadProcessMemory(process, offset, &buffer[0], uintptr(nSize), nil)
	if err != nil {
		return "", err
	}
	text := ""
	for _, v := range buffer {
		if v == 0 {
			break
		}
		text += string(v)
	}
	return text, nil
}

// 获取微信key：入参为微信进程句柄，偏移地址，返回值为key和错误信息
func GetWeChatKey(process windows.Handle, offset uintptr) (string, error) {
	var buffer = make([]byte, 4)
	err := windows.ReadProcessMemory(process, offset, &buffer[0], 4, nil)
	if err != nil {
		return "", err
	}
	var num = 32
	var buffer2 = make([]byte, num)
	// c# 代码(IntPtr)(((int)array[3] << 24) + ((int)array[2] << 16) + ((int)array[1] << 8) + (int)array[0]);转成go代码如下
	offset2 := uintptr((int(buffer[3]) << 24) + (int(buffer[2]) << 16) + (int(buffer[1]) << 8) + int(buffer[0]))
	err = windows.ReadProcessMemory(process, offset2, &buffer2[0], uintptr(num), nil)
	if err != nil {
		return "", err
	}
	// 将byte数组转成hex字符串，并转成大写
	key := hex.EncodeToString(buffer2)
	key = strings.ToUpper(key)
	return key, nil

}

func main() {
	processAllAccess := uint32(
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
	process, err := GetWeChatProcess()
	if err != nil {
		fmt.Println("GetWeChatProcess error: ", err)
		return
	}
	wechatProcessHandle, err := windows.OpenProcess(processAllAccess, false, process.ProcessID)
	if err != nil {
		fmt.Println("OpenProcess error: ", err)
		return
	}
	module, err := GetWeChatWinModule(process)
	if err != nil {
		fmt.Println("GetWeChatWinModule error: ", err)
		return
	}
	version, err := GetVersion(module)
	if err != nil {
		fmt.Println("GetVersion error: ", err)
		return
	}
	fmt.Println("WeChat Version: ", version)
	nickName, err := GetWeChatData(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][0]), 100)
	if err != nil {
		fmt.Println("GetWeChatNickName error: ", err)
		return
	}
	account, err := GetWeChatData(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][1]), 100)
	if err != nil {
		fmt.Println("GetWeChatAccount error: ", err)
		return
	}
	mobile, err := GetWeChatData(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][2]), 100)
	if err != nil {
		fmt.Println("GetWeChatMobile error: ", err)
		return
	}
	key, err := GetWeChatKey(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][4]))
	if err != nil {
		fmt.Println("GetWeChatKey error: ", err)
		return
	}
	fmt.Println("WeChat NickName: ", nickName)
	fmt.Println("WeChat Account: ", account)
	fmt.Println("WeChat Mobile: ", mobile)
	fmt.Println("WeChat Key: ", key)

}
