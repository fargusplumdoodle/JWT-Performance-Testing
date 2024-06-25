package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	secretKey = "I think maybe python is over"
	role      = "sub"
)

func generateAndVerifyJWT() error {
	sub := "1234567890"
	claims := jwt.MapClaims{
		"sub":  sub,
		"exp":  time.Now().Add(10 * time.Minute).Unix(),
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return err
}

func runTest(count, parallelism int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, parallelism)

	for i := 0; i < count; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()
			err := generateAndVerifyJWT()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}()
	}

	wg.Wait()
	return time.Since(start)
}

func main() {
	count := 100000
	parallelisms := []int{1, 2, 4, 8, 16, 32, 64, 128}

	for _, p := range parallelisms {
		duration := runTest(count, p)
		fmt.Printf("Parallelism: %d, Time: %v\n", p, duration)
	}
}
