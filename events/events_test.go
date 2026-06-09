package events

import (
	"bunker3000/player"
	"testing"
)

func TestConstructEvent_Valid(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error { return nil },
			func(p *player.Player) error { return nil },
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
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error { return nil },
			func(p *player.Player) error { return nil },
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
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error { return nil },
			func(p *player.Player) error { return nil },
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
		Actions:  [2]func(p *player.Player) error{nil, nil},
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
	if len(pool) < 2 {
		t.Errorf("Expected at least 2 events, got %d", len(pool))
	}
}

func TestEventExecute_ValidChoice(t *testing.T) {
	executed := false
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error { executed = true; return nil },
			func(p *player.Player) error { return nil },
		},
	}
	p := player.CreatePlayer()

	err := event.Execute(1, &p)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !executed {
		t.Error("Action was not executed")
	}
}

func TestEventExecute_InvalidChoice(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions: [2]func(p *player.Player) error{
			func(p *player.Player) error { return nil },
			func(p *player.Player) error { return nil },
		},
	}
	p := player.CreatePlayer()

	err := event.Execute(3, &p)

	if err == nil {
		t.Error("Expected error for invalid choice, got nil")
	}
}

func TestEventExecute_NilAction(t *testing.T) {
	event := Event{
		Message:  "Test event",
		Variants: [2]string{"Option 1", "Option 2"},
		Actions:  [2]func(p *player.Player) error{nil, nil},
	}
	p := player.CreatePlayer()

	err := event.Execute(1, &p)

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

// Интеграционный тест для конкретных событий
func TestEvent_AddWaterEvent(t *testing.T) {
	pool := ConstructEventsPool()
	p := player.CreatePlayer()
	initialWater := p.Water

	// Находим событие с водой и выполняем его
	for _, event := range pool {
		if event.Message == "Вы нашли чистый ручей среди скал." {
			err := event.Execute(1, &p)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if p.Water <= initialWater {
				t.Errorf("Expected water to increase from %d, got %d", initialWater, p.Water)
			}
			return
		}
	}
	t.Error("Water event not found in pool")
}

func TestEvent_AttackEvent(t *testing.T) {
	pool := ConstructEventsPool()
	p := player.CreatePlayer()
	p.Health = 80

	// Находим событие с атакой и выполняем его
	for _, event := range pool {
		if event.Message == "На вас напал бродяга с ножом!" {
			err := event.Execute(1, &p)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if p.Health >= 80 {
				t.Errorf("Expected health to decrease from 80, got %d", p.Health)
			}
			return
		}
	}
	t.Error("Attack event not found in pool")
}
