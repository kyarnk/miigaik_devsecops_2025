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
