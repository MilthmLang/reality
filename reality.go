package reality

import (
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

	bestScore := score.GetScore()
	switch {
	case bestScore >= 1000000:
		return difficulty + 1.5, nil
	case bestScore >= 850000:
		return difficulty + float64(bestScore-850000)/100000.0, nil
	case bestScore >= 700000:
		result := difficulty*(0.5+float64(bestScore-700000)/300000.0) + float64(bestScore-850000)/100000.0
		if result < 0 {
			return 0, nil
		}
		return result, nil
	case bestScore >= 600000:
		result := (difficulty - 3) * float64(bestScore-600000) / 200000.0
		if result < 0 {
			return 0, nil
		}
		return result, nil
	default:
		return 0, nil
	}
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
