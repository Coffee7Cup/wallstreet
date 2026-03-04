import os
from pathlib import Path
import psycopg2
from dotenv import load_dotenv

# ----------------------------
# Path & ENV setup
# ----------------------------
BASE_DIR = Path(__file__).resolve().parent.parent.parent
load_dotenv(BASE_DIR / ".env")

DATABASE_URL = os.getenv("DB_URL")
if not DATABASE_URL:
    raise RuntimeError("DB_URL not found in .env")

DATA_DIR = BASE_DIR / "data/company-data"

# ----------------------------
# Loader
# ----------------------------
def load_profit_loss():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        csv_path = company_dir / "profit_loss.csv"
        if not csv_path.exists():
            print(f"⚠ Skipping {company_dir.name} (no profit_loss.csv)")
            continue

        company_symbol = company_dir.name.upper()
        print(f"Loading profit & loss for: {company_symbol}")
        with open(csv_path, "r") as f:
            cur.copy_expert("""
                COPY profit_loss_temp(
                    year,
                    sales,
                    expenses,
                    operating_profit,
                    opm_percent,
                    other_income,
                    interest,
                    depreciation,
                    profit_before_tax,
                    tax_percent,
                    net_profit,
                    eps,
                    dividend_payout
                )
                FROM STDIN WITH CSV HEADER
            """, f)

        cur.execute("""
            INSERT INTO profit_loss (
                company_id, year, sales, expenses, operating_profit, opm_percent,
                other_income, interest, depreciation, profit_before_tax, tax_percent,
                net_profit, eps, dividend_payout
            )
            SELECT
                c.id, p.year, p.sales, p.expenses, p.operating_profit, p.opm_percent,
                p.other_income, p.interest, p.depreciation, p.profit_before_tax, p.tax_percent,
                p.net_profit, p.eps, p.dividend_payout
            FROM profit_loss_temp p
            JOIN companies c ON c.symbol = %s
            ON CONFLICT (company_id, year)
            DO UPDATE SET
                sales = EXCLUDED.sales,
                expenses = EXCLUDED.expenses,
                operating_profit = EXCLUDED.operating_profit,
                opm_percent = EXCLUDED.opm_percent,
                other_income = EXCLUDED.other_income,
                interest = EXCLUDED.interest,
                depreciation = EXCLUDED.depreciation,
                profit_before_tax = EXCLUDED.profit_before_tax,
                tax_percent = EXCLUDED.tax_percent,
                net_profit = EXCLUDED.net_profit,
                eps = EXCLUDED.eps,
                dividend_payout = EXCLUDED.dividend_payout
        """, (company_symbol,))

        cur.execute("TRUNCATE profit_loss_temp")
        print("  ✔ Done")

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ All profit & loss data inserted successfully")

if __name__ == "__main__":
    load_profit_loss()
