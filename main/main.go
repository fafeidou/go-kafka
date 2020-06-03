package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	// 获取程序目录
	strDir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
	if errDir != nil {
		log.Fatal(errDir)
	}

	strFile := strDir + "/testFile.txt"
	fmt.Println("file to open is ", strFile)

	// 打开文件，如果没有，那么创建，设置为读写操作，文件权限为755
	file, errOpenFile := os.OpenFile(strFile, os.O_RDWR|os.O_CREATE, 755)
	if errOpenFile != nil {
		fmt.Println("open file fail")
		log.Fatal(errOpenFile)
	}

	// 保证文件能够被关闭
	defer file.Close()

	// 写文件
	for i := 1; i < 10; i++ {
		strToWrite := "Line " + strconv.Itoa(i) + "\r\n"
		file.WriteString(strToWrite)
	}

	szbuf := make([]byte, 1024)

	// 读文件
	// 由于对file进行的写操作已经将文件指针偏移到文件末尾
	// 如果此处用Read，会导致失败，因此使用ReadAt，且从文件开头读取
	nRead, errRead := file.ReadAt(szbuf, 0)
	if errRead != nil && errRead != io.EOF {
		fmt.Println("read file fail")
		log.Fatal(errRead)
	}

	fmt.Printf("read size:%d, content：%s", nRead, szbuf)
}
