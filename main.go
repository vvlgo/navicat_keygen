package main

import (
	"fmt"
	"navicat_keygen/keygen"
)

func main() {
	var (
		name     string
		organize string
		time     int
		code     string
	)

	fmt.Println("********************************************")
	fmt.Println("*       Navicat Keygen by @fei             *")
	fmt.Println("*             version 1.0                  *")
	fmt.Println("********************************************")

STEP1:
	fmt.Print("名称(随便填):")
	length, _ := fmt.Scanf("%s", &name)
	if length == 0 {
		fmt.Println("请输入名称")
		goto STEP1
	}
STEP2:
	fmt.Print("组织(随便填):")
	length, _ = fmt.Scanf("%s", &organize)
	if length == 0 {
		fmt.Println("请输入组织")
		goto STEP2
	}
STEP3:
	fmt.Print("时间戳(当前时间戳或-1天时间戳):")
	length, _ = fmt.Scanf("%d", &time)
	if length == 0 {
		fmt.Println("请输入时间戳")
		goto STEP3
	}
STEP4:
	fmt.Print("请求码(注意不要换行):")
	length, _ = fmt.Scanf("%s", &code)
	if length == 0 {
		fmt.Println("请输入请求码")
		goto STEP4
	}

	kg := keygen.NewKeygen("./public.pem", "./private.pem", name, organize, time)
	kg.GetActivationCode(code)

	select {}
}
