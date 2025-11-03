package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	const url = "http://srv.msk01.gigacorp.local/_stats"

	errorCount := 0 // счётчик неудачных попыток

	for {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic.")
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		// Читаем тело ответа
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic.")
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		// Разделяем строку по запятым
		fields := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(fields) != 7 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic.")
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		// Сбрасываем счётчик ошибок при успешном парсинге
		errorCount = 0

		// Преобразуем значения
		values := make([]float64, 7)
		for i, f := range fields {
			v, err := strconv.ParseFloat(f, 64)
			if err != nil {
				errorCount++
				if errorCount >= 3 {
					fmt.Println("Unable to fetch server statistic.")
					return
				}
				continue
			}
			values[i] = v
		}

		load := values[0]
		memTotal := values[1]
		memUsed := values[2]
		diskTotal := values[3]
		diskUsed := values[4]
		netTotal := values[5]
		netUsed := values[6]

		// Проверяем пороги
		if load > 30 {
			fmt.Printf("Load Average is too high: %.0f\n", load)
		}

		if memTotal > 0 {
			memUsage := (memUsed / memTotal) * 100
			if memUsage > 80 {
				fmt.Printf("Memory usage too high: %.0f%%\n", memUsage)
			}
		}

		if diskTotal > 0 {
			diskFree := diskTotal - diskUsed
			if diskUsed/diskTotal > 0.9 {
				mbLeft := diskFree / 1024 / 1024
				fmt.Printf("Free disk space is too low: %.0f Mb left\n", mbLeft)
			}
		}

		if netTotal > 0 {
			netUsage := netUsed / netTotal
			if netUsage > 0.9 {
				// свободная полоса в мегабитах в секунду
				freeMbit := (netTotal - netUsed) * 8 / 1024 / 1024
				fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeMbit)
			}
		}

		time.Sleep(1 * time.Second) // пауза между запросами
	}
}
