package events

import (
	"errors"
	"fmt"
	"main/player"
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

	ev1, _ := ConstructEvent(Event{
		Message:  "Вы нашли чистый ручей среди скал.",
		Variants: [2]string{"Напиться и набрать воды с собой", "Пройти мимо, опасаясь засады"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				p.AddWater(5)
				return nil
			},
			func(p *player.Player) error {
				fmt.Println("Вы осторожно обошли ручей.")
				return nil
			},
		},
	})
	pool = append(pool, ev1)

	ev2, _ := ConstructEvent(Event{
		Message:  "На вас напал бродяга с ножом!",
		Variants: [2]string{"Принять бой", "Отдать ему часть припасов и убежать"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error {
				return p.ChangeHealth(-30)
			},
			func(p *player.Player) error {
				if p.Eat >= 3 {
					p.Eat -= 3
				} else {
					p.Eat = 0
				}
				fmt.Println("Вы бросили еду и убежали.")
				return nil
			},
		},
	})
	pool = append(pool, ev2)

	return pool
}

func (e Event) Print() {
	fmt.Println("\n========================================")
	fmt.Printf("СОБЫТИЕ: %s\n", e.Message)
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
