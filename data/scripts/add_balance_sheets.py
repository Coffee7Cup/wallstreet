import os
from pathlib import Path
import psycopg2
from dotenv import load_dotenv

BASE_DIR = Path(__file__).resolve().parent.parent.parent
load_dotenv(BASE_DIR / ".env")

DATABASE_URL = os.getenv("DB_URL")
if not DATABASE_URL:
    raise RuntimeError("DB_URL not found in .env")

DATA_DIR = BASE_DIR / "data/company-data"

def load_balance_sheets():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        csv_path = company_dir / "balance_sheets.csv"
        if not csv_path.exists():
            print(f"⚠ Skipping {company_dir.name} (no balance_sheets.csv)")
            continue

        company_symbol = company_dir.name.upper()
        print(f"Loading balance sheets for: {company_symbol}")
        with open(csv_path, "r") as f:
            cur.copy_expert("""
                COPY balance_sheets_temp(
                    year,
                    equity_capital,
                    reserves,
                    borrowings,
                    other_liabilities,
                    total_liabilities,
                    fixed_assets,
                    cwip,
                    investments,
                    other_assets,
                    total_assets
                )
                FROM STDIN WITH CSV HEADER
            """, f)

        cur.execute("""
            INSERT INTO balance_sheets (
                company_id, year, equity_capital, reserves, borrowings, other_liabilities, total_liabilities,
                fixed_assets, cwip, investments, other_assets, total_assets
            )
            SELECT
                c.id, b.year, b.equity_capital, b.reserves, b.borrowings, b.other_liabilities, b.total_liabilities,
                b.fixed_assets, b.cwip, b.investments, b.other_assets, b.total_assets
            FROM balance_sheets_temp b
            JOIN companies c ON c.symbol = %s
            ON CONFLICT (company_id, year)
            DO UPDATE SET
                equity_capital = EXCLUDED.equity_capital,
                reserves = EXCLUDED.reserves,
                borrowings = EXCLUDED.borrowings,
                other_liabilities = EXCLUDED.other_liabilities,
                total_liabilities = EXCLUDED.total_liabilities,
                fixed_assets = EXCLUDED.fixed_assets,
                cwip = EXCLUDED.cwip,
                investments = EXCLUDED.investments,
                other_assets = EXCLUDED.other_assets,
                total_assets = EXCLUDED.total_assets
        """, (company_symbol,))

        cur.execute("TRUNCATE balance_sheets_temp")
        print("  ✔ Done")

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ All balance sheets data inserted successfully")

if __name__ == "__main__":
    load_balance_sheets()
