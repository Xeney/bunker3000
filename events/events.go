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
	Actions  [2]func(p *player.Player) string // Теперь возвращает string, а не error
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
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(5)
				return "💧 Вы набрали воды из ручья!"
			},
			func(p *player.Player) string {
				return "👣 Вы осторожно обошли ручей, не рискуя."
			},
		},
	})
	pool = append(pool, ev1)

	// 2. Ягодная поляна
	ev2, _ := ConstructEvent(Event{
		Message:  "🍓 Вы наткнулись на поляну с дикими ягодами.",
		Variants: [2]string{"Собрать ягоды", "Проверить, не ядовиты ли они сначала"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				return "🍓 Вы собрали полную корзину сладких ягод!"
			},
			func(p *player.Player) string {
				p.AddEat(1)
				return "🔍 Вы потратили время на проверку, но нашли только половину урожая."
			},
		},
	})
	pool = append(pool, ev2)

	// 3. Заброшенный склад
	ev3, _ := ConstructEvent(Event{
		Message:  "🏚️ Вы обнаружили заброшенный склад с консервами.",
		Variants: [2]string{"Взять все, что можно", "Взять немного, чтобы не перегружаться"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(8)
				return "🥫 Вы набрали полный рюкзак консервов! (Еда: +8)"
			},
			func(p *player.Player) string {
				p.AddEat(4)
				return "🥫 Вы взяли только часть консервов, чтобы не перегружаться. (Еда: +4)"
			},
		},
	})
	pool = append(pool, ev3)

	// 4. Дождь
	ev4, _ := ConstructEvent(Event{
		Message:  "☔ Начался сильный дождь.",
		Variants: [2]string{"Собрать дождевую воду", "Укрыться в пещере и переждать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(4)
				return "💧 Вы наполнили все емкости дождевой водой! (Вода: +4)"
			},
			func(p *player.Player) string {
				return "🏔️ Вы переждали дождь в пещере в полной безопасности."
			},
		},
	})
	pool = append(pool, ev4)

	// ========== БОЕВЫЕ СОБЫТИЯ ==========

	// 5. Нападение бродяги
	ev5, _ := ConstructEvent(Event{
		Message:  "⚔️ На вас напал бродяга с ножом!",
		Variants: [2]string{"Принять бой", "Отдать ему часть припасов и убежать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				err := p.ChangeHealth(-30)
				if err != nil {
					return "💀 Вы сражались, но получили смертельные раны..."
				}
				return "🤕 Вы победили бродягу, но получили серьезные раны. (-30% здоровья)"
			},
			func(p *player.Player) string {
				loss := int8(3)
				if p.Eat >= loss {
					p.Eat -= loss
				} else {
					p.Eat = 0
				}
				return fmt.Sprintf("🏃 Вы бросили %d еды и убежали. Живы, но голодны.", loss)
			},
		},
	})
	pool = append(pool, ev5)

	// 6. Дикое животное
	ev6, _ := ConstructEvent(Event{
		Message:  "🐺 На вас вышла голодная волчья стая!",
		Variants: [2]string{"Сражаться с волками", "Забраться на дерево"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				err := p.ChangeHealth(-40)
				if err != nil {
					return "💀 Волки оказались сильнее... Вы не выжили."
				}
				return "🗡️ Вы отбились от волков, но получили серьезные раны. (-40% здоровья)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "🌳 Вы просидели на дереве несколько часов. Волки ушли, но вы замерзли и устали. (-5% здоровья)"
			},
		},
	})
	pool = append(pool, ev6)

	// 7. Ловушка
	ev7, _ := ConstructEvent(Event{
		Message:  "🎣 Вы наткнулись на охотничью ловушку.",
		Variants: [2]string{"Осторожно обойти", "Попытаться забрать приманку"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(2)
				return "🍄 Вы аккуратно обошли ловушку и нашли неподалеку грибы. (Еда: +2)"
			},
			func(p *player.Player) string {
				err := p.ChangeHealth(-15)
				if err != nil {
					return "💀 Ловушка сработала смертельно..."
				}
				return "💢 Ловушка сработала! Вы получили травму. (-15% здоровья)"
			},
		},
	})
	pool = append(pool, ev7)

	// ========== ПОМОЩЬ И ВСТРЕЧИ ==========

	// 8. Странник
	ev8, _ := ConstructEvent(Event{
		Message:  "🧙 Встретился загадочный странник.",
		Variants: [2]string{"Попросить помощи", "Обменять припасы на информацию"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				p.AddWater(3)
				return "🎁 Странник поделился с вами едой и водой! (Еда: +3, Вода: +3)"
			},
			func(p *player.Player) string {
				if p.Eat >= 2 {
					p.Eat -= 2
					p.AddEat(5)
					p.AddWater(5)
					return "🗺️ Странник рассказал о тайнике с припасами! (Еда: +3, Вода: +5)"
				}
				return "😔 У вас нечем платить, странник уходит ни с чем."
			},
		},
	})
	pool = append(pool, ev8)

	// 9. Торговец
	ev9, _ := ConstructEvent(Event{
		Message:  "💰 Вы встретили торговца, который ищет еду.",
		Variants: [2]string{"Продать часть еды за воду", "Отказаться, еда нужнее"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				if p.Eat >= 3 {
					p.Eat -= 3
					p.AddWater(6)
					return "🔄 Вы обменяли 3 еды на 6 воды! Выгодная сделка."
				}
				return "😞 У вас недостаточно еды для обмена. Торговец уходит."
			},
			func(p *player.Player) string {
				p.AddWater(2)
				return "🎁 Торговец уходит, но оставляет вам небольшой подарок за честность. (Вода: +2)"
			},
		},
	})
	pool = append(pool, ev9)

	// ========== БОЛЕЗНИ И ОПАСНОСТИ ==========

	// 10. Испорченная вода
	ev10, _ := ConstructEvent(Event{
		Message:  "🤢 Вы выпили воду из сомнительного источника.",
		Variants: [2]string{"Попытаться очистить воду", "Искать другой источник"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(2)
				return "💧 Вы прокипятили воду и получили немного чистой воды. (Вода: +2)"
			},
			func(p *player.Player) string {
				p.AddWater(4)
				return "💧 Вы потратили время на поиски и нашли чистый ручей. (Вода: +4)"
			},
		},
	})
	pool = append(pool, ev10)

	// 11. Отравленная еда
	ev11, _ := ConstructEvent(Event{
		Message:  "🍄 Вы нашли грибы, но сомневаетесь в их съедобности.",
		Variants: [2]string{"Съесть грибы", "Выбросить их"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 30 {
					err := p.ChangeHealth(-20)
					if err != nil {
						return "💀 Грибы оказались смертельно ядовитыми..."
					}
					return "🤢 Грибы оказались ядовитыми! Вас тошнит. (-20% здоровья)"
				} else {
					p.AddEat(4)
					return "🍄 Грибы оказались съедобными и очень питательными! (Еда: +4)"
				}
			},
			func(p *player.Player) string {
				p.AddEat(2)
				return "🍓 Вы выбросили грибы, но нашли поблизости ягоды. (Еда: +2)"
			},
		},
	})
	pool = append(pool, ev11)

	// ========== ПОГОДА И ПРИРОДА ==========

	// 12. Жара
	ev12, _ := ConstructEvent(Event{
		Message:  "🌡️ Аномальная жара! Вы быстро теряете влагу.",
		Variants: [2]string{"Искать тень и отдыхать", "Продолжать путь"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "😓 Вы переждали жару в тени, но немного ослабли. (-5% здоровья)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-15)
				if p.Water >= 2 {
					p.Water -= 2
				}
				return "🥵 Вы продолжали идти и сильно обезвожились. (-15% здоровья, -2 воды)"
			},
		},
	})
	pool = append(pool, ev12)

	// 13. Туман
	ev13, _ := ConstructEvent(Event{
		Message:  "🌫️ Густой туман спустился на долину.",
		Variants: [2]string{"Идти осторожно", "Остановиться и ждать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				p.AddWater(2)
				return "🎁 Вы наткнулись на кладовку в тумане! (Еда: +3, Вода: +2)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "⏳ Вы переждали туман, но потеряли время. (-5% здоровья)"
			},
		},
	})
	pool = append(pool, ev13)

	// ========== НАХОДКИ И СОКРОВИЩА ==========

	// 14. Рюкзак с припасами
	ev14, _ := ConstructEvent(Event{
		Message:  "🎒 Вы нашли заброшенный рюкзак с припасами!",
		Variants: [2]string{"Забрать все", "Оставить на месте"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(4)
				p.AddWater(4)
				return "🎒 В рюкзаке были консервы и вода! (Еда: +4, Вода: +4)"
			},
			func(p *player.Player) string {
				return "🙏 Вы решили не трогать чужое. Возможно, это правильный выбор."
			},
		},
	})
	pool = append(pool, ev14)

	// 15. Аптечка
	ev15, _ := ConstructEvent(Event{
		Message:  "💊 Вы нашли медицинскую аптечку!",
		Variants: [2]string{"Использовать лекарства", "Сохранить на потом"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.ChangeHealth(30)
				return "💊 Вы перевязали раны и приняли лекарства! (+30% здоровья)"
			},
			func(p *player.Player) string {
				return "🎒 Вы положили аптечку в рюкзак на черный день."
			},
		},
	})
	pool = append(pool, ev15)

	// 16. Рыбалка
	ev16, _ := ConstructEvent(Event{
		Message:  "🎣 Вы нашли хорошее место для рыбалки.",
		Variants: [2]string{"Порыбачить", "Пройти мимо"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				success := rand.Intn(100)
				if success > 40 {
					p.AddEat(5)
					return "🐟 Удачная рыбалка! Вы поймали несколько рыб. (Еда: +5)"
				}
				p.ChangeHealth(-5)
				return "😞 Рыба не клюет. Вы потратили время впустую и ослабли. (-5% здоровья)"
			},
			func(p *player.Player) string {
				return "🚶 Вы не стали тратить время на рыбалку."
			},
		},
	})
	pool = append(pool, ev16)

	// 17. Дружелюбные поселенцы
	ev17, _ := ConstructEvent(Event{
		Message:  "🏘️ Вы наткнулись на поселок выживших.",
		Variants: [2]string{"Попросить убежища", "Обменяться ресурсами"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				p.AddWater(3)
				p.ChangeHealth(10)
				return "🏠 Вас накормили, дали воды и позволили отдохнуть! (Еда: +3, Вода: +3, Здоровье: +10)"
			},
			func(p *player.Player) string {
				if p.Eat >= 2 {
					p.Eat -= 2
					p.AddWater(5)
					return "🔄 Вы обменяли еду на воду и инструменты. (Еда: -2, Вода: +5)"
				}
				return "😔 У вас нечего предложить для обмена."
			},
		},
	})
	pool = append(pool, ev17)

	// 18. Землетрясение
	ev18, _ := ConstructEvent(Event{
		Message:  "🌋 Началось землетрясение!",
		Variants: [2]string{"Искать укрытие", "Бежать на открытое место"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				return "🏔️ Вы спрятались в пещере и не пострадали. Удачно!"
			},
			func(p *player.Player) string {
				err := p.ChangeHealth(-35)
				if err != nil {
					return "💀 На вас упал камень... Это был фатальный исход."
				}
				return "🪨 На вас упал камень! Вы серьезно ранены. (-35% здоровья)"
			},
		},
	})
	pool = append(pool, ev18)

	// 19. Пчелы и мед
	ev19, _ := ConstructEvent(Event{
		Message:  "🍯 Вы нашли улей с диким медом!",
		Variants: [2]string{"Попытаться забрать мед", "Оставить улей в покое"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 50 {
					p.ChangeHealth(-15)
					p.AddEat(5)
					return "🐝 Пчелы атаковали! Вы получили укусы, но мед забрали. (-15% здоровья, Еда: +5)"
				} else {
					p.AddEat(6)
					return "🍯 Вы удачно забрали мед, не потревожив пчел! (Еда: +6)"
				}
			},
			func(p *player.Player) string {
				return "🍯 Вы решили не рисковать и ушли. Мед остался пчелам."
			},
		},
	})
	pool = append(pool, ev19)

	// 20. Старая карта
	ev20, _ := ConstructEvent(Event{
		Message:  "🗺️ Вы нашли старую карту с отметкой тайника.",
		Variants: [2]string{"Пойти искать тайник", "Игнорировать карту"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(7)
				p.AddWater(7)
				p.ChangeHealth(10)
				return "💰 Вы нашли тайник с большим количеством припасов! (Еда: +7, Вода: +7, Здоровье: +10)"
			},
			func(p *player.Player) string {
				return "🗺️ Вы выбросили карту, решив не рисковать."
			},
		},
	})
	pool = append(pool, ev20)

	return pool
}

func (e Event) Print() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Printf("║ 📖 СОБЫТИЕ: %s\n", e.Message)
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Printf("\n 1. %s\n", e.Variants[0])
	fmt.Printf(" 2. %s\n", e.Variants[1])
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func (e Event) Execute(choice int, p *player.Player) (string, error) {
	if choice != 1 && choice != 2 {
		return "", errors.New("неверный выбор: введите 1 или 2")
	}

	index := choice - 1
	if e.Actions[index] == nil {
		return "", errors.New("критическая ошибка: действие для этого выбора не настроено")
	}

	// Вызываем действие и получаем сообщение
	resultMessage := e.Actions[index](p)
	return resultMessage, nil
}

func GetRandomEvent(pool []Event) (Event, error) {
	if len(pool) == 0 {
		return Event{}, errors.New("список событий пуст")
	}

	randomIndex := rand.Intn(len(pool))
	return pool[randomIndex], nil
}
