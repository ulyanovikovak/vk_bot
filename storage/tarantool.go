package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"vk_bot/poll"

	"github.com/tarantool/go-tarantool"
)

type TarantoolStorage struct {
	conn *tarantool.Connection
}

func NewTarantoolStorage(conn *tarantool.Connection) *TarantoolStorage {
	return &TarantoolStorage{conn: conn}
}

func (s *TarantoolStorage) SavePoll(p *poll.Poll) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	log.Printf("Сохраняем голосование %s в Tarantool", p.ID)
	_, err = s.conn.Insert("polls", []interface{}{p.ID, string(data)})
	return err
}

func (s *TarantoolStorage) GetPoll(id string) (*poll.Poll, error) {
	log.Printf("Получаем голосование %s из Tarantool", id)
	resp, err := s.conn.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{id})
	if err != nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("голосование не найдено")
	}

	data := resp.Data[0].([]interface{})[1].(string)
	var p poll.Poll
	err = json.Unmarshal([]byte(data), &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *TarantoolStorage) DeletePoll(id string) error {
	log.Printf("Удаляем голосование %s из Tarantool", id)
	_, err := s.conn.Delete("polls", "primary", []interface{}{id})
	return err
}
