Вот детальное Техническое Задание (ТЗ), структурированное для разработки системы в микросервисной архитектуре с использованием Golang, gRPC и брокеров сообщений.

Техническое Задание: Симулятор Банковской Системы (Microservices)
1. Общее описание
Разработать бэкенд-систему, симулирующую работу банка. Система должна состоять из микросервисов, взаимодействующих через gRPC и RabbitMQ, обеспечивать управление счетами, обработку транзакций и кэширование данных.
2. Стек технологий
Язык: Golang (1.21+)
База данных: PostgreSQL
Кэш: Redis
Брокер сообщений: RabbitMQ
Протокол межсервисного взаимодействия: gRPC (Protobuf)
Конфигурация: YAML
Миграции: (на выбор: golang-migrate, goose)

3. Архитектура системы
Система делится на два основных микросервиса:
Account Service (Сервис Счетов): Отвечает за создание/удаление счетов, хранение баланса, предоставление HTTP API.
Transaction Service (Сервис Транзакций): Отвечает за логику проведения операций (Producer/Consumer), хранение истории транзакций.
3.1 Схема взаимодействия
HTTP Requests поступают в Account Service (выступает в роли API Gateway для клиента).
Чтение данных (баланс, история) идет через Redis. Если данных нет — запрос в БД.
Изменение баланса (депозит, снятие, перевод) происходит асинхронно через RabbitMQ.
Синхронная коммуникация (валидация существования счета) происходит через gRPC.

4. База данных (PostgreSQL)
Необходимо создать две базы данных (или две схемы) для изоляции сервисов.
4.1 Таблица accounts (Владелец: Account Service)
Поле
Тип
Описание
id
UUID / BIGSERIAL
Первичный ключ
balance
DECIMAL / BIGINT
Текущий баланс (копейки/центы)
currency
VARCHAR(3)
Валюта (напр. USD, RUB)
is_locked
BOOLEAN
Блокировка счета при операциях
created_at
TIMESTAMP
Дата создания
deleted_at
TIMESTAMP
Дата удаления (Soft Delete)

4.2 Таблица transactions (Владелец: Transaction Service)
Поле
Тип
Описание
id
UUID / BIGSERIAL
Первичный ключ
account_id
UUID / BIGINT
Ссылка на счет
to_account_id
UUID / BIGINT
Ссылка на счет получателя (для переводов)
amount
DECIMAL / BIGINT
Сумма операции
transaction_type
VARCHAR
deposit, withdraw, transfer
created_at
TIMESTAMP
Дата проведения


5. API Интерфейс (REST)
Все запросы принимает Account Service.
5.1 Управление счетами
POST /api/accounts — Создание счета.
Действие: Запись в БД, запись в Redis.
DELETE /api/accounts/{id} — Удаление счета.
Действие: Soft delete в БД, удаление из Redis.
GET /api/accounts/{id} — Получение счета.
Логика: Искать в Redis -> если нет, искать в БД -> сохранить в Redis.
GET /api/accounts — Список всех счетов.
5.2 Транзакции (Асинхронные)
POST /api/accounts/{id}/deposit — Тело: {amount}
POST /api/accounts/{id}/withdraw — Тело: {amount}
POST /api/accounts/{id}/transfer — Тело: {amount, to_account_id}
GET /api/accounts/{id}/transactions — История операций.

6. RabbitMQ и Обработка событий
6.1 Producers (В Account Service)
При получении POST-запроса на транзакцию, сервис не пишет в БД сразу, а отправляет сообщение в RabbitMQ.
Топики (Queues):
account-deposit
account-withdraw
account-transfer
Формат сообщения (JSON):
JSON
{
  "transaction_id": "uuid",
  "account_id": 1,
  "to_account_id": 2, 
  "amount": 1000,
  "type": "transfer",
  "timestamp": "..."
}


6.2 Consumers (В Transaction Service)
Сервис подписывается на очереди и обрабатывает логику:
Deposit Consumer:
Создает запись в таблице transactions.
Через gRPC вызывает Account Service для обновления баланса (UpdateBalance).
Обновляет кэш Redis.
Withdraw Consumer:
Через gRPC проверяет баланс в Account Service.
Если средств достаточно -> Снимает средства (gRPC) -> Пишет транзакцию -> Обновляет Redis.
Если недостаточно -> Логирует ошибку (или отправляет в очередь DLQ).
Transfer Consumer:
Использует транзакцию БД (Saga или Two-phase commit упрощенно):
Снять с User A (gRPC) -> Зачислить User B (gRPC) -> Записать транзакцию.

7. gRPC Коммуникация
Определить .proto файл для взаимодействия между Transaction Service (Client) и Account Service (Server).
Методы gRPC:
GetAccountBalance(AccountID) — возвращает текущий баланс.
UpdateAccountBalance(AccountID, Amount, OperationType) — атомарно обновляет баланс (+ или -).
CheckAccountExists(AccountID) — валидация существования.

8. Кэширование (Redis)
8.1 Стратегия
TTL: 24 часа для всех ключей.
Паттерн: Write-through (обновление кэша при записи в БД) или Cache-aside.
8.2 Структура данных
Аккаунт:
Key: account:{id}
Value: Hash { "id": 1, "balance": 1000, "currency": "USD" }
История:
Key: transactions:{account_id}
Value: List (JSON objects) — хранить последние N операций или все за сутки.

9. Требования к реализации
Конфигурация: Использовать config.yaml для настройки портов, хостов БД/Redis/RabbitMQ.
Docker: Написать docker-compose.yml, поднимающий:
Postgres
Redis
RabbitMQ
Account Service
Transaction Service
Graceful Shutdown: Корректное завершение работы консюмеров и закрытие соединений.
Код: Чистая архитектура (Clean Architecture) или слоистая (Handler -> Service -> Repository).
10. План сдачи (Definition of Done)
Запущен docker-compose up, все контейнеры работают.
Применены миграции БД.
Через Postman/cURL создан аккаунт.
Выполнен депозит -> Проверка: баланс в Redis и БД увеличился.
Выполнен перевод -> Проверка: у отправителя убыло, у получателя прибыло.
В логах RabbitMQ видно прохождение сообщений.
