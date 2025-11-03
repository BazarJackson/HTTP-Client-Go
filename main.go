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

	errorCount := 0
	for {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		parts := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(parts) != 7 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		errorCount = 0 // успешное чтение — сбрасываем счётчик ошибок

		load, _ := strconv.ParseFloat(parts[0], 64)
		memTotal, _ := strconv.ParseFloat(parts[1], 64)
		memUsed, _ := strconv.ParseFloat(parts[2], 64)
		diskTotal, _ := strconv.ParseFloat(parts[3], 64)
		diskUsed, _ := strconv.ParseFloat(parts[4], 64)
		netTotal, _ := strconv.ParseFloat(parts[5], 64)
		netUsed, _ := strconv.ParseFloat(parts[6], 64)

		// Load Average
		if load > 30 {
			fmt.Printf("Load Average is too high: %.0f\n", load)
		}

		// Memory usage
		memUsage := memUsed / memTotal * 100
		if memUsage > 80 {
			fmt.Printf("Memory usage too high: %.0f%%\n", memUsage)
		}

		// Disk usage
		freeDisk := (diskTotal - diskUsed) / (1024 * 1024)
		diskUsage := diskUsed / diskTotal * 100
		if diskUsage > 90 {
			fmt.Printf("Free disk space is too low: %.0f Mb left\n", freeDisk)
		}

		// Network usage
		netUsage := netUsed / netTotal
		if netUsage > 0.9 {
			freeBandwidth := (netTotal - netUsed) * 8 / (1024 * 1024) // в мегабитах
			fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeBandwidth)
		}

		time.Sleep(1 * time.Second)
	}
}
