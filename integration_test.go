package main

import (
	"bunker3000/events"
	"bunker3000/player"
	"testing"
)

func TestIntegration_FullDayCycle(t *testing.T) {
	p := player.CreatePlayer()
	eventPool := events.ConstructEventsPool()

	// Симулируем полный день
	err := p.StartNewDay()
	if err != nil {
		t.Fatalf("Failed to start day: %v", err)
	}

	event, err := events.GetRandomEvent(eventPool)
	if err != nil {
		t.Fatalf("Failed to get random event: %v", err)
	}

	// Выполняем первое действие события
	err = event.Execute(1, &p)
	if err != nil {
		t.Logf("Event execution had error (may be expected): %v", err)
	}

	// Проверяем, что игра не крашнулась
	if p.Lock && p.Health <= 0 {
		t.Log("Player died, which is acceptable in test")
	}
}

func TestIntegration_MultipleDays(t *testing.T) {
	p := player.CreatePlayer()
	eventPool := events.ConstructEventsPool()

	// Симулируем несколько дней
	for i := 0; i < 3; i++ {
		if p.Lock {
			break
		}

		err := p.StartNewDay()
		if err != nil && !p.Lock {
			t.Errorf("Unexpected error on day %d: %v", i+1, err)
		}

		if !p.Lock {
			event, err := events.GetRandomEvent(eventPool)
			if err != nil {
				t.Errorf("Failed to get event: %v", err)
				continue
			}
			event.Execute(1, &p)
		}
	}

	// Проверяем, что все поля в допустимых пределах
	if p.Health < 0 || p.Health > player.MaxHealth {
		t.Errorf("Health out of bounds: %d", p.Health)
	}
	if p.Eat < 0 || p.Eat > player.MaxResourceLimit {
		t.Errorf("Eat out of bounds: %d", p.Eat)
	}
	if p.Water < 0 || p.Water > player.MaxResourceLimit {
		t.Errorf("Water out of bounds: %d", p.Water)
	}
}

func TestIntegration_DeathScenario(t *testing.T) {
	p := player.CreatePlayer()

	// Доводим игрока до смерти через недостаток еды
	p.Eat = 0
	p.Health = 20

	err := p.StartNewDay()

	if err == nil {
		t.Error("Expected death error, got nil")
	}
	if !p.Lock {
		t.Error("Expected game to be locked after death")
	}
	if p.Health != 0 {
		t.Errorf("Expected health 0, got %d", p.Health)
	}
}

func TestIntegration_ResourceLimits(t *testing.T) {
	p := player.CreatePlayer()

	// Проверяем лимиты ресурсов
	p.AddEat(100)
	if p.Eat > player.MaxResourceLimit {
		t.Errorf("Eat exceeded limit: %d > %d", p.Eat, player.MaxResourceLimit)
	}

	p.AddWater(100)
	if p.Water > player.MaxResourceLimit {
		t.Errorf("Water exceeded limit: %d > %d", p.Water, player.MaxResourceLimit)
	}

	// Проверяем лимит здоровья - используем корректное значение int8
	// Максимальное значение, которое можно добавить к здоровью
	maxAddHealth := int8(player.MaxHealth - p.Health)
	err := p.ChangeHealth(maxAddHealth)
	if err != nil {
		t.Errorf("Unexpected error when adding health: %v", err)
	}

	// Проверяем, что здоровье не превышает максимум
	if p.Health > player.MaxHealth {
		t.Errorf("Health exceeded limit: %d > %d", p.Health, player.MaxHealth)
	}

	// Проверяем, что здоровье не может стать больше MaxHealth
	p.ChangeHealth(10) // Пытаемся добавить еще
	if p.Health > player.MaxHealth {
		t.Errorf("Health exceeded limit after additional add: %d > %d", p.Health, player.MaxHealth)
	}

	// Проверяем корректное ограничение здоровья
	if p.Health != player.MaxHealth {
		t.Errorf("Expected health to be capped at %d, got %d", player.MaxHealth, p.Health)
	}
}

// Дополнительный тест для проверки граничных значений int8
func TestIntegration_HealthBoundaries(t *testing.T) {
	p := player.CreatePlayer()

	// Устанавливаем здоровье в 0
	p.Health = 0

	// Пытаемся вылечиться
	err := p.ChangeHealth(50)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if p.Health != 50 {
		t.Errorf("Expected health 50, got %d", p.Health)
	}

	// Пытаемся добавить слишком много за один раз (но в пределах int8)
	err = p.ChangeHealth(70)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Должно быть ограничено MaxHealth (100)
	if p.Health != player.MaxHealth {
		t.Errorf("Expected health capped at %d, got %d", player.MaxHealth, p.Health)
	}
}
