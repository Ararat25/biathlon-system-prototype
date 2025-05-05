package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	TimeLayout      = "15:04:05.000"
	StartTimeLayout = "15:04:05"
)

// Config - структура для хранения переменных конфигурации
type Config struct {
	Laps        int    `json:"laps"`        // Количество кругов на основной дистанции
	LapLen      int    `json:"lapLen"`      // Длина каждого основного круга
	PenaltyLen  int    `json:"penaltyLen"`  // Продолжительность каждого штрафного круга
	FiringLines int    `json:"firingLines"` // Количество огневых рубежей на круге
	Start       string `json:"start"`       // Запланированное время старта первого участника
	StartDelta  string `json:"startDelta"`  // Планируемый интервал между запусками
}

// LoadConfig считывает данные из файла по указанному пути и возвращает в виде структуры Config
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия config файла: %v", err)
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения config файла: %v", err)
	}

	_, err = time.Parse(StartTimeLayout, config.Start)
	if err != nil {
		panic(fmt.Sprintf("некорректный формат запланированного времени старта первого участника %s: %v\n", config.Start, err))
	}

	_, err = time.Parse(StartTimeLayout, config.StartDelta)
	if err != nil {
		panic(fmt.Sprintf("некорректный формат интервала между стартами участников %s: %v\n", config.StartDelta, err))
	}

	return &config, nil
}
