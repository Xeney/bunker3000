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
	Actions  [2]func(p *player.Player) string
	Category string
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
	if event.Category == "" {
		event.Category = "other"
	}

	return event, nil
}

func categoryColor(cat string) string {
	switch cat {
	case "resource":
		return "\033[92m"
	case "combat":
		return "\033[91m"
	case "helper":
		return "\033[94m"
	case "weather":
		return "\033[96m"
	case "hazard":
		return "\033[93m"
	case "find":
		return "\033[95m"
	default:
		return "\033[97m"
	}
}

func categoryLabel(cat string) string {
	switch cat {
	case "resource":
		return "[РЕСУРСЫ]"
	case "combat":
		return "[БИТВА]"
	case "helper":
		return "[ПОМОЩЬ]"
	case "weather":
		return "[ПОГОДА]"
	case "hazard":
		return "[ОПАСНОСТЬ]"
	case "find":
		return "[НАХОДКА]"
	default:
		return "[СОБЫТИЕ]"
	}
}

func ConstructEventsPool() []Event {
	var pool []Event

	// ========== РЕСУРСНЫЕ СОБЫТИЯ ==========

	// 1. Ручей с водой
	ev1, _ := ConstructEvent(Event{
		Message:  "[РУЧЕЙ] Вы нашли чистый ручей среди скал.",
		Variants: [2]string{"Напиться и набрать воды с собой", "Пройти мимо, опасаясь засады"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(5)
				return "[ВОДА] Вы набрали воды из ручья!"
			},
			func(p *player.Player) string {
				return "[ОСТОРОЖНО] Вы осторожно обошли ручей, не рискуя."
			},
		},
		Category: "resource",
	})
	pool = append(pool, ev1)

	// 2. Ягодная поляна
	ev2, _ := ConstructEvent(Event{
		Message:  "[ЯГОДЫ] Вы наткнулись на поляну с дикими ягодами.",
		Variants: [2]string{"Собрать ягоды", "Проверить, не ядовиты ли они сначала"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				return "[ЕДА] Вы собрали полную корзину сладких ягод!"
			},
			func(p *player.Player) string {
				p.AddEat(1)
				return "[ПРОВЕРКА] Вы потратили время на проверку, но нашли только половину урожая."
			},
		},
		Category: "resource",
	})
	pool = append(pool, ev2)

	// 3. Заброшенный склад
	ev3, _ := ConstructEvent(Event{
		Message:  "[СКЛАД] Вы обнаружили заброшенный склад с консервами.",
		Variants: [2]string{"Взять все, что можно", "Взять немного, чтобы не перегружаться"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(8)
				return "[ЕДА] Вы набрали полный рюкзак консервов! (Еда: +8)"
			},
			func(p *player.Player) string {
				p.AddEat(4)
				return "[ЕДА] Вы взяли только часть консервов, чтобы не перегружаться. (Еда: +4)"
			},
		},
		Category: "resource",
	})
	pool = append(pool, ev3)

	// 4. Дождь
	ev4, _ := ConstructEvent(Event{
		Message:  "[ДОЖДЬ] Начался сильный дождь.",
		Variants: [2]string{"Собрать дождевую воду", "Укрыться в пещере и переждать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(4)
				return "[ВОДА] Вы наполнили все емкости дождевой водой! (Вода: +4)"
			},
			func(p *player.Player) string {
				return "[УКРЫТИЕ] Вы переждали дождь в пещере в полной безопасности."
			},
		},
		Category: "weather",
	})
	pool = append(pool, ev4)

	// ========== БОЕВЫЕ СОБЫТИЯ ==========

	// 5. Нападение бродяги
	ev5, _ := ConstructEvent(Event{
		Message:  "[НАПАДЕНИЕ] На вас напал бродяга с ножом!",
		Variants: [2]string{"Принять бой", "Отдать ему часть припасов и убежать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				err := p.ChangeHealth(-30)
				if err != nil {
					return "[СМЕРТЬ] Вы сражались, но получили смертельные раны..."
				}
				return "[БОЙ] Вы победили бродягу, но получили серьезные раны. (-30 здоровья)"
			},
			func(p *player.Player) string {
				loss := int8(3)
				if p.Eat >= loss {
					p.Eat -= loss
				} else {
					p.Eat = 0
				}
				return fmt.Sprintf("[БЕГСТВО] Вы бросили %d еды и убежали.", loss)
			},
		},
		Category: "combat",
	})
	pool = append(pool, ev5)

	// 6. Дикое животное
	ev6, _ := ConstructEvent(Event{
		Message:  "[ВОЛКИ] На вас вышла голодная волчья стая!",
		Variants: [2]string{"Сражаться с волками", "Забраться на дерево"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				err := p.ChangeHealth(-40)
				if err != nil {
					return "[СМЕРТЬ] Волки оказались сильнее... Вы не выжили."
				}
				return "[БОЙ] Вы отбились от волков, но получили серьезные раны. (-40 здоровья)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "[СПАСЕНИЕ] Вы просидели на дереве несколько часов. Волки ушли. (-5 здоровья)"
			},
		},
		Category: "combat",
	})
	pool = append(pool, ev6)

	// 7. Ловушка
	ev7, _ := ConstructEvent(Event{
		Message:  "[ЛОВУШКА] Вы наткнулись на охотничью ловушку.",
		Variants: [2]string{"Осторожно обойти", "Попытаться забрать приманку"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(2)
				return "[НАХОДКА] Вы обошли ловушку и нашли грибы. (Еда: +2)"
			},
			func(p *player.Player) string {
				err := p.ChangeHealth(-15)
				if err != nil {
					return "[СМЕРТЬ] Ловушка сработала смертельно..."
				}
				return "[ТРАВМА] Ловушка сработала! Вы получили травму. (-15 здоровья)"
			},
		},
		Category: "hazard",
	})
	pool = append(pool, ev7)

	// ========== ПОМОЩЬ И ВСТРЕЧИ ==========

	// 8. Странник
	ev8, _ := ConstructEvent(Event{
		Message:  "[СТРАННИК] Встретился загадочный странник.",
		Variants: [2]string{"Попросить помощи", "Обменять припасы на информацию"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				p.AddWater(3)
				return "[ПОМОЩЬ] Странник поделился с вами едой и водой! (Еда: +3, Вода: +3)"
			},
			func(p *player.Player) string {
				if p.Eat >= 2 {
					p.Eat -= 2
					p.AddEat(5)
					p.AddWater(5)
					return "[СДЕЛКА] Странник рассказал о тайнике с припасами! (Еда: +3, Вода: +5)"
				}
				return "[НЕУДАЧА] У вас нечем платить, странник уходит."
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev8)

	// 9. Торговец
	ev9, _ := ConstructEvent(Event{
		Message:  "[ТОРГОВЕЦ] Вы встретили торговца.",
		Variants: [2]string{"Продать часть еды за воду", "Отказаться"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				if p.Eat >= 3 {
					p.Eat -= 3
					p.AddWater(6)
					return "[СДЕЛКА] Вы обменяли 3 еды на 6 воды!"
				}
				return "[НЕУДАЧА] Недостаточно еды для обмена. Торговец уходит."
			},
			func(p *player.Player) string {
				p.AddWater(2)
				return "[ПОДАРОК] Торговец оставил воду за честность. (Вода: +2)"
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev9)

	// ========== БОЛЕЗНИ И ОПАСНОСТИ ==========

	// 10. Испорченная вода
	ev10, _ := ConstructEvent(Event{
		Message:  "[ВОДА] Вы выпили воду из сомнительного источника.",
		Variants: [2]string{"Попытаться очистить воду", "Искать другой источник"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(2)
				return "[ВОДА] Вы прокипятили воду и получили чистую. (Вода: +2)"
			},
			func(p *player.Player) string {
				p.AddWater(4)
				return "[ВОДА] Вы нашли чистый ручей. (Вода: +4)"
			},
		},
		Category: "hazard",
	})
	pool = append(pool, ev10)

	// 11. Отравленная еда
	ev11, _ := ConstructEvent(Event{
		Message:  "[ГРИБЫ] Вы нашли грибы, но сомневаетесь в их съедобности.",
		Variants: [2]string{"Съесть грибы", "Выбросить их"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 30 {
					err := p.ChangeHealth(-20)
					if err != nil {
						return "[СМЕРТЬ] Грибы оказались смертельно ядовитыми..."
					}
					return "[ОТРАВЛЕНИЕ] Грибы ядовитые! Вас тошнит. (-20 здоровья)"
				} else {
					p.AddEat(4)
					return "[ЕДА] Грибы оказались съедобными и питательными! (Еда: +4)"
				}
			},
			func(p *player.Player) string {
				p.AddEat(2)
				return "[ЕДА] Вы нашли ягоды рядом. (Еда: +2)"
			},
		},
		Category: "hazard",
	})
	pool = append(pool, ev11)

	// ========== ПОГОДА И ПРИРОДА ==========

	// 12. Жара
	ev12, _ := ConstructEvent(Event{
		Message:  "[ЖАРА] Аномальная жара! Вы быстро теряете влагу.",
		Variants: [2]string{"Искать тень и отдыхать", "Продолжать путь"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "[ОТДЫХ] Вы переждали жару в тени. (-5 здоровья)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-15)
				if p.Water >= 2 {
					p.Water -= 2
				}
				return "[ОПАСНОСТЬ] Вы обезвожились. (-15 здоровья, -2 воды)"
			},
		},
		Category: "weather",
	})
	pool = append(pool, ev12)

	// 13. Туман
	ev13, _ := ConstructEvent(Event{
		Message:  "[ТУМАН] Густой туман спустился на долину.",
		Variants: [2]string{"Идти осторожно", "Остановиться и ждать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				p.AddWater(2)
				return "[НАХОДКА] Вы наткнулись на кладовку в тумане! (Еда: +3, Вода: +2)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "[ОЖИДАНИЕ] Вы переждали туман. (-5 здоровья)"
			},
		},
		Category: "weather",
	})
	pool = append(pool, ev13)

	// ========== НАХОДКИ И СОКРОВИЩА ==========

	// 14. Рюкзак с припасами
	ev14, _ := ConstructEvent(Event{
		Message:  "[РЮКЗАК] Вы нашли заброшенный рюкзак!",
		Variants: [2]string{"Забрать все", "Оставить на месте"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(4)
				p.AddWater(4)
				return "[НАХОДКА] В рюкзаке консервы и вода! (Еда: +4, Вода: +4)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы решили не трогать чужое."
			},
		},
		Category: "find",
	})
	pool = append(pool, ev14)

	// 15. Аптечка
	ev15, _ := ConstructEvent(Event{
		Message:  "[АПТЕЧКА] Вы нашли медицинскую аптечку!",
		Variants: [2]string{"Использовать лекарства", "Сохранить на потом"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.ChangeHealth(30)
				return "[ЛЕЧЕНИЕ] Вы перевязали раны! (+30 здоровья)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы положили аптечку в рюкзак."
			},
		},
		Category: "find",
	})
	pool = append(pool, ev15)

	// 16. Рыбалка
	ev16, _ := ConstructEvent(Event{
		Message:  "[РЫБАЛКА] Вы нашли хорошее место для рыбалки.",
		Variants: [2]string{"Порыбачить", "Пройти мимо"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				success := rand.Intn(100)
				if success > 40 {
					p.AddEat(5)
					return "[ЕДА] Удачная рыбалка! (Еда: +5)"
				}
				p.ChangeHealth(-5)
				return "[НЕУДАЧА] Рыба не клюет. (-5 здоровья)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы не стали тратить время."
			},
		},
		Category: "resource",
	})
	pool = append(pool, ev16)

	// 17. Дружелюбные поселенцы
	ev17, _ := ConstructEvent(Event{
		Message:  "[ПОСЕЛЕНИЕ] Вы наткнулись на поселок выживших.",
		Variants: [2]string{"Попросить убежища", "Обменяться ресурсами"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				p.AddWater(3)
				p.ChangeHealth(10)
				return "[ПОМОЩЬ] Вас накормили и приютили! (Еда: +3, Вода: +3, Здоровье: +10)"
			},
			func(p *player.Player) string {
				if p.Eat >= 2 {
					p.Eat -= 2
					p.AddWater(5)
					return "[СДЕЛКА] Вы обменяли еду на воду. (Еда: -2, Вода: +5)"
				}
				return "[НЕУДАЧА] У вас нечего предложить."
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev17)

	// 18. Землетрясение
	ev18, _ := ConstructEvent(Event{
		Message:  "[ЗЕМЛЕТРЯСЕНИЕ] Началось землетрясение!",
		Variants: [2]string{"Искать укрытие", "Бежать на открытое место"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				return "[УКРЫТИЕ] Вы спрятались в пещере и не пострадали."
			},
			func(p *player.Player) string {
				err := p.ChangeHealth(-35)
				if err != nil {
					return "[СМЕРТЬ] Камень упал... Фатальный исход."
				}
				return "[ТРАВМА] Камень упал! (-35 здоровья)"
			},
		},
		Category: "hazard",
	})
	pool = append(pool, ev18)

	// 19. Пчелы и мед
	ev19, _ := ConstructEvent(Event{
		Message:  "[МЕД] Вы нашли улей с диким медом!",
		Variants: [2]string{"Забрать мед", "Оставить улей"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 50 {
					p.ChangeHealth(-15)
					p.AddEat(5)
					return "[РИСК] Пчелы покусали, но мед ваш! (-15 здоровья, Еда: +5)"
				} else {
					p.AddEat(6)
					return "[УДАЧА] Вы забрали мед без потерь! (Еда: +6)"
				}
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы решили не рисковать."
			},
		},
		Category: "resource",
	})
	pool = append(pool, ev19)

	// 20. Старая карта
	ev20, _ := ConstructEvent(Event{
		Message:  "[КАРТА] Вы нашли старую карту с отметкой тайника.",
		Variants: [2]string{"Искать тайник", "Игнорировать карту"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(7)
				p.AddWater(7)
				p.ChangeHealth(10)
				return "[ТАЙНИК] Вы нашли припасы! (Еда: +7, Вода: +7, Здоровье: +10)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы выбросили карту."
			},
		},
		Category: "find",
	})
	pool = append(pool, ev20)

	// ========== НОВЫЕ СОБЫТИЯ ==========

	// 21. Старый колодец
	ev21, _ := ConstructEvent(Event{
		Message:  "[КОЛОДЕЦ] Вы нашли старый колодец.",
		Variants: [2]string{"Набрать воды", "Проверить качество воды"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(3)
				return "[ВОДА] Вода чистая! (Вода: +3)"
			},
			func(p *player.Player) string {
				p.AddWater(2)
				risk := rand.Intn(100)
				if risk < 20 {
					p.ChangeHealth(-10)
					return "[ОПАСНОСТЬ] Вода оказалась зараженной. (-10 здоровья, Вода: +2)"
				}
				return "[ВОДА] Вода хорошая! (Вода: +2)"
			},
		},
		Category: "resource",
	})
	pool = append(pool, ev21)

	// 22. Засада
	ev22, _ := ConstructEvent(Event{
		Message:  "[ЗАСАДА] Вы попали в засаду бандитов!",
		Variants: [2]string{"Прорываться с боем", "Спрятаться и ждать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				err := p.ChangeHealth(-25)
				if err != nil {
					return "[СМЕРТЬ] Бандиты одолели..."
				}
				return "[БОЙ] Вы прорвались! (-25 здоровья)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "[ОЖИДАНИЕ] Бандиты ушли, вы потеряли время. (-5 здоровья)"
			},
		},
		Category: "combat",
	})
	pool = append(pool, ev22)

	// 23. Снежная буря
	ev23, _ := ConstructEvent(Event{
		Message:  "[БУРЯ] Началась сильная снежная буря!",
		Variants: [2]string{"Укрыться в снежной пещере", "Продолжать идти"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				return "[УКРЫТИЕ] Пурга стихла, вы в порядке. (-5 здоровья)"
			},
			func(p *player.Player) string {
				err := p.ChangeHealth(-30)
				if err != nil {
					return "[СМЕРТЬ] Буря была слишком сильной..."
				}
				return "[ОПАСНОСТЬ] Вы чуть не замерзли! (-30 здоровья)"
			},
		},
		Category: "weather",
	})
	pool = append(pool, ev23)

	// 24. Заброшенный дом
	ev24, _ := ConstructEvent(Event{
		Message:  "[ДОМ] Вы нашли заброшенный дом.",
		Variants: [2]string{"Обыскать дом", "Пройти мимо"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				find := rand.Intn(3)
				switch find {
				case 0:
					p.AddEat(3)
					return "[ЕДА] В доме были консервы! (Еда: +3)"
				case 1:
					p.AddWater(3)
					return "[ВОДА] Вы нашли бутылки с водой! (Вода: +3)"
				default:
					p.AddEat(2)
					p.AddWater(2)
					return "[НАХОДКА] В доме есть припасы! (Еда: +2, Вода: +2)"
				}
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы прошли мимо."
			},
		},
		Category: "find",
	})
	pool = append(pool, ev24)

	// 25. Раненый путник
	ev25, _ := ConstructEvent(Event{
		Message:  "[ПУТНИК] Вы нашли раненого путника.",
		Variants: [2]string{"Помочь путнику", "Пройти мимо"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(-2)
				p.ChangeHealth(-5)
				if p.Eat < 0 {
					p.Eat = 0
				}
				reward := rand.Intn(100)
				if reward > 50 {
					p.AddWater(4)
					return "[НАГРАДА] Путник оказался торговцем и отблагодарил! (Еда: -2, Вода: +4)"
				}
				return "[ПОМОЩЬ] Путник ушел, забрав вашу еду. (Еда: -2, Здоровье: -5)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы прошли мимо."
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev25)

	// 26. Лесной пожар
	ev26, _ := ConstructEvent(Event{
		Message:  "[ПОЖАР] В лесу начался пожар!",
		Variants: [2]string{"Бежать от огня", "Укрыться в ручье"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				err := p.ChangeHealth(-20)
				if err != nil {
					return "[СМЕРТЬ] Огонь настиг вас..."
				}
				p.AddWater(-2)
				if p.Water < 0 {
					p.Water = 0
				}
				return "[БЕГСТВО] Вы убежали от огня. (-20 здоровья, -2 воды)"
			},
			func(p *player.Player) string {
				p.AddWater(2)
				return "[СПАСЕНИЕ] Ручей спас вас! (Вода: +2)"
			},
		},
		Category: "hazard",
	})
	pool = append(pool, ev26)

	// 27. Затопленный подвал
	ev27, _ := ConstructEvent(Event{
		Message:  "[ПОДВАЛ] Вы нашли затопленный подвал.",
		Variants: [2]string{"Обыскать подвал", "Не рисковать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 40 {
					p.ChangeHealth(-15)
					return "[ТРАВМА] Пол провалился! (-15 здоровья)"
				}
				p.AddEat(5)
				p.AddWater(3)
				return "[НАХОДКА] В подвале много припасов! (Еда: +5, Вода: +3)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы прошли мимо."
			},
		},
		Category: "find",
	})
	pool = append(pool, ev27)

	// 28. Сигнальный костер
	ev28, _ := ConstructEvent(Event{
		Message:  "[КОСТЕР] Вы заметили сигнальный дым.",
		Variants: [2]string{"Пойти на дым", "Игнорировать"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(4)
				p.ChangeHealth(5)
				return "[ПОМОЩЬ] Вы нашли группу выживших! (Еда: +4, Здоровье: +5)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы проигнорировали сигнал."
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev28)

	// 29. Горная тропа
	ev29, _ := ConstructEvent(Event{
		Message:  "[ТРОПА] Вы нашли горную тропу.",
		Variants: [2]string{"Пойти по тропе", "Вернуться на дорогу"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 30 {
					p.ChangeHealth(-20)
					return "[ПАДЕНИЕ] Вы сорвались с тропы! (-20 здоровья)"
				}
				p.AddEat(3)
				return "[ПУТЬ] Тропа вывела к поляне с ягодами! (Еда: +3)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы вернулись на безопасную дорогу."
			},
		},
		Category: "hazard",
	})
	pool = append(pool, ev29)

	// 30. Брошенная машина
	ev30, _ := ConstructEvent(Event{
		Message:  "[МАШИНА] Вы нашли брошенный автомобиль.",
		Variants: [2]string{"Обыскать машину", "Пройти мимо"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				find := rand.Intn(4)
				switch find {
				case 0:
					p.AddEat(4)
					return "[ЕДА] В багажнике консервы! (Еда: +4)"
				case 1:
					p.AddWater(4)
					return "[ВОДА] Вы нашли бутылки с водой! (Вода: +4)"
				case 2:
					return "[ПУСТО] В машине ничего полезного."
				default:
					p.AddEat(2)
					p.AddWater(2)
					return "[НАХОДКА] Немного припасов. (Еда: +2, Вода: +2)"
				}
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы не стали тратить время."
			},
		},
		Category: "find",
	})
	pool = append(pool, ev30)

	// 31. Ядовитый источник
	ev31, _ := ConstructEvent(Event{
		Message:  "[ИСТОЧНИК] Вы нашли подземный источник.",
		Variants: [2]string{"Попить воды", "Проверить воду"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddWater(3)
				risk := rand.Intn(100)
				if risk < 25 {
					p.ChangeHealth(-15)
					return "[ОТРАВЛЕНИЕ] Вода заражена! (Вода: +3, -15 здоровья)"
				}
				return "[ВОДА] Вода чистая и вкусная! (Вода: +3)"
			},
			func(p *player.Player) string {
				p.AddWater(2)
				return "[ВОДА] Вы проверили - вода безопасна. (Вода: +2)"
			},
		},
		Category: "hazard",
	})
	pool = append(pool, ev31)

	// 32. Ночная охота
	ev32, _ := ConstructEvent(Event{
		Message:  "[ОХОТА] Ночью вы слышите звуки зверей.",
		Variants: [2]string{"Пойти на охоту", "Остаться в лагере"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				success := rand.Intn(100)
				if success > 60 {
					p.AddEat(6)
					return "[ДОБЫЧА] Вы поймали крупную дичь! (Еда: +6)"
				}
				p.ChangeHealth(-10)
				return "[НЕУДАЧА] Охота провалилась, вы ранены. (-10 здоровья)"
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы остались в лагере."
			},
		},
		Category: "resource",
	})
	pool = append(pool, ev32)

	// 33. Крик о помощи
	ev33, _ := ConstructEvent(Event{
		Message:  "[КРИК] Вы слышите крик о помощи из заброшенных руин.",
		Variants: [2]string{"Пойти на крик и помочь", "Пройти мимо — это может быть ловушкой"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 25 {
					p.ChangeHealth(-10)
					p.Flags["helped_cries"] = true
					return "[ЗАСАДА] Это была ловушка! Вы ранены. (-10 здоровья)\nНо вы заметили следы настоящих припасов."
				}
				p.AddEat(3)
				p.AddWater(2)
				p.Flags["helped_cries"] = true
				return "[ПОМОЩЬ] Вы отбили атаку и нашли припасы в руинах! (Еда: +3, Вода: +2)"
			},
			func(p *player.Player) string {
				p.Flags["ignored_cries"] = true
				return "[ОСТОРОЖНО] Вы прошли мимо. Крики стихли."
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev33)

	// 34. Раненый бандит
	ev34, _ := ConstructEvent(Event{
		Message:  "[БАНДИТ] Вы нашли раненого бандита у дороги. Он просит пощады.",
		Variants: [2]string{"Добить бандита", "Перевязать и отпустить"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.AddEat(3)
				p.Flags["killed_bandit"] = true
				return "[БЕЗЖАЛОСТНОСТЬ] Вы обыскали тело. (Еда: +3)"
			},
			func(p *player.Player) string {
				p.ChangeHealth(-5)
				p.Flags["spared_bandit"] = true
				return "[МИЛОСЕРДИЕ] Бандит благодарит вас и уходит. (-5 здоровья на бинты)"
			},
		},
		Category: "combat",
	})
	pool = append(pool, ev34)

	// 35. Таинственная карта
	ev35, _ := ConstructEvent(Event{
		Message:  "[КАРТА] В старом тайнике вы нашли потрёпанную карту с пометками.",
		Variants: [2]string{"Взять карту и изучить", "Оставить — мало ли что"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				p.Flags["found_map"] = true
				return "[НАХОДКА] На карте отмечен тайник! Нужно будет поискать."
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы оставили карту."
			},
		},
		Category: "find",
	})
	pool = append(pool, ev35)

	// 36. Старик просит еду
	ev36, _ := ConstructEvent(Event{
		Message:  "[СТАРИК] Дряхлый старик сидит у дороги и просит еды.",
		Variants: [2]string{"Поделиться едой ( -3 еды)", "Пройти мимо"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				if p.Eat < 3 {
					p.Flags["refused_oldman"] = true
					return "[НЕУДАЧА] У вас самих не хватает еды. Старик понимающе кивает."
				}
				p.AddEat(-3)
				p.Flags["fed_oldman"] = true
				return "[ДОБРОТА] Старик благословляет вас: «Добро вернётся». (Еда: -3)"
			},
			func(p *player.Player) string {
				p.Flags["refused_oldman"] = true
				return "[ВЫБОР] Вы прошли мимо. Старик смотрит вслед."
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev36)

	// 37. Запертый склад
	ev37, _ := ConstructEvent(Event{
		Message:  "[СКЛАД] Вы нашли запертый склад с довоенными припасами.",
		Variants: [2]string{"Взломать замок (шумно)", "Поискать другой вход"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 35 {
					p.ChangeHealth(-15)
					p.Flags["broke_storage"] = true
					return "[ОПАСНОСТЬ] На шум пришли бандиты! Вы ранены. (-15 здоровья)"
				}
				p.AddEat(5)
				p.AddWater(5)
				p.Flags["broke_storage"] = true
				return "[НАХОДКА] Внутри ящики с консервами и водой! (Еда: +5, Вода: +5)"
			},
			func(p *player.Player) string {
				p.AddEat(2)
				p.AddWater(2)
				p.Flags["left_storage"] = true
				return "[ПОИСК] Вы нашли чёрный вход и взяли немного припасов. (Еда: +2, Вода: +2)"
			},
		},
		Category: "find",
	})
	pool = append(pool, ev37)

	// 38. Заброшенная рация
	ev38, _ := ConstructEvent(Event{
		Message:  "[РАЦИЯ] В заброшенном лагере вы нашли рабочую рацию.",
		Variants: [2]string{"Попытаться вызвать помощь", "Не рисковать — сигнал могут перехватить"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string {
				risk := rand.Intn(100)
				if risk < 30 {
					p.ChangeHealth(-10)
					return "[ПЕРЕХВАТ] Сигнал перехватили бандиты! Они выследили вас. (-10 здоровья)"
				}
				p.Flags["answered_signal"] = true
				return "[УДАЧА] Вы связались с дружественным караваном! Они обещали помочь при встрече."
			},
			func(p *player.Player) string {
				return "[ВЫБОР] Вы оставили рацию."
			},
		},
		Category: "helper",
	})
	pool = append(pool, ev38)

	return pool
}

func (e Event) Print() {
	clr := categoryColor(e.Category)
	label := categoryLabel(e.Category)
	reset := "\033[0m"

	fmt.Printf("\n%s+------------------------------------------+%s\n", clr, reset)
	fmt.Printf("%s| %s %s%s\n", clr, label, e.Message, reset)
	fmt.Printf("%s+------------------------------------------+%s\n", clr, reset)
	fmt.Printf("\n%s 1.%s %s\n", "\033[93m", reset, e.Variants[0])
	fmt.Printf("%s 2.%s %s\n", "\033[93m", reset, e.Variants[1])
	fmt.Println("\n--------------------------------------------")
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
