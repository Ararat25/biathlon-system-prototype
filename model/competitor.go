package model

import (
	"biathlon-system-prototype/config"
	"fmt"
	"strings"
	"time"
)

const numberTargets = 5

type Competitor struct {
	Id             int       // id участника
	Registered     bool      // флаг для указания зарегистрированности участника
	Finished       bool      // для отслеживания статуса участника
	StartPlanned   time.Time // планируемое время старта
	StartActual    time.Time // реальное время старта
	Laps           []LapInfo // информация о пройденных основных кругах
	PenaltyLaps    LapInfo   // информация о пройденных штрафных кругах
	Hits           int       // количество попаданий
	NotFinishedMsg string    // причина из-за которой участник не может продолжить участие
}

func BuildCompetitors(events map[int][]Event, conf *config.Config) []Competitor {
	var competitors []Competitor

	for competitorId, eventList := range events {
		competitor := Competitor{Id: competitorId}

		var (
			lastPenaltyLapEnter time.Time
			lastMainLapEnd      time.Time
			numberPenaltyLaps   int
		)
	loop:
		for _, event := range eventList {
			switch event.Id {
			case Registered:
				competitor.Registered = true
				break
			case StartTimeWasSet:
				parseTime, err := time.Parse(config.TimeLayout, event.ExtraParams[0])
				if err != nil {
					fmt.Printf("некорректный формат времени старта для участника %d: %v\n", competitorId, err)
					continue loop
				}

				startDelta, _ := time.Parse(config.StartTimeLayout, conf.StartDelta)
				delta := startDelta.Sub(time.Time{})

				competitor.StartPlanned = parseTime.Add(delta)
				break
			case Started:
				competitor.StartActual = event.Time
				lastMainLapEnd = event.Time
				break
			case TargetHit:
				competitor.Hits++
				break
			case EnteredPenaltyLaps:
				lastPenaltyLapEnter = event.Time
				break
			case LeftPenaltyLaps:
				numberPenaltyLaps++
				distance := numberPenaltyLaps * conf.PenaltyLen

				competitor.PenaltyLaps.Duration += event.Time.Sub(lastPenaltyLapEnter)

				speed := float64(distance) / competitor.PenaltyLaps.Duration.Seconds()

				competitor.PenaltyLaps.Speed = speed
				break
			case EndedMainLap:
				duration := event.Time.Sub(lastMainLapEnd)
				competitor.Laps = append(competitor.Laps, LapInfo{Duration: duration, Speed: float64(conf.LapLen) / duration.Seconds()})
				lastMainLapEnd = event.Time

				if len(competitor.Laps) == conf.Laps {
					event.Id = Finished
				}
				break
			case Finished:
				competitor.Finished = true
				break
			case CanNotContinue:
				competitor.NotFinishedMsg = strings.Join(event.ExtraParams, " ")
				competitor.Finished = false
			default:
				continue loop
			}
		}

		competitors = append(competitors, competitor)
	}

	return competitors
}
