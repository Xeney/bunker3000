package player

import (
	"testing"
)

func BenchmarkStartNewDay(b *testing.B) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.StartNewDay()
	}
}

func BenchmarkChangeHealth(b *testing.B) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.ChangeHealth(-5)
		p.ChangeHealth(5)
	}
}

func BenchmarkAddResource(b *testing.B) {
	p := CreatePlayer(DifficultyNormal, ClassSurvivor)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.AddEat(1)
		p.AddWater(1)
	}
}
