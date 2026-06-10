package player

import (
	"errors"
	"fmt"
)

type Difficulty int8

const (
	DifficultyEasy     Difficulty = 0
	DifficultyNormal   Difficulty = 1
	DifficultyHard     Difficulty = 2
	DifficultyEndless  Difficulty = 3
)

func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "Лёгкая"
	case DifficultyNormal:
		return "Обычная"
	case DifficultyHard:
		return "Сложная"
	case DifficultyEndless:
		return "Бесконечный"
	}
	return "Неизвестно"
}

type ClassType int8

const (
	ClassSurvivor ClassType = 0
	ClassMedic    ClassType = 1
	ClassScout    ClassType = 2
	ClassBrawler  ClassType = 3
	ClassHauler   ClassType = 4
	ClassMechanic ClassType = 5
	ClassTrader   ClassType = 6
)

func (c ClassType) String() string {
	switch c {
	case ClassSurvivor:
		return "Выживальщик"
	case ClassMedic:
		return "Медик"
	case ClassScout:
		return "Разведчик"
	case ClassBrawler:
		return "Боец"
	case ClassHauler:
		return "Грузчик"
	case ClassMechanic:
		return "Механик"
	case ClassTrader:
		return "Торговец"
	}
	return "Неизвестно"
}

func (c ClassType) Description() string {
	switch c {
	case ClassSurvivor:
		return "Сбалансированный персонаж без особых бонусов"
	case ClassMedic:
		return "+20 к макс. здоровью, двойное лечение"
	case ClassScout:
		return "Двойная добыча ресурсов"
	case ClassBrawler:
		return "Половина урона от всех источников"
	case ClassHauler:
		return "+6 к вместимости, +2 стартовых ресурса"
	case ClassMechanic:
		return "+15 HP, −25% урона отовсюду"
	case ClassTrader:
		return "+2 еды/воды из событий помощи"
	}
	return ""
}

type DifficultyConfig struct {
	Name            string
	MaxDays         int8
	ConsumeFood     int8
	ConsumeWater    int8
	StartEat        int8
	StartWater      int8
	MaxResource     int8
	StarvationDmg   int8
	DehydrationDmg  int8
}

var DifficultyConfigs = map[Difficulty]DifficultyConfig{
	DifficultyEasy: {
		Name:           "Лёгкая",
		MaxDays:        10,
		ConsumeFood:    1,
		ConsumeWater:   1,
		StartEat:       15,
		StartWater:     15,
		MaxResource:    25,
		StarvationDmg:  10,
		DehydrationDmg: 10,
	},
	DifficultyNormal: {
		Name:           "Обычная",
		MaxDays:        7,
		ConsumeFood:    1,
		ConsumeWater:   2,
		StartEat:       10,
		StartWater:     10,
		MaxResource:    20,
		StarvationDmg:  20,
		DehydrationDmg: 20,
	},
	DifficultyHard: {
		Name:           "Сложная",
		MaxDays:        5,
		ConsumeFood:    2,
		ConsumeWater:   3,
		StartEat:       5,
		StartWater:     5,
		MaxResource:    15,
		StarvationDmg:  30,
		DehydrationDmg: 30,
	},
	DifficultyEndless: {
		Name:           "Бесконечный",
		MaxDays:        -1,
		ConsumeFood:    2,
		ConsumeWater:   2,
		StartEat:       12,
		StartWater:     12,
		MaxResource:    20,
		StarvationDmg:  15,
		DehydrationDmg: 15,
	},
}

type Player struct {
	Health     int8
	Eat        int8
	Water      int8
	ThisDay    int8
	Lock       bool
	Difficulty Difficulty
	Class      ClassType
	Karma      int8
	XP         int
	Level      int8
	Flags      map[string]bool
}

var xpThresholds = []int{0, 20, 50, 90, 140, 200, 270, 350, 440, 540, 650}

func XPForLevel(level int8) int {
	idx := int(level)
	if idx >= len(xpThresholds) {
		return 650
	}
	return xpThresholds[idx]
}

func CreatePlayer(diff Difficulty, class ClassType) Player {
	cfg := DifficultyConfigs[diff]
	health := int8(100)
	switch class {
	case ClassMedic:
		health = 120
	case ClassMechanic:
		health = 115
	}

	eat := cfg.StartEat
	water := cfg.StartWater
	if class == ClassHauler {
		eat += 2
		water += 2
	}

	return Player{
		Health:     health,
		Eat:        eat,
		Water:      water,
		ThisDay:    1,
		Lock:       false,
		Difficulty: diff,
		Class:      class,
		Karma:      0,
		XP:         0,
		Level:      0,
		Flags:      make(map[string]bool),
	}
}

func (p *Player) AddKarma(delta int8) {
	p.Karma += delta
	if p.Karma > 100 {
		p.Karma = 100
	}
	if p.Karma < -100 {
		p.Karma = -100
	}
}

func (p *Player) AddXP(amount int) bool {
	p.XP += amount
	if p.XP < 0 {
		p.XP = 0
	}
	nextLevel := int(p.Level) + 1
	if nextLevel < len(xpThresholds) && p.XP >= xpThresholds[nextLevel] {
		p.Level++
		return true
	}
	return false
}

func (p *Player) GetMaxHP() int8 {
	base := MaxHealth(p.Class)
	return base + p.Level*5
}

func (p *Player) GetMaxResource() int8 {
	base := p.GetCfg().MaxResource
	if p.Class == ClassHauler {
		base += 6
	}
	return base + p.Level
}

func (p *Player) GetCfg() DifficultyConfig {
	return DifficultyConfigs[p.Difficulty]
}

func (p *Player) StartNewDay() error {
	if p.Lock {
		return errors.New("игра уже завершена")
	}

	cfg := p.GetCfg()

	if cfg.MaxDays < 0 {
		fmt.Printf("========================================\n")
		fmt.Printf("          ДЕНЬ %d\n", p.ThisDay)
		fmt.Printf("========================================\n")
	} else {
		fmt.Printf("========================================\n")
		fmt.Printf("          ДЕНЬ %d ИЗ %d\n", p.ThisDay, cfg.MaxDays)
		fmt.Printf("========================================\n")
	}

	var (
		consumeFood = cfg.ConsumeFood
		consumeWater = cfg.ConsumeWater
		starvationDmg = cfg.StarvationDmg
		dehydrationDmg = cfg.DehydrationDmg
	)

	var damage int8 = 0

	if p.Eat >= consumeFood {
		p.Eat -= consumeFood
		fmt.Printf("[-] Еда: -%d (осталось: %d)\n", consumeFood, p.Eat)
	} else {
		p.Eat = 0
		fmt.Println("[!] НЕТ ЕДЫ! Здоровье уменьшается.")
		damage += starvationDmg
	}

	if p.Water >= consumeWater {
		p.Water -= consumeWater
		fmt.Printf("[-] Вода: -%d (осталось: %d)\n", consumeWater, p.Water)
	} else {
		p.Water = 0
		fmt.Println("[!] НЕТ ВОДЫ! Здоровье уменьшается.")
		damage += dehydrationDmg
	}

	if damage > 0 {
		if p.Health <= damage {
			p.Health = 0
			p.Lock = true
			return errors.New("персонаж погиб от истощения")
		}
		p.Health -= damage
		fmt.Printf("[-] Потеря здоровья от истощения: -%d%% (осталось: %d%%)\n", damage, p.Health)
	}

	if cfg.MaxDays >= 0 && p.ThisDay >= cfg.MaxDays {
		p.Lock = true
		fmt.Println("\n=== ПОБЕДА! Вы успешно продержались все дни! ===")
		return nil
	}

	return nil
}

func (p *Player) AddEat(amount int8) {
	if p.Class == ClassScout && amount > 0 {
		amount *= 2
	}
	p.Eat += amount
	maxRes := p.GetMaxResource()
	if p.Eat > maxRes {
		p.Eat = maxRes
	}
	if p.Eat < 0 {
		p.Eat = 0
	}
}

func (p *Player) AddWater(amount int8) {
	if p.Class == ClassScout && amount > 0 {
		amount *= 2
	}
	p.Water += amount
	maxRes := p.GetMaxResource()
	if p.Water > maxRes {
		p.Water = maxRes
	}
	if p.Water < 0 {
		p.Water = 0
	}
}

func (p *Player) PrintStatus() {
	cfg := p.GetCfg()
	fmt.Println("\n+------------------------------------------+")
	fmt.Println("|         ТЕКУЩИЙ СТАТУС ИГРОКА            |")
	fmt.Println("+------------------------------------------+")
	fmt.Printf("| Здоровье: %d%%\n", p.Health)
	fmt.Printf("| Еда:      %d / %d\n", p.Eat, p.GetMaxResource())
	fmt.Printf("| Вода:     %d / %d\n", p.Water, p.GetMaxResource())
	if cfg.MaxDays < 0 {
		fmt.Printf("| День:     %d (∞)\n", p.ThisDay)
	} else {
		fmt.Printf("| День:     %d / %d\n", p.ThisDay, cfg.MaxDays)
	}
	fmt.Printf("| Сложность: %s\n", p.Difficulty.String())
	fmt.Printf("| Класс:     %s\n", p.Class.String())
	fmt.Printf("| Уровень:   %d (XP: %d)\n", p.Level, p.XP)
	fmt.Println("+------------------------------------------+")
}

func (p *Player) ChangeHealth(amount int8) error {
	if p.Lock {
		return errors.New("игра уже завершена")
	}

	if amount < 0 {
		damage := -amount
		switch p.Class {
		case ClassBrawler:
			damage = damage / 2
		case ClassMechanic:
			damage = damage * 3 / 4
		}
		if damage < 1 {
			damage = 1
		}
		if p.Health <= damage {
			p.Health = 0
			p.Lock = true
			return errors.New("персонаж погиб от полученных ран")
		}
		p.Health -= damage
	} else {
		if p.Class == ClassMedic {
			amount *= 2
		}
		maxHP := p.GetMaxHP()
		if int(p.Health)+int(amount) > int(maxHP) {
			p.Health = maxHP
		} else {
			p.Health += amount
		}
	}

	return nil
}

func MaxHealth(class ClassType) int8 {
	switch class {
	case ClassMedic:
		return 120
	case ClassMechanic:
		return 115
	default:
		return 100
	}
}

func (p *Player) PrintDeathMessage() {
	fmt.Println("\n+------------------------------------------+")
	fmt.Println("|              GAME OVER                   |")
	fmt.Println("+------------------------------------------+")

	if p.Health <= 0 {
		fmt.Println("| Причина смерти: Полученные травмы")
	} else if p.Eat <= 0 && p.Water <= 0 {
		fmt.Println("| Причина смерти: Голод и обезвоживание")
	} else if p.Eat <= 0 {
		fmt.Println("| Причина смерти: Голод")
	} else if p.Water <= 0 {
		fmt.Println("| Причина смерти: Обезвоживание")
	}

	cfg := p.GetCfg()
	if cfg.MaxDays < 0 {
		fmt.Printf("| Дней продержался: %d (∞)\n", p.ThisDay)
	} else {
		fmt.Printf("| Дней продержался: %d / %d\n", p.ThisDay, cfg.MaxDays)
	}
}
