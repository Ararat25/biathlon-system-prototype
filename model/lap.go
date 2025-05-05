package model

import "time"

type LapInfo struct {
	Duration time.Duration // длительность преодоления круга
	Speed    float64       // средняя скорость на круге
}
