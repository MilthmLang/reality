package reality

type ChartInformationRepository interface {
	GetDifficulty(chartID string) (float64, error)
}

type ScoreRecord interface {
	GetChartID() string
	GetScore() float64
}
