package player

import (
	"errors"
	"fmt"
)

type Difficulty int8

const (
	DifficultyEasy   Difficulty = 0
	DifficultyNormal Difficulty = 1
	DifficultyHard   Difficulty = 2
)

func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "Лёгкая"
	case DifficultyNormal:
		return "Обычная"
	case DifficultyHard:
		return "Сложная"
	}
	return "Неизвестно"
}

type ClassType int8

const (
	ClassSurvivor ClassType = 0
	ClassMedic    ClassType = 1
	ClassScout    ClassType = 2
	ClassBrawler  ClassType = 3
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
}

type Player struct {
	Health     int8
	Eat        int8
	Water      int8
	ThisDay    int8
	Lock       bool
	Difficulty Difficulty
	Class      ClassType
}

func CreatePlayer(diff Difficulty, class ClassType) Player {
	cfg := DifficultyConfigs[diff]
	health := int8(100)
	if class == ClassMedic {
		health = 120
	}

	return Player{
		Health:     health,
		Eat:        cfg.StartEat,
		Water:      cfg.StartWater,
		ThisDay:    1,
		Lock:       false,
		Difficulty: diff,
		Class:      class,
	}
}

func (p *Player) GetCfg() DifficultyConfig {
	return DifficultyConfigs[p.Difficulty]
}

func (p *Player) StartNewDay() error {
	if p.Lock {
		return errors.New("игра уже завершена")
	}

	cfg := p.GetCfg()

	fmt.Printf("========================================\n")
	fmt.Printf("          ДЕНЬ %d ИЗ %d\n", p.ThisDay, cfg.MaxDays)
	fmt.Printf("========================================\n")

	var damage int8 = 0

	if p.Eat >= cfg.ConsumeFood {
		p.Eat -= cfg.ConsumeFood
		fmt.Printf("[-] Еда: -%d (осталось: %d)\n", cfg.ConsumeFood, p.Eat)
	} else {
		p.Eat = 0
		fmt.Println("[!] НЕТ ЕДЫ! Здоровье уменьшается.")
		damage += cfg.StarvationDmg
	}

	if p.Water >= cfg.ConsumeWater {
		p.Water -= cfg.ConsumeWater
		fmt.Printf("[-] Вода: -%d (осталось: %d)\n", cfg.ConsumeWater, p.Water)
	} else {
		p.Water = 0
		fmt.Println("[!] НЕТ ВОДЫ! Здоровье уменьшается.")
		damage += cfg.DehydrationDmg
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

	if p.ThisDay >= cfg.MaxDays {
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
	cfg := p.GetCfg()
	if p.Eat > cfg.MaxResource {
		p.Eat = cfg.MaxResource
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
	cfg := p.GetCfg()
	if p.Water > cfg.MaxResource {
		p.Water = cfg.MaxResource
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
	fmt.Printf("| Еда:      %d / %d\n", p.Eat, cfg.MaxResource)
	fmt.Printf("| Вода:     %d / %d\n", p.Water, cfg.MaxResource)
	fmt.Printf("| День:     %d / %d\n", p.ThisDay, cfg.MaxDays)
	fmt.Printf("| Сложность: %s\n", p.Difficulty.String())
	fmt.Printf("| Класс:     %s\n", p.Class.String())
	fmt.Println("+------------------------------------------+")
}

func (p *Player) ChangeHealth(amount int8) error {
	if p.Lock {
		return errors.New("игра уже завершена")
	}

	if amount < 0 {
		damage := -amount
		if p.Class == ClassBrawler {
			damage = damage / 2
			if damage < 1 {
				damage = 1
			}
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
		if int(p.Health)+int(amount) > int(MaxHealth(p.Class)) {
			p.Health = MaxHealth(p.Class)
		} else {
			p.Health += amount
		}
	}

	return nil
}

func MaxHealth(class ClassType) int8 {
	if class == ClassMedic {
		return 120
	}
	return 100
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
	fmt.Printf("| Дней продержался: %d / %d\n", p.ThisDay, cfg.MaxDays)
}
