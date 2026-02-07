import psycopg2
import csv
import os

# Database configuration
DB_CONFIG = {
    "dbname": "wallstreet",
    "user": "postgres",
    "password": "yash123",
    "host": "localhost",
    "port": 5432,
}

def load_users_admins():
    conn = psycopg2.connect(**DB_CONFIG)
    cur = conn.cursor()

    # Seed data (since we don't have CSVs for these yet, let's provide a way to load them)
    # Note: password_hash is no longer in the schema as per user update.
    
    users = []
    admins = []

    # Path handling to ensure script works from any location
    base_dir = os.path.dirname(os.path.abspath(__file__))
    user_data_path = os.path.join(base_dir, "..", "dummy-data", "user_data.csv")
    admin_data_path = os.path.join(base_dir, "..", "dummy-data", "admin_data.csv")

    with open(user_data_path, "r") as f:
        reader = csv.reader(f)
        next(reader)  # Skip header row
        for row in reader:
            users.append(row)

    with open(admin_data_path, "r") as f:
        reader = csv.reader(f)
        next(reader)  # Skip header row
        for row in reader:
            admins.append(row)

    print("Inserting users...")
    for user in users:
        try:
            cur.execute("""
                INSERT INTO users (username, email, cash_balance)
                VALUES (%s, %s, %s)
                ON CONFLICT (username) DO UPDATE SET 
                    email = EXCLUDED.email,
                    cash_balance = EXCLUDED.cash_balance
            """, user)
        except Exception as e:
            print(f"Error inserting user {user[0]}: {e}")
            conn.rollback()
            return

    print("Inserting admins...")
    for admin in admins:
        try:
            cur.execute("""
                INSERT INTO admins (username, email)
                VALUES (%s, %s)
                ON CONFLICT (username) DO UPDATE SET
                    email = EXCLUDED.email
            """, admin)
        except Exception as e:
            print(f"Error inserting admin {admin[0]}: {e}")
            conn.rollback()
            return

    conn.commit()
    cur.close()
    conn.close()
    print("Users and Admins insertion complete.")

if __name__ == "__main__":
    try:
        load_users_admins()
    except Exception as e:
        print(f"General Error: {e}")
