package internal

import (
	"biathlon-system-prototype/config"
	"biathlon-system-prototype/model"
	"fmt"
	"os"
	"slices"
)

// PrintReport выводит участников с их результатами
func PrintReport(config *config.Config, file *os.File) error {
	events, err := model.ReadEvents(file)
	if err != nil {
		return err
	}

	competitors := model.BuildCompetitors(events, config)

	sortCompetitors(competitors)

	for _, comp := range *competitors {
		if comp.Registered {
			fmt.Println(comp.String(config))
		}
	}

	return nil
}

// sortCompetitors сортирует участников по возрастанию общего времени
func sortCompetitors(competitors *[]model.Competitor) {
	slices.SortFunc(*competitors, func(a, b model.Competitor) int {
		if a.TotalTime < b.TotalTime {
			return -1
		} else if a.TotalTime > b.TotalTime {
			return 1
		}
		return 0
	})
}
