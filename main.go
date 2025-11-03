package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	serverURL = "http://srv.msk01.gigacorp.local/_stats"

	loadAvgThreshold      = 30.0
	memUsageThreshold     = 0.8
	diskUsageThreshold    = 0.9
	networkUsageThreshold = 0.9
)

func main() {
	errorCount := 0

	for {
		resp, err := http.Get(serverURL)
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic.")
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic.")
				return
			}
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}

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

		errorCount = 0

		values := make([]float64, len(fields))
		for i, f := range fields {
			v, err := strconv.ParseFloat(strings.TrimSpace(f), 64)
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

		loadAvg := values[0]
		totalMem := values[1]
		usedMem := values[2]
		totalDisk := values[3]
		usedDisk := values[4]
		totalNet := values[5]
		usedNet := values[6]

		if loadAvg > loadAvgThreshold {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}

		if totalMem > 0 {
			memUsage := usedMem / totalMem
			if memUsage > memUsageThreshold {
				fmt.Printf("Memory usage too high: %.0f%%\n", memUsage*100)
			}
		}

		if totalDisk > 0 {
			diskUsage := usedDisk / totalDisk
			if diskUsage > diskUsageThreshold {
				freeBytes := totalDisk - usedDisk
				freeMb := int(freeBytes / 1024 / 1024)
				fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
			}
		}

		if totalNet > 0 {
			netUsage := usedNet / totalNet
			if netUsage > networkUsageThreshold {
				// ✅ без умножения на 8, округляем вниз
				freeMbitPerSec := int((totalNet - usedNet) / 1_000_000)
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbitPerSec)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
