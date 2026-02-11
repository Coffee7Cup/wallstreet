import os
import csv
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
def load_news():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        csv_path = company_dir / "news.csv"
        if not csv_path.exists():
            continue

        company_symbol = company_dir.name.upper()
        print(f"Loading news for: {company_symbol}")

        with open(csv_path, mode="r", encoding="utf-8") as f:
            reader = csv.DictReader(f)

            for row in reader:
                cur.execute("""
                    INSERT INTO news (
                        release_date,
                        title,
                        content,
                        tick,
                        company_id

                    )
                    SELECT
                        TO_DATE(%s, 'MM/DD/YYYY'),
                        %s,
                        %s,
                        %s,
                        c.id
                    FROM companies c
                    WHERE c.symbol = %s
                """, (
                    row["release_date"],
                    row["title"],
                    row.get("content", ""),
                    row.get("tick"),
                    company_symbol
                ))

        print("  ✔ Done")

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ News insertion complete")


if __name__ == "__main__":
    load_news()
