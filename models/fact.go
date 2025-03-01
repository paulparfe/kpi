package models

import (
	"fmt"
)

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

// GenerateFacts Генерируем count фактов для отправки.
func GenerateFacts(count int) []Fact {
	facts := make([]Fact, count)
	for i := 0; i < count; i++ {
		facts[i] = Fact{
			PeriodStart:         "2024-12-01",
			PeriodEnd:           "2024-12-31",
			PeriodKey:           "month",
			IndicatorToMoID:     "227373",
			IndicatorToMoFactID: "0",
			Value:               "1",
			FactTime:            "2024-12-31",
			IsPlan:              "0",
			AuthUserID:          "40",
			Comment:             fmt.Sprintf("buffer PaulParfe %d", i+1),
		}
	}
	return facts
}
