package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/paulparfe/kpi/buffer"
	"net/http"
	"net/url"
	"time"
)

var Token = string([]byte{52, 56, 97, 98, 51, 52, 52, 54, 52, 97, 53, 53, 55, 51, 53, 49, 57, 55, 50, 53, 100, 101, 98, 53, 56, 54, 53, 99, 99, 55, 52, 99})
var URL = string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 100, 101, 118, 101, 108, 111, 112, 109, 101, 110, 116, 46, 107, 112, 105, 45, 100, 114, 105, 118, 101, 46, 114, 117, 47, 95, 97, 112, 105, 47, 102, 97, 99, 116, 115, 47, 115, 97, 118, 101, 95, 102, 97, 99, 116})

func SendFact(ctx context.Context, fact *buffer.Fact) (*buffer.Fact, error) {
	// Устанавливаем таймаут 10 секунд для запроса
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Подготовка данных к отправке.
	data := url.Values{}
	data.Set("period_start", fact.PeriodStart)
	data.Set("period_end", fact.PeriodEnd)
	data.Set("period_key", fact.PeriodKey)
	data.Set("indicator_to_mo_id", fact.IndicatorToMoID)
	data.Set("indicator_to_mo_fact_id", fact.IndicatorToMoFactID)
	data.Set("value", fact.Value)
	data.Set("fact_time", fact.FactTime)
	data.Set("is_plan", fact.IsPlan)
	data.Set("auth_user_id", fact.AuthUserID)
	data.Set("comment", fact.Comment)

	// Формирование API-запроса.
	req, err := http.NewRequestWithContext(ctx, "POST", URL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fact, fmt.Errorf("ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Отправка API-запроса.
	client := &http.Client{}
	resp, err := client.Do(req)

	// Проверяем, не истек ли таймаут.
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fact, fmt.Errorf("превышено время ожидания отправки факта %v", fact.ID)
	}

	// Обработка ошибки.
	if err != nil {
		return fact, fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	// Если ответ был не 200.
	if resp.StatusCode != http.StatusOK {
		return fact, fmt.Errorf("ошибка: %d %v", resp.StatusCode, resp.Status)
	}

	// Факт успешно отправлен.
	return fact, nil
}
