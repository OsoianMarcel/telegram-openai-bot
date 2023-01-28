package stats

import (
	"encoding/json"
	"sync"
)

// Common stats properties.
type CommonStats struct {
	AiAllMessages   uint32 `json:"ai_all_messages"`
	AiResponses     uint32 `json:"ai_responses"`
	AiInvalidErrors uint32 `json:"ai_invalid_errors"`
	AiErrors        uint32 `json:"ai_errors"`
	AiTimeoutErrors uint32 `json:"ai_timeout_errors"`
}

// User data.
type User struct {
	CommonStats
}

// Stats data.
type Stats struct {
	CommonStats
	Users map[string]User `json:"users"`
	mutex sync.RWMutex
}

// Generate string format of all stats.
func (s *Stats) GenString() (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	jb, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "", err
	}

	return string(jb), nil
}

// Increase the counter of AI responses.
func (s *Stats) IncrAiResponses(userId int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashId := HashUserId(userId)
	s.AiAllMessages++
	s.AiResponses++

	user, ok := s.Users[hashId]
	if !ok {
		user = User{}
	}

	user.AiAllMessages++
	user.AiResponses++

	s.Users[hashId] = user
}

// Increase the counter of AI invalid errors.
func (s *Stats) IncrAiInvalidErrors(userId int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashId := HashUserId(userId)
	s.AiAllMessages++
	s.AiInvalidErrors++

	user, ok := s.Users[hashId]
	if !ok {
		user = User{}
	}

	user.AiAllMessages++
	user.AiInvalidErrors++

	s.Users[hashId] = user
}

// Increase the counter of AI general errors.
func (s *Stats) IncrAiErrors(userId int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashId := HashUserId(userId)
	s.AiAllMessages++
	s.AiErrors++

	user, ok := s.Users[hashId]
	if !ok {
		user = User{}
	}

	user.AiAllMessages++
	user.AiErrors++

	s.Users[hashId] = user
}

// Increase the counter of AI timeout errors.
func (s *Stats) IncrAiTimeoutErrors(userId int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashId := HashUserId(userId)

	s.AiAllMessages++
	s.AiTimeoutErrors++

	user, ok := s.Users[hashId]
	if !ok {
		user = User{}
	}

	user.AiAllMessages++
	user.AiTimeoutErrors++

	s.Users[hashId] = user
}
