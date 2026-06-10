package save

import (
	"bunker3000/player"
	"encoding/json"
	"os"
	"path/filepath"
)

const SaveFileName = "bunker3000_save.json"

type SaveData struct {
	Player       player.Player `json:"player"`
	Achievements []SaveAchieve `json:"achievements"`
}

type SaveAchieve struct {
	ID       string `json:"id"`
	Unlocked bool   `json:"unlocked"`
	TotalDmg int8   `json:"total_dmg"`
}

func Save(p player.Player, achieves []struct {
	ID       string
	Unlocked bool
}, totalDmg int8) error {
	data := SaveData{
		Player: p,
	}
	for _, a := range achieves {
		data.Achievements = append(data.Achievements, SaveAchieve{
			ID:       a.ID,
			Unlocked: a.Unlocked,
			TotalDmg: totalDmg,
		})
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(SaveFileName, bytes, 0644)
}

func Load() (*SaveData, error) {
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		path := filepath.Join(dir, SaveFileName)
		if _, err := os.Stat(path); err == nil {
			bytes, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			var data SaveData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			return &data, nil
		}
	}

	bytes, err := os.ReadFile(SaveFileName)
	if err != nil {
		return nil, err
	}
	var data SaveData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func HasSave() bool {
	if _, err := os.Stat(SaveFileName); err == nil {
		return true
	}
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	dir := filepath.Dir(exe)
	path := filepath.Join(dir, SaveFileName)
	_, err = os.Stat(path)
	return err == nil
}

func DeleteSave() error {
	if err := os.Remove(SaveFileName); err != nil && !os.IsNotExist(err) {
		exe, err2 := os.Executable()
		if err2 == nil {
			dir := filepath.Dir(exe)
			path := filepath.Join(dir, SaveFileName)
			return os.Remove(path)
		}
		return err
	}
	return nil
}
