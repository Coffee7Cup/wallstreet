# WallStreet Simulation

A real-time stock market simulation with a Go backend and a SvelteKit frontend.

## Project Structure

This project is a monorepo consisting of:

- `backend/`: Go (Fiber) backend and simulation engine.
- `frontend/`: SvelteKit frontend with real-time updates via WebSockets.

## Getting Started

### Prerequisites

- Go 1.25.5
- Node.js 24.10.0
- PostgreSQL 16
- python 

### SetUp
1. You can use the sample data already give or collect data by yourself and put it in the /backend/data/company-data folder - each folder can have fundamentals.csv, prices.csv, ratios.csv, news.csv

#### You need the following packeges installed in your python
```bash
pip install psycopg2
pip install python-dotenv
```

2. Run the /backend/data/scripts/init.sql ot copy paste it

3. Run the following commands to populate the database
```bash
cd backend/data/scripts
python add_fundamentals.py
python add_ratios.py
python add_stocks.py
python add_users_admins.py

```

### Database

```bash
# Create a new database
CREATE DATABASE wallstreet;

# Connect to the database
\c wallstreet;

# Run the migration

```

### Backend

1. Navigate to the `backend` directory:
   ```bash
   cd backend
   ```
2. Copy the sample environment file and configure your database:
   ```bash
   cp .env.sample .env
   ```
3. Run the backend server:
   ```bash
   go run cmd/server/main.go
   ```

### Frontend

1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the development server:
   ```bash
   npm run dev
   ```
4. Open your browser to `http://localhost:5173`.

## Features

- Real-time stock price updates via WebSockets.
- Admin panel for simulation control (Start, Stop, Pause, Resume).
- Detailed company fundamentals and performance ratios.
- Responsive dashboard with price history graphs.
- Structured logging for system monitoring.
