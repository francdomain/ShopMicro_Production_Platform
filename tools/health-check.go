package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkEndpoint(url string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}

func main() {
	endpoints := []string{
		"http://localhost:3001/health",
		"http://localhost:8080",
		"http://localhost:3002/health",
		"http://localhost:3000/login",
	}

	fmt.Println("ShopMicro Health Check")
	fmt.Println("======================")

	allHealthy := true
	for _, endpoint := range endpoints {
		fmt.Printf("Checking %s... ", endpoint)
		if err := checkEndpoint(endpoint); err != nil {
			fmt.Printf("FAIL: %v\n", err)
			allHealthy = false
		} else {
			fmt.Println("OK")
		}
	}

	if allHealthy {
		fmt.Println("\nAll services are healthy!")
		os.Exit(0)
	} else {
		fmt.Println("\nSome services are unhealthy!")
		os.Exit(1)
	}
}
