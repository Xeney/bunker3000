package player

import (
	"errors"
	"fmt"
)

const (
	MaxHealth  int8 = 100
	StartEat   int8 = 10
	StartWater int8 = 10
	MaxDays    int8 = 7
)

type Player struct {
	Health  int8
	Eat     int8
	Water   int8
	ThisDay int8
	Lock    bool
}

func CreatePlayer() Player {
	return Player{
		Health:  MaxHealth,
		Eat:     StartEat,
		Water:   StartWater,
		ThisDay: 1,
		Lock:    false,
	}
}

// Каждое утро вызываем этот метод для симуляции начала дня
func (p *Player) StartNewDay() error {
	if p.Lock {
		return errors.New("игра уже завершена")
	}

	fmt.Printf("--- День %d из %d ---\n", p.ThisDay, MaxDays)

	p.Eat -= 1
	p.Water -= 2

	var damage int8
	if p.Eat <= 0 {
		fmt.Println("Предупреждение: Нет еды! Здоровье уменьшается.")
		damage += 20
	}
	if p.Water <= 0 {
		fmt.Println("Предупреждение: Нет воды! Здоровье уменьшается.")
		damage += 20
	}

	if damage > 0 {
		p.Health -= damage
		if p.Health <= 0 {
			p.Health = 0
			p.Lock = true
			return errors.New("персонаж погиб от истощения")
		}
		fmt.Printf("Текущее здоровье: %d%%\n", p.Health)
	}

	if p.ThisDay >= MaxDays {
		p.Lock = true
		fmt.Printf("Победа! Вы успешно продержались %d дней!\n", MaxDays)
		return nil
	}

	p.ThisDay++
	return nil
}
