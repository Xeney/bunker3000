package achievements

import (
	"bunker3000/player"
	"fmt"
)

type Achievement struct {
	ID       string
	Name     string
	Desc     string
	Unlocked bool
}

type Tracker struct {
	Achievements []*Achievement
	TotalDamage  int8
}

func NewTracker() *Tracker {
	return &Tracker{
		Achievements: []*Achievement{
			{ID: "first_day", Name: "ПЕРВЫЙ ШАГ", Desc: "Пережить первый день"},
			{ID: "veteran", Name: "ВЕТЕРАН", Desc: "Пережить 3 дня"},
			{ID: "survivor", Name: "ВЫЖИВАЛЬЩИК", Desc: "Пережить 7 дней"},
			{ID: "marathon", Name: "МАРАФОНЕЦ", Desc: "Пережить 10 дней"},
			{ID: "tough", Name: "КРЕПКИЙ ОРЕШЕК", Desc: "Выдержать суммарно 50+ урона и выжить"},
			{ID: "collector", Name: "КОЛЛЕКЦИОНЕР", Desc: "Заполнить ресурс до максимума"},
			{ID: "speedster", Name: "СКОРОСТНОЙ ФИНИШ", Desc: "Победить на сложной сложности"},
			{ID: "lucky", Name: "СЧАСТЛИВЧИК", Desc: "Победить, ни разу не потеряв всё здоровье"},
			{ID: "medic", Name: "ВРАЧЕВАТЕЛЬ", Desc: "Победить в роли Медика"},
			{ID: "scout", Name: "СЛЕДОПЫТ", Desc: "Победить в роли Разведчика"},
			{ID: "brawler", Name: "ГЛАДИАТОР", Desc: "Победить в роли Бойца"},
			{ID: "hauler", Name: "ТЯЖЕЛОВОЗ", Desc: "Победить в роли Грузчика"},
			{ID: "mechanic", Name: "РЕМОНТНИК", Desc: "Победить в роли Механика"},
			{ID: "trader", Name: "БАРЫГА", Desc: "Победить в роли Торговца"},
			{ID: "maxlevel", Name: "МАКСИМАЛИСТ", Desc: "Достичь 10-го уровня"},
		},
	}
}

func (t *Tracker) Check(p *player.Player) {
	for _, a := range t.Achievements {
		if a.Unlocked {
			continue
		}
		switch a.ID {
		case "first_day":
			if p.ThisDay >= 2 || p.Lock {
				a.Unlocked = true
			}
		case "veteran":
			if p.ThisDay >= 3 {
				a.Unlocked = true
			}
		case "survivor":
			if p.ThisDay >= 7 {
				a.Unlocked = true
			}
		case "marathon":
			if p.ThisDay >= 10 {
				a.Unlocked = true
			}
		case "tough":
			if t.TotalDamage >= 50 && p.Health > 0 {
				a.Unlocked = true
			}
		case "collector":
			cfg := p.GetCfg()
			if p.Eat >= cfg.MaxResource || p.Water >= cfg.MaxResource {
				a.Unlocked = true
			}
		case "speedster":
			if p.Health > 0 && p.Lock && p.Difficulty == player.DifficultyHard {
				a.Unlocked = true
			}
		case "lucky":
			if p.Health > 0 && p.Lock && t.TotalDamage < 100 {
				a.Unlocked = true
			}
		case "medic":
			if p.Health > 0 && p.Lock && p.Class == player.ClassMedic {
				a.Unlocked = true
			}
		case "scout":
			if p.Health > 0 && p.Lock && p.Class == player.ClassScout {
				a.Unlocked = true
			}
		case "brawler":
			if p.Health > 0 && p.Lock && p.Class == player.ClassBrawler {
				a.Unlocked = true
			}
		case "hauler":
			if p.Health > 0 && p.Lock && p.Class == player.ClassHauler {
				a.Unlocked = true
			}
		case "mechanic":
			if p.Health > 0 && p.Lock && p.Class == player.ClassMechanic {
				a.Unlocked = true
			}
		case "trader":
			if p.Health > 0 && p.Lock && p.Class == player.ClassTrader {
				a.Unlocked = true
			}
		case "maxlevel":
			if p.Level >= 10 {
				a.Unlocked = true
			}
		}
	}
}

func (t *Tracker) RecordDamage(dmg int8) {
	t.TotalDamage += dmg
}

func (t *Tracker) PrintAll() {
	unlocked := 0
	for _, a := range t.Achievements {
		if a.Unlocked {
			unlocked++
		}
	}
	if unlocked == 0 {
		fmt.Println("\n  Достижения не получены.")
		return
	}
	fmt.Println("\n+------------------------------------------+")
	fmt.Println("|              ДОСТИЖЕНИЯ                  |")
	fmt.Println("+------------------------------------------+")
	for _, a := range t.Achievements {
		if a.Unlocked {
			fmt.Printf("| [X] %s\n", a.Name)
			fmt.Printf("|     %s\n", a.Desc)
		}
	}
	fmt.Printf("| Получено: %d / %d\n", unlocked, len(t.Achievements))
	fmt.Println("+------------------------------------------+")
}
