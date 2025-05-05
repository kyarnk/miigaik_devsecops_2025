# Без докера 

```
запуск одной командой 
go run main.go

и на http://localhost:8080
будет фронт

```

# С докером

```
docker build -t task-manager .
docker run -d -p 8080:8080 --name task-manager-container task-manager

```
