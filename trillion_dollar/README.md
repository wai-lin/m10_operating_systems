# NVDA Watch

A small Go service that prints the NVDA price every minute during regular NASDAQ hours (9:30-16:00 US/Eastern), shows the delta from the previous tick, and prints end-of-day (EOD) summary stats. Daily stats persist in SQLite so restarts do not lose the EOD summary.

## Requirements

- Go (1.20+ recommended)
- Alpha Vantage API key
- Linux with systemd

## Run Locally

Set your API key:

Create free API Key at [alpha vantage](https://www.alphavantage.co/support/#api-key).

```bash
export ALPHA_VANTAGE_API_KEY="your_key_here"
```

Run:

```bash
go run .
```

## Build

```bash
go build -o nvda-watch .
```

## Install as a systemd service

1) Install binary and working directory:

```bash
sudo mkdir -p /opt/nvda-watch
sudo cp nvda-watch /opt/nvda-watch/
```

2) Create environment file:

```bash
echo 'ALPHA_VANTAGE_API_KEY=your_key_here' | sudo tee /etc/nvda-watch.env
sudo chmod 600 /etc/nvda-watch.env
```

3) Install service unit:

```bash
sudo cp nvda-watch.service /etc/systemd/system/nvda-watch.service
sudo systemctl daemon-reload
sudo systemctl enable --now nvda-watch
```

4) View logs:

```bash
sudo journalctl -u nvda-watch -f
```

## EOD Summary and Restart Safety

- Each minute tick is stored in SQLite under a row keyed by trading date.
- The program updates open, close, min, max, and last price in that row.
- After market close (or after a restart when the market is closed), it prints the EOD summary for the most recent day that has not been printed and marks it as printed.
- Because the printed flag is stored in SQLite, the summary is printed exactly once even across restarts.

## Files

- main.go: main application logic
- nvda-watch.service: systemd unit file
- nvda.db: SQLite database created at runtime
