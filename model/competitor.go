package model

import (
	"biathlon-system-prototype/config"
	"fmt"
	"math"
	"strings"
	"time"
)

// Competitor - структура для хранения данных участника
type Competitor struct {
	Id             int           // id участника
	Registered     bool          // флаг для указания зарегистрированности участника
	StartPlanned   time.Time     // планируемое время старта
	StartActual    time.Time     // реальное время старта
	FinishedTime   time.Time     // время финиширования
	Laps           []LapInfo     // информация о пройденных основных кругах
	PenaltyLaps    LapInfo       // информация о пройденных штрафных кругах
	Hits           int           // количество попаданий
	NotFinishedMsg string        // причина из-за которой участник не может продолжить участие
	Status         int           // статус участника
	TotalTime      time.Duration // общее время, затраченное на маршрут
}

// String выводит данные участника в форматрированном виде
func (comp *Competitor) String(conf *config.Config) string {
	var total string

	switch comp.Status {
	case Disqualified:
		total = "NotStarted"
	case CanNotContinue:
		total = "NotFinished"
	case Finished:
		total = time.Time{}.Add(comp.TotalTime).Format(config.TimeLayout)
	}

	var lapsStr string
	for _, lap := range comp.Laps {
		if lap.Speed == 0 {
			lapsStr += "{,}, "
		} else {
			lapsStr += fmt.Sprintf("{%s, %.3f}, ", formatDuration(lap.Duration), math.Floor(lap.Speed*1000)/1000)
		}
	}

	lapsStr = strings.TrimSpace(lapsStr)
	if len(lapsStr) > 0 {
		lapsStr = lapsStr[:len(lapsStr)-1]
	}

	var penalty string
	if comp.PenaltyLaps.Speed == 0 {
		penalty = "{,}"
	} else {
		penalty = fmt.Sprintf("{%s, %.3f}", formatDuration(comp.PenaltyLaps.Duration), comp.PenaltyLaps.Speed)
	}

	return fmt.Sprintf("[%s] %d [%s] %s %d/%d", total, comp.Id, lapsStr, penalty, comp.Hits, config.NumberTargets*conf.FiringLines)
}

// BuildCompetitors на основе событий создает участников и возвращает их в виде списка
func BuildCompetitors(events map[int][]*Event, conf *config.Config) *[]Competitor {
	var competitors []Competitor

	for competitorId, eventList := range events {
		competitor := Competitor{Id: competitorId, TotalTime: time.Duration(math.MaxInt64)}

		sortEvents(eventList)

		fillCompetitorFields(&competitor, eventList, conf)

		completeLapsIfNeeded(&competitor, conf)

		competitors = append(competitors, competitor)
	}

	return &competitors
}

// fillCompetitorFields заполняет поля участника данными на основе списка событий
func fillCompetitorFields(competitor *Competitor, eventList []*Event, conf *config.Config) {
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
				fmt.Printf("некорректный формат времени старта для участника %d: %v\n", competitor.Id, err)
				continue loop
			}

			competitor.StartPlanned = parseTime
			break
		case Started:
			competitor.StartActual = event.Time
			lastMainLapEnd = competitor.StartPlanned

			startDelta, _ := time.Parse(config.StartTimeLayout, conf.StartDelta)
			delta := startDelta.Sub(time.Time{})

			if toTimeOfDay(competitor.StartActual).After(toTimeOfDay(competitor.StartPlanned.Add(delta))) {
				competitor.Status = Disqualified
				break loop
			}
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

			speed = math.Floor(speed*1000) / 1000

			competitor.PenaltyLaps.Speed = speed
			break
		case EndedMainLap:
			duration := event.Time.Sub(lastMainLapEnd)
			competitor.Laps = append(competitor.Laps, LapInfo{Duration: duration, Speed: float64(conf.LapLen) / duration.Seconds()})
			lastMainLapEnd = event.Time

			if len(competitor.Laps) == conf.Laps {
				competitor.Status = Finished
				competitor.FinishedTime = event.Time

				competitor.TotalTime = competitor.FinishedTime.Sub(competitor.StartPlanned)
			}
			break
		case CanNotContinue:
			competitor.Status = CanNotContinue
			competitor.NotFinishedMsg = strings.Join(event.ExtraParams, " ")
		default:
			continue loop
		}
	}
}

// completeLapsIfNeeded добавляет пустые круги, если необходимо
func completeLapsIfNeeded(competitor *Competitor, conf *config.Config) {
	if competitor.Status != Finished && len(competitor.Laps) < conf.Laps {
		for i := len(competitor.Laps); i < conf.Laps; i++ {
			competitor.Laps = append(competitor.Laps, LapInfo{})
		}
	}
}

// formatDuration переводит time.Duration в форматированный time.Time
func formatDuration(d time.Duration) string {
	t := time.Time{}.Add(d)
	return t.Format(config.TimeLayout)
}

// toTimeOfDay обнуляет дату оставляя только время
func toTimeOfDay(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}
