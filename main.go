package main

import (
	"bunker3000/events"
	"bunker3000/player"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	p := player.CreatePlayer()
	eventPool := events.ConstructEventsPool()

	fmt.Println("=== ДОБРО ПОЖАЛОВАТЬ В СИМУЛЯТОР ВЫЖИВАНИЯ ===")
	p.PrintStatus()

	for !p.Lock {
		err := p.StartNewDay()
		if err != nil {
			fmt.Printf("\n💀 Игра окончена: %s\n", err.Error())
			break
		}

		if p.Lock {
			break
		}

		ev, err := events.GetRandomEvent(eventPool)
		if err != nil {
			fmt.Printf("❌ Критическая ошибка пула событий: %s\n", err.Error())
			break
		}

		ev.Print()

		var choice int
		for {
			fmt.Print("Введите ваш выбор (1 или 2): ")
			_, scanErr := fmt.Scan(&choice)

			if scanErr != nil {
				fmt.Println("Ошибка: Пожалуйста, введите корректное число.")
				var discard string
				fmt.Scanln(&discard)
				continue
			}

			executeErr := ev.Execute(choice, &p)
			if executeErr != nil {
				fmt.Printf("❌ %s. Попробуйте еще раз.\n", executeErr.Error())
				continue
			}

			break
		}

		p.PrintStatus()
		if p.Lock && p.Health == 0 {
			fmt.Println("\n💀 Игра окончена: Вы погибли от полученных в событии повреждений.")
			break
		}
	}

	fmt.Println("\n=== Спасибо за игру! ===")
}
