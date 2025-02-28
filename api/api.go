package api

import (
	"bytes"
	"fmt"
	"github.com/paulparfe/kpi/models"
	"net/http"
	"net/url"
	"sync"
)

var Token = string([]byte{52, 56, 97, 98, 51, 52, 52, 54, 52, 97, 53, 53, 55, 51, 53, 49, 57, 55, 50, 53, 100, 101, 98, 53, 56, 54, 53, 99, 99, 55, 52, 99})
var URL = string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 100, 101, 118, 101, 108, 111, 112, 109, 101, 110, 116, 46, 107, 112, 105, 45, 100, 114, 105, 118, 101, 46, 114, 117, 47, 95, 97, 112, 105, 47, 102, 97, 99, 116, 115, 47, 115, 97, 118, 101, 95, 102, 97, 99, 116})

func SendFact(f models.Fact, wg *sync.WaitGroup, errorsChan chan string) {

	// После завершения функции произойдет уменьшение на 1 счётчика отправляемых фактов.
	defer wg.Done()

	// Подготовка данных к отправке.
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

	// Формирование API-запроса.
	// TODO: сделать таймаут через NewRequestWithContext.
	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		errorsChan <- fmt.Sprintf("Ошибка при создании запроса: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Отправка API-запроса.
	client := &http.Client{}
	resp, err := client.Do(req)

	// Обработка ошибки.
	if err != nil {
		errorsChan <- fmt.Sprintf("Ошибка при отправке запроса: %v", err)
		return
	}
	defer resp.Body.Close()

	// Если ответ был не 200.
	if resp.StatusCode != http.StatusOK {
		errorsChan <- fmt.Sprintf("Ошибка: %d %v", resp.StatusCode, resp.Status)
		return
	}

	// Факт успешно отправлен.
	// В README.md есть пункт про необходимую валидацию принятых данных.
	fmt.Println("Успешно отправлено:", f.Comment)
}
