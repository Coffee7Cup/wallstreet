import os
import psycopg2
import csv
from pathlib import Path
from dotenv import load_dotenv

# Setup paths
BASE_DIR = Path(__file__).resolve().parent.parent.parent
load_dotenv(BASE_DIR / ".env")

DATABASE_URL = os.getenv("DB_URL")
if not DATABASE_URL:
    print("❌ DB_URL not found in .env")
    exit(1)

NEWS_CSV = BASE_DIR / "_data" / "news.csv"

def import_news():
    if not NEWS_CSV.exists():
        print(f"❌ News CSV not found at {NEWS_CSV}")
        return

    try:
        conn = psycopg2.connect(DATABASE_URL)
        cur = conn.cursor()
        
        print(f"🚀 Starting news import from {NEWS_CSV.name}")

        # 1. Get company mapping (symbol -> id)
        cur.execute("SELECT symbol, id FROM companies")
        company_map = {row[0]: row[1] for row in cur.fetchall()}

        # 2. Read and insert news
        with open(NEWS_CSV, "r", encoding="utf-8") as f:
            reader = csv.DictReader(f)
            count = 0
            for row in reader:
                symbol = row.get("company_symbol")
                release_date = row.get("release_date")
                title = row.get("title")
                # Handle the 'contan' typo in CSV
                content = row.get("contan") or row.get("content") 

                if not symbol or not release_date or not title:
                    continue

                company_id = company_map.get(symbol)
                if not company_id:
                    print(f"⚠️ Warning: Company symbol '{symbol}' not found in database. Skipping row.")
                    continue

                # Insert into news table
                cur.execute("""
                    INSERT INTO news (release_date, title, content, company_id)
                    VALUES (%s, %s, %s, %s)
                """, (release_date, title, content, company_id))
                count += 1

        conn.commit()
        cur.close()
        conn.close()
        print(f"✅ Successfully imported {count} news items.")

    except Exception as e:
        print(f"❌ Error during import: {e}")

if __name__ == "__main__":
    import_news()
