# Vote Bot

##  Функциональноые требования

- **Создание голосования**: `/vote Вопрос | Вариант1 | Вариант2 | ...`
- **Голосование**: `/vote poll_id option_id` (если пользователь уже голосовал, то голос перезаписывается)
- **Просмотр результатов**: `/vote results poll_id`
- **Завершение голосования**: `/vote close poll_id`
- **Удаление голосования**: `/vote delete poll_id`


## Запуск


git clone https://github.com/ulyanovikovak/vk_bot.git
cd vk_bot
docker compose build --no-cache
docker compose up

Бот будет доступен на http://localhost:8080/vote