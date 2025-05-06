package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	TimeLayout      = "15:04:05.000" // шаблон для формата времени в событиях
	StartTimeLayout = "15:04:05"     // шаблон для формата времени старта
	NumberTargets   = 5              // количество мишеней на огневом рубеже
)

var (
	ConfigPath = "../config.json" // путь до файла конфигурации
	EventsPath = "../src/events"  // путь до файла с событиями
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

// LoadEnvVars загружает переменные окружения из файла .env
func LoadEnvVars() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Ошибка загрузки .env файла, будут использованы значения по умолчанию: %v\n", err)
	}

	DefaultConfigPath := ConfigPath
	ConfigPath = os.Getenv("CONFIG_PATH")
	if ConfigPath == "" {
		ConfigPath = DefaultConfigPath
		log.Println("CONFIG_PATH не задан, используется значение по умолчанию:", ConfigPath)
	}

	DefaultEventsPath := EventsPath
	EventsPath = os.Getenv("EVENTS_PATH")
	if EventsPath == "" {
		EventsPath = DefaultEventsPath
		log.Println("EVENTS_PATH не задан, используется значение по умолчанию:", EventsPath)
	}
}
