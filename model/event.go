package model

import (
	"biathlon-system-prototype/config"
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	Registered         = iota + 1 // участник зарегистрировался
	StartTimeWasSet               // время старта было определено жеребьевкой
	OnStarLine                    // участник находится на линии старта
	Started                       // участник стартовал
	OnFiringRange                 // участник находится на стрельбище
	TargetHit                     // мишень поражена
	LeftFiringRange               // участник покинул огневой рубеж
	EnteredPenaltyLaps            // участник начал штрафной круг
	LeftPenaltyLaps               // участник покинул штрафной круг
	EndedMainLap                  // участник завершил основной круг
	CanNotContinue                // участник не может продолжить
)

const (
	Disqualified = 32 // участник дисквалифицирован
	Finished     = 33 // участник финишировал
)

// Event - структура для хранения данных о событии
type Event struct {
	Time         time.Time // время, когда произошло событие
	Id           int       // id события
	CompetitorID int       // id участника
	ExtraParams  []string  // дополнительные параметры
}

// ParseEvent - преобразует строку в объект структуры Event
func ParseEvent(str string) (*Event, error) {
	if !strings.HasPrefix(str, "[") {
		return nil, fmt.Errorf("некорректный формат строки события: %s", str)
	}

	timeEnd := strings.Index(str, "]")
	if timeEnd == -1 {
		return nil, fmt.Errorf("пропущена закрывающаяся скобка ] в строке события: %s", str)
	}

	parts := strings.Fields(str)

	if len(parts) < 3 {
		return nil, fmt.Errorf("некорректный формат строки события: %s", parts)
	}

	timeStr := strings.Trim(parts[0], "[]")

	eventTime, err := time.Parse(config.TimeLayout, timeStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат времени: %v", err)
	}

	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("некорректный формат id события: %v", err)
	}

	competitorID, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("некорректный формат id участника: %v", err)
	}

	var params []string
	if len(parts) > 3 {
		params = parts[3:]
	}

	return &Event{
		Time:         eventTime,
		Id:           eventID,
		CompetitorID: competitorID,
		ExtraParams:  params,
	}, nil
}

// ReadEvents - считывает события из файла и выводит логи
func ReadEvents(file *os.File) (map[int][]*Event, error) {
	scanner := bufio.NewScanner(file)

	events := make(map[int][]*Event)
	var eventsSlice []*Event

	for scanner.Scan() {
		t := scanner.Text()

		func(t string) {
			event, err := ParseEvent(t)
			if err != nil {
				fmt.Println(err)
				return
			}

			events[event.CompetitorID] = append(events[event.CompetitorID], event)
			eventsSlice = append(eventsSlice, event)
		}(t)
	}

	sortEvents(eventsSlice)

	for _, ev := range eventsSlice {
		printEventLog(ev)
	}

	return events, nil
}

// printEventLog - выводит лог события
func printEventLog(event *Event) {
	var message string

	switch event.Id {
	case Registered:
		message = fmt.Sprintf("The competitor(%d) registered", event.CompetitorID)
		break
	case StartTimeWasSet:
		message = fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", event.CompetitorID, event.ExtraParams[0])
		break
	case OnStarLine:
		message = fmt.Sprintf("The competitor(%d) is on the start line", event.CompetitorID)
		break
	case Started:
		message = fmt.Sprintf("The competitor(%d) has started", event.CompetitorID)
		break
	case OnFiringRange:
		message = fmt.Sprintf("The competitor(%d) is on the firing range(%s)", event.CompetitorID, event.ExtraParams[0])
		break
	case TargetHit:
		message = fmt.Sprintf("The target(%s) has been hit by competitor(%d)", event.ExtraParams[0], event.CompetitorID)
		break
	case LeftFiringRange:
		message = fmt.Sprintf("The competitor(%d) left the firing range", event.CompetitorID)
		break
	case EnteredPenaltyLaps:
		message = fmt.Sprintf("The competitor(%d) entered the penalty laps", event.CompetitorID)
		break
	case LeftPenaltyLaps:
		message = fmt.Sprintf("The competitor(%d) left the penalty laps", event.CompetitorID)
		break
	case EndedMainLap:
		message = fmt.Sprintf("The competitor(%d) ended the main lap", event.CompetitorID)
		break
	case CanNotContinue:
		message = fmt.Sprintf("The competitor(%d) can't continue: %s", event.CompetitorID, strings.Join(event.ExtraParams, " "))
		break
	default:
		message = fmt.Sprintf("invalid event %d for competitor(%d)", event.Id, event.CompetitorID)
		break
	}

	fmt.Printf("[%v] %v\n", event.Time.Format(config.TimeLayout), message)
}

// sortEvents сортирует события по возрастанию времени
func sortEvents(events []*Event) {
	slices.SortFunc(events, func(a, b *Event) int {
		if a.Time.Before(b.Time) {
			return -1
		} else if a.Time.After(b.Time) {
			return 1
		}
		return 0
	})
}
