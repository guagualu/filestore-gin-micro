package util

import (
	"bytes"
	"os"
	"os/exec"
)

// 执行 linux shell command
func ExecLinuxShell(s string) (string, error) {
	//函数返回一个io.Writer类型的*Cmd
	cmd := exec.Command("/bin/bash", "-c", s)

	//通过bytes.Buffer将byte类型转化为string类型
	var result bytes.Buffer
	cmd.Stdout = &result

	//Run执行cmd包含的命令，并阻塞直至完成
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return result.String(), err
}

// 原理：1、程序创建一个sh文件 2、需要git自带的 git bash，并设置为bash环境变量 3、执行这个shell程序 执行文件合并命令
func ExecWinShell(s string, shText string) (string, string, error) {
	// 创建文件
	file, err := os.Create("merge.sh")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// 写入内容
	_, err = file.WriteString("#!/bin/bash\n")
	if err != nil {
		return "", "", err
	}
	_, err = file.WriteString(shText)
	if err != nil {
		return "", "", err
	}
	//函数返回一个io.Writer类型的*Cmd
	cmd := exec.Command("bash", s)

	//通过bytes.Buffer将byte类型转化为string类型
	var result bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &result
	cmd.Stderr = &stderr
	//
	////Run执行cmd包含的命令，并阻塞直至完成
	err = cmd.Run()
	//out, err := cmd.CombinedOutput()
	if err != nil {
		return "", stderr.String(), err
	}

	return result.String(), stderr.String(), err
	//return string(out), nil
}
