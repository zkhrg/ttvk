package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

type IPData struct {
	ID          string    `json:"id,omitempty"`
	IPAddress   string    `json:"ip_address"`
	PingTime    int       `json:"ping_time"`
	LastSuccess time.Time `json:"last_success"`
}

func getEnvURL() string {
	url := os.Getenv("API_URL")
	if url == "" {
		fmt.Println("API_URL environment variable is not set")
		os.Exit(1)
	}
	return url
}

func pingIPAddress(ip string) (int, error) {
	cmd := exec.Command("ping", "-c", "1", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Используем регулярное выражение для извлечения времени пинга
	re := regexp.MustCompile(`time[=](\d+)\.?(\d*)\s*ms`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return 0, fmt.Errorf("ping output does not contain time")
	}

	// Объединяем части времени (целая и дробная части) в одно число
	var pingTime int
	_, err = fmt.Sscanf(matches[1], "%d", &pingTime)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ping time: %v", err)
	}

	// Если есть дробная часть, можно немного уточнить значение
	if len(matches) > 2 {
		// Дополняем дробную часть (например, если вывод был "time=10.123 ms", то 123 миллисекунды)
		var fraction int
		_, err = fmt.Sscanf(matches[2], "%d", &fraction)
		if err == nil {
			pingTime = pingTime*1000 + fraction
		}
	}

	return pingTime, nil
}

func fetchIPData(url string) ([]IPData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []IPData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func updateIPData(url string, ipData IPData) error {
	client := &http.Client{}
	data, err := json.Marshal(ipData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update IP data: %s", resp.Status)
	}

	return nil
}

func processIPData(url string, ipData IPData, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	// Пингуем IP
	pingTime, err := pingIPAddress(ipData.IPAddress)
	if err != nil {
		fmt.Printf("Error pinging %s: %v\n", ipData.IPAddress, err)
		return
	}

	ipData.PingTime = pingTime
	ipData.LastSuccess = time.Now()

	// Отправляем обновленные данные на сервер
	err = updateIPData(url, ipData)
	if err != nil {
		fmt.Printf("Error updating IP data for %s: %v\n", ipData.IPAddress, err)
	} else {
		fmt.Printf("Successfully updated IP data for %s\n", ipData.IPAddress)
	}

	// Проверка на контекст (например, для graceful shutdown)
	select {
	case <-ctx.Done():
		fmt.Println("Received shutdown signal, terminating processing.")
		return
	default:
	}
}

func main() {
	url := getEnvURL()

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для получения сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Запускаем цикл для выполнения задачи до получения сигнала
	for {
		select {
		case <-sigChan:
			// Получили сигнал завершения (например, SIGINT)
			fmt.Println("Received shutdown signal, shutting down gracefully...")
			cancel()  // Отменяем все горутины
			wg.Wait() // Ожидаем завершения всех горутин
			return
		default:
			// Получаем данные с API
			ipDataList, err := fetchIPData(url)
			if err != nil {
				fmt.Printf("Error fetching IP data: %v\n", err)
				time.Sleep(5 * time.Second) // Задержка перед новой попыткой
				continue
			}

			// Запускаем горутины для каждого IP-адреса
			for _, ipData := range ipDataList {
				wg.Add(1)
				go processIPData(url, ipData, &wg, ctx)
			}

			// Пауза перед следующим циклом
			time.Sleep(10 * time.Second) // Ждем 10 секунд перед повторной проверкой
		}
	}
}
