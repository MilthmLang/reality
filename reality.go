package reality

import (
	"math"
	"sort"
)

// CalculateSingleEntryReality is used for calculating how much reality this record donated
func CalculateSingleEntryReality(score ScoreRecord, chartRepo ChartInformationRepository) (float64, error) {
	chartId := score.GetChartID()

	difficulty, err := chartRepo.GetDifficulty(chartId)
	if err != nil {
		return 0, err
	}

	if difficulty < 1e-3 {
		return 0, nil
	}

	const (
		b          = 3.1
		k          = 3.65
		range1     = 0.8
		range2     = 0.7
		limitScore = 100_5000.0
		offset     = -0.5
	)

	s := score.GetScore()
	if s > limitScore {
		s = limitScore
	}

	var value float64
	switch {
	case s < 70_0000:
		return 0, nil

	case s < 98_0000:
		// Linear mapping from [700k..980k] onto [-1..0]
		value = s/(98_0000-700000) + (-1 - ((0-(-1))/(98_0000-float64(70_0000)))*70_0000)

	case s < 99_5000:
		x1 := (s - 98_0000) / (99_5000 - 98_0000)
		top := math.Exp(b*x1) - 1
		bottom := math.Exp(b) - 1
		value = top / bottom * range1

	case s < 100_5000:
		x2 := (s - 99_5000) / (100_5000 - 99_5000)
		bot := 1 + math.Exp(-k*x2)
		value = (1/bot-0.5)*2*range2 + range1

	case s <= 101_0000:
		value = 1.5
	}

	ret := difficulty + value + offset
	if ret > 0 {
		return ret, nil
	}
	return 0, nil
}

func calculateTotalRealityFromValues(values []float64) float64 {
	const size = 20
	ws := make([]float64, size)

	for i := 0; i < size; i++ {
		if i < len(values) {
			ws[i] = values[i] / float64(size)
		}
	}

	psize := size
	for psize > 1 {
		halfSize := psize / 2
		for i := 0; i < halfSize; i++ {
			ws[i] += ws[psize-i-1]
		}
		psize -= halfSize
	}

	return ws[0]
}

// CalculateReality is used for calculating reality from given scores list and charts information
func CalculateReality(scores []ScoreRecord, chartRepo ChartInformationRepository) (float64, error) {
	var ratings []float64

	for _, score := range scores {
		r, err := CalculateSingleEntryReality(score, chartRepo)
		if err != nil {
			continue
		}
		ratings = append(ratings, r)
	}

	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i] > ratings[j]
	})

	if len(ratings) > 20 {
		ratings = ratings[:20]
	}

	total := calculateTotalRealityFromValues(ratings)
	return total, nil
}
