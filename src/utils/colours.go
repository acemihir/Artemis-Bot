package utils

import (
	"fmt"
	"os"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

func Cout(text string, colour string, params ...interface{}) {
	if len(params) == 0 {
		fmt.Println(colour + text + Reset)
		if strings.Contains(text, "[ERROR]") {
			os.Exit(1)
		}
	} else {
		fmt.Printf(colour+text+Reset+"\n", params...)
		if strings.Contains(text, "[ERROR]") {
			os.Exit(1)
		}
	}
}
