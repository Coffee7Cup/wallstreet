import os
import psycopg2
import csv

# Database configuration
DB_CONFIG = {
    "dbname": "wallstreet",
    "user": "postgres",
    "password": "yash123",
    "host": "localhost",
    "port": 5432,
}

def load_news():
    conn = psycopg2.connect(**DB_CONFIG)
    cur = conn.cursor()

    # Paths to news files (assuming news.csv in each company dir or a global one)
    # The user request mentioned "news" without specifying format. 
    # Let's assume a central data/news.csv or per-company news.csv
    
    news_dirs = ["../airtel", "../infosys", "../itc", "../relience-industries", "../tcs"]
    
    for news_dir in news_dirs:
        path = os.path.join(news_dir, "news.csv")
        if not os.path.exists(path):
            continue
            
        print(f"Loading news from: {path}")
        with open(path, mode='r', encoding='utf-8') as f:
            reader = csv.DictReader(f)
            for row in reader:
                cur.execute("""
                    INSERT INTO news (release_date, title, content, news_type, impact_factor)
                    VALUES (%s, %s, %s, %s, %s)
                """, (row['release_date'], row['title'], row.get('content', ''), row.get('news_type', ''), row.get('impact_factor', 0)))

    conn.commit()
    cur.close()
    conn.close()
    print("News insertion complete.")

if __name__ == "__main__":
    load_news()
