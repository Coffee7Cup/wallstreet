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
def load_ratios():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    # Ensure temp table exists
    cur.execute("""
        CREATE TABLE IF NOT EXISTS ratios_temp (
            year DATE,
            roe NUMERIC,
            debt_equity NUMERIC,
            opm NUMERIC,
            intrinsic_value NUMERIC
        )
    """)

    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        csv_path = company_dir / "ratios.csv"
        if not csv_path.exists():
            continue

        company_symbol = company_dir.name.upper()
        print(f"Loading ratios for: {company_symbol}")

        with open(csv_path, "r") as f:
            cur.copy_expert("""
                COPY ratios_temp(
                    year,
                    roe,
                    debt_equity,
                    opm,
                    intrinsic_value
                )
                FROM STDIN WITH CSV HEADER
            """, f)

        cur.execute("""
            INSERT INTO ratios (
                company_id,
                year,
                roe,
                debt_equity,
                opm,
                intrinsic_value
            )
            SELECT
                c.id,
                r.year,
                r.roe,
                r.debt_equity,
                r.opm,
                r.intrinsic_value
            FROM ratios_temp r
            JOIN companies c
              ON c.symbol = %s
            ON CONFLICT (company_id, year)
            DO UPDATE SET
                roe = EXCLUDED.roe,
                debt_equity = EXCLUDED.debt_equity,
                opm = EXCLUDED.opm,
                intrinsic_value = EXCLUDED.intrinsic_value
        """, (company_symbol,))

        cur.execute("TRUNCATE ratios_temp")
        print("  ✔ Done")

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ Ratios insertion complete")


if __name__ == "__main__":
    load_ratios()
