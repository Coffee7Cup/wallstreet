import os
import json
from pathlib import Path
import psycopg2
from dotenv import load_dotenv

# ----------------------------
# ENV & PATH
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
def load_companies():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    for company_dir in DATA_DIR.iterdir():
        if not company_dir.is_dir():
            continue

        meta_file = company_dir / "company.json"
        if not meta_file.exists():
            print(f"⚠ Skipping {company_dir.name} (no company.json)")
            continue

        with open(meta_file, "r", encoding="utf-8") as f:
            data = json.load(f)

        print(f"Inserting company: {data['symbol']}")

        cur.execute("""
            INSERT INTO companies (symbol, name, sector, total_shares)
            VALUES (%s, %s, %s, %s)
            ON CONFLICT (symbol)
            DO UPDATE SET
                name = EXCLUDED.name,
                sector = EXCLUDED.sector,
                total_shares = EXCLUDED.total_shares
        """, (
            data["symbol"],
            data["name"],
            data.get("sector"),
            data["total_shares"],
        ))

    conn.commit()
    cur.close()
    conn.close()
    print("\n✅ Companies insertion complete")


if __name__ == "__main__":
    load_companies()
