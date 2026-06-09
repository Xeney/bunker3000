package main

import (
	"bunker3000/events"
	"bunker3000/player"
	"bunker3000/utils"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	utils.Clear()

	p := player.CreatePlayer()
	eventPool := events.ConstructEventsPool()

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║     ДОБРО ПОЖАЛОВАТЬ В БУНКЕР-3000     ║")
	fmt.Println("╚════════════════════════════════════════╝")
	p.PrintStatus()
	utils.WaitEnter()

	for !p.Lock {
		utils.Clear()

		err := p.StartNewDay()
		if err != nil {
			fmt.Printf("\n💀 Игра окончена: %s\n", err.Error())
			utils.WaitEnter()
			break
		}

		if p.Lock {
			break
		}

		ev, err := events.GetRandomEvent(eventPool)
		if err != nil {
			fmt.Printf("❌ Ошибка: %s\n", err.Error())
			utils.WaitEnter()
			break
		}

		ev.Print()

		var choice int
		for {
			fmt.Print("\n👉 Введите ваш выбор (1 или 2): ")
			_, scanErr := fmt.Scan(&choice)

			if scanErr != nil {
				fmt.Println("❌ Ошибка: введите 1 или 2")
				var discard string
				fmt.Scanln(&discard)
				continue
			}

			resultMessage, executeErr := ev.Execute(choice, &p)
			if executeErr != nil {
				fmt.Printf("❌ %s\n", executeErr.Error())
				continue
			}

			utils.Clear()
			fmt.Println("╔════════════════════════════════════════╗")
			fmt.Println("║           РЕЗУЛЬТАТ ДЕЙСТВИЯ          ║")
			fmt.Println("╚════════════════════════════════════════╝")
			fmt.Printf("\n📌 %s\n", resultMessage)

			break
		}

		if p.Lock {
			if p.Health == 0 {
				fmt.Println("\n💀 Игра окончена: Вы погибли")
			}
			utils.WaitEnter()
			break
		}

		p.PrintStatus()
		utils.WaitEnter()
	}

	utils.Clear()
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║        СПАСИБО ЗА ИГРУ!                ║")
	fmt.Println("╚════════════════════════════════════════╝")

	fmt.Printf("\n📊 ИТОГОВАЯ СТАТИСТИКА:\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if p.Health > 0 && p.Lock {
		fmt.Printf("🏆 ПОБЕДА! %d из %d дней\n", p.ThisDay-1, player.MaxDays)
	} else {
		fmt.Printf("💀 ПОРАЖЕНИЕ на %d дне\n", p.ThisDay-1)
	}

	fmt.Printf("❤️ Здоровье: %d%%\n", p.Health)
	fmt.Printf("🍗 Еда: %d / %d\n", p.Eat, player.MaxResourceLimit)
	fmt.Printf("💧 Вода: %d / %d\n", p.Water, player.MaxResourceLimit)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Println("\n📌 Нажмите Enter для выхода...")
	fmt.Scanln()
}
