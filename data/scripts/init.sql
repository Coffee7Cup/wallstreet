-- Initial Schema for WallStreet project

-- 1. Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    cash_balance NUMERIC(15, 2) DEFAULT 1000000.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_cash CHECK (cash_balance >= 0)
);

-- 2. Admins table
CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Companies table
CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    sector VARCHAR(100),
    total_shares BIGINT NOT NULL
);

-- 4. Stock Prices table
CREATE TABLE IF NOT EXISTS stock_prices (
    id SERIAL PRIMARY KEY,
    company_id INT NOT NULL REFERENCES companies(id),
    date DATE NOT NULL,
    open_price NUMERIC(15, 2) NOT NULL,
    close_price NUMERIC(15, 2) NOT NULL,
    high_price NUMERIC(15, 2) NOT NULL,
    low_price NUMERIC(15, 2) NOT NULL,
    no_of_shares BIGINT NOT NULL,
    no_of_trades INT NOT NULL,
    total_turnover NUMERIC(15, 2) NOT NULL,
    UNIQUE(company_id, date)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_stock_prices_date ON stock_prices(date);
CREATE INDEX IF NOT EXISTS idx_stock_prices_company_date ON stock_prices(company_id, date);

-- 5. Ratios table
CREATE TABLE IF NOT EXISTS ratios (
    id SERIAL PRIMARY KEY,
    company_id INT NOT NULL REFERENCES companies(id),
    year DATE NOT NULL,
    roe NUMERIC(10, 4),
    debt_equity NUMERIC(10, 4),
    opm NUMERIC(10, 4),
    intrinsic_value NUMERIC(15, 2),
    UNIQUE(company_id, year)
);

-- 6. Fundamentals table
CREATE TABLE IF NOT EXISTS fundamentals (
    id SERIAL PRIMARY KEY,
    company_id INT NOT NULL REFERENCES companies(id),
    year DATE NOT NULL,
    sales NUMERIC(15, 2),
    operating_profit NUMERIC(15, 2),
    net_profit NUMERIC(15, 2),
    eps NUMERIC(10, 4),
    equity_capital NUMERIC(15, 2),
    reserves NUMERIC(15, 2),
    borrowings NUMERIC(15, 2),
    total_assets NUMERIC(15, 2),
    UNIQUE(company_id, year)
);

-- 7. Trades table
CREATE TABLE IF NOT EXISTS trades (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    company_id INT NOT NULL REFERENCES companies(id),
    trade_type VARCHAR(10) NOT NULL, -- 'BUY' or 'SELL'
    quantity INT NOT NULL,
    date DATE NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trades_user_id ON trades(user_id);

-- 8. Portfolio Entries table
CREATE TABLE IF NOT EXISTS portfolio_entries (
    user_id INT NOT NULL REFERENCES users(id),
    company_id INT NOT NULL REFERENCES companies(id),
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY(user_id, company_id),
    CONSTRAINT positive_quantity CHECK (quantity >= 0)
);

-- 9. News table
CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    release_date DATE NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    news_type VARCHAR(50),
    impact_factor NUMERIC(5, 2),
    tick INT -- Manual tick override
);

CREATE INDEX IF NOT EXISTS idx_news_date ON news(release_date);

-- 10. Simulation State table
CREATE TABLE IF NOT EXISTS simulation_state (
    id SERIAL PRIMARY KEY,
    tick INT DEFAULT 0,
    is_active BOOLEAN DEFAULT FALSE,
    is_paused BOOLEAN DEFAULT FALSE,
    start_time TIMESTAMP NOT NULL,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT one_row CHECK (id = 1)
);

-- Compatibility views (optional but helpful)
CREATE OR REPLACE VIEW prices AS SELECT * FROM stock_prices;
CREATE OR REPLACE VIEW trade AS SELECT * FROM trades;
CREATE OR REPLACE VIEW portfolio AS SELECT * FROM portfolio_entries;
CREATE OR REPLACE VIEW stock_prices_with_ticks AS
SELECT *, ROW_NUMBER() OVER (PARTITION BY company_id ORDER BY date) - 1 as tick_idx
FROM stock_prices;

CREATE OR REPLACE VIEW news_with_ticks AS
WITH sim_start AS (
    SELECT COALESCE(MIN(date), CURRENT_DATE) as start_date FROM stock_prices
)
SELECT n.*, 
       COALESCE(n.tick, (n.release_date - s.start_date)) as tick_idx
FROM news n, sim_start s;

-- Staging tables for scripts (must be regular tables to use with COPY in some drivers)
CREATE TABLE IF NOT EXISTS prices_temp (
    date DATE,
    open_price NUMERIC,
    close_price NUMERIC,
    low_price NUMERIC,
    high_price NUMERIC,
    no_of_shares BIGINT,
    no_of_trades INT,
    total_turnover NUMERIC
);

CREATE TABLE IF NOT EXISTS ratios_temp (
    year DATE,
    roe NUMERIC,
    debt_equity NUMERIC,
    opm NUMERIC,
    intrinsic_value NUMERIC
);

-- Seed basic company data
INSERT INTO companies (symbol, name, sector, total_shares) VALUES
('AIRTEL', 'Bharti Airtel', 'TELECOM', 1000000000),
('INFOSYS', 'Infosys Limited', 'IT', 1000000000),
('ITC', 'ITC Limited', 'FMCG', 1000000000),
('RELIANCE', 'Reliance Industries', 'CONGLOMERATE', 1000000000),
('TCS', 'Tata Consultancy Services', 'IT', 1000000000)
ON CONFLICT (symbol) DO NOTHING;
