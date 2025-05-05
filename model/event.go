package model

import (
	"biathlon-system-prototype/config"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	Registered = iota + 1
	StartTimeWasSet
	OnStarLine
	Started
	OnFiringRange
	TargetHit
	LeftFiringRange
	EnteredPenaltyLaps
	LeftPenaltyLaps
	EndedMainLap
	CanNotContinue
)

const (
	Disqualified = 32
	Finished     = 33
)

type Event struct {
	Time         time.Time // время, когда произошло событие
	Id           int       // id события
	CompetitorID int       // id участника
	ExtraParams  []string  // дополнительные параметры
}

// ParseEvent - преобразует строку в объект структуры Event
func ParseEvent(str string) (Event, error) {
	if !strings.HasPrefix(str, "[") {
		return Event{}, fmt.Errorf("некорректный формат строки: %s", str)
	}

	timeEnd := strings.Index(str, "]")
	if timeEnd == -1 {
		return Event{}, fmt.Errorf("пропущена закрывающаяся скобка ] в строке: %s", str)
	}

	parts := strings.Fields(str)

	if len(parts) < 3 {
		return Event{}, fmt.Errorf("некорректный формат строки события: %s", parts)
	}

	timeStr := strings.Trim(parts[0], "[]")

	eventTime, err := time.Parse(config.TimeLayout, timeStr)
	if err != nil {
		return Event{}, fmt.Errorf("некорректный формат времени: %v", err)
	}

	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return Event{}, fmt.Errorf("некорректный формат id события: %v", err)
	}

	competitorID, err := strconv.Atoi(parts[2])
	if err != nil {
		return Event{}, fmt.Errorf("некорректный формат id участника: %v", err)
	}

	var params []string
	if len(parts) > 3 {
		params = parts[3:]
	}

	return Event{
		Time:         eventTime,
		Id:           eventID,
		CompetitorID: competitorID,
		ExtraParams:  params,
	}, nil
}
