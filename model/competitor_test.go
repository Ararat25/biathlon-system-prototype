package model

import (
	"biathlon-system-prototype/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockConfig возвращает тестовую конфигурацию соревнования
func mockConfig() *config.Config {
	return &config.Config{
		Laps:        2,
		LapLen:      1000,
		PenaltyLen:  200,
		FiringLines: 1,
		Start:       "12:00:00",
		StartDelta:  "00:00:30",
	}
}

// TestBuildCompetitors_Finished проверяет, что участник, завершивший все круги, получает статус Finished
func TestBuildCompetitors_Finished(t *testing.T) {
	conf := mockConfig()
	startPlanned, _ := time.Parse(config.StartTimeLayout, "12:00:00")
	startActual := startPlanned
	lap1End := startPlanned.Add(3 * time.Minute)
	lap2End := lap1End.Add(3 * time.Minute)

	regTime, _ := time.Parse(config.TimeLayout, "11:59:00.000")

	events := map[int][]*Event{
		1: {
			{Id: Registered, Time: regTime},
			{Id: StartTimeWasSet, ExtraParams: []string{"12:00:00.000"}, Time: regTime},
			{Id: Started, Time: startActual},
			{Id: EndedMainLap, Time: lap1End},
			{Id: EndedMainLap, Time: lap2End},
		},
	}

	competitors := BuildCompetitors(events, conf)
	comp := (*competitors)[0]

	assert.Equal(t, 1, comp.Id)
	assert.Equal(t, Finished, comp.Status)
	assert.Len(t, comp.Laps, 2)
	assert.WithinDuration(t, lap2End, comp.FinishedTime, time.Second)
	assert.Equal(t, 2, len(comp.Laps))
}

// TestBuildCompetitors_Disqualified проверяет, что участник, стартовавший с нарушением времени, получает статус Disqualified.
func TestBuildCompetitors_Disqualified(t *testing.T) {
	conf := mockConfig()
	startPlanned, _ := time.Parse(config.StartTimeLayout, "12:00:00")
	startActual := startPlanned.Add(2 * time.Minute) // больше StartDelta

	events := map[int][]*Event{
		2: {
			{Id: Registered},
			{Id: StartTimeWasSet, ExtraParams: []string{"12:00:00.000"}},
			{Id: Started, Time: startActual},
		},
	}

	competitors := BuildCompetitors(events, conf)
	comp := (*competitors)[0]

	assert.Equal(t, Disqualified, comp.Status)
}

// TestBuildCompetitors_CannotContinue проверяет, что участник, не сумевший продолжить гонку, получает статус CanNotContinue и сообщение о причине.
func TestBuildCompetitors_CannotContinue(t *testing.T) {
	conf := mockConfig()
	startPlanned, _ := time.Parse(config.StartTimeLayout, "12:00:00")
	startActual := startPlanned

	eventTime, _ := time.Parse(config.TimeLayout, "11:00:00.000")

	events := map[int][]*Event{
		3: {
			{Id: Registered, Time: eventTime},
			{Id: StartTimeWasSet, Time: eventTime.Add(10 * time.Minute), ExtraParams: []string{"12:00:00.000"}},
			{Id: Started, Time: startActual},
			{Id: CanNotContinue, Time: startActual.Add(30 * time.Minute), ExtraParams: []string{"упал, сломал лыжи"}},
		},
	}

	competitors := BuildCompetitors(events, conf)
	comp := (*competitors)[0]

	assert.Equal(t, CanNotContinue, comp.Status)
	assert.Contains(t, comp.NotFinishedMsg, "упал, сломал лыжи")
}
