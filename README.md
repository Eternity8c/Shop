# Shop

Микросервисный бэкенд интернет-магазина на Go. Включает аутентификацию, каталог товаров, корзину и API Gateway как единую точку входа для клиентов.

## Возможности

- Регистрация и вход пользователей, получение JWT токена
- Просмотр и поиск товаров в каталоге
- Добавление товаров в личную корзину
- Ролевой доступ: администраторы могут создавать товары
- Хранение корзины в Redis (TTL 24ч)
- Метрики Prometheus на стороне API Gateway
- Централизованная агрегация логов через Loki + Promtail

## Установка

Необходимы **Go 1.21+** и **Docker Engine 24+** с Compose v2.

### API Gateway

```bash
cd api-gateway
make run
```

Запускает шлюз на порту `8080`. Подключается к `auth-service:44044` и `product-service:8081`.

### Auth Service

```bash
cd auth-service
make run
```

Запускает gRPC сервер на порту `44044`. Требует запущенный PostgreSQL.

### Product Service

```bash
cd product-service
make run
```

Запускает REST сервер на порту `8081`. Требует запущенный PostgreSQL.

### Cart Service

```bash
cd cart-service
make run
```

Запускает REST сервер на порту `8082`. Требует запущенный Redis.

> **Важно:** при локальном запуске в `configs/.env` каждого сервиса укажите `localhost` в качестве хоста. В Docker используйте имя контейнера (`postgres`, `redis`).

Каждый сервис читает конфигурацию через переменную окружения `CONFIG_PATH`, которая должна указывать на `.env` файл внутри `configs/`.

Примеры конфигов:

**api-gateway/configs/.env**
```env
GATEWAY_PORT=8080
AUTH_SERVICES_ADDR=localhost:44044
PRODUCT_SERVICE_ADDR=localhost:8081
```

**auth-service/configs/.env**
```env
PORT=44044
STORAGE_PATH=postgres://admin:1234@localhost:5432/auth?sslmode=disable
TOKEN_TTL=1h
SECRET=your-jwt-secret-here
```

**product-service/configs/.env**
```env
PORT=8081
STORAGE_PATH=postgres://admin:1234@localhost:5432/product?sslmode=disable
```

**cart-service/configs/.env**
```env
PORT=8082
STORAGE_PATH=redis://localhost:6379/0
```

## Docker

Запустить все сервисы одной командой:

```bash
docker compose up --build
```

Поднимает все четыре сервиса вместе с PostgreSQL, Redis, Kafka и стеком мониторинга (Prometheus, Grafana, Loki, Promtail).

> Перед запуском в Docker замените `localhost` на имена контейнеров (`postgres`, `redis`) в каждом `configs/.env`.

## Использование

1. Зарегистрируйте пользователя через `POST /register`.
2. Войдите через `POST /login` — в ответе придёт JWT токен.
3. Передавайте токен в заголовке `Authorization: Bearer <token>` для защищённых маршрутов.
4. Просматривайте товары через `GET /products` или ищите по названию `GET /products/search?name=`.
5. Добавляйте товары в корзину через `POST /cart`.
6. Администраторы могут создавать товары через `POST /products`.

### API

```
POST   /register              — создать аккаунт (email, password, fullName)
POST   /login                 — войти, возвращает JWT
GET    /products              — список всех товаров
GET    /products/search?name= — поиск по названию
POST   /products              — создать товар (только admin)
POST   /cart                  — добавить товар в корзину (требует авторизации)
```

## Мониторинг

```bash
docker compose -f infra/monitoring.yml up
```

| Сервис | URL |
|---|---|
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |

## Нагрузочное тестирование

```bash
docker compose -f k6/k6.yml up
```
## Архитектура
![Архитектура](imgs/image.png)