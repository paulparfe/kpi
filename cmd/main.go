package main

import (
	"fmt"
	"github.com/paulparfe/kpi/api"
	"github.com/paulparfe/kpi/models"
	"sync"
	"time"
)

// Генерируем count фактов для отправки
func generateFacts(count int) []models.Fact {
	facts := make([]models.Fact, count)
	for i := 0; i < count; i++ {
		facts[i] = models.Fact{
			PeriodStart:         "2024-12-01",
			PeriodEnd:           "2024-12-31",
			PeriodKey:           "month",
			IndicatorToMoID:     "227373",
			IndicatorToMoFactID: "0",
			Value:               "1",
			FactTime:            "2024-12-31",
			IsPlan:              "0",
			AuthUserID:          "40",
			Comment:             fmt.Sprintf("PaulParfe %d", i+1),
		}
	}
	return facts
}

func main() {
	maxBatchSize := 2         // Размер пачки
	facts := generateFacts(7) // Факты для отправки

	for i := 0; i < len(facts); i += maxBatchSize {
		var wg sync.WaitGroup
		errorsChan := make(chan string, maxBatchSize)

		end := i + maxBatchSize
		if end > len(facts) {
			end = len(facts)
		}

		// Параллельная отправка фактов из пачки
		for _, fact := range facts[i:end] {
			wg.Add(1)
			go api.SendFact(fact, &wg, errorsChan)
		}

		// Ожидаем когда отправка пачки завершится
		wg.Wait()
		close(errorsChan)

		// Пауза чтобы пачки отличалсь по времени ("post_time": "26.02.2025 08:09:21",)
		time.Sleep(2 * time.Second)

		// Если произошла ошибка, то завершаем обработку фактов
		for err := range errorsChan {
			fmt.Println("Ошибка в пачке:", err)
			return
		}
	}

	fmt.Println("Все данные успешно отправлены!")
}
