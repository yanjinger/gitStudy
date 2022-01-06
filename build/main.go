package main

import (
	"fmt"
	"github.com/yanjinger/gitStudy/build/sayArch"
	"github.com/yanjinger/gitStudy/build/sayOs"
	"runtime"
)

func main() {
	fmt.Println("Os=", runtime.GOOS)
	sayOs.SayOs()
	fmt.Println("Arch=", runtime.GOARCH)
	sayArch.SayArch()
}
