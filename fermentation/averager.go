package fermentation

import (
	"sort"
	"time"
)

type averager struct {
	readings []averagerReading
	window   *time.Duration
	outliers *float64
}

type averagerReading struct {
	reading   float64
	timestamp time.Time
}

func (a *averager) prepare() {
	if a.window == nil {
		defaultWindow := 5 * time.Minute
		a.window = &defaultWindow
	}
	if a.readings == nil {
		a.readings = make([]averagerReading, 0, 1000)
	}
	if a.outliers == nil {
		defaultOutlier := 0.1
		a.outliers = &defaultOutlier
	}
}

func (a *averager) addReading(reading float64) float64 {
	a.prepare()

	curTime := time.Now()

	curReadings := make([]float64, 0, 1000)
	curReadings = append(curReadings, reading)
	inserted := false

	for i, r := range a.readings {
		if r.timestamp.Add(*a.window).After(curTime) {
			curReadings = append(curReadings, r.reading)
		} else {
			a.readings[i] = averagerReading{reading, curTime}
		}
	}
	if !inserted {
		a.readings = append(a.readings, averagerReading{reading, curTime})
	}

	return avg(removeOutliers(curReadings, *a.outliers))
}

func removeOutliers(vals []float64, fraction float64) []float64 {
	sort.Float64s(vals)
	cut := int(float64(len(vals)) * fraction)
	if cut == 0 && len(vals) >= 3 {
		cut = 1
	}

	return vals[cut : len(vals)-cut]
}

func avg(vals []float64) (ret float64) {
	for _, v := range vals {
		ret += v
	}
	ret /= float64(len(vals))
	return
}
