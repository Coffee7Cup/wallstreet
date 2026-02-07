import os
import psycopg2

conn = psycopg2.connect(
    dbname="wallstreet",
    user="postgres",
    password="yash123",
    host="localhost",
    port=5432,
)
cur = conn.cursor()

names = ["airtel", "infosys", "itc", "relience-industries", "tcs"]

names_symbols = {
    "airtel": "AIRTEL",
    "infosys": "INFOSYS",
    "itc": "ITC",
    "relience-industries": "RELIANCE",
    "tcs": "TCS",
}

for name in names:
    path = f"../{name}/prices.csv"

    if not os.path.exists(path):
        continue

    with open(path) as f:
        cur.copy_expert(
            "COPY prices_temp(date, open_price, close_price, low_price, high_price, no_of_shares, no_of_trades, total_turnover) FROM STDIN HEADER CSV",
            f,
        )

        print("in the function")

        cur.execute(
            """
            INSERT INTO stock_prices (company_id, date,open_price, close_price, low_price, high_price, no_of_shares, no_of_trades, total_turnover)
            SELECT
                c.id,
                p.date,
                p.open_price,
                p.close_price,
                p.high_price,
                p.low_price,
                p.no_of_shares,
                p.no_of_trades,
                p.total_turnover
            FROM prices_temp p
            JOIN companies c
            ON c.symbol = %s
        """,
            (names_symbols[name],),
        )
        cur.execute("TRUNCATE TABLE prices_temp")

        print("function complete")

conn.commit()
cur.close()
conn.close()
