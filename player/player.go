package player

import (
	"errors"
	"fmt"
)

const (
	MaxHealth  uint8 = 100
	StartEat   uint8 = 10
	StartWater uint8 = 10
	MaxDays    uint8 = 7
)

type Player struct {
	Health  uint8
	Eat     uint8
	Water   uint8
	ThisDay uint8 // 1-7
	Lock    bool  // true, если игра окончена (победа или поражение)
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

	// 1. Выводим текущий день
	fmt.Printf("--- День %d из %d ---\n", p.ThisDay, MaxDays)

	// 2. Тратим ресурсы на жизнеобеспечение
	p.Eat -= 1
	p.Water -= 2

	// 3. Проверяем нехватку ресурсов и наносим урон
	var damage uint8
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

	// 4. Переходим на следующий день или объявляем победу
	if p.ThisDay >= MaxDays {
		p.Lock = true
		fmt.Printf("Победа! Вы успешно продержались %d дней!\n", MaxDays)
		return nil
	}

	p.ThisDay++
	return nil
}
