package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type AttackConfig struct {
	IPs        []string
	Endpoint   string
	Method     string
	Body       string
	Requests   int
	RPS        int
	Duration   time.Duration
	RandomBody bool
	LogToFile  bool
}

var logMutex sync.Mutex

func generateRandomBody() string {
	username := fmt.Sprintf("user_%d", rand.Intn(10000))
	password := fmt.Sprintf("pass_%d", rand.Intn(10000))
	body := map[string]string{
		"username": username,
		"password": password,
	}
	b, _ := json.Marshal(body)
	return string(b)
}

func sendRequest(cfg AttackConfig, ip string, wg *sync.WaitGroup, logCh chan<- []string) {
	defer wg.Done()

	body := cfg.Body
	if cfg.RandomBody {
		body = generateRandomBody()
	}

	url := fmt.Sprintf("http://%s%s", ip, cfg.Endpoint)
	var req *http.Request
	var err error

	if cfg.Method == "POST" || cfg.Method == "PUT" {
		req, err = http.NewRequest(cfg.Method, url, bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(cfg.Method, url, nil)
	}

	if err != nil {
		fmt.Printf("[ERR] Create req to %s: %v\n", url, err)
		return
	}

	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		fmt.Printf("[ERR] Send req to %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[%s] %s -> %d (%v)\n", cfg.Method, url, resp.StatusCode, latency)

	if cfg.LogToFile {
		logCh <- []string{time.Now().Format(time.RFC3339), ip, cfg.Endpoint, cfg.Method, fmt.Sprintf("%d", resp.StatusCode), latency.String()}
	}
}

func startAttack(cfg AttackConfig) {
	logCh := make(chan []string, 1000)
	var wg sync.WaitGroup
	var ticker *time.Ticker
	limit := cfg.Requests

	if cfg.RPS > 0 {
		ticker = time.NewTicker(time.Second / time.Duration(cfg.RPS))
		defer ticker.Stop()
	}

	endTime := time.Now().Add(cfg.Duration)

	go func() {
		if cfg.LogToFile {
			f, err := os.Create("log.csv")
			if err != nil {
				fmt.Println("Ошибка создания log.csv:", err)
				return
			}
			defer f.Close()
			w := csv.NewWriter(f)
			defer w.Flush()
			w.Write([]string{"time", "ip", "endpoint", "method", "status", "latency"})
			for row := range logCh {
				w.Write(row)
			}
		}
	}()

	sent := 0
	for {
		if cfg.Duration > 0 && time.Now().After(endTime) {
			break
		}
		if limit > 0 && sent >= limit {
			break
		}

		if cfg.RPS > 0 {
			<-ticker.C
		}

		ip := cfg.IPs[rand.Intn(len(cfg.IPs))]
		wg.Add(1)
		go sendRequest(cfg, ip, &wg, logCh)
		sent++
	}

	wg.Wait()
	if cfg.LogToFile {
		close(logCh)
	}
	fmt.Println("\nГенерация трафика завершена.")
}

func main() {
	endpoint := flag.String("endpoint", "/", "Endpoint пути (например, /user/login)")
	method := flag.String("method", "GET", "HTTP-метод (GET, POST, PUT, DELETE)")
	body := flag.String("body", "", "JSON-данные запроса (для POST/PUT)")
	ips := flag.String("ips", "127.0.0.1", "Список IP через запятую")
	rps := flag.Int("rps", 0, "Запросов в секунду (0 = без ограничения)")
	requests := flag.Int("requests", 100, "Количество запросов")
	duration := flag.Int("duration", 0, "Длительность атаки в секундах (0 = по количеству)")
	random := flag.Bool("random", false, "Случайные JSON-данные (POST/PUT)")
	log := flag.Bool("log", false, "Вести лог в файл log.csv")

	flag.Parse()

	cfg := AttackConfig{
		IPs:        strings.Split(*ips, ","),
		Endpoint:   *endpoint,
		Method:     strings.ToUpper(*method),
		Body:       *body,
		Requests:   *requests,
		RPS:        *rps,
		Duration:   time.Duration(*duration) * time.Second,
		RandomBody: *random,
		LogToFile:  *log,
	}

	rand.Seed(time.Now().UnixNano())
	startAttack(cfg)
}


