package events

import (
	"bunker3000/player"
	"testing"
)

func TestConstructEvent_Valid(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string { return "OK" },
			func(p *player.Player) string { return "OK" },
		},
	}

	result, err := ConstructEvent(event)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Message != event.Message {
		t.Errorf("Expected message %s, got %s", event.Message, result.Message)
	}
}

func TestConstructEvent_EmptyMessage(t *testing.T) {
	event := Event{
		Message:  "",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string { return "OK" },
			func(p *player.Player) string { return "OK" },
		},
	}

	_, err := ConstructEvent(event)

	if err == nil {
		t.Error("Expected error for empty message, got nil")
	}
}

func TestConstructEvent_EmptyVariant(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"", "Option 2"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string { return "OK" },
			func(p *player.Player) string { return "OK" },
		},
	}

	_, err := ConstructEvent(event)

	if err == nil {
		t.Error("Expected error for empty variant, got nil")
	}
}

func TestConstructEvent_NilAction(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions:  [2]func(p *player.Player) string{nil, nil},
	}

	_, err := ConstructEvent(event)

	if err == nil {
		t.Error("Expected error for nil action, got nil")
	}
}

func TestConstructEventsPool(t *testing.T) {
	pool := ConstructEventsPool()

	if len(pool) == 0 {
		t.Error("Expected non-empty event pool")
	}
	if len(pool) < 30 {
		t.Errorf("Expected at least 30 events, got %d", len(pool))
	}
}

func TestEventExecute_ValidChoice(t *testing.T) {
	executed := false
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string { executed = true; return "OK" },
			func(p *player.Player) string { return "OK" },
		},
	}
	p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)

	resultMessage, err := event.Execute(1, &p)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !executed {
		t.Error("Action was not executed")
	}
	if resultMessage != "OK" {
		t.Errorf("Expected message 'OK', got '%s'", resultMessage)
	}
}

func TestEventExecute_InvalidChoice(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions: [2]func(p *player.Player) string{
			func(p *player.Player) string { return "OK" },
			func(p *player.Player) string { return "OK" },
		},
	}
	p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)

	_, err := event.Execute(3, &p)

	if err == nil {
		t.Error("Expected error for invalid choice, got nil")
	}
}

func TestEventExecute_NilAction(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions:  [2]func(p *player.Player) string{nil, nil},
	}
	p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)

	_, err := event.Execute(1, &p)

	if err == nil {
		t.Error("Expected error for nil action, got nil")
	}
}

func TestGetRandomEvent(t *testing.T) {
	pool := []Event{
		{Message: "Event 1"},
		{Message: "Event 2"},
		{Message: "Event 3"},
	}

	event, err := GetRandomEvent(pool)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if event.Message == "" {
		t.Error("Expected non-empty event")
	}
}

func TestGetRandomEvent_EmptyPool(t *testing.T) {
	pool := []Event{}

	_, err := GetRandomEvent(pool)

	if err == nil {
		t.Error("Expected error for empty pool, got nil")
	}
}

func TestEvent_AddWaterEvent(t *testing.T) {
	pool := ConstructEventsPool()
	p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)
	initialWater := p.Water

	for _, event := range pool {
		if event.Message == "[РУЧЕЙ] Среди серых скал вы слышите журчание воды. Небольшой ручей с кристально чистой водой пробивается сквозь расщелину." {
			resultMessage, err := event.Execute(1, &p)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if p.Water <= initialWater {
				t.Errorf("Expected water to increase from %d, got %d", initialWater, p.Water)
			}
			if resultMessage == "" {
				t.Error("Expected non-empty result message")
			}
			return
		}
	}
	t.Error("Water event not found in pool")
}

func TestEvent_AttackEvent(t *testing.T) {
	pool := ConstructEventsPool()
	p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)
	p.Health = 80

	for _, event := range pool {
		if event.Message == "[НАПАДЕНИЕ] Из-за ржавого остова грузовика выскакивает человек в грязной одежде. В его руке сверкает заточенный нож, глаза горят безумной решимостью." {
			resultMessage, err := event.Execute(1, &p)
			if err != nil && p.Health > 0 {
				if p.Health > 0 {
					t.Errorf("Unexpected error: %v", err)
				}
			}
			if p.Health >= 80 && p.Health > 0 {
				t.Errorf("Expected health to decrease from 80, got %d", p.Health)
			}
			if resultMessage == "" {
				t.Error("Expected non-empty result message")
			}
			return
		}
	}
	t.Error("Attack event not found in pool")
}

func TestAllEvents_NoPanic(t *testing.T) {
	pool := ConstructEventsPool()

	for i, event := range pool {
		p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)

		for choice := 1; choice <= 2; choice++ {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Event %d, choice %d panicked: %v", i, choice, r)
					}
				}()

				resultMessage, err := event.Execute(choice, &p)
				_ = resultMessage
				_ = err
			}()
		}
	}
}

func TestEventsCount(t *testing.T) {
	pool := ConstructEventsPool()
	expectedCount := 238

	if len(pool) != expectedCount {
		t.Errorf("Expected %d events, got %d", expectedCount, len(pool))
	}
}

func TestUniqueEventMessages(t *testing.T) {
	pool := ConstructEventsPool()
	messages := make(map[string]bool)

	for _, event := range pool {
		if messages[event.Message] {
			t.Errorf("Duplicate event message: %s", event.Message)
		}
		messages[event.Message] = true
	}
}

func TestAllEventsExecution(t *testing.T) {
	pool := ConstructEventsPool()

	for i, event := range pool {
		for choice := 1; choice <= 2; choice++ {
			p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)
			p.Health = 100
			p.Eat = 10
			p.Water = 10
			cfg := p.GetCfg()

			resultMessage, err := event.Execute(choice, &p)

			if err != nil {
				if p.Health > 0 {
					t.Errorf("Event %d, choice %d unexpected error: %v", i, choice, err)
				}
			}

			if err == nil && resultMessage == "" {
				t.Errorf("Event %d, choice %d returned empty message", i, choice)
			}

			if p.Health < 0 || p.Health > player.MaxHealth(p.Class) {
				t.Errorf("Event %d, choice %d health out of bounds: %d", i, choice, p.Health)
			}
			if p.Eat < 0 || p.Eat > cfg.MaxResource {
				t.Errorf("Event %d, choice %d eat out of bounds: %d", i, choice, p.Eat)
			}
			if p.Water < 0 || p.Water > cfg.MaxResource {
				t.Errorf("Event %d, choice %d water out of bounds: %d", i, choice, p.Water)
			}
		}
	}
}

func TestAllEventsReturnMessage(t *testing.T) {
	pool := ConstructEventsPool()

	for i, event := range pool {
		for choice := 1; choice <= 2; choice++ {
			p := player.CreatePlayer(player.DifficultyNormal, player.ClassSurvivor)

			resultMessage, err := event.Execute(choice, &p)

			if err == nil && resultMessage == "" {
				t.Errorf("Event %d, choice %d returned empty message but no error", i, choice)
			}

			if err != nil && resultMessage == "" {
				t.Errorf("Event %d, choice %d returned error but empty message: %v", i, choice, err)
			}
		}
	}
}
