package internal

import (
	"biathlon-system-prototype/config"
	"biathlon-system-prototype/model"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PrintReport - выводит участников с их результатами
func PrintReport(config *config.Config, file *os.File) error {
	events, err := readEvents(file)
	if err != nil {
		return err
	}

	competitors := model.BuildCompetitors(events, config)

	fmt.Println(competitors)

	return nil
}

// readEvents - считывает события из файла и выводит логи
func readEvents(file *os.File) (map[int][]model.Event, error) {
	scanner := bufio.NewScanner(file)

	events := make(map[int][]model.Event)

	for scanner.Scan() {
		t := scanner.Text()

		func(t string) {
			event, err := model.ParseEvent(t)
			if err != nil {
				fmt.Println(err)
				return
			}

			printEventLog(event)

			events[event.CompetitorID] = append(events[event.CompetitorID], event)
		}(t)
	}

	return events, nil
}

// printEventLog - выводит логи событий
func printEventLog(event model.Event) {
	var message string

	switch event.Id {
	case model.Registered:
		message = fmt.Sprintf("The competitor(%d) registered", event.CompetitorID)
		break
	case model.StartTimeWasSet:
		message = fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", event.CompetitorID, event.ExtraParams[0])
		break
	case model.OnStarLine:
		message = fmt.Sprintf("The competitor(%d) is on the start line", event.CompetitorID)
		break
	case model.Started:
		message = fmt.Sprintf("The competitor(%d) has started", event.CompetitorID)
		break
	case model.OnFiringRange:
		message = fmt.Sprintf("The competitor(%d) is on the firing range(%s)", event.CompetitorID, event.ExtraParams[0])
		break
	case model.TargetHit:
		message = fmt.Sprintf("The target(%s) has been hit by competitor(%d)", event.ExtraParams[0], event.CompetitorID)
		break
	case model.LeftFiringRange:
		message = fmt.Sprintf("The competitor(%d) left the firing range", event.CompetitorID)
		break
	case model.EnteredPenaltyLaps:
		message = fmt.Sprintf("The competitor(%d) entered the penalty laps", event.CompetitorID)
		break
	case model.LeftPenaltyLaps:
		message = fmt.Sprintf("The competitor(%d) left the penalty laps", event.CompetitorID)
		break
	case model.EndedMainLap:
		message = fmt.Sprintf("The competitor(%d) ended the main lap", event.CompetitorID)
		break
	case model.CanNotContinue:
		message = fmt.Sprintf("The competitor(%d) can't continue: %s", event.CompetitorID, strings.Join(event.ExtraParams, " "))
		break
	default:
		message = fmt.Sprintf("invalid event %d for competitor(%d)", event.Id, event.CompetitorID)
		break
	}

	fmt.Printf("[%v] %v\n", event.Time.Format(config.TimeLayout), message)
}
