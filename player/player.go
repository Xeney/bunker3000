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

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("          ДЕНЬ %d ИЗ %d\n", p.ThisDay, MaxDays)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	var damage int8 = 0

	if p.Eat >= 1 {
		p.Eat -= 1
		fmt.Printf("🍗 Потрачено еды: 1 (осталось: %d)\n", p.Eat)
	} else {
		p.Eat = 0
		fmt.Println("⚠️ ПРЕДУПРЕЖДЕНИЕ: Нет еды! Здоровье уменьшается.")
		damage += 20
	}

	if p.Water >= 2 {
		p.Water -= 2
		fmt.Printf("💧 Потрачено воды: 2 (осталось: %d)\n", p.Water)
	} else {
		p.Water = 0
		fmt.Println("⚠️ ПРЕДУПРЕЖДЕНИЕ: Нет воды! Здоровье уменьшается.")
		damage += 20
	}

	if damage > 0 {
		if p.Health <= damage {
			p.Health = 0
			p.Lock = true
			return errors.New("персонаж погиб от истощения")
		}
		p.Health -= damage
		fmt.Printf("❤️ Потеря здоровья от истощения: -%d%% (осталось: %d%%)\n", damage, p.Health)
	}

	if p.ThisDay >= MaxDays {
		p.Lock = true
		fmt.Println("\n🎉 ПОБЕДА! Вы успешно продержались 7 дней! 🎉")
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
	fmt.Printf("🍗 Получено еды: +%d (Всего: %d)\n", amount, p.Eat)
}

func (p *Player) AddWater(amount int8) {
	p.Water += amount
	if p.Water > MaxResourceLimit {
		p.Water = MaxResourceLimit
	}
	fmt.Printf("💧 Получено воды: +%d (Всего: %d)\n", amount, p.Water)
}

func (p *Player) PrintStatus() {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║        ТЕКУЩИЙ СТАТУС ИГРОКА       ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Printf("❤️  Здоровье: %d%%\n", p.Health)
	fmt.Printf("🍗  Еда:      %d / %d\n", p.Eat, MaxResourceLimit)
	fmt.Printf("💧  Вода:     %d / %d\n", p.Water, MaxResourceLimit)
	fmt.Printf("📅  День:     %d / %d\n", p.ThisDay-1, MaxDays)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
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
		fmt.Printf("❤️ Потеря здоровья: -%d%% (осталось: %d%%)\n", damage, p.Health)
	} else {
		// Исправлено: проверяем переполнение при сложении
		newHealth := int16(p.Health) + int16(amount)
		if newHealth > int16(MaxHealth) {
			p.Health = MaxHealth
		} else {
			p.Health = int8(newHealth)
		}
		fmt.Printf("❤️ Восстановление здоровья: +%d%% (текущее: %d%%)\n", amount, p.Health)
	}

	return nil
}

func (p *Player) RandomResource(amount int8, resourceType string) {
	switch resourceType {
	case "food":
		p.AddEat(amount)
	case "water":
		p.AddWater(amount)
	case "both":
		p.AddEat(amount)
		p.AddWater(amount)
	}
}

// Добавьте в конец файла player.go
func (p *Player) PrintDeathMessage() {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║            💀 GAME OVER 💀         ║")
	fmt.Println("╚════════════════════════════════════╝")

	if p.Health <= 0 {
		fmt.Println("Причина смерти: Полученные травмы")
	} else if p.Eat <= 0 && p.Water <= 0 {
		fmt.Println("Причина смерти: Голод и обезвоживание")
	} else if p.Eat <= 0 {
		fmt.Println("Причина смерти: Голод")
	} else if p.Water <= 0 {
		fmt.Println("Причина смерти: Обезвоживание")
	}

	fmt.Printf("Дней продержался: %d\n", p.ThisDay-1)
}
