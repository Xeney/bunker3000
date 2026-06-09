package utils

import (
	"testing"
)

func TestGetStringInput(t *testing.T) {
	// Примечание: Этот тест требует взаимодействия с пользователем
	// В реальном проекте лучше использовать моки или интеграционные тесты
	// Это пример того, как можно структурировать тест

	// Для автоматического тестирования рекомендуется рефакторинг
	// с использованием интерфейсов для ввода/вывода

	t.Skip("Skipping interactive test - requires manual input or mocking")
}

func TestClear(t *testing.T) {
	// Тестирование очистки консоли - сложно автоматизировать
	// Рекомендуется проверить, что функция не паникует

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Clear function panicked: %v", r)
		}
	}()

	// Просто проверяем, что функция выполняется без паники
	Clear()
}

func TestWaitEnter(t *testing.T) {
	// Аналогично, тест требует взаимодействия с пользователем
	t.Skip("Skipping interactive test - requires manual input")
}
