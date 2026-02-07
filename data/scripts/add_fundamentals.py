import os
import psycopg2
import pandas as pd

# Database configuration
DB_CONFIG = {
    "dbname": "wallstreet",
    "user": "postgres",
    "password": "yash123",
    "host": "localhost",
    "port": 5432,
}

COMPANIES = ["airtel", "infosys", "itc", "relience-industries", "tcs"]
SYMBOL_MAP = {
    "airtel": "AIRTEL",
    "infosys": "INFOSYS",
    "itc": "ITC",
    "relience-industries": "RELIANCE",
    "tcs": "TCS",
}

def load_fundamentals():
    conn = psycopg2.connect(**DB_CONFIG)
    cur = conn.cursor()

    # Create staging table for fundamentals if not exists
    cur.execute("""
        CREATE TABLE IF NOT EXISTS fundamentals_temp (
            year DATE,
            sales NUMERIC,
            operating_profit NUMERIC,
            net_profit NUMERIC,
            eps_in_rs NUMERIC,
            equity_capital NUMERIC,
            reserves NUMERIC,
            borrowings NUMERIC,
            total_assets NUMERIC
        )
    """)

    for company in COMPANIES:
        path = f"../{company}/fundamentals.csv"
        if not os.path.exists(path):
            print(f"File not found: {path}")
            continue

        print(f"Inserting fundamentals for: {company}")
        
        with open(path, "r") as f:
            # COPY into temp table
            cur.copy_expert(
                "COPY fundamentals_temp(year, sales, operating_profit, net_profit, eps_in_rs, equity_capital, reserves, borrowings, total_assets) FROM STDIN WITH CSV HEADER",
                f
            )

        # Move from temp to main table
        cur.execute("""
            INSERT INTO fundamentals (company_id, year, sales, operating_profit, net_profit, eps, equity_capital, reserves, borrowings, total_assets)
            SELECT
                c.id,
                f.year,
                f.sales,
                f.operating_profit,
                f.net_profit,
                f.eps_in_rs,
                f.equity_capital,
                f.reserves,
                f.borrowings,
                f.total_assets
            FROM fundamentals_temp f
            JOIN companies c ON c.symbol = %s
            ON CONFLICT (company_id, year) DO UPDATE SET
                sales = EXCLUDED.sales,
                operating_profit = EXCLUDED.operating_profit,
                net_profit = EXCLUDED.net_profit,
                eps = EXCLUDED.eps,
                equity_capital = EXCLUDED.equity_capital,
                reserves = EXCLUDED.reserves,
                borrowings = EXCLUDED.borrowings,
                total_assets = EXCLUDED.total_assets
        """, (SYMBOL_MAP[company],))
        
        cur.execute("TRUNCATE fundamentals_temp")
        print(f"  > Done.")

    conn.commit()
    cur.close()
    conn.close()
    print("\nFundamentals insertion complete.")

if __name__ == "__main__":
    load_fundamentals()
