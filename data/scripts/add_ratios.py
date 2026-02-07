import os
import psycopg2


# TODO: complete this code

conn = psycopg2.connect(
    dbname="wallstreet",
    user="postgres",
    password="yash123",
    host="localhost",
    port=5432,
)
cur = conn.cursor()

names = ["airtel", "infosys", "itc", "relience-industries", "tcs"]
names_symbols = {"airtel" : "AIRTEL", "infosys" : "INFOSYS", "itc" : "ITC", "relience-industries" : "RELIANCE", "tcs":"TCS"}

for name in names:
    path = f"../{name}/ratios.csv"
    if not os.path.exists(path):
        continue

    with open(path, "r") as f:
        cur.copy_expert(
            "COPY ratios_temp(year, roe, debt_equity, opm, intrinsic_value) FROM STDIN WITH CSV HEADER",
            f
        )
    
    print(f"Inserting ratios for: {name}")
    cur.execute("""
        INSERT INTO ratios (company_id, year, roe, debt_equity, opm, intrinsic_value)
        SELECT
            c.id,
            r.year,
            r.roe,
            r.debt_equity,
            r.opm,
            r.intrinsic_value
        FROM ratios_temp r
        JOIN companies c ON c.symbol = %s
        ON CONFLICT (company_id, year) DO UPDATE SET
            roe = EXCLUDED.roe,
            debt_equity = EXCLUDED.debt_equity,
            opm = EXCLUDED.opm,
            intrinsic_value = EXCLUDED.intrinsic_value
        """, (names_symbols[name],))

    cur.execute("TRUNCATE ratios_temp")

conn.commit()
cur.close()
conn.close()
print("\nRatios insertion complete.")
