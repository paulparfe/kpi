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
	// Количество одновременно отправляемых фактов.
	// Возможно лучше будет сделать несколько одновременно работающих "отправлятелей" фактов.
	maxBatchSize := 2

	// Сейчас генерируется фиксированное количество фактов для отправки.
	// В реальности нужна функция для добавления фактов в буфер.
	// В реальности нужна валидация принятых данных.
	facts := generateFacts(7) // Все факты для отправки

	// Цикл заполнения батча, отправки, проверки результатов отправки.
	// Должен быть бесконечный цикл, с возможностью выхода по graceful shutdown.
	for i := 0; i < len(facts); i += maxBatchSize {

		// Для ожидания завершения отправки всех фактов из пачки.
		var wg sync.WaitGroup

		// Для приёма сообщений об ошибках из каждого отправленного факта.
		errorsChan := make(chan string, maxBatchSize)

		// Возможно в буфере остался 1 факт, а размер пачки = 2.
		// Тогда нужно в пачку положить этот 1 факт.
		end := i + maxBatchSize
		if end > len(facts) {
			end = len(facts)
		}

		// Считывать из буфера факты по-одному.
		// В отдельных горутинах отправлять факт за фактом параллельно.
		for _, fact := range facts[i:end] {
			// Увеличим счётчик отправляемых фактов на 1.
			wg.Add(1)

			// Отправляем факт.
			// После завершения функции произойдет уменьшение на 1 счётчика отправляемых фактов.
			go api.SendFact(fact, &wg, errorsChan)
		}

		// Ожидаем когда отправка всех фактов из пачки завершится с ошибками или без.
		// Счётчик отправляемых фактов будет равен 0.
		wg.Wait()

		// Закроем канал сообщений об ошибках.
		close(errorsChan)

		// Если произошла ошибка, то выводим её в логи.
		// Факты с ошибками, возможно нужно оставить в буфере и попытаться отправить позже.
		for err := range errorsChan {
			fmt.Println("Ошибка в пачке:", err)
		}

		// Если не было ошибки, то факты нужно удалить из буфера.

		// Пауза чтобы пачки отличались по времени ("post_time": "26.02.2025 08:09:21",)
		time.Sleep(2 * time.Second)
	}

}
