import os
import psycopg2
import csv
from pathlib import Path
from dotenv import load_dotenv

BASE_DIR = Path(__file__).resolve().parent.parent.parent
load_dotenv(BASE_DIR / ".env")

DATABASE_URL = os.getenv("DB_URL")
if not DATABASE_URL:
    raise RuntimeError("DB_URL not found in .env")

CSV_PATH = BASE_DIR / "_data" / "user_data.csv"

def main():
    if not CSV_PATH.exists():
        print(f"❌ CSV file not found: {CSV_PATH}")
        return

    try:
        conn = psycopg2.connect(DATABASE_URL)
        cur = conn.cursor()
        
        print(f"🚀 Importing users from {CSV_PATH.name}...")
        
        with open(CSV_PATH, "r", encoding="utf-8") as f:
            reader = csv.DictReader(f)
            users_added = 0
            
            for row in reader:
                raw_username = row.get("username", "").strip()
                raw_email = row.get("email", "").strip()
                
                if not raw_username or not raw_email:
                    continue
                
                # Normalization: lower case, spaces to underscores
                username = raw_username.lower().replace(" ", "_").replace(".", "_").replace("__", "_")
                # Email: append domain
                email = f"{raw_email.strip()}@gprec.ac.in".lower() # lowercasing email for consistency
                
                try:
                    cur.execute("""
                        INSERT INTO users (username, email)
                        VALUES (%s, %s)
                        ON CONFLICT DO NOTHING
                    """, (username, email))
                    if cur.rowcount > 0:
                        users_added += 1
                except Exception as e:
                    print(f"⚠️ Error adding user {username}: {e}")
                    conn.rollback()
                    continue
            
        conn.commit()
        cur.close()
        conn.close()
        print(f"✅ Successfully added {users_added} new users.")
        
    except Exception as e:
        print(f"❌ Database error: {e}")

if __name__ == "__main__":
    main()
