package models

import (
	"context"
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

// StreamFacts Возвращает канал фактов и заполняет ими.
func StreamFacts(ctx context.Context) <-chan Fact {
	ch := make(chan Fact)

	go func() {
		defer close(ch)

		i := 1
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Функция FactsChannel завершает работу")
				return
			default:
			}

			ch <- getFact(i)
			i++

			// Временный выход после 10 фактов. Только для тестового задания.
			if i == 10 {
				return
			}
		}
	}()

	return ch
}

func getFact(i int) Fact {
	return Fact{
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
