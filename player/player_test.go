package player

import (
	"testing"
)

func TestCreatePlayer(t *testing.T) {
	p := CreatePlayer()

	if p.Health != MaxHealth {
		t.Errorf("Expected Health %d, got %d", MaxHealth, p.Health)
	}
	if p.Eat != StartEat {
		t.Errorf("Expected Eat %d, got %d", StartEat, p.Eat)
	}
	if p.Water != StartWater {
		t.Errorf("Expected Water %d, got %d", StartWater, p.Water)
	}
	if p.ThisDay != 1 {
		t.Errorf("Expected ThisDay 1, got %d", p.ThisDay)
	}
	if p.Lock {
		t.Error("Expected Lock to be false")
	}
}

func TestStartNewDay_Success(t *testing.T) {
	p := CreatePlayer()
	err := p.StartNewDay()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Eat != StartEat-1 {
		t.Errorf("Expected Eat %d, got %d", StartEat-1, p.Eat)
	}
	if p.Water != StartWater-2 {
		t.Errorf("Expected Water %d, got %d", StartWater-2, p.Water)
	}
	if p.ThisDay != 2 {
		t.Errorf("Expected ThisDay 2, got %d", p.ThisDay)
	}
}

func TestStartNewDay_NoFood(t *testing.T) {
	p := CreatePlayer()
	p.Eat = 0
	p.Health = 80

	err := p.StartNewDay()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Health != 60 {
		t.Errorf("Expected Health 60, got %d", p.Health)
	}
	if p.Lock {
		t.Error("Expected Lock to be false")
	}
}

func TestStartNewDay_NoFood_Death(t *testing.T) {
	p := CreatePlayer()
	p.Eat = 0
	p.Health = 15

	err := p.StartNewDay()

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if p.Health != 0 {
		t.Errorf("Expected Health 0, got %d", p.Health)
	}
	if !p.Lock {
		t.Error("Expected Lock to be true")
	}
}

func TestStartNewDay_NoWater(t *testing.T) {
	p := CreatePlayer()
	p.Water = 1
	p.Health = 80

	err := p.StartNewDay()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Health != 60 {
		t.Errorf("Expected Health 60, got %d", p.Health)
	}
}

func TestStartNewDay_NoWater_Death(t *testing.T) {
	p := CreatePlayer()
	p.Water = 1
	p.Health = 15

	err := p.StartNewDay()

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if p.Health != 0 {
		t.Errorf("Expected Health 0, got %d", p.Health)
	}
	if !p.Lock {
		t.Error("Expected Lock to be true")
	}
}

func TestStartNewDay_MaxDaysReached(t *testing.T) {
	p := CreatePlayer()
	p.ThisDay = MaxDays

	err := p.StartNewDay()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !p.Lock {
		t.Error("Expected Lock to be true after max days")
	}
}

func TestStartNewDay_Locked(t *testing.T) {
	p := CreatePlayer()
	p.Lock = true

	err := p.StartNewDay()

	if err == nil {
		t.Error("Expected error when game is locked, got nil")
	}
}

func TestAddEat(t *testing.T) {
	p := CreatePlayer()
	initialEat := p.Eat

	p.AddEat(5)

	if p.Eat != initialEat+5 {
		t.Errorf("Expected Eat %d, got %d", initialEat+5, p.Eat)
	}
}

func TestAddEat_OverLimit(t *testing.T) {
	p := CreatePlayer()
	p.AddEat(MaxResourceLimit + 10)

	if p.Eat != MaxResourceLimit {
		t.Errorf("Expected Eat capped at %d, got %d", MaxResourceLimit, p.Eat)
	}
}

func TestAddWater(t *testing.T) {
	p := CreatePlayer()
	initialWater := p.Water

	p.AddWater(5)

	if p.Water != initialWater+5 {
		t.Errorf("Expected Water %d, got %d", initialWater+5, p.Water)
	}
}

func TestAddWater_OverLimit(t *testing.T) {
	p := CreatePlayer()
	p.AddWater(MaxResourceLimit + 10)

	if p.Water != MaxResourceLimit {
		t.Errorf("Expected Water capped at %d, got %d", MaxResourceLimit, p.Water)
	}
}

func TestChangeHealth_Negative(t *testing.T) {
	p := CreatePlayer()
	p.Health = 80

	err := p.ChangeHealth(-30)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Health != 50 {
		t.Errorf("Expected Health 50, got %d", p.Health)
	}
}

func TestChangeHealth_Negative_Death(t *testing.T) {
	p := CreatePlayer()
	p.Health = 20

	err := p.ChangeHealth(-30)

	if err == nil {
		t.Error("Expected error on death, got nil")
	}
	if p.Health != 0 {
		t.Errorf("Expected Health 0, got %d", p.Health)
	}
	if !p.Lock {
		t.Error("Expected Lock to be true")
	}
}

func TestChangeHealth_Positive(t *testing.T) {
	p := CreatePlayer()
	p.Health = 50

	err := p.ChangeHealth(30)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Health != 80 {
		t.Errorf("Expected Health 80, got %d", p.Health)
	}
}

func TestChangeHealth_Positive_OverMax(t *testing.T) {
	p := CreatePlayer()
	p.Health = 90

	err := p.ChangeHealth(30)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Health != MaxHealth {
		t.Errorf("Expected Health capped at %d, got %d", MaxHealth, p.Health)
	}
}

func TestChangeHealth_Locked(t *testing.T) {
	p := CreatePlayer()
	p.Lock = true

	err := p.ChangeHealth(10)

	if err == nil {
		t.Error("Expected error when game is locked, got nil")
	}
}
