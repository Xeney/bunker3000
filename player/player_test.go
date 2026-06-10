package player

import (
	"testing"
)

func TestCreatePlayer(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)

	if p.Health != MaxHealth(ClassSurvivor) {
		t.Errorf("Expected Health %d, got %d", MaxHealth(ClassSurvivor), p.Health)
	}
	cfg := DifficultyConfigs[DifficultyNormal]
	if p.Eat != cfg.StartEat {
		t.Errorf("Expected Eat %d, got %d", cfg.StartEat, p.Eat)
	}
	if p.Water != cfg.StartWater {
		t.Errorf("Expected Water %d, got %d", cfg.StartWater, p.Water)
	}
	if p.ThisDay != 1 {
		t.Errorf("Expected ThisDay 1, got %d", p.ThisDay)
	}
	if p.Lock {
		t.Error("Expected Lock to be false")
	}
}

func TestStartNewDay_Success(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	cfg := DifficultyConfigs[DifficultyNormal]
	err := p.StartNewDay()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Eat != cfg.StartEat-cfg.ConsumeFood {
		t.Errorf("Expected Eat %d, got %d", cfg.StartEat-cfg.ConsumeFood, p.Eat)
	}
	if p.Water != cfg.StartWater-cfg.ConsumeWater {
		t.Errorf("Expected Water %d, got %d", cfg.StartWater-cfg.ConsumeWater, p.Water)
	}
	if p.ThisDay != 1 {
		t.Errorf("Expected ThisDay 1 (unchanged), got %d", p.ThisDay)
	}
}

func TestStartNewDay_NoFood(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
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
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
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
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
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
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
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
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	cfg := DifficultyConfigs[DifficultyNormal]
	p.ThisDay = cfg.MaxDays

	err := p.StartNewDay()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !p.Lock {
		t.Error("Expected Lock to be true after max days")
	}
}

func TestStartNewDay_Locked(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	p.Lock = true

	err := p.StartNewDay()

	if err == nil {
		t.Error("Expected error when game is locked, got nil")
	}
}

func TestAddEat(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	initialEat := p.Eat

	p.AddEat(5)

	if p.Eat != initialEat+5 {
		t.Errorf("Expected Eat %d, got %d", initialEat+5, p.Eat)
	}
}

func TestAddEat_OverLimit(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	cfg := DifficultyConfigs[DifficultyNormal]
	p.AddEat(cfg.MaxResource + 10)

	if p.Eat != cfg.MaxResource {
		t.Errorf("Expected Eat capped at %d, got %d", cfg.MaxResource, p.Eat)
	}
}

func TestAddWater(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	initialWater := p.Water

	p.AddWater(5)

	if p.Water != initialWater+5 {
		t.Errorf("Expected Water %d, got %d", initialWater+5, p.Water)
	}
}

func TestAddWater_OverLimit(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	cfg := DifficultyConfigs[DifficultyNormal]
	p.AddWater(cfg.MaxResource + 10)

	if p.Water != cfg.MaxResource {
		t.Errorf("Expected Water capped at %d, got %d", cfg.MaxResource, p.Water)
	}
}

func TestChangeHealth_Negative(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
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
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
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
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
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
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	p.Health = 90

	err := p.ChangeHealth(30)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Health != MaxHealth(ClassSurvivor) {
		t.Errorf("Expected Health capped at %d, got %d", MaxHealth(ClassSurvivor), p.Health)
	}
}

func TestChangeHealth_Locked(t *testing.T) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)
	p.Lock = true

	err := p.ChangeHealth(10)

	if err == nil {
		t.Error("Expected error when game is locked, got nil")
	}
}
