package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"
	"xiaoLong/could"
	"xiaoLong/xc"
)

func type1() {
	fmt.Println("请输入保存cid的文件名，不需要输入.txt")
	var saveFileName string
	_, err := fmt.Scanln(&saveFileName)
	if err != nil {
		fmt.Println("输入错误", err)
		ScanQuit()
	}
	file, err := os.Create(saveFileName + ".txt")
	if err != nil {
		fmt.Println("创建文件失败", err)
		ScanQuit()
	}
	fmt.Println("开始运行")
	startTime := time.Now()
	defer file.Close()
	xc.Wg.Add(1)
	go could.XiaoLong.GetFileSuperList(file, "0", "", 0)
	xc.Wg.Wait()
	for _, i := range xc.FileXc {
		i.Close()
	}
	if could.XiaoLong.Count == 0 {
		fmt.Println("获取文件列表失败")
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("运行完毕，共%d个资源，耗时%s\n", could.XiaoLong.Count, elapsedTime)
}
func type2() {
	_, err := os.Stat("cid")
	if os.IsNotExist(err) {
		err := os.Mkdir("cid", 755)
		if err != nil {
			fmt.Println("创建文件夹失败")
			ScanQuit()
		}
	}
	file, err := os.Create("cid\\root.txt")
	if err != nil {
		fmt.Println("文件创建失败", err)
		ScanQuit()
	}
	fmt.Println("开始运行")
	startTime := time.Now()
	defer file.Close()
	xc.Wg.Add(1)
	go could.XiaoLong.GetFileSuperList(file, "0", "", 0)
	xc.Wg.Wait()
	for _, i := range xc.FileXc {
		i.Close()
	}
	if could.XiaoLong.Count == 0 {
		fmt.Println("获取文件列表失败")
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("运行完毕，共%d个资源，耗时%s\n", could.XiaoLong.Count, elapsedTime)
}

func main() {
	fmt.Println("作者：BY易仝\nQQ：1944876825\n功能：批量保存小龙云盘cid\n注意：防止闪退看不到日志，建议用命令行打开")
	fmt.Println("请选择保存模式\n1.所有cid保存在一个txt中\n2.一个文件夹对应一个txt")
	var saveType int
	_, err := fmt.Scanln(&saveType)
	if err != nil {
		fmt.Println("输入错误", err)
		ScanQuit()
	}
	could.XiaoLong.SaveType = saveType
	fmt.Println("请输入爬取线程，单线程输入1，（指的是同时爬取多少个文件夹）")
	var maxPros int
	_, err = fmt.Scanln(&maxPros)
	if err != nil {
		fmt.Println("输入错误", err)
		ScanQuit()
	}
	runtime.GOMAXPROCS(maxPros)
	var forMat string
	fmt.Println("请输入数据保存格式，不填走默认，可选参数如下，参数两边必须加{}，如{FileName}，暂时没有做保存为Excell，不过你可以直接复制生成的内容，粘贴在Excell中，前提是分隔符需要是\\t")
	var arr = []could.FormatItem{
		{
			Key:   "FileName",
			Title: "文件名称",
		},
		{
			Key:   "FileCid",
			Title: "文件cid",
		},
		{
			Key:   "FileDir",
			Title: "文件路径",
		},
		{
			Key:   "FileSize",
			Title: "文件大小",
		},
		{
			Key:   "Cover",
			Title: "文件封面cid",
		},
		{
			Key:   "Suffix",
			Title: "扩展名",
		},
	}
	for _, i := range arr {
		fmt.Print(i.Key + "\t")
	}
	fmt.Print("\n")
	for _, i := range arr {
		fmt.Print(i.Title + "\t")
	}
	fmt.Print("\n")
	Scanf(&forMat)
	if forMat != "" {
		could.XiaoLong.ForMatText = forMat
	}
	fmt.Println("数据保存格式：", could.XiaoLong.ForMatText)
	var Ak string
	fmt.Println("请输入您的小龙云密钥")
	_, err = fmt.Scanln(&Ak)
	if err != nil {
		fmt.Println("输入错误", err)
		ScanQuit()
	}
	if could.XiaoLong.Login(Ak) == false {
		ScanQuit()
	}
	could.XiaoLong.Count = 0

	if saveType == 1 {
		type1()
	}
	if saveType == 2 {
		type2()
	}
	ScanQuit()
}
func ScanQuit() {
	fmt.Println("按下任意建退出...")
	_, _ = fmt.Scanln()
	os.Exit(0)
}

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}
