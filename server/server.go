package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Currency struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type CurrencyQuote struct {
	USDBRL Currency `json:"USDBRL"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type Response struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func sendJSONErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorMessage := ErrorMessage{Message: message}
	json.NewEncoder(w).Encode(errorMessage)
}

func successResponse(w http.ResponseWriter, rate CurrencyQuote) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{
		Bid: rate.USDBRL.Bid,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to encode response: ", err)
		sendJSONErrorMessage(w, http.StatusInternalServerError, "Failed to process the request due to an internal error.")
		return
	}
	log.Println("Request processed successfully")
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond) //time.Microsecond)
	defer cancel()

	log.Println("Request started")
	defer log.Println("Request finished")

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			sendJSONErrorMessage(w, http.StatusRequestTimeout, "Request timeout, please try again.")
			log.Println("Request timeout")
		} else {
			log.Println("Request cancelled by the client")
		}
		return
	default:
		LatestUSDBRLRate, error := getLatestUSDBRLRate()
		if error != nil {
			sendJSONErrorMessage(w, http.StatusInternalServerError, "Failed to get latest USD to BRL rate.")
			return
		}

		err := insertCurrencyIntoDB(LatestUSDBRLRate)
		if err != nil {
			sendJSONErrorMessage(w, http.StatusInternalServerError, "Failed to insert currency into DB.")
			return
		}
		successResponse(w, *LatestUSDBRLRate)
	}
}

func getLatestUSDBRLRate() (*CurrencyQuote, error) {
	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c CurrencyQuote
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}

func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data/rates.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS rates (
		code TEXT NOT NULL,
		codeIn TEXT NOT NULL,
		name TEXT NOT NULL,
		high TEXT NOT NULL,
		low TEXT NOT NULL,
		varBid TEXT NOT NULL,
		pctChange TEXT NOT NULL,
		bid TEXT NOT NULL,
		ask TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		createDate TEXT NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

func insertCurrencyIntoDB(currency *CurrencyQuote) error {
	db, err := dbConnect()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		return err
	}

	timestamp := int64(0)
	if currency.USDBRL.Timestamp != "" {
		timestamp, err = strconv.ParseInt(currency.USDBRL.Timestamp, 10, 64)
		if err != nil {
			log.Fatalf("Failed to parse timestamp: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO rates(code, codeIn, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		currency.USDBRL.Code,
		currency.USDBRL.CodeIn,
		currency.USDBRL.Name,
		currency.USDBRL.High,
		currency.USDBRL.Low,
		currency.USDBRL.VarBid,
		currency.USDBRL.PctChange,
		currency.USDBRL.Bid,
		currency.USDBRL.Ask,
		timestamp,
		currency.USDBRL.CreateDate)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
