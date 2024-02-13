package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func createRequest(ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	return req
}

func retrieveResponse(req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	return res
}

func saveToFile(res *http.Response) {
	file, err := os.Create("data/cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var result map[string]string
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		panic(err)
	}

	formattedString := fmt.Sprintf("Dólar: %s", result["bid"])
	_, err = file.WriteString(formattedString)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, res.Body)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req := createRequest(ctx)
	res := retrieveResponse(req)
	saveToFile(res)
}
