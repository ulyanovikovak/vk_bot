package poll

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type PollService struct {
	storage Storage
}

type Storage interface {
	SavePoll(poll *Poll) error
	GetPoll(id string) (*Poll, error)
	DeletePoll(id string) error
}

func NewPollService(storage Storage) *PollService {
	return &PollService{storage: storage}
}

func (s *PollService) CreatePoll(question string, options []string, userID string) (*Poll, error) {
	id := uuid.New().String()

	optionMap := make(map[string]string)
	for _, opt := range options {
		optionMap[uuid.New().String()] = opt
	}

	poll := &Poll{
		ID:        id,
		Question:  question,
		Options:   optionMap,
		Votes:     make(map[string]string),
		CreatedBy: userID,
		IsClosed:  false,
	}

	err := s.storage.SavePoll(poll)
	if err != nil {
		return nil, err
	}

	log.Printf("Создано голосование: %s от пользователя %s", id, userID)
	return poll, nil
}

func (s *PollService) Vote(pollID, optionID, userID string) error {
	poll, err := s.storage.GetPoll(pollID)
	if err != nil {
		return err
	}
	if poll.IsClosed {
		return fmt.Errorf("голосование закрыто")
	}

	if _, ok := poll.Options[optionID]; !ok {
		return fmt.Errorf("неверный вариант")
	}

	poll.Votes[userID] = optionID
	log.Printf("Пользователь %s проголосовал за %s в голосовании %s", userID, optionID, pollID)
	return s.storage.SavePoll(poll)
}

func (s *PollService) GetResults(pollID string) (string, error) {
	poll, err := s.storage.GetPoll(pollID)
	if err != nil {
		return "", err
	}

	counts := make(map[string]int)
	for _, optionID := range poll.Votes {
		counts[optionID]++
	}

	result := fmt.Sprintf("Результаты голосования: %s\n", poll.Question)
	for id, text := range poll.Options {
		result += fmt.Sprintf("- %s: %d голос(ов)\n", text, counts[id])
	}

	log.Printf("Пользователь запросил результаты голосования %s", pollID)
	return result, nil
}

func (s *PollService) ClosePoll(pollID, userID string) error {
	p, err := s.storage.GetPoll(pollID)
	if err != nil {
		return err
	}

	if p.CreatedBy != userID {
		return fmt.Errorf("только создатель может завершить голосование")
	}

	p.IsClosed = true
	log.Printf("Пользователь %s завершил голосование %s", userID, pollID)
	return s.storage.SavePoll(p)
}

func (s *PollService) DeletePoll(pollID, userID string) error {
	p, err := s.storage.GetPoll(pollID)
	if err != nil {
		return err
	}

	if p.CreatedBy != userID {
		return fmt.Errorf("только создатель может удалить голосование")
	}

	log.Printf("Пользователь %s удалил голосование %s", userID, pollID)
	return s.storage.DeletePoll(pollID)
}
