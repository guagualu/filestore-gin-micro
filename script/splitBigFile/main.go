package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	// 读取命令行参数

	inputFile := "./script/tmp/ai融合_2023-03-26.mp4"
	chunkSize := 1024 * 1024
	// 打开输入文件
	in, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer in.Close()
	SplitFile(in, chunkSize)

}
func SplitFile(file *os.File, size int) {
	finfo, err := file.Stat()
	if err != nil {
		fmt.Println("get file info failed:", file, size)
	}
	fmt.Println(finfo, size)
	//每次最多拷贝1m
	bufsize := 1024 * 1024
	if size < bufsize {
		bufsize = size
	}
	buf := make([]byte, bufsize)
	num := (int(finfo.Size()) + size - 1) / size
	fmt.Println(num, len(buf))
	for i := 0; i < num; i++ {
		copylen := 0
		newfilename := strconv.Itoa(i) + finfo.Name()
		newfile, err1 := os.Create(newfilename)
		if err1 != nil {
			fmt.Println("failed to create file", newfilename)
		} else {
			fmt.Println("create file:", newfilename)
		}
		for copylen < size {
			n, err2 := file.Read(buf)
			if err2 != nil && err2 != io.EOF {
				fmt.Println(err2, "failed to read from:", file)
				break
			}
			if n <= 0 {
				break
			}
			//fmt.Println(n, len(buf))
			//写文件
			w_buf := buf[:n]
			newfile.Write(w_buf)
			copylen += n
		}
	}
	return
	//////每次最多拷贝1m
	//chunkSize := 1024 * 1024
	//finfo, _ := file.Stat()
	//for i := 0; ; i++ {
	//	// 每次都需要重新初始化，防止文件内容重复（需要优化）
	//	var buf = make([]byte, chunkSize) //此处chunkSize是读取分片的核心设置
	//	_, err := file.Read(buf)
	//	newfilename := strconv.Itoa(i) + finfo.Name()
	//	newfile, _ := os.Create(newfilename)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		panic(err)
	//	}
	//	newfile.Write(buf)
	//}
}
