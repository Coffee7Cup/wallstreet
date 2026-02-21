import os
import pandas as pd
import psycopg2
from pathlib import Path
from dotenv import load_dotenv
import re
from datetime import datetime

# ----------------------------
# CONFIG & PATHS
# ----------------------------
BASE_DIR = Path(__file__).resolve().parent.parent.parent
load_dotenv(BASE_DIR / ".env")

DATABASE_URL = os.getenv("DB_URL")
if not DATABASE_URL:
    raise RuntimeError("DB_URL not found in .env")

DATA_DIR = BASE_DIR / "data/company-data"

# Mapping for specific folder names to logical symbols if needed
SYMBOL_MAPPING = {
    "relience-industries": "RELIANCE",
    "airtel": "AIRTEL",
    "infosys": "INFOSYS",
    "itc": "ITC",
    "tcs": "TCS"
}

def normalize_name(name):
    # Convert 'relience-industries' to 'Relience Industries'
    return name.replace("-", " ").title()

def get_symbol(name):
    if name.lower() in SYMBOL_MAPPING:
        return SYMBOL_MAPPING[name.lower()]
    return name.replace("-", "").upper()[:10]

def clean_numeric(val):
    if pd.isna(val) or val == "":
        return None
    cleaned = re.sub(r'[^0-9.\-]', '', str(val))
    return float(cleaned) if cleaned else None

# ----------------------------
# DB CONNECTION
# ----------------------------
def get_db_connection():
    return psycopg2.connect(DATABASE_URL)

# ----------------------------
# INGESTION LOGIC
# ----------------------------
def ingest_all():
    conn = get_db_connection()
    cur = conn.cursor()

    try:
        if not DATA_DIR.exists():
            print(f"❌ Data directory not found: {DATA_DIR}")
            return

        for company_dir in DATA_DIR.iterdir():
            if not company_dir.is_dir():
                continue

            folder_name = company_dir.name
            company_name = normalize_name(folder_name)
            symbol = get_symbol(folder_name)

            print(f"\n🚀 Processing {company_name} ({symbol})...")

            # 1. Upsert Company
            cur.execute("""
                INSERT INTO companies (symbol, name, sector, total_shares)
                VALUES (%s, %s, %s, %s)
                ON CONFLICT (symbol) DO UPDATE SET
                    name = EXCLUDED.name,
                    total_shares = EXCLUDED.total_shares
                RETURNING id;
            """, (symbol, company_name, "FIXED_SECTOR", 1000000000))
            company_id = cur.fetchone()[0]

            # 2. Fundamentals
            fund_file = company_dir / "fundamentals.csv"
            if fund_file.exists():
                print(f"  - Loading fundamentals...")
                df = pd.read_csv(fund_file)
                for _, row in df.iterrows():
                    cur.execute("""
                        INSERT INTO fundamentals (company_id, year, sales, operating_profit, net_profit, eps, equity_capital, reserves, borrowings, total_assets)
                        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                        ON CONFLICT (company_id, year) DO UPDATE SET
                            sales = EXCLUDED.sales,
                            net_profit = EXCLUDED.net_profit,
                            eps = EXCLUDED.eps
                    """, (
                        company_id, row['year'], 
                        clean_numeric(row.get('sales')), clean_numeric(row.get('operating_profit')), 
                        clean_numeric(row.get('net_profit')), clean_numeric(row.get('eps_in_rs')),
                        clean_numeric(row.get('equity_capital')), clean_numeric(row.get('reserves')),
                        clean_numeric(row.get('borrowings')), clean_numeric(row.get('total_assets'))
                    ))

            # 3. Ratios
            ratio_file = company_dir / "ratios.csv"
            if ratio_file.exists():
                print(f"  - Loading ratios...")
                df = pd.read_csv(ratio_file)
                for _, row in df.iterrows():
                    cur.execute("""
                        INSERT INTO ratios (company_id, year, roe, debt_equity, opm, intrinsic_value)
                        VALUES (%s, %s, %s, %s, %s, %s)
                        ON CONFLICT (company_id, year) DO UPDATE SET
                            roe = EXCLUDED.roe,
                            intrinsic_value = EXCLUDED.intrinsic_value
                    """, (company_id, row['year'], clean_numeric(row.get('roe')), clean_numeric(row.get('debt_equity')), clean_numeric(row.get('opm')), clean_numeric(row.get('intrinsic_value'))))

            # 4. Prices
            price_file = company_dir / "prices.csv"
            if price_file.exists():
                print(f"  - Loading prices...")
                df = pd.read_csv(price_file)
                # Map CSV columns if different
                # date,open_price,close_price,low_price,high_price,no_of_shares,no_of_trades,total_turnover
                for _, row in df.iterrows():
                    cur.execute("""
                        INSERT INTO stock_prices (company_id, date, open_price, close_price, high_price, low_price, no_of_shares, no_of_trades, total_turnover)
                        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
                        ON CONFLICT (company_id, date) DO UPDATE SET
                            close_price = EXCLUDED.close_price
                    """, (
                        company_id, row['date'], row['open_price'], row['close_price'],
                        row['high_price'], row['low_price'], row['no_of_shares'],
                        row['no_of_trades'], row['total_turnover']
                    ))

            # 5. News
            news_file = company_dir / "news.csv"
            if news_file.exists():
                print(f"  - Loading news...")
                df = pd.read_csv(news_file)
                for _, row in df.iterrows():
                    # Handle multiple date formats if needed
                    date_val = row['release_date']
                    try:
                        # Try MM/DD/YYYY
                        dt = datetime.strptime(date_val, "%m/%d/%Y").date()
                    except:
                        try:
                            # Try YYYY-MM-DD
                            dt = datetime.strptime(date_val, "%Y-%m-%d").date()
                        except:
                            dt = date_val

                    cur.execute("""
                        INSERT INTO news (release_date, title, content, company_id)
                        VALUES (%s, %s, %s, %s)
                    """, (dt, row['title'], row.get('content', ''), company_id))

        conn.commit()
        print("\n✅ Ingestion finished successfully!")

    except Exception as e:
        conn.rollback()
        print(f"❌ Error during ingestion: {e}")
    finally:
        cur.close()
        conn.close()

if __name__ == "__main__":
    ingest_all()
