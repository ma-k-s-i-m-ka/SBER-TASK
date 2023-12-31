# SBER-TASK
# Демонстрация работы приложения
[Пример работы и запуска приложения](https://drive.google.com/file/d/1fzC7PREYBYbiMnP6tAz34D3apPQg9pEe/view?usp=sharing)
# Инструкция к запуску приложения 
Изначально приложение запускается на http://localhost:3003/ и представляет из себя вебсайт с возможность просмотра задач

Документация Swagger UI открывается по адресу http://localhost:3003/docs/

Измените конфигурацию приложения для запуска сервера и подключения к бд в следующих файлах:

[.env]
```sh
ENV=dev
USE_HTTP=true
DATABASE_DSN= your postgres patch
```
[config.yml]
```sh
http:
  host:            0.0.0.0    # your host
  port:            3003       # your port
  read_timeout:    30         # Seconds
  write_timeout:   30         # Seconds
```
Таблица с который работает приложение находится в файле [db.sql]. Перед запуском приложения ее необходимо создать заранее
```sh
DROP TABLE IF EXISTS Task;

CREATE TABLE IF NOT EXISTS Task (
id              serial       primary key,
title           text         not null,
description     text         not null,
date            timestamptz  not null,
status          bool         not null
);
```
Приложение готово к запуску! Команда "go run main.go" и вперед!

Онсновные команды для работы с приложением прописаны в Makefile:
```sh
test:
  go test ./app/internal/test 
  
run:
  go run main.go
  
docker-build:
  docker build -t sber_task:local .
  
docker-compose-up:
  docker compose -f docker-compose.yml up
```
# Инструкция к запуску приложения через Docker
В выше упомянутом [.env] файле необходимо изменить host на host.docker.internal
```sh
ENV=dev
USE_HTTP=true
DATABASE_DSN=postgres://youruser:yourpassword@host.docker.internal:5432/yourdatabase
```
Если была изменена конфигурация приложения в файле [config.yml], то необходимо прокинуть новые порты в файле [docker-compose.yml]:
```sh
services:
  app:
    image: sber_task:local
    container_name: sber
    ports:
    - "3003:3003" (your ports)
```
Приложение готово к запуску через Docker!
Выполнить следующие команды и вперед!
```sh
docker build -t sber_task:local .
docker compose -f docker-compose.yml up
 ```
# Тестирование
Протестировать все запрос можно через Swagger UI по адресу http://localhost:3003/docs/

Или воспользоваться POSTMAN с моей [коллекцией запросов](https://drive.google.com/file/d/1gz1qN-e3l_tvE_8Ucw4n_ThvseVUXpZl/view?usp=sharing)
