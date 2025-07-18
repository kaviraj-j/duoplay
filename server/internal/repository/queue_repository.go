package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kaviraj-j/duoplay/internal/model"
)

type QueueRepository interface {
	AddToQueue(ctx context.Context, userID string, conn *websocket.Conn) error
	RemoveFromQueue(ctx context.Context, userID string) error
	GetWaitingPlayers(ctx context.Context) ([]string, error)
	GetPlayerConnection(ctx context.Context, userID string) (*websocket.Conn, error)
	SetMatchCallback(callback func(ctx context.Context, player1ID, player2ID string) (*model.Room, error))
	PlayerExistsInQueue(ctx context.Context, userID string) (bool, error)
}

type inMemoryQueueRepository struct {
	waitingPlayers map[string]*websocket.Conn
	mu             sync.RWMutex
	matchCallback  func(ctx context.Context, player1ID, player2ID string) (*model.Room, error)
}

func NewQueueRepository() QueueRepository {
	return &inMemoryQueueRepository{
		waitingPlayers: make(map[string]*websocket.Conn),
	}
}

func (q *inMemoryQueueRepository) SetMatchCallback(callback func(ctx context.Context, player1ID, player2ID string) (*model.Room, error)) {
	q.matchCallback = callback
}

func (q *inMemoryQueueRepository) AddToQueue(ctx context.Context, userID string, conn *websocket.Conn) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.waitingPlayers[userID] = conn

	return nil
}

func (q *inMemoryQueueRepository) RemoveFromQueue(ctx context.Context, userID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if conn, exists := q.waitingPlayers[userID]; exists {
		conn.Close()
		delete(q.waitingPlayers, userID)
	}

	return nil
}

func (q *inMemoryQueueRepository) GetWaitingPlayers(ctx context.Context) ([]string, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	players := make([]string, 0, len(q.waitingPlayers))
	for userID := range q.waitingPlayers {
		players = append(players, userID)
	}

	return players, nil
}

func (q *inMemoryQueueRepository) GetPlayerConnection(ctx context.Context, userID string) (*websocket.Conn, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	conn, exists := q.waitingPlayers[userID]
	if !exists {
		return nil, fmt.Errorf("player not found in queue")
	}

	return conn, nil
}

func (q *inMemoryQueueRepository) PlayerExistsInQueue(ctx context.Context, userID string) (bool, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	_, exists := q.waitingPlayers[userID]
	return exists, nil
}
