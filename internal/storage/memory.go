package storage

import (
	"errors"
	"github.com/K1la/delayed-notifier/internal/models"
	"sync"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type MemoryStorage struct {
	data map[string]models.Notification
	mu   sync.RWMutex
}

func New() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]models.Notification),
	}
}

func (m *MemoryStorage) Save(n models.Notification) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[n.ID] = n
}

func (m *MemoryStorage) Get(id string) (models.Notification, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n, ok := m.data[id]
	if !ok {
		return models.Notification{}, ErrNotificationNotFound
	}
	return n, nil
}

func (s *MemoryStorage) GetAll() []models.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]models.Notification, 0, len(s.data))
	for _, n := range s.data {
		list = append(list, n)
	}
	return list
}

func (m *MemoryStorage) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.data[id]
	if !ok {
		return ErrNotificationNotFound
	}
	delete(m.data, id)
	return nil
}
