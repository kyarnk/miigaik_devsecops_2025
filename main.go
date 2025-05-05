package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func sendRequest(endpoint, method, jsonData, ip string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("http://%s%s", ip, endpoint)
	var req *http.Request
	var err error

	// Если метод требует тело — используем его
	if method == "POST" || method == "PUT" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(jsonData)))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		fmt.Printf("Ошибка при создании запроса на %s: %v\n", url, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при отправке запроса на %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Ответ от %s [%s]: %d\n", url, method, resp.StatusCode)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Использование:")
		fmt.Println("  GET/DELETE: go run main.go <endpoint> <method> <ip_address> [requests_count]")
		fmt.Println("  POST/PUT:   go run main.go <endpoint> <method> <json_data> <ip_address> [requests_count]")
		return
	}

	endpoint := os.Args[1]
	method := strings.ToUpper(os.Args[2])
	var jsonData, ip string
	requestsCount := 1000

	switch method {
	case "GET", "DELETE":
		if len(os.Args) < 4 {
			fmt.Println("Недостаточно аргументов для метода", method)
			return
		}
		ip = os.Args[3]
		if len(os.Args) >= 5 {
			fmt.Sscanf(os.Args[4], "%d", &requestsCount)
		}
	case "POST", "PUT":
		if len(os.Args) < 5 {
			fmt.Println("Недостаточно аргументов для метода", method)
			return
		}
		jsonData = os.Args[3]
		ip = os.Args[4]
		if len(os.Args) >= 6 {
			fmt.Sscanf(os.Args[5], "%d", &requestsCount)
		}
	default:
		fmt.Println("Неподдерживаемый метод:", method)
		return
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < requestsCount; i++ {
		wg.Add(1)
		go sendRequest(endpoint, method, jsonData, ip, &wg)
	}

	wg.Wait()
	elapsedTime := time.Since(startTime)
	fmt.Printf("Время выполнения: %v\n", elapsedTime)
}
