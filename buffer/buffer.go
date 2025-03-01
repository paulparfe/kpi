package buffer

import (
	"github.com/paulparfe/kpi/models"
	"sync"
)

type FactStatus string

const (
	StatusNew        FactStatus = "new"        // Новый факт, ожидающий обработки.
	StatusProcessing FactStatus = "processing" // Факт в процессе обработки.
	StatusCompleted  FactStatus = "completed"  // Факт успешно обработан.
	StatusFailed     FactStatus = "failed"     // Ошибка обработки факта.
)

type Fact struct {
	ID     int        // Уникальный идентификатор факта.
	Status FactStatus // Текущий статус факта.

	models.Fact
}

// Buffer представляет собой хранилище фактов с синхронизацией доступа.
type Buffer struct {
	lastFactID int           // Последний использованный ID для фактов.
	facts      map[int]*Fact // Хранилище фактов.
	mu         sync.Mutex
}

// NewBuffer создает и возвращает новый экземпляр буфера.
func NewBuffer() *Buffer {
	return &Buffer{
		facts: make(map[int]*Fact),
	}
}

// Add добавляет новый факт в буфер, назначая ему уникальный идентификатор.
func (b *Buffer) Add(modelFact models.Fact) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lastFactID++ // Увеличиваем счетчик ID.

	bufferFact := &Fact{
		ID:     b.lastFactID,
		Status: StatusNew,
		Fact:   modelFact,
	}

	b.facts[b.lastFactID] = bufferFact
}

// GetFactForProcessing находит и возвращает факт для обработки, меняя его статус на "processing".
func (b *Buffer) GetFactForProcessing() *Fact {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, fact := range b.facts {
		if fact.Status != StatusProcessing { // Ищем факт, который еще не в обработке.
			fact.Status = StatusProcessing
			return fact
		}
	}

	return nil // Если нет фактов для обработки, возвращаем nil.
}

// UpdateStatus обновляет статус указанного факта.
func (b *Buffer) UpdateStatus(fact *Fact, status FactStatus) {
	b.mu.Lock()
	defer b.mu.Unlock()

	fact.Status = status
}

// Remove удаляет факт из буфера.
func (b *Buffer) Remove(fact *Fact) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.facts, fact.ID)
}

// IsEmpty проверяет, пуст ли буфер (нет ли фактов).
func (b *Buffer) IsEmpty() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.facts) == 0
}

// IsNothingToProcess нужен для проверки, есть ли факты для обработки.
func (b *Buffer) IsNothingToProcess() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, fact := range b.facts {
		if fact.Status != StatusProcessing {
			return false
		}
	}

	return true
}

// GetAllFacts возвращает все факты, хранящиеся в буфере.
func (b *Buffer) GetAllFacts() map[int]*Fact {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.facts
}
