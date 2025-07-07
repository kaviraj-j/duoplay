package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

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

	// Start a goroutine to check for matches
	go q.checkForMatches(ctx)

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

func (q *inMemoryQueueRepository) checkForMatches(ctx context.Context) {
	// Wait a bit to allow more players to join
	time.Sleep(100 * time.Millisecond)

	q.mu.Lock()
	defer q.mu.Unlock()

	// If we have 2 or more players, create a match
	waitingPlayers := make([]string, 0, len(q.waitingPlayers))
	for userID := range q.waitingPlayers {
		waitingPlayers = append(waitingPlayers, userID)
	}

	if len(waitingPlayers) >= 2 {
		// Take the first two players
		player1ID := waitingPlayers[0]
		player2ID := waitingPlayers[1]

		// Remove them from queue
		delete(q.waitingPlayers, player1ID)
		delete(q.waitingPlayers, player2ID)

		// Call the match callback if set
		if q.matchCallback != nil {
			go q.matchCallback(ctx, player1ID, player2ID)
		}
	}
}

func (q *inMemoryQueueRepository) PlayerExistsInQueue(ctx context.Context, userID string) (bool, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	_, exists := q.waitingPlayers[userID]
	return exists, nil
}
