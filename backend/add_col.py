import psycopg2

try:
    conn = psycopg2.connect("postgresql://postgres:yash123@localhost:5432/wallstreet")
    conn.autocommit = True
    cur = conn.cursor()
    cur.execute("ALTER TABLE trades ADD COLUMN IF NOT EXISTS price NUMERIC(15, 2) NOT NULL DEFAULT 0")
    print("Successfully added price column to trades table!")
    cur.close()
    conn.close()
except psycopg2.Error as e:
    print(f"Database error: {e}")
