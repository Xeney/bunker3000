package player

import (
	"errors"
	"fmt"
)

const (
	MaxHealth        int8 = 100
	StartEat         int8 = 10
	StartWater       int8 = 10
	MaxDays          int8 = 7
	MaxResourceLimit int8 = 20
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

	var damage int8 = 0

	if p.Eat >= 1 {
		p.Eat -= 1
	} else {
		p.Eat = 0
		fmt.Println("Предупреждение: Нет еды! Здоровье уменьшается.")
		damage += 20
	}

	if p.Water >= 2 {
		p.Water -= 2
	} else {
		p.Water = 0
		fmt.Println("Предупреждение: Нет воды! Здоровье уменьшается.")
		damage += 20
	}

	if damage > 0 {
		if p.Health <= damage {
			p.Health = 0
			p.Lock = true
			return errors.New("персонаж погиб от истощения")
		}
		p.Health -= damage
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

func (p *Player) AddEat(amount int8) {
	p.Eat += amount
	if p.Eat > MaxResourceLimit {
		p.Eat = MaxResourceLimit
	}
	fmt.Printf("Получено еды: +%d (Всего: %d)\n", amount, p.Eat)
}

func (p *Player) AddWater(amount int8) {
	p.Water += amount
	if p.Water > MaxResourceLimit {
		p.Water = MaxResourceLimit
	}
	fmt.Printf("Получено воды: +%d (Всего: %d)\n", amount, p.Water)
}

func (p *Player) PrintStatus() {
	fmt.Println("\n📊 СТАТУС ИГРОКА:")
	fmt.Printf("❤️ Здоровье: %d%%\n", p.Health)
	fmt.Printf("🍗 Еда:      %d единиц\n", p.Eat)
	fmt.Printf("💧 Вода:     %d единиц\n", p.Water)
	fmt.Println("----------------------------------------")
}

func (p *Player) ChangeHealth(amount int8) error {
	if p.Lock {
		return errors.New("игра уже завершена")
	}

	if amount < 0 {
		damage := -amount
		if p.Health <= damage {
			p.Health = 0
			p.Lock = true
			return errors.New("персонаж погиб от полученных ран")
		}
		p.Health -= damage
		fmt.Printf("Вы потеряли %d%% здоровья. Осталось: %d%%\n", damage, p.Health)
	} else {
		p.Health += amount
		if p.Health > MaxHealth {
			p.Health = MaxHealth
		}
		fmt.Printf("Вы восстановили %d%% здоровья. Текущее: %d%%\n", amount, p.Health)
	}

	return nil
}
