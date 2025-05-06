package model

import "time"

// LapInfo - структура для хранения данных о круге
type LapInfo struct {
	Duration time.Duration // длительность преодоления круга
	Speed    float64       // средняя скорость на круге
}
