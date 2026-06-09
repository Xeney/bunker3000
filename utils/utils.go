package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func WaitEnter() {
	fmt.Println("\nНажмите Enter, чтобы продолжить...")
	_, _ = reader.ReadString('\n')
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func GetStringInput(message string) string {
	fmt.Printf("%s :> ", message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
