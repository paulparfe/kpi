package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/paulparfe/kpi/api"
	"github.com/paulparfe/kpi/buffer"
	"github.com/paulparfe/kpi/models"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Создаем контекст, который автоматически завершится при получении SIGTERM/SIGINT.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Канал фактов для отправки.
	facts := models.StreamFacts(ctx)

	// Добавляем сгенерированные факты в буфер.
	buff := buffer.NewBuffer()
	for fact := range facts {
		buff.Add(fact)
	}

	// Количество одновременно отправляемых фактов.
	maxBatchSize := 2

	// Цикл отправки фактов пачками.
mainLoop:
	for {
		select {
		case <-ctx.Done(): // Если контекст завершен (например, получен SIGTERM)
			fmt.Println("Выход из цикла обработки")
			break mainLoop
		default:
		}

		// Если буфер пуст, делаем паузу и продолжаем цикл.
		if buff.IsNothingToProcess() {
			time.Sleep(time.Second)
			continue
		}

		// Создаем WaitGroup для ожидания завершения отправки всех фактов из текущей пачки.
		var wg sync.WaitGroup

		// Запускаем несколько горутин для отправки фактов.
		for i := 0; i < maxBatchSize; i++ {
			wg.Add(1)

			// В отдельных горутинах отправляем факт за фактом параллельно.
			go worker(ctx, &wg, buff)
		}

		// Ожидаем когда отправка всех фактов из пачки завершится с ошибками или без.
		wg.Wait()

		// Пауза для тестового задания чтобы пачки отличались по времени ("post_time": "26.02.2025 08:09:21",)
		time.Sleep(2 * time.Second)
	}

	// Перед завершением программы сохраняем буфер в файл.
	saveBufferToFile(buff)
}

// Worker обрабатывает отправку одного факта из буфера
func worker(ctx context.Context, wg *sync.WaitGroup, buff *buffer.Buffer) {
	defer wg.Done() // Уменьшаем счетчик горутин при завершении работы.

	select {
	case <-ctx.Done(): // Если контекст завершен, выходим
		fmt.Println("Worker завершает работу (прервано)")
		return
	default:
	}

	// Извлекаем факт для отправки.
	fact := buff.GetFactForProcessing()
	if fact == nil {
		return // Если фактов нет, выходим.
	}

	// Пытаемся отправить факт через API.
	sentFact, err := api.SendFact(ctx, fact)
	if err != nil {
		// В случае ошибки меняем статус факта на "failed" (чтобы отправить повторно).
		buff.UpdateStatus(sentFact, buffer.StatusFailed)
		fmt.Println("Ошибка отправки:", err)
		return
	}

	// Если отправка успешна, удаляем факт из буфера.
	fmt.Printf("Успешно отправлен %v\n", sentFact.Comment)
	buff.Remove(sentFact)
}

// saveBufferToFile сохраняет оставшиеся факты в файл для возможности продолжить отправку и анализа.
func saveBufferToFile(buff *buffer.Buffer) {
	// Преобразуем буфер в JSON-формат.
	data, err := json.MarshalIndent(buff.GetAllFacts(), "", "  ")
	if err != nil {
		fmt.Println("Ошибка сериализации буфера:", err)
		return
	}

	// Записываем данные в файл buffer_backup.json.
	err = os.WriteFile("buffer_backup.json", data, 0644)
	if err != nil {
		fmt.Println("Ошибка сохранения буфера в файл:", err)
		return
	}

	fmt.Println("Буфер успешно сохранен в buffer_backup.json")
}
