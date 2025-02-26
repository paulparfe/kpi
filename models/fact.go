package models

type Fact struct {
	PeriodStart         string
	PeriodEnd           string
	PeriodKey           string
	IndicatorToMoID     string
	IndicatorToMoFactID string
	Value               string
	FactTime            string
	IsPlan              string
	AuthUserID          string
	Comment             string
}
