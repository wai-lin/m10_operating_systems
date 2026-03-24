package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type quoteResp struct {
	GlobalQuote struct {
		Price string `json:"05. price"`
	} `json:"Global Quote"`
	Note string `json:"Note"`
}

func main() {
	key := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if key == "" {
		fmt.Println("Missing ALPHA_VANTAGE_API_KEY")
		return
	}

	db, err := initDB("nvda.db")
	if err != nil {
		fmt.Println("DB Error:", err)
		return
	}
	defer db.Close()

	for {
		if !inMarketHours(time.Now()) {
			if err := printEODIfNeeded(db); err != nil {
				fmt.Println("EOD error:", err)
			}
			fmt.Println("Market closed; Sleeping 5 minutes.")
			time.Sleep(5 * time.Minute)
			continue
		}

		priceStr, err := fetchPrice(key)
		if err != nil {
			fmt.Println("Error:", err)
			time.Sleep(60 * time.Second)
			continue
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Println("Parse Error:", err)
			time.Sleep(60 * time.Second)
			continue
		}

		date := tradeDate(time.Now())
		prev, hadPrev, err := upsertStats(db, date, price)
		if err != nil {
			fmt.Println("DB Update Error:", err)
			time.Sleep(60 * time.Second)
			continue
		}

		if hadPrev {
			fmt.Printf("NVDA price: %.4f (diff %+0.4f)\n", price, price-prev)
		} else {
			fmt.Printf("NVDA Price: %.4f (diff N/A)\n", price)
		}
		time.Sleep(60 * time.Second)
	}
}

func fetchPrice(key string) (string, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=NVDA&apikey=%s", key)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("HTTP Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	var qr quoteResp
	if err := json.NewDecoder(resp.Body).Decode(&qr); err != nil {
		fmt.Println("Decode Error:", err)
		return "", err
	}

	if qr.Note != "" {
		return "", fmt.Errorf("API Note: %s", qr.Note)
	}
	if qr.GlobalQuote.Price == "" {
		return "", fmt.Errorf("Missing price in response.")
	}
	return qr.GlobalQuote.Price, nil
}

func inMarketHours(now time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.Local // fallback to local timezone
	}
	t := now.In(loc)

	// Mon-Fri only
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return false
	}

	// 9:30-16:00
	h, m, _ := t.Clock()
	if h < 9 || h > 16 {
		return false
	}
	if h == 9 && m < 30 {
		return false
	}
	if h == 16 && m > 0 {
		return false
	}
	return true
}

func initDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	schema := `CREATE TABLE IF NOT EXISTS daily_stats (
	trade_date TEXT PRIMARY KEY,
	open REAL,
	close REAL,
	min REAL,
	max REAL,
	last REAL,
	updated_at TEXT,
	summary_printed INTEGER DEFAULT 0
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return db, nil
}

func tradeDate(now time.Time) string {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.Local
	}
	return now.In(loc).Format("2006-01-02")
}

func upsertStats(db *sql.DB, date string, price float64) (last float64, hadLast bool, err error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var open, min, max, prevLast float64
	row := tx.QueryRow("SELECT open, min, max, last FROM daily_stats WHERE trade_date = ?", date)
	switch errScan := row.Scan(&open, &min, &max, &prevLast); errScan {
	case sql.ErrNoRows:
		// First tick of the day
		_, err = tx.Exec(`INSERT INTO daily_stats
		(trade_date, open, close, min, max, last, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'))`,
			date, price, price, price, price, price)
		if err != nil {
			return 0, false, err
		}
		if err = tx.Commit(); err != nil {
			return 0, false, err
		}
		return 0, false, nil
	case nil:
		// Update
		newMin := min
		if price < min {
			newMin = price
		}
		newMax := max
		if price > max {
			newMax = price
		}
		_, err = tx.Exec(`UPDATE daily_stats SET
		close = ?, min = ?, max = ?, last = ?, updated_at = datetime('now')
		WHERE trade_date = ?`,
			price, newMin, newMax, price, date)
		if err != nil {
			return 0, false, err
		}
		if err = tx.Commit(); err != nil {
			return 0, false, err
		}
		return prevLast, true, nil
	default:
		return 0, false, errScan
	}
}

func printEODIfNeeded(db *sql.DB) error {
	var date string
	var open, close, min, max float64

	row := db.QueryRow(`SELECT trade_date, open, close, min, max
	FROM daily_stats
	WHERE summary_printed = 0
	ORDER BY trade_date DESC
	LIMIT 1`)
	if err := row.Scan(&date, &open, &close, &min, &max); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	fmt.Printf("EOD %s: open %.4f close %.4f min %.4f max %.4f\n",
		date, open, close, min, max)

	_, err := db.Exec(`UPDATE daily_stats SET summary_printed = 1 WHERE trade_date = ?`, date)
	return err
}
