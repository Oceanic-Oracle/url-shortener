## Запуск
Убедитесь, что Docker запущен, затем соберите и запустите сервис:
```bash
docker compose up --build
```
После запуска:
- Сервис будет доступен на http://localhost:8080

## Тесты
Тесты находятся в папке [test](url-shortener/test/) \
Запуск:
```bash
make e2e
```

## Нагрузочное тестирование
Тестирование проводилось с помощью `hey` при следующих условиях:
- 1 контейнер 
- 30 параллельных соединений
- 10 минут нагрузки

#### 1. Создание короткой ссылки
```bash
POST /shorten
```
<div>
  <img src="images/create.png" alt="test" style="width:60%; height: auto;">
</div>

#### 2. Редирект по короткому коду
```bash
GET /{code}
```
<div>
  <img src="images/redirect.png" alt="test" style="width:60%; height: auto;">
</div>

## Линтер
В проекте используется `golangci-lint` для статического анализа кода. \
Конфигурация: [.golangci.yml](pr/.golangci.yml).\
Использовать можно через: [lint.go](url-shortener/lint.go) \
Запуск:
```bash
make lint
```