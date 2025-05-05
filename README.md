# Инструкция по установке Go и запуску программы для тестирования нагрузки

## Шаг 1: Установка Go

### Для Для Linux

```
sudo apt update
sudo apt install golang-go
go version -- вывод go version go1.19.4 linux/amd64

```

## Шаг 2: Запуск
```
go run main.go <endpoint> <method> <json data> <ip_address> <requests_count> 

```

Параметры:

<endpoint> — путь до REST API, например, /api/v1/test.
<method> - метод запроса REST API.
<json data> - данные для запросаов
<ip_address> — IP-адрес или доменное имя вашего сервера.
<requests_count> (необязательный параметр) — количество HTTP-запросов, которые программа отправит (по умолчанию 1000).


Пример:
```
для ГЕТ

go run main.go /api/test GET 127.0.0.1:8080 100

для ПОСТ

go run main.go /api/test POST '{"name":"Alice"}' 127.0.0.1:8080 50

для ПУТ

go run main.go /api/test PUT '{"id":1,"value":"updated"}' 127.0.0.1:8080 10

для ДЕЛИТ

go run main.go /api/test DELETE 127.0.0.1:8080 5
```


Пример вывода:
```
Ответ от http://192.168.1.100/api/v1/test: 200
Ответ от http://192.168.1.100/api/v1/test: 200
...
Время выполнения: 5.625
```

# Нововведения

## 1. Простая нагрузка:
```
go run main.go --endpoint /tasks --method GET --ips 127.0.0.1 --requests 100

```

## 2. DDoS-стиль (1000 RPS на 3 IP):
```
go run main.go --endpoint /user/login --method POST --random --ips 192.168.0.10,192.168.0.11,192.168.0.12 --rps 1000 --duration 10 --log


```

## 3. Тест регистрации с реальными JSON-данными:
```
go run main.go --endpoint /user/register --method POST --body '{\"username\":\"test\",\"password\":\"pass\"}' --ips 127.0.0.1 --requests 500

```

```
В терминале: логи каждого запроса и время отклика.
В log.csv: журнал всех запросов с IP, статусами, задержкой.
``` 
