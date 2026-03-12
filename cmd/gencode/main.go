package main

import (
	"ginp-api/cmd/gencode/cmd"
	"os"
)

func main() {
	// 检测是否无参数运行，如果是则进入交互模式
	if len(os.Args) == 1 {
		cmd.ShowMainMenu()
		return
	}
	cmd.Execute()
}
