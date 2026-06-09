package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var reader = bufio.NewReader(os.Stdin)

// WaitEnter ожидает нажатия Enter
func WaitEnter() {
	fmt.Println("\n📌 Нажмите Enter, чтобы продолжить...")
	_, _ = reader.ReadString('\n')
}

// Clear очищает экран (работает на Windows, Linux, MacOS)
func Clear() {
	// Определяем операционную систему
	switch runtime.GOOS {
	case "windows":
		// Для Windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		// Для Linux, MacOS и других Unix-like систем
		fmt.Print("\033[H\033[2J")
	}
}

// ClearLine очищает текущую строку
func ClearLine() {
	fmt.Print("\r\033[K")
}

// GetStringInput получает строковый ввод от пользователя
func GetStringInput(message string) string {
	fmt.Printf("%s :> ", message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// GetIntInput получает числовой ввод от пользователя
func GetIntInput(message string) int {
	var result int
	fmt.Printf("%s :> ", message)
	fmt.Scan(&result)
	return result
}

// PrintSlow печатает текст с эффектом печати (для атмосферы)
func PrintSlow(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(delay)
	}
	fmt.Println()
}

// ConfirmAction запрашивает подтверждение действия
func ConfirmAction(message string) bool {
	fmt.Printf("%s (y/n): ", message)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes" || input == "да"
}
