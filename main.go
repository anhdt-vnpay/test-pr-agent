package main

import (
	"fmt"
	"os"

	"github.com/blcvn/corev4-explorer/cmd"
)

func main() {
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}