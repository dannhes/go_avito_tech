Сервис для назначений Pr`ов
Инструкция к использованию приложения:
1) Склонировать репозиторий
```
git clone git@github.com:dannhes/go_avito_tech.git
```
2) Зайти в Docker
3) Поднять приложение в докере
   ```
   docker-compose up —build  
   ```
4) Создать еще один терминал для команд (Post, Get и тд)
5) Для запуска тестов
   ```
   docker-compose -f docker-compose.test.yaml up --build -d 
   go test ./tests/e2e -v
   docker-compose -f docker-compose.test.yaml down -v  
   ```
6) Для запуска линтера (должна быть версия 1.25
   ```
   golangci-lint run
   ```
