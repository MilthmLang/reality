package reality

import (
	"fmt"
	"math"
	"testing"
)

type SimpleChartRepo struct {
	data map[string]float64
}

func (s SimpleChartRepo) GetDifficulty(chartID string) (float64, error) {
	ret, ok := s.data[chartID]
	if !ok {
		return 0, fmt.Errorf("chart not exist")
	}
	return ret, nil
}

var simpleChartRepo = SimpleChartRepo{
	data: make(map[string]float64),
}

func init() {
	simpleChartRepo.data["1"] = 2
}

type SimpleScoreRecord struct {
	ChartID string
	Score   float64
}

func (s SimpleScoreRecord) GetChartID() string {
	return s.ChartID
}

func (s SimpleScoreRecord) GetScore() float64 {
	return s.Score
}

func TestSingleRecordValue(t *testing.T) {
	score := SimpleScoreRecord{
		ChartID: "1",
		Score:   1005000,
	}
	result, err := CalculateSingleEntryReality(score, simpleChartRepo)
	if err != nil {
		t.Error(err)
	}
	if math.Abs(result-3.5) >= 1e8 {
		t.Errorf("result should be 3.5")
	}
}

func TestEmptyRecordReality(t *testing.T) {
	result, err := CalculateReality([]ScoreRecord{}, simpleChartRepo)
	if err != nil {
		t.Error(err)
	}
	if math.Abs(result) >= 1e8 {
		t.Errorf("result should be 3.5")
	}
}

func TestSingleRecordReality(t *testing.T) {
	score := SimpleScoreRecord{
		ChartID: "1",
		Score:   1005000,
	}
	result, err := CalculateReality([]ScoreRecord{score}, simpleChartRepo)
	if err != nil {
		t.Error(err)
	}
	if math.Abs(result-3.5/20) >= 1e8 {
		t.Errorf("result should be 3.5")
	}
}

func TestRecordRealityWithSingleValue(t *testing.T) {
	score := SimpleScoreRecord{
		ChartID: "1",
		Score:   1005000,
	}
	scores := make([]ScoreRecord, 0)
	for i := 0; i < 50; i++ {
		scores = append(scores, score)
	}

	result, err := CalculateReality(scores, simpleChartRepo)
	if err != nil {
		t.Error(err)
	}
	if math.Abs(result-3.5) >= 1e8 {
		t.Errorf("result should be 3.5")
	}
}
