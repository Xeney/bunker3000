package events

import (
	"bunker3000/player"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type Event struct {
	Message  string
	Variants [2]string
	Actions  [2]func(p *player.Player) error
}

func ConstructEvent(event Event) (Event, error) {
	if strings.TrimSpace(event.Message) == "" {
		return Event{}, errors.New("сообщение события не может быть пустым")
	}
	if strings.TrimSpace(event.Variants[0]) == "" || strings.TrimSpace(event.Variants[1]) == "" {
		return Event{}, errors.New("оба варианта ответа должны быть заполнены")
	}
	if event.Actions[0] == nil || event.Actions[1] == nil {
		return Event{}, errors.New("для каждого варианта ответа должно быть назначено действие (функция не может быть nil)")
	}

	return event, nil
}

func ConstructEventsPool() []Event {
	var pool []Event

	// ========== РЕСУРСНЫЕ СОБЫТИЯ ==========

	// 1. Ручей с водой
	ev1, _ := ConstructEvent(Event{
		Message:  "🌊 Вы нашли чистый ручей среди скал.",
		Variants: [2]string{"Напиться и набрать воды с собой", "Пройти мимо, опасаясь засады"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.AddWater(5)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы осторожно обошли ручей, не рискуя.")
				return nil
			},
		},
	})
	pool = append(pool, ev1)

	// 2. Ягодная поляна
	ev2, _ := ConstructEvent(Event{
		Message:  "🍓 Вы наткнулись на поляну с дикими ягодами.",
		Variants: [2]string{"Собрать ягоды", "Проверить, не ядовиты ли они сначала"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.AddEat(3)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы потратили время на проверку, но нашли только половину урожая.")
				p.AddEat(1)
				return nil
			},
		},
	})
	pool = append(pool, ev2)

	// 3. Заброшенный склад
	ev3, _ := ConstructEvent(Event{
		Message:  "🏚️ Вы обнаружили заброшенный склад с консервами.",
		Variants: [2]string{"Взять все, что можно", "Взять немного, чтобы не перегружаться"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.AddEat(8)
				fmt.Println("Вы набрали полный рюкзак консервов, но немного замедлились.")
				return nil
			},
			func(p *player.Player) error {
				p.AddEat(4)
				return nil
			},
		},
	})
	pool = append(pool, ev3)

	// 4. Дождь
	ev4, _ := ConstructEvent(Event{
		Message:  "☔ Начался сильный дождь.",
		Variants: [2]string{"Собрать дождевую воду", "Укрыться в пещере и переждать"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.AddWater(4)
				fmt.Println("Вы наполнили все емкости дождевой водой!")
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы переждали дождь в безопасности, но потеряли время.")
				return nil
			},
		},
	})
	pool = append(pool, ev4)

	// ========== БОЕВЫЕ СОБЫТИЯ ==========

	// 5. Нападение бродяги
	ev5, _ := ConstructEvent(Event{
		Message:  "⚔️ На вас напал бродяга с ножом!",
		Variants: [2]string{"Принять бой", "Отдать ему часть припасов и убежать"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				return p.ChangeHealth(-30)
			},
			func(p *player.Player) error {
				loss := int8(3)
				if p.Eat >= loss {
					p.Eat -= loss
				} else {
					p.Eat = 0
				}
				fmt.Printf("Вы бросили %d еды и убежали.\n", loss)
				return nil
			},
		},
	})
	pool = append(pool, ev5)

	// 6. Дикое животное
	ev6, _ := ConstructEvent(Event{
		Message:  "🐺 На вас вышла голодная волчья стая!",
		Variants: [2]string{"Сражаться с волками", "Забраться на дерево"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				err := p.ChangeHealth(-40)
				fmt.Println("Вы отбились, но получили серьезные раны.")
				return err
			},
			func(p *player.Player) error {
				fmt.Println("Вы просидели на дереве несколько часов, теряя время.")
				p.ChangeHealth(-5) // Замерзли и устали
				return nil
			},
		},
	})
	pool = append(pool, ev6)

	// 7. Ловушка
	ev7, _ := ConstructEvent(Event{
		Message:  "🎣 Вы наткнулись на охотничью ловушку.",
		Variants: [2]string{"Осторожно обойти", "Попытаться забрать приманку"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				fmt.Println("Вы аккуратно обошли ловушку и нашли неподалеку грибы.")
				p.AddEat(2)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Ловушка сработала! Вы получили травму.")
				return p.ChangeHealth(-15)
			},
		},
	})
	pool = append(pool, ev7)

	// ========== ПОМОЩЬ И ВСТРЕЧИ ==========

	// 8. Странник
	ev8, _ := ConstructEvent(Event{
		Message:  "🧙 Встретился загадочный странник.",
		Variants: [2]string{"Попросить помощи", "Обменять припасы на информацию"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.AddEat(3)
				p.AddWater(3)
				fmt.Println("Странник поделился с вами едой и водой!")
				return nil
			},
			func(p *player.Player) error {
				if p.Eat >= 2 {
					p.Eat -= 2
					fmt.Println("Странник рассказал о тайнике с припасами неподалеку.")
					p.AddEat(5)
					p.AddWater(5)
				} else {
					fmt.Println("У вас нечем платить, странник уходит.")
				}
				return nil
			},
		},
	})
	pool = append(pool, ev8)

	// 9. Торговец
	ev9, _ := ConstructEvent(Event{
		Message:  "💰 Вы встретили торговца, который ищет еду.",
		Variants: [2]string{"Продать часть еды за воду", "Отказаться, еда нужнее"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				if p.Eat >= 3 {
					p.Eat -= 3
					p.AddWater(6)
					fmt.Println("Вы обменяли 3 еды на 6 воды!")
				} else {
					fmt.Println("У вас недостаточно еды для обмена.")
				}
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Торговец уходит, но оставляет вам небольшой подарок за честность.")
				p.AddWater(2)
				return nil
			},
		},
	})
	pool = append(pool, ev9)

	// ========== БОЛЕЗНИ И ОПАСНОСТИ ==========

	// 10. Испорченная вода
	ev10, _ := ConstructEvent(Event{
		Message:  "🤢 Вы выпили воду из сомнительного источника.",
		Variants: [2]string{"Попытаться очистить воду", "Искать другой источник"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				fmt.Println("Вы прокипятили воду, но потеряли часть.")
				p.AddWater(2)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы потратили время на поиски и нашли чистый ручей.")
				p.AddWater(4)
				return nil
			},
		},
	})
	pool = append(pool, ev10)

	// 11. Отравленная еда
	ev11, _ := ConstructEvent(Event{
		Message:  "🍄 Вы нашли грибы, но сомневаетесь в их съедобности.",
		Variants: [2]string{"Съесть грибы", "Выбросить их"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				risk := rand.Intn(100)
				if risk < 30 {
					fmt.Println("Грибы оказались ядовитыми! Вас тошнит.")
					return p.ChangeHealth(-20)
				} else {
					fmt.Println("Грибы оказались съедобными и очень питательными!")
					p.AddEat(4)
					return nil
				}
			},
			func(p *player.Player) error {
				fmt.Println("Вы выбросили грибы, но нашли поблизости ягоды.")
				p.AddEat(2)
				return nil
			},
		},
	})
	pool = append(pool, ev11)

	// ========== ПОГОДА И ПРИРОДА ==========

	// 12. Жара
	ev12, _ := ConstructEvent(Event{
		Message:  "🌡️ Аномальная жара! Вы быстро теряете влагу.",
		Variants: [2]string{"Искать тень и отдыхать", "Продолжать путь"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				fmt.Println("Вы переждали жару в тени, но потеряли время.")
				p.ChangeHealth(-5)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы продолжали идти и сильно обезвожились.")
				p.ChangeHealth(-15)
				if p.Water >= 2 {
					p.Water -= 2
				}
				return nil
			},
		},
	})
	pool = append(pool, ev12)

	// 13. Туман
	ev13, _ := ConstructEvent(Event{
		Message:  "🌫️ Густой туман спустился на долину.",
		Variants: [2]string{"Идти осторожно", "Остановиться и ждать"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				fmt.Println("Вы наткнулись на кладовку в тумане!")
				p.AddEat(3)
				p.AddWater(2)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы переждали туман, но потеряли день.")
				p.ChangeHealth(-5)
				return nil
			},
		},
	})
	pool = append(pool, ev13)

	// ========== НАХОДКИ И СОКРОВИЩА ==========

	// 14. Рюкзак с припасами
	ev14, _ := ConstructEvent(Event{
		Message:  "🎒 Вы нашли заброшенный рюкзак с припасами!",
		Variants: [2]string{"Забрать все", "Оставить на месте"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.AddEat(4)
				p.AddWater(4)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы решили не трогать чужое, но нашли монетку неподалеку.")
				return nil
			},
		},
	})
	pool = append(pool, ev14)

	// 15. Аптечка
	ev15, _ := ConstructEvent(Event{
		Message:  "💊 Вы нашли медицинскую аптечку!",
		Variants: [2]string{"Использовать лекарства", "Сохранить на потом"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.ChangeHealth(30)
				fmt.Println("Вы перевязали раны и приняли лекарства!")
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы положили аптечку в рюкзак на черный день.")
				return nil
			},
		},
	})
	pool = append(pool, ev15)

	// 16. Рыбалка
	ev16, _ := ConstructEvent(Event{
		Message:  "🎣 Вы нашли хорошее место для рыбалки.",
		Variants: [2]string{"Порыбачить", "Пройти мимо"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				success := rand.Intn(100)
				if success > 40 {
					fmt.Println("Удачная рыбалка! Вы поймали несколько рыб.")
					p.AddEat(5)
				} else {
					fmt.Println("Рыба не клюет. Вы потратили время впустую.")
					p.ChangeHealth(-5)
				}
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы не стали тратить время на рыбалку.")
				return nil
			},
		},
	})
	pool = append(pool, ev16)

	// 17. Дружелюбные поселенцы
	ev17, _ := ConstructEvent(Event{
		Message:  "🏘️ Вы наткнулись на поселок выживших.",
		Variants: [2]string{"Попросить убежища", "Обменяться ресурсами"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				fmt.Println("Вас накормили и дали отдохнуть.")
				p.AddEat(3)
				p.AddWater(3)
				p.ChangeHealth(10)
				return nil
			},
			func(p *player.Player) error {
				if p.Eat >= 2 {
					p.Eat -= 2
					p.AddWater(5)
					fmt.Println("Вы обменяли еду на воду и инструменты.")
				} else {
					fmt.Println("У вас нечего предложить для обмена.")
				}
				return nil
			},
		},
	})
	pool = append(pool, ev17)

	// 18. Землетрясение
	ev18, _ := ConstructEvent(Event{
		Message:  "🌋 Началось землетрясение!",
		Variants: [2]string{"Искать укрытие", "Бежать на открытое место"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				fmt.Println("Вы спрятались в пещере и не пострадали.")
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("На вас упал камень! Вы серьезно ранены.")
				return p.ChangeHealth(-35)
			},
		},
	})
	pool = append(pool, ev18)

	// 19. Пчелы и мед
	ev19, _ := ConstructEvent(Event{
		Message:  "🍯 Вы нашли улей с диким медом!",
		Variants: [2]string{"Попытаться забрать мед", "Оставить улей в покое"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				risk := rand.Intn(100)
				if risk < 50 {
					fmt.Println("Пчелы атаковали! Вы получили множество укусов.")
					p.ChangeHealth(-15)
					p.AddEat(5)
				} else {
					fmt.Println("Вы удачно забрали мед, не потревожив пчел!")
					p.AddEat(6)
				}
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы решили не рисковать и ушли.")
				return nil
			},
		},
	})
	pool = append(pool, ev19)

	// 20. Старая карта
	ev20, _ := ConstructEvent(Event{
		Message:  "🗺️ Вы нашли старую карту с отметкой тайника.",
		Variants: [2]string{"Пойти искать тайник", "Игнорировать карту"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				fmt.Println("Вы нашли тайник с большим количеством припасов!")
				p.AddEat(7)
				p.AddWater(7)
				p.ChangeHealth(10)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы выбросили карту, решив не рисковать.")
				return nil
			},
		},
	})
	pool = append(pool, ev20)

	return pool
}

func (e Event) Print() {
	fmt.Println("\n========================================")
	fmt.Printf("📖 СОБЫТИЕ: %s\n", e.Message)
	fmt.Println("========================================")
	fmt.Printf("1.  %s\n", e.Variants[0])
	fmt.Printf("2.  %s\n", e.Variants[1])
	fmt.Println("----------------------------------------")
}

func (e Event) Execute(choice int, p *player.Player) error {
	if choice != 1 && choice != 2 {
		return errors.New("неверный выбор: введите 1 или 2")
	}

	index := choice - 1
	if e.Actions[index] == nil {
		return errors.New("критическая ошибка: действие для этого выбора не настроено")
	}

	return e.Actions[index](p)
}

func GetRandomEvent(pool []Event) (Event, error) {
	if len(pool) == 0 {
		return Event{}, errors.New("список событий пуст")
	}

	randomIndex := rand.Intn(len(pool))
	return pool[randomIndex], nil
}
