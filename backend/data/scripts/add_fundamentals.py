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
def load_fundamentals():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    # Temp table
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

    # Iterate all directories inside company-data
    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        csv_path = company_dir / "fundamentals.csv"
        if not csv_path.exists():
            print(f"⚠ Skipping {company_dir.name} (no fundamentals.csv)")
            continue

        company_symbol = company_dir.name.upper()
        print(f"Loading fundamentals for: {company_symbol}")
        with open(csv_path, "r") as f:
            cur.copy_expert("""
                COPY fundamentals_temp(
                    year,
                    sales,
                    operating_profit,
                    net_profit,
                    eps_in_rs,
                    equity_capital,
                    reserves,
                    borrowings,
                    total_assets
                )
                FROM STDIN WITH CSV HEADER
            """, f)

        # Map folder → company via symbol column
        cur.execute("""
            INSERT INTO fundamentals (
                company_id,
                year,
                sales,
                operating_profit,
                net_profit,
                eps,
                equity_capital,
                reserves,
                borrowings,
                total_assets
            )
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
            JOIN companies c
              ON c.symbol = %s
            ON CONFLICT (company_id, year)
            DO UPDATE SET
                sales = EXCLUDED.sales,
                operating_profit = EXCLUDED.operating_profit,
                net_profit = EXCLUDED.net_profit,
                eps = EXCLUDED.eps,
                equity_capital = EXCLUDED.equity_capital,
                reserves = EXCLUDED.reserves,
                borrowings = EXCLUDED.borrowings,
                total_assets = EXCLUDED.total_assets
        """, (company_symbol,))

        cur.execute("TRUNCATE fundamentals_temp")
        print("  ✔ Done")

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ All fundamentals inserted successfully")


if __name__ == "__main__":
    load_fundamentals()
