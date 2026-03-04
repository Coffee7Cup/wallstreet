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

def load_cash_flows():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        csv_path = company_dir / "cash_flows.csv"
        if not csv_path.exists():
            print(f"⚠ Skipping {company_dir.name} (no cash_flows.csv)")
            continue

        company_symbol = company_dir.name.upper()
        print(f"Loading cash flows for: {company_symbol}")
        with open(csv_path, "r") as f:
            cur.copy_expert("""
                COPY cash_flows_temp(
                    year,
                    cash_from_operating_activity,
                    cash_from_investing_activity,
                    cash_from_financing_activity,
                    net_cash_flow
                )
                FROM STDIN WITH CSV HEADER
            """, f)

        cur.execute("""
            INSERT INTO cash_flows (
                company_id, year, cash_from_operating_activity, cash_from_investing_activity,
                cash_from_financing_activity, net_cash_flow
            )
            SELECT
                c.id, cf.year, cf.cash_from_operating_activity, cf.cash_from_investing_activity,
                cf.cash_from_financing_activity, cf.net_cash_flow
            FROM cash_flows_temp cf
            JOIN companies c ON c.symbol = %s
            ON CONFLICT (company_id, year)
            DO UPDATE SET
                cash_from_operating_activity = EXCLUDED.cash_from_operating_activity,
                cash_from_investing_activity = EXCLUDED.cash_from_investing_activity,
                cash_from_financing_activity = EXCLUDED.cash_from_financing_activity,
                net_cash_flow = EXCLUDED.net_cash_flow
        """, (company_symbol,))

        cur.execute("TRUNCATE cash_flows_temp")
        print("  ✔ Done")

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ All cash flows data inserted successfully")

if __name__ == "__main__":
    load_cash_flows()
