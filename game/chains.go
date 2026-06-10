package game

import (
	"bunker3000/events"
	"bunker3000/player"
	"math/rand"
)

type ChainDef struct {
	ID          string
	TriggerFlag string
	Steps       []events.Event
	FinalReward func(p *player.Player) string
}

func findChain(id string) *ChainDef {
	for i := range AllChains {
		if AllChains[i].ID == id {
			return &AllChains[i]
		}
	}
	return nil
}

func activateChain(gs *GameState, chainID string) {
	gs.Player.Flags[chainID+"_started"] = true
	gs.ActiveChainID = chainID
	gs.ChainStep = 0
}

var AllChains = []ChainDef{
	// Цепочка: Крик о помощи → Тайник
	{
		ID:          "chain_cries",
		TriggerFlag: "helped_cries",
		Steps: []events.Event{
			{
				Message:  "[ВСТРЕЧА] Спасённый догоняет вас: «Я помню вашу помощь! Вот карта тайника.»",
				Variants: [2]string{"Взять карту и искать тайник", "Вежливо отказаться — некогда"},
				Actions: [2]func(p *player.Player) string{
					func(p *player.Player) string {
						p.Flags["chain_cries_map"] = true
						return "[КАРТА] Вы взяли карту. Тайник отмечен крестом."
					},
					func(p *player.Player) string {
						return "[ВЫБОР] Вы отказались и продолжили путь."
					},
				},
				Category: "helper",
			},
			{
				Message:  "[ТАЙНИК] Вы нашли место, отмеченное на карте!",
				Variants: [2]string{"Раскопать тайник", "Осторожно проверить на ловушки"},
				Actions: [2]func(p *player.Player) string{
					func(p *player.Player) string {
						p.AddEat(5)
						p.AddWater(5)
						return "[НАХОДКА] В тайнике консервы и вода! (Еда: +5, Вода: +5)"
					},
					func(p *player.Player) string {
						p.AddEat(3)
						p.AddWater(3)
						return "[ОСТОРОЖНО] Ловушек нет, но вы потеряли время. (Еда: +3, Вода: +3)"
					},
				},
				Category: "find",
			},
		},
		FinalReward: func(p *player.Player) string {
			p.AddEat(3)
			p.AddWater(3)
			p.ChangeHealth(10)
			return "[НАГРАДА] Спасённый оказался торговцем. Он оставил вам подарок! (Еда: +3, Вода: +3, Здоровье: +10)"
		},
	},

	// Цепочка: Рация → Караван
	{
		ID:          "chain_signal",
		TriggerFlag: "answered_signal",
		Steps: []events.Event{
			{
				Message:  "[РАДИО] Рация оживает: «Говорит караван «Свобода». Нас атаковали! Нужна помощь!»",
				Variants: [2]string{"Идти на помощь каравану", "Ответить, что ничем не можете помочь"},
				Actions: [2]func(p *player.Player) string{
					func(p *player.Player) string {
						p.Flags["chain_signal_help"] = true
						return "[ПОМОЩЬ] Вы идёте на помощь каравану."
					},
					func(p *player.Player) string {
						return "[ВЫБОР] Вы оставляете рацию."
					},
				},
				Category: "helper",
			},
			{
				Message:  "[БОЙ] Вы прибыли к каравану. На них нападают бандиты!",
				Variants: [2]string{"Вступить в бой", "Попытаться договориться"},
				Actions: [2]func(p *player.Player) string{
					func(p *player.Player) string {
						err := p.ChangeHealth(-15)
						if err != nil {
							return "[СМЕРТЬ] Вы погибли в бою с бандитами..."
						}
						return "[БОЙ] Вместе с караваном вы отбили атаку! (-15 здоровья)"
					},
					func(p *player.Player) string {
						if p.Eat >= 3 {
							p.AddEat(-3)
							p.Flags["chain_signal_diplomacy"] = true
							return "[ДИПЛОМАТИЯ] Вы откупились частью припасов. Бандиты ушли. (Еда: -3)"
						}
						p.ChangeHealth(-10)
						return "[НЕУДАЧА] Торговаться нечем. Бандиты атаковали. (-10 здоровья)"
					},
				},
				Category: "combat",
			},
		},
		FinalReward: func(p *player.Player) string {
			p.AddEat(4)
			p.AddWater(4)
			return "[БЛАГОДАРНОСТЬ] Караванщики делятся припасами! (Еда: +4, Вода: +4)"
		},
	},

	// Цепочка: Карта → Сокровище
	{
		ID:          "chain_map",
		TriggerFlag: "found_map",
		Steps: []events.Event{
			{
				Message:  "[ПОИСК] Вы пытаетесь разобрать карту. Нужно найти ориентир.",
				Variants: [2]string{"Искать по компасу", "Искать по приметным местам"},
				Actions: [2]func(p *player.Player) string{
					func(p *player.Player) string {
						risk := rand.Intn(100)
						if risk < 30 {
							p.ChangeHealth(-10)
							return "[НЕУДАЧА] Вы заблудились и напоролись на терновник. (-10 здоровья)"
						}
						p.Flags["chain_map_orienteer"] = true
						return "[УСПЕХ] Компас вывел вас к цели!"
					},
					func(p *player.Player) string {
						p.AddEat(-1)
						if p.Eat < 0 {
							p.Eat = 0
						}
						return "[ПОИСК] Вы потратили день на поиски, но нашли верный путь. (Еда: -1)"
					},
				},
				Category: "find",
			},
			{
				Message:  "[СОКРОВИЩЕ] Вы нашли тайник, отмеченный на карте!",
				Variants: [2]string{"Открыть тайник", "Осмотреть на ловушки"},
				Actions: [2]func(p *player.Player) string{
					func(p *player.Player) string {
						p.AddEat(7)
						p.AddWater(7)
						p.ChangeHealth(10)
						return "[СОКРОВИЩЕ] Припасы и медикаменты! (Еда: +7, Вода: +7, Здоровье: +10)"
					},
					func(p *player.Player) string {
						p.AddEat(4)
						p.AddWater(4)
						return "[ОСТОРОЖНО] Тайник заминирован, но вы успели обезвредить. (Еда: +4, Вода: +4)"
					},
				},
				Category: "find",
			},
		},
		FinalReward: func(p *player.Player) string {
			p.Flags["treasure_found"] = true
			return "[СОКРОВИЩЕ] Вы нашли легендарный тайник выживших!"
		},
	},
}
