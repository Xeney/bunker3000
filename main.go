package main

import (
	"bunker3000/achievements"
	"bunker3000/colors"
	"bunker3000/events"
	"bunker3000/player"
	"bunker3000/save"
	"bunker3000/utils"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var reset = "\033[0m"

func main() {
	rand.Seed(time.Now().UnixNano())
	colors.Init()
	utils.Clear()

	for {
		utils.Clear()
		printTitle()
		choice := showMainMenu()

		switch choice {
		case 1:
			startNewGame()
		case 2:
			loadAndContinue()
		case 3:
			showHallOfFame()
		case 4:
			utils.Clear()
			printBox("СПАСИБО ЗА ИГРУ!", "\033[36m")
			fmt.Println()
			utils.WaitEnter()
			return
		}
	}
}

func printTitle() {
	c := "\033[96m"
	fmt.Printf("%s", c)
	fmt.Println("  ==========================================")
	fmt.Println("  |           БУНКЕР-3000                  |")
	fmt.Println("  |      Симулятор выживания в пустоши     |")
	fmt.Println("  ==========================================")
	fmt.Printf("%s", reset)
	fmt.Println()
}

func printBox(title string, color string) {
	n := len(title) + 4
	fmt.Printf("%s", color)
	fmt.Print("+" + strings.Repeat("-", n-2) + "+\n")
	fmt.Printf("| %s |\n", title)
	fmt.Print("+" + strings.Repeat("-", n-2) + "+\n")
	fmt.Printf("%s", reset)
}

func showMainMenu() int {
	fmt.Println("  1. Новая игра")
	fmt.Println("  2. Загрузить игру")
	fmt.Println("  3. Достижения и рекорды")
	fmt.Println("  4. Выход")
	return utils.GetIntInput(1, 4)
}

func startNewGame() {
	utils.Clear()
	printTitle()

	fmt.Println("\n--- ВЫБОР СЛОЖНОСТИ ---")
	fmt.Println("  1. Лёгкая (10 дней, меньше расход)")
	fmt.Println("  2. Обычная (7 дней, стандартный расход)")
	fmt.Println("  3. Сложная (5 дней, большой расход)")
	diff := utils.GetIntInput(1, 3) - 1

	utils.Clear()
	printTitle()

	fmt.Println("\n--- ВЫБОР КЛАССА ---")
	classes := []player.ClassType{
		player.ClassSurvivor,
		player.ClassMedic,
		player.ClassScout,
		player.ClassBrawler,
	}
	for i, c := range classes {
		fmt.Printf("  %d. %s - %s\n", i+1, c.String(), c.Description())
	}
	classIdx := utils.GetIntInput(1, 4) - 1

	p := player.CreatePlayer(player.Difficulty(diff), classes[classIdx])
	eventPool := events.ConstructEventsPool()
	ach := achievements.NewTracker()

	utils.Clear()
	fmt.Printf("%s", "\033[92m")
	printBox("ДОБРО ПОЖАЛОВАТЬ В БУНКЕР-3000!", "\033[92m")
	fmt.Printf("%s", reset)
	p.PrintStatus()
	utils.WaitEnter()

	for !p.Lock {
		utils.Clear()

		err := p.StartNewDay()
		if err != nil {
			fmt.Printf("\n%s[СМЕРТЬ]%s %s\n", "\033[91m", reset, err.Error())
			utils.WaitEnter()
			break
		}

		if p.Lock {
			break
		}

		ev, err := events.GetRandomEvent(eventPool)
		if err != nil {
			fmt.Printf("%s[ОШИБКА]%s %s\n", "\033[91m", reset, err.Error())
			utils.WaitEnter()
			break
		}

		ev.Print()

		var choice int
		var resultMessage string

		for {
			choice = utils.GetIntInput(1, 2)

			resultMessage, err = ev.Execute(choice, &p)
			if err != nil {
				fmt.Printf("\n%s[ОШИБКА]%s %s\n", "\033[91m", reset, err.Error())
				continue
			}
			break
		}

		utils.Clear()

		fmt.Printf("%s", "\033[93m")
		printBox("РЕЗУЛЬТАТ ДЕЙСТВИЯ", "\033[93m")
		fmt.Printf("%s", reset)
		fmt.Printf("\n%s\n", resultMessage)

		if resultMessage != "" {
			if len(resultMessage) > 5 {
				tag := resultMessage[:6]
				switch {
				case tag == "[СМЕРТ" || tag == "[ТРАВМ":
					fmt.Printf("  [%s-%s%s%s]\n", "\033[91m", reset, "\033[91m", reset)
				case tag == "[ЕДА]" || tag == "[ВОДА]":
					fmt.Printf("  [%s+РЕСУРСЫ%s]\n", "\033[92m", reset)
				}
			}
		}

		if p.Lock {
			ach.RecordDamage(calcDamage(&p))
			ach.Check(&p)
			utils.WaitEnter()
			break
		}

		p.PrintStatus()

		_ = save.Save(
			p,
			[]struct {
				ID       string
				Unlocked bool
			}(nil),
			0,
		)

		ach.Check(&p)
		utils.WaitEnter()

		if !p.Lock {
			p.ThisDay++
		}
	}

	showGameOver(&p, ach)
}

func calcDamage(p *player.Player) int8 {
	return 100 - p.Health
}

func loadAndContinue() {
	sd, err := save.Load()
	if err != nil {
		utils.Clear()
		fmt.Printf("\n%s[ОШИБКА]%s Сохранение не найдено.\n", "\033[91m", reset)
		utils.WaitEnter()
		return
	}

	p := sd.Player
	if p.Lock {
		fmt.Printf("\n%s[ИНФО]%s Эта игра уже завершена.\n", "\033[93m", reset)
		utils.WaitEnter()
		return
	}

	eventPool := events.ConstructEventsPool()
	ach := achievements.NewTracker()
	for _, sa := range sd.Achievements {
		for _, a := range ach.Achievements {
			if a.ID == sa.ID && sa.Unlocked {
				a.Unlocked = true
			}
		}
	}

	utils.Clear()
	fmt.Printf("%s", "\033[92m")
	printBox("ИГРА ЗАГРУЖЕНА", "\033[92m")
	fmt.Printf("%s", reset)
	p.PrintStatus()
	utils.WaitEnter()

	for !p.Lock {
		utils.Clear()
		err := p.StartNewDay()
		if err != nil {
			fmt.Printf("\n%s[СМЕРТЬ]%s %s\n", "\033[91m", reset, err.Error())
			utils.WaitEnter()
			break
		}
		if p.Lock {
			break
		}

		ev, err := events.GetRandomEvent(eventPool)
		if err != nil {
			fmt.Printf("%s[ОШИБКА]%s %s\n", "\033[91m", reset, err.Error())
			utils.WaitEnter()
			break
		}

		ev.Print()

		var resultMessage string
		for {
			choice := utils.GetIntInput(1, 2)
			resultMessage, err = ev.Execute(choice, &p)
			if err != nil {
				fmt.Printf("\n%s[ОШИБКА]%s %s\n", "\033[91m", reset, err.Error())
				continue
			}
			break
		}

		utils.Clear()

		fmt.Printf("%s", "\033[93m")
		printBox("РЕЗУЛЬТАТ ДЕЙСТВИЯ", "\033[93m")
		fmt.Printf("%s", reset)
		fmt.Printf("\n%s\n", resultMessage)

		if p.Lock {
			ach.Check(&p)
			utils.WaitEnter()
			break
		}

		p.PrintStatus()
		ach.Check(&p)
		utils.WaitEnter()

		if !p.Lock {
			p.ThisDay++
		}
	}

	showGameOver(&p, ach)
}

func showGameOver(p *player.Player, ach *achievements.Tracker) {
	utils.Clear()
	cfg := p.GetCfg()

	fmt.Println()
	if p.Health > 0 && p.Lock {
		fmt.Printf("%s", "\033[92m")
		printBox("ПОБЕДА!", "\033[92m")
	} else {
		fmt.Printf("%s", "\033[91m")
		printBox("ПОРАЖЕНИЕ", "\033[91m")
	}
	fmt.Printf("%s", reset)

	fmt.Println("\n--- ИТОГОВАЯ СТАТИСТИКА ---")
	if p.Health > 0 && p.Lock {
		fmt.Printf("  Статус:     ПОБЕДА! %d из %d дней\n", p.ThisDay, cfg.MaxDays)
	} else if p.Health == 0 {
		fmt.Printf("  Статус:     ПОРАЖЕНИЕ на %d дне\n", p.ThisDay)
	} else {
		fmt.Printf("  Статус:     Игра завершена на %d дне\n", p.ThisDay)
	}
	fmt.Printf("  Здоровье:   %d%%\n", p.Health)
	fmt.Printf("  Еда:        %d / %d\n", p.Eat, cfg.MaxResource)
	fmt.Printf("  Вода:       %d / %d\n", p.Water, cfg.MaxResource)
	fmt.Printf("  Сложность:  %s\n", p.Difficulty.String())
	fmt.Printf("  Класс:      %s\n", p.Class.String())

	ach.Check(p)
	ach.PrintAll()

	// Save player stats for records saving
	_ = save.DeleteSave()

	utils.WaitEnter()
}

func showHallOfFame() {
	utils.Clear()
	printTitle()
	fmt.Println("\n--- РЕКОРДЫ И ДОСТИЖЕНИЯ ---")
	fmt.Println("\nСохранение рекордов пока недоступно.")
	fmt.Println("\nПодсказка: для просмотра достижений завершите игру.")
	fmt.Println()

	ach := achievements.NewTracker()
	// Try to load saved achievements from save data
	sd, err := save.Load()
	if err == nil {
		for _, sa := range sd.Achievements {
			for _, a := range ach.Achievements {
				if a.ID == sa.ID && sa.Unlocked {
					a.Unlocked = true
				}
			}
		}
		ach.PrintAll()
	} else {
		fmt.Println("  Достижения не найдены. Сыграйте игру, чтобы получить их.")
	}

	utils.WaitEnter()
}


