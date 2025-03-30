package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"vk_bot/poll"
)

func Handler(service *poll.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
			return
		}

		text := r.FormValue("text")
		userID := r.FormValue("user_id")
		log.Printf("Запрос от пользователя %s: %s", userID, text)

		args := strings.Fields(text)

		if len(args) == 2 && args[0] == "results" {
			pollID := args[1]
			log.Printf("Пользователь %s запросил результаты голосования %s", userID, pollID)
			resultText, err := service.GetResults(pollID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp := map[string]string{"response_type": "in_channel", "text": resultText}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if len(args) == 2 && args[0] == "close" {
			pollID := args[1]
			log.Printf("Пользователь %s хочет завершить голосование %s", userID, pollID)
			err := service.ClosePoll(pollID, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			resp := map[string]string{"response_type": "in_channel", "text": "Голосование завершено."}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if len(args) == 2 && args[0] == "delete" {
			pollID := args[1]
			log.Printf("Пользователь %s хочет удалить голосование %s", userID, pollID)
			err := service.DeletePoll(pollID, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			resp := map[string]string{"response_type": "in_channel", "text": "Голосование удалено."}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if len(args) == 2 {
			pollID := args[0]
			optionID := args[1]
			log.Printf("Пользователь %s голосует за %s в голосовании %s", userID, optionID, pollID)
			err = service.Vote(pollID, optionID, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp := map[string]string{"response_type": "in_channel", "text": "Голос принят!"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		args = strings.Split(text, "|")
		if len(args) < 2 {
			http.Error(w, "Неверный формат. Пример: вопрос | вариант1 | вариант2", http.StatusBadRequest)
			return
		}

		question := strings.TrimSpace(args[0])
		var options []string
		for _, opt := range args[1:] {
			options = append(options, strings.TrimSpace(opt))
		}

		p, err := service.CreatePoll(question, options, userID)
		if err != nil {
			http.Error(w, "Ошибка создания голосования", http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf("Голосование создано: %s\n", p.Question)
		for id, opt := range p.Options {
			msg += fmt.Sprintf("- [%s] %s\n", id[:6], opt)
		}

		log.Printf("Голосование %s создано пользователем %s", p.ID, userID)

		resp := map[string]string{"response_type": "in_channel", "text": msg}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
