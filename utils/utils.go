package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func WaitEnter() {
	fmt.Println("\nНажмите Enter, чтобы продолжить...")
	_, _ = reader.ReadString('\n')
}

func Clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func GetStringInput(message string) string {
	fmt.Printf("%s :> ", message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func GetIntInput(minVal, maxVal int) int {
	for {
		fmt.Printf("\n=> Введите число (%d-%d): ", minVal, maxVal)
		var val int
		_, err := fmt.Scan(&val)
		var trailing string
		fmt.Scanln(&trailing)
		if err == nil && val >= minVal && val <= maxVal {
			return val
		}
		fmt.Println("Ошибка: введите число от", minVal, "до", maxVal)
	}
}
