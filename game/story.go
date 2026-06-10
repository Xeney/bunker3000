package game

import (
	"bunker3000/player"
	"fmt"
	"sort"
	"strings"
)

func GenerateEpilogue(gs *GameState) string {
	var b strings.Builder
	p := &gs.Player

	// 1. Opening
	if p.Health > 0 {
		if gs.Endless {
			b.WriteString("Пустоши не смогли вас сломить. Вы продержались до самого конца.\n\n")
		} else {
			b.WriteString("Вы выжили. Пустоши не сломили вас.\n\n")
		}
	} else {
		b.WriteString("Пустоши забрали ещё одну жизнь...\n\n")
		switch {
		case p.Eat <= 0 && p.Water <= 0:
			b.WriteString("Голод и жажда — самая страшная смерть в пустошах. Ресурсы кончились, и организм не выдержал.\n\n")
		case p.Eat <= 0:
			b.WriteString("Голод оказался сильнее воли к жизни. Еда кончилась, и силы покинули вас.\n\n")
		case p.Water <= 0:
			b.WriteString("Обезвоживание забрало последние силы. Вода — главное сокровище пустоши.\n\n")
		default:
			b.WriteString("Раны, полученные в пути, оказались смертельными. Пустоши не прощают ошибок.\n\n")
		}
	}

	// 2. Key decisions from flags
	storyMap := map[string]string{
		"helped_cries":  "Вы не прошли мимо чужих криков — и обрели благодарного союзника.",
		"ignored_cries": "Чужие крики остались без ответа. Пустоши не прощают равнодушия.",
		"killed_bandit": "Раненый бандит пал от вашей руки. Кто-то в пустошах запомнит это.",
		"spared_bandit": "Вы пощадили раненого бандита. Милосердие — редкая роскошь в эти дни.",
		"fed_oldman":    "Последняя еда, отданная незнакомому старику, не была забыта.",
		"refused_oldman": "Старик ушёл в пустоши с пустыми руками. Вы запомнили его взгляд.",
		"found_map":     "Карта, найденная в руинах, вела к сокровищу. Но хватило ли смелости пойти по ней?",
		"broke_storage": "Вы взломали довоенный склад. Шум привлёк внимание — но добыча стоила риска.",
		"left_storage":  "Вы не стали рисковать со складом. Осторожность сохранила вам жизнь.",
		"answered_signal": "Сигнал через рацию привёл к контакту с другими выжившими.",
		"chain_cries_done":  "Спасённый из руин указал путь к тайнику. Его благодарность была искренней.",
		"chain_signal_done": "Караван, с которым вы связались по рации, стал надёжным союзником.",
		"chain_map_done":    "Тайник, отмеченный на карте, оказался настоящим сокровищем.",
		"treasure_found":    "Легендарный тайник выживших найден. История запомнит ваше имя.",
	}

	// Collect matching flag keys sorted for deterministic order
	keys := make([]string, 0, len(p.Flags))
	for k := range p.Flags {
		if p.Flags[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		if sentence, ok := storyMap[k]; ok {
			b.WriteString(sentence + "\n")
		}
	}

	// 3. Stats
	b.WriteString(fmt.Sprintf("\nВы продержались %d дней.\n", p.ThisDay))
	if gs.TotalDamage > 0 {
		b.WriteString(fmt.Sprintf("За это время вы получили %d%% суммарного урона.\n", gs.TotalDamage))
	}
	b.WriteString(fmt.Sprintf("Максимальный уровень здоровья: %d%%.\n", player.MaxHealth(p.Class)))
	b.WriteString(fmt.Sprintf("Минимальный: вы были на волосок от смерти.\n"))

	// Karma block
	switch {
	case p.Karma >= 50:
		b.WriteString("\nВаша карма сияет чистотой. Вы были светом в этом тёмном мире.\n")
	case p.Karma >= 20:
		b.WriteString("\nВаша карма склоняется к добру. Вы старались помогать другим.\n")
	case p.Karma >= -20:
		b.WriteString("\nВаша карма в равновесии. Вы просто выживали, как могли.\n")
	case p.Karma >= -50:
		b.WriteString("\nВаша карма омрачена. Пустоши сделали вас жестоким.\n")
	default:
		b.WriteString("\nВаша карма чернее бездны. Вы сеяли смерть и хаос.\n")
	}

	// 5. Ending
	b.WriteString(fmt.Sprintf("\n%s на сложности «%s» — ", p.Class.String(), p.Difficulty.String()))
	if p.Health > 0 {
		b.WriteString("достойная история выживания в мире, где выживают единицы.")
	} else {
		b.WriteString("трагическая, но поучительная история. Пустоши не прощают слабости.")
	}

	return b.String()
}
