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
# Исправления уязвимостей 

## Semgrep
```markdown
# Отчёт по исправлению уязвимостей из Semgrep

## 1. Уязвимость: dockerfile.security.missing-user.missing-user

### Проблема
**Описание**: Контейнер запускается от имени root, что нарушает принцип минимальных привилегий.

**Место в коде**:
```dockerfile
FROM alpine:latest
...
CMD ["./main"]  # Запуск от root
```

**Риски**:
- Полный контроль над контейнером при компрометации
- Возможность эскалации привилегий на хосте

### Исправление
```dockerfile
FROM alpine:latest

# Добавляем непривилегированного пользователя
RUN adduser -D -S -H -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser  # Переключаем пользователя
CMD ["./main"]
```

**Что изменилось**:
1. Создан отдельный пользователь `appuser` без прав root
2. Рекурсивно изменены владельцы файлов
3. Приложение гарантированно запускается с минимальными правами

## 2. Уязвимость: go.lang.security.audit.net.use-tls.use-tls

### Проблема
**Описание**: Используется незашифрованное HTTP-соединение.

**Место в коде**:
```go
http.ListenAndServe(":8080", nil)
```

**Риски**:
- Перехват данных (логины, пароли, токены)
- Возможность MITM-атак

### Исправление
```go
func main() {
    server := &http.Server{
        Addr:    ":8443",
        Handler: nil,
        TLSConfig: &tls.Config{
            MinVersion:   tls.VersionTLS12,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            },
        },
    }
    
    log.Fatal(server.ListenAndServeTLS(
        "cert.pem", 
        "key.pem",
    ))
}
```

**Дополнительные меры**:
1. Генерация тестового сертификата:
   ```bash
   openssl req -x509 -newkey rsa:4096 -nodes -keyout key.pem -out cert.pem -days 365
   ```
2. В продакшене - использование Let's Encrypt

## 3. Уязвимость: go.lang.security.audit.xss.no-direct-write-to-responsewriter

### Проблема
**Описание**: Прямой вывод в ResponseWriter без экранирования.

**Место в коде**:
```go
w.Write([]byte("<div>" + userInput + "</div>"))
```

**Риски**:
- Возможность XSS-атак
- Внедрение произвольного HTML/JS

### Исправление
```go
import "html/template"

func safeHandler(w http.ResponseWriter, r *http.Request) {
    const tpl = `<div>{{ . }}</div>`
    
    tmpl := template.Must(template.New("safe").Parse(tpl))
    tmpl.Execute(w, userInput)  // Автоматическое экранирование
}
```

**Что изменилось**:
1. Использование html/template вместо ручного форматирования
2. Автоматическое экранирование опасных символов

## 4. Уязвимость: go.lang.security.audit.xss.no-printf-in-responsewriter

### Проблема
**Описание**: Использование fmt.Fprintf для вывода пользовательских данных.

**Место в коде**:
```go
fmt.Fprintf(w, "Hello, %s!", r.FormValue("name"))
```

**Риски**:
- Обход механизмов экранирования
- Потенциальные XSS-уязвимости

### Исправление
```go
func safeGreeting(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.New("greet").Parse(
        `Hello, {{ . }}!`,
    ))
    
    if err := tmpl.Execute(w, r.FormValue("name")); err != nil {
        http.Error(w, "Render error", http.StatusInternalServerError)
    }
}
```
