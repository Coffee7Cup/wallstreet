import os
from pathlib import Path
import psycopg2
from dotenv import load_dotenv

# ----------------------------
# Path & ENV setup
# ----------------------------
BASE_DIR = Path(__file__).resolve().parent.parent.parent
load_dotenv(BASE_DIR / ".env")

print(BASE_DIR)

DATABASE_URL = os.getenv("DB_URL")
if not DATABASE_URL:
    raise RuntimeError("DB_URL not found in .env")

DATA_DIR = BASE_DIR / "data/company-data"
# ----------------------------
# Loader
# ----------------------------
def load_prices():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    # Ensure temp table exists
    cur.execute("""
        CREATE TABLE IF NOT EXISTS prices_temp (
            date DATE,
            open_price NUMERIC,
            close_price NUMERIC,
            low_price NUMERIC,
            high_price NUMERIC,
            no_of_shares BIGINT,
            no_of_trades INT,
            total_turnover NUMERIC
        )
    """)

    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        csv_path = company_dir / "prices.csv"
        if not csv_path.exists():
            continue

        company_symbol = company_dir.name.upper()
        print(f"Loading prices for: {company_symbol}")

        with open(csv_path, "r") as f:
            cur.copy_expert("""
                COPY prices_temp(
                    date,
                    open_price,
                    close_price,
                    low_price,
                    high_price,
                    no_of_shares,
                    no_of_trades,
                    total_turnover
                )
                FROM STDIN WITH CSV HEADER
            """, f)

        cur.execute("""
            INSERT INTO stock_prices (
                company_id,
                date,
                open_price,
                close_price,
                high_price,
                low_price,
                no_of_shares,
                no_of_trades,
                total_turnover
            )
            SELECT
                c.id,
                p.date,
                p.open_price,
                p.close_price,
                p.high_price,
                p.low_price,
                p.no_of_shares,
                p.no_of_trades,
                p.total_turnover
            FROM prices_temp p
            JOIN companies c
              ON c.symbol = %s
            ON CONFLICT (company_id, date)
            DO UPDATE SET
                open_price = EXCLUDED.open_price,
                close_price = EXCLUDED.close_price,
                high_price = EXCLUDED.high_price,
                low_price = EXCLUDED.low_price,
                no_of_shares = EXCLUDED.no_of_shares,
                no_of_trades = EXCLUDED.no_of_trades,
                total_turnover = EXCLUDED.total_turnover
        """, (company_symbol,))

        cur.execute("TRUNCATE prices_temp")
        print("  ✔ Done")

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ Prices insertion complete")


if __name__ == "__main__":
    load_prices()
