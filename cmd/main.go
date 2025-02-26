package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var token = string([]byte{52, 56, 97, 98, 51, 52, 52, 54, 52, 97, 53, 53, 55, 51, 53, 49, 57, 55, 50, 53, 100, 101, 98, 53, 56, 54, 53, 99, 99, 55, 52, 99})
var apiURL = string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 100, 101, 118, 101, 108, 111, 112, 109, 101, 110, 116, 46, 107, 112, 105, 45, 100, 114, 105, 118, 101, 46, 114, 117, 47, 95, 97, 112, 105, 47, 102, 97, 99, 116, 115, 47, 115, 97, 118, 101, 95, 102, 97, 99, 116})

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

// Генерируем count фактов для отправки
func generateFacts(count int) []Fact {
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
			Comment:             fmt.Sprintf("PaulParfe %d", i+1),
		}
	}
	return facts
}

func sendFact(f Fact, wg *sync.WaitGroup, errorsChan chan string) {
	defer wg.Done()

	data := url.Values{}
	data.Set("period_start", f.PeriodStart)
	data.Set("period_end", f.PeriodEnd)
	data.Set("period_key", f.PeriodKey)
	data.Set("indicator_to_mo_id", f.IndicatorToMoID)
	data.Set("indicator_to_mo_fact_id", f.IndicatorToMoFactID)
	data.Set("value", f.Value)
	data.Set("fact_time", f.FactTime)
	data.Set("is_plan", f.IsPlan)
	data.Set("auth_user_id", f.AuthUserID)
	data.Set("comment", f.Comment)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		errorsChan <- fmt.Sprintf("Ошибка при создании запроса: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorsChan <- fmt.Sprintf("Ошибка при отправке запроса: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorsChan <- fmt.Sprintf("Ошибка: %d %v", resp.StatusCode, resp.Status)
		return
	}

	fmt.Println("Успешно отправлено:", f.Comment)
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
			go sendFact(fact, &wg, errorsChan)
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
