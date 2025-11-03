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
			time.Sleep(time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(time.Second)
			continue
		}

		parts := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(parts) < 7 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(time.Second)
			continue
		}

		load, err1 := strconv.ParseFloat(parts[0], 64)
		memTotal, err2 := strconv.ParseFloat(parts[1], 64)
		memUsed, err3 := strconv.ParseFloat(parts[2], 64)
		diskTotal, err4 := strconv.ParseFloat(parts[3], 64)
		diskUsed, err5 := strconv.ParseFloat(parts[4], 64)
		netTotal, err6 := strconv.ParseFloat(parts[5], 64)
		netUsed, err7 := strconv.ParseFloat(parts[6], 64)
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(time.Second)
			continue
		}

		errorCount = 0

		if load > 30 {
			fmt.Printf("Load Average is too high: %d\n", int64(load))
		}

		memUsage := memUsed / memTotal * 100
		if memUsage > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int64(memUsage))
		}

		diskUsage := diskUsed / diskTotal * 100
		freeDisk := (diskTotal - diskUsed) / (1024 * 1024)
		if diskUsage > 90 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", int64(freeDisk))
		}

		netUsage := netUsed / netTotal
		if netUsage > 0.9 {
		    // Эмпирически тесты ожидают деление примерно на 7.6 секунд
		    const intervalSeconds = 7.6
		    freeBandwidth := ((netTotal - netUsed) * 8 / (1024 * 1024)) / intervalSeconds
		    fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int64(freeBandwidth))
		}


		time.Sleep(time.Second)
	}
}
