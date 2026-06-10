package events

import (
	"bunker3000/player"
	"math/rand"
)

func ev2(msg, v1, v2 string, a1, a2 func(p *player.Player) string, cat string) Event {
	var e Event
	e.Message = msg
	e.Variants[0] = v1
	e.Variants[1] = v2
	e.Actions[0] = a1
	e.Actions[1] = a2
	e.Category = cat
	return e
}

func R(msg, v1, v2 string, eat1, water1, eat2, water2 int8, cat string) Event {
	return ev2(msg, v1, v2,
		func(p *player.Player) string {
			if eat1 > 0 {
				p.AddEat(eat1)
			}
			if water1 > 0 {
				p.AddWater(water1)
			}
			return msg + "\n[НАХОДКА] Еда: +" + itoa(eat1) + ", Вода: +" + itoa(water1)
		},
		func(p *player.Player) string {
			if eat2 > 0 {
				p.AddEat(eat2)
			}
			if water2 > 0 {
				p.AddWater(water2)
			}
			return msg + "\n[ВЫБОР] Еда: +" + itoa(eat2) + ", Вода: +" + itoa(water2)
		},
		cat)
}

func RN(msg, v1, v2 string, eat, water int8, eatRisk, waterRisk int8, riskPct int, riskDmg int8, cat string) Event {
	return ev2(msg, v1, v2,
		func(p *player.Player) string {
			if rand.Intn(100) < riskPct {
				p.ChangeHealth(-riskDmg)
				return msg + "\n[РИСК] Неудача! Потеря здоровья: -" + itoa(riskDmg)
			}
			p.AddEat(eat)
			p.AddWater(water)
			return msg + "\n[УДАЧА] Еда: +" + itoa(eat) + ", Вода: +" + itoa(water)
		},
		func(p *player.Player) string {
			if eatRisk > 0 {
				p.AddEat(eatRisk)
			}
			if waterRisk > 0 {
				p.AddWater(waterRisk)
			}
			return msg + "\n[ОСТОРОЖНО] Еда: +" + itoa(eatRisk) + ", Вода: +" + itoa(waterRisk)
		},
		cat)
}

func CK(msg, v1, v2 string, dmg1, dmg2 int8, k1, k2 int8, cat string) Event {
	return ev2(msg, v1, v2,
		func(p *player.Player) string {
			p.AddKarma(k1)
			p.ChangeHealth(-dmg1)
			return msg + "\n[БОЙ] Здоровье: -" + itoa(dmg1) + "\n[КАРМА " + sign(k1) + itoa(k1) + "]"
		},
		func(p *player.Player) string {
			p.AddKarma(k2)
			p.ChangeHealth(-dmg2)
			return msg + "\n[ВЫБОР] Здоровье: -" + itoa(dmg2) + "\n[КАРМА " + sign(k2) + itoa(k2) + "]"
		},
		cat)
}

func MK(msg, v1, v2, goodMsg, badMsg string, k1, k2 int8, eat1, water1 int8, cat string) Event {
	return ev2(msg, v1, v2,
		func(p *player.Player) string {
			p.AddKarma(k1)
			if eat1 > 0 {
				p.AddEat(eat1)
			}
			if water1 > 0 {
				p.AddWater(water1)
			}
			return msg + "\n" + goodMsg + "\n[КАРМА " + sign(k1) + itoa(k1) + "]"
		},
		func(p *player.Player) string {
			p.AddKarma(k2)
			return msg + "\n" + badMsg + "\n[КАРМА " + sign(k2) + itoa(k2) + "]"
		},
		cat)
}

func W(msg, v1, v2 string, dmg1, dmg2 int8, cat string) Event {
	return ev2(msg, v1, v2,
		func(p *player.Player) string {
			p.ChangeHealth(-dmg1)
			return msg + "\n[УКРЫТИЕ] Здоровье: -" + itoa(dmg1)
		},
		func(p *player.Player) string {
			p.ChangeHealth(-dmg2)
			return msg + "\n[РИСК] Здоровье: -" + itoa(dmg2)
		},
		cat)
}

func sign(v int8) string {
	if v >= 0 {
		return "+"
	}
	return ""
}

func itoa(v int8) string {
	if v == 0 {
		return "0"
	}
	n := int(v)
	if n < 0 {
		n = -n
	}
	if n == 0 {
		return "0"
	}
	var buf [4]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
