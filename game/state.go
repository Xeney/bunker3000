package game

import (
	"bunker3000/achievements"
	"bunker3000/events"
	"bunker3000/player"
	"bunker3000/save"
	"errors"
	"math/rand"
	"time"
)

type GamePhase int

const (
	PhaseMenu GamePhase = iota
	PhaseDifficulty
	PhaseClass
	PhasePlaying
	PhaseResult
	PhaseGameOver
	PhaseAchievements
)

type GameState struct {
	Player        player.Player
	Achievements  *achievements.Tracker
	EventPool     []events.Event
	CurrentEvent  *events.Event
	LastResult    string
	Phase         GamePhase
	TotalDamage   int8
	DayMessage    string
	DayError      error
	Endless       bool
	StoryBlocks   []string
	ActiveChainID string
	ChainStep     int
	UsedEvents    map[int]bool
	LeveledUp     bool
}

func NewGame(diff player.Difficulty, class player.ClassType) *GameState {
	rand.Seed(time.Now().UnixNano())
	p := player.CreatePlayer(diff, class)
	eventPool := events.ConstructEventsPool()
	ach := achievements.NewTracker()

	return &GameState{
		Player:       p,
		Achievements: ach,
		EventPool:    eventPool,
		Phase:        PhasePlaying,
		Endless:      diff == player.DifficultyEndless,
		UsedEvents:   make(map[int]bool),
	}
}

func (gs *GameState) StartDay() {
	gs.DayError = nil
	gs.DayMessage = ""
	gs.CurrentEvent = nil
	gs.LastResult = ""
	gs.LeveledUp = false

	err := gs.Player.StartNewDay()

	if gs.Player.Lock && gs.Player.Health > 0 {
		gs.Player.AddXP(10)
		gs.DayMessage = "ПОБЕДА! Вы успешно продержались все дни!"
		gs.Achievements.Check(&gs.Player)
		return
	}

	if err != nil {
		gs.DayError = err
		return
	}

	if gs.Player.AddXP(5) {
		gs.LeveledUp = true
	}

	// Active chain step
	if gs.ActiveChainID != "" {
		chain := findChain(gs.ActiveChainID)
		if chain != nil && gs.ChainStep < len(chain.Steps) {
			gs.CurrentEvent = &chain.Steps[gs.ChainStep]
			return
		}
		gs.ActiveChainID = ""
		gs.ChainStep = 0
	}

	// Check flag triggers for new chains
	if gs.ActiveChainID == "" {
		for i := range AllChains {
			chain := &AllChains[i]
			if gs.Player.Flags[chain.TriggerFlag] && !gs.Player.Flags[chain.ID+"_started"] {
				activateChain(gs, chain.ID)
				gs.CurrentEvent = &chain.Steps[0]
				return
			}
		}
	}

	ev, err := gs.pickEvent()
	if err != nil {
		gs.DayError = err
		return
	}
	gs.CurrentEvent = ev
}

func (gs *GameState) ExecuteChoice(choice int) (string, error) {
	if gs.CurrentEvent == nil {
		return "", errors.New("нет активного события")
	}

	resultMessage, err := gs.CurrentEvent.Execute(choice, &gs.Player)
	gs.UsedEvents[gs.CurrentEvent.ID] = true
	gs.LastResult = resultMessage

	// Trader class bonus
	if gs.Player.Class == player.ClassTrader &&
		(gs.CurrentEvent.Category == "helper") {
		gs.Player.AddEat(2)
		gs.Player.AddWater(2)
		resultMessage += "\n[ТОРГОВЕЦ] Бонус от сделки: еда +2, вода +2"
		gs.LastResult = resultMessage
	}

	// XP for karma-positive choice
	if gs.Player.Karma > 0 {
		if gs.Player.AddXP(3) {
			gs.LeveledUp = true
		}
	}

	// Advance active chain
	if gs.ActiveChainID != "" {
		gs.ChainStep++
		chain := findChain(gs.ActiveChainID)
		if chain != nil && gs.ChainStep >= len(chain.Steps) {
			gs.Player.AddXP(15)
			if chain.FinalReward != nil {
				rewardMsg := chain.FinalReward(&gs.Player)
				resultMessage = rewardMsg + "\n" + resultMessage
				gs.LastResult = resultMessage
			}
			gs.StoryBlocks = append(gs.StoryBlocks, gs.ActiveChainID)
			gs.Player.Flags[gs.ActiveChainID+"_done"] = true
			gs.ActiveChainID = ""
			gs.ChainStep = 0
		}
	}

	gs.Achievements.Check(&gs.Player)

	if gs.Player.Lock {
		gs.TotalDamage = 100 - gs.Player.Health
		if gs.TotalDamage < 0 {
			gs.TotalDamage = 0
		}
		gs.Achievements.RecordDamage(gs.TotalDamage)
		gs.Achievements.Check(&gs.Player)
	}

	return resultMessage, err
}

func (gs *GameState) pickEvent() (*events.Event, error) {
	available := make([]events.Event, 0, len(gs.EventPool))
	for _, e := range gs.EventPool {
		if !gs.UsedEvents[e.ID] {
			available = append(available, e)
		}
	}

	if len(available) == 0 {
		gs.UsedEvents = make(map[int]bool)
		available = gs.EventPool
	}

	k := gs.Player.Karma
	type weightEntry struct {
		ev     events.Event
		weight int
	}
	entries := make([]weightEntry, 0, len(available))

	for _, e := range available {
		w := 10
		switch {
		case k > 30 && (e.Category == "resource" || e.Category == "helper" || e.Category == "find"):
			w += 15
		case k < -30 && (e.Category == "combat" || e.Category == "hazard"):
			w += 15
		case k > 30 && (e.Category == "combat" || e.Category == "hazard"):
			w -= 5
		case k < -30 && (e.Category == "resource" || e.Category == "helper" || e.Category == "find"):
			w -= 5
		}
		if w < 1 {
			w = 1
		}
		entries = append(entries, weightEntry{e, w})
	}

	total := 0
	for _, e := range entries {
		total += e.weight
	}
	roll := rand.Intn(total)
	cum := 0
	for _, e := range entries {
		cum += e.weight
		if roll < cum {
			return &e.ev, nil
		}
	}
	return &entries[len(entries)-1].ev, nil
}

func (gs *GameState) NextDay() {
	gs.Player.ThisDay++
	gs.Achievements.Check(&gs.Player)
}

func (gs *GameState) IsGameOver() bool {
	return gs.Player.Lock
}

func (gs *GameState) Save() error {
	// Load old achievements for persistence between runs
	oldAchieves, _ := save.LoadAchievements()

	// Merge: keep any previously unlocked achievement
	unlocked := make(map[string]bool)
	for _, sa := range oldAchieves {
		if sa.Unlocked {
			unlocked[sa.ID] = true
		}
	}
	for _, a := range gs.Achievements.Achievements {
		if a.Unlocked {
			unlocked[a.ID] = true
		}
	}

	// Build merged save list
	saveAchieves := make([]struct {
		ID       string
		Unlocked bool
	}, len(gs.Achievements.Achievements))
	persistAchieves := make([]save.SaveAchieve, len(gs.Achievements.Achievements))
	for i, a := range gs.Achievements.Achievements {
		state := unlocked[a.ID]
		saveAchieves[i] = struct {
			ID       string
			Unlocked bool
		}{ID: a.ID, Unlocked: state}
		persistAchieves[i] = save.SaveAchieve{
			ID:       a.ID,
			Unlocked: state,
			TotalDmg: gs.TotalDamage,
		}
	}

	// Persist achievements separately (never deleted)
	save.SaveAchievements(persistAchieves)

	// Save game state
	return save.Save(gs.Player, saveAchieves, gs.TotalDamage, gs.Endless, gs.StoryBlocks, gs.ActiveChainID, gs.ChainStep)
}

func LoadFromSave(saveData *save.SaveData) *GameState {
	eventPool := events.ConstructEventsPool()
	ach := achievements.NewTracker()
	for _, sa := range saveData.Achievements {
		for _, a := range ach.Achievements {
			if a.ID == sa.ID && sa.Unlocked {
				a.Unlocked = true
			}
		}
	}
	return &GameState{
		Player:        saveData.Player,
		Achievements:  ach,
		EventPool:     eventPool,
		Phase:         PhasePlaying,
		TotalDamage:   saveData.TotalDamage,
		Endless:       saveData.Endless,
		StoryBlocks:   saveData.StoryBlocks,
		ActiveChainID: saveData.ActiveChainID,
		ChainStep:     saveData.ChainStep,
	}
}
