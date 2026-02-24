import os
import psycopg2
import csv
from psycopg2.extras import execute_values
from pathlib import Path
from dotenv import load_dotenv

BASE_DIR = Path(__file__).resolve().parent.parent.parent
load_dotenv(BASE_DIR / ".env")

DATABASE_URL = os.getenv("DB_URL")
if not DATABASE_URL:
    raise RuntimeError("DB_URL not found in .env")

DATA_DIR = BASE_DIR / "_data"

def process_table(cur, csv_path, temp_table, main_table, company_symbol, cols):
    if not csv_path.exists():
        return

    cur.execute(f"TRUNCATE {temp_table}")
    
    with open(csv_path, "r", encoding="utf-8") as f:
        # Some CSVs may have blank lines, Postgres COPY STDIN handles standard CSVs nicely
        cur.copy_expert(f"""
            COPY {temp_table}({', '.join(cols)})
            FROM STDIN WITH CSV HEADER
        """, f)

    # Insert into main table
    conflict_cols = "company_id, date" if main_table == "stock_prices" else "company_id, year"
    
    update_clause = ",\n".join([f"{c} = EXCLUDED.{c}" for c in cols if c not in ('year', 'date')])
    if not update_clause:
        update_clause = f"{cols[0]} = EXCLUDED.{cols[0]}" # Fallback if we only have the key, unlikely
        
    date_col = 'date' if 'date' in cols else 'year'

    query = f"""
        INSERT INTO {main_table} (company_id, {', '.join(cols)})
        SELECT c.id, {', '.join([f't.{c}' for c in cols])}
        FROM {temp_table} t
        JOIN companies c ON c.symbol = %s
        ON CONFLICT ({conflict_cols})
        DO UPDATE SET
            {update_clause}
    """
    cur.execute(query, (company_symbol,))
    cur.execute(f"TRUNCATE {temp_table}")


def main():
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    print("🚀 Starting import of financial data from _data/")
    
    # Iterate through sectors
    for sector_dir in DATA_DIR.iterdir():
        if not sector_dir.is_dir() or sector_dir.name == "scripts":
            continue
            
        sector = sector_dir.name.upper()
        
        # Iterate through companies in sector
        for company_dir in sector_dir.iterdir():
            if not company_dir.is_dir():
                continue
                
            symbol = company_dir.name.upper()
            mapping = {
                "MAHINDRA-AND-MAHINDRA": "M&M",
                "MARUTI-SUZUKI": "MARUTI",
                "MARUTHI-SUZUKI": "MARUTI",
                "TATA-MOTORS": "TATAMOTORS",
                "BAJAJ-AUTO": "BAJAJ-AUTO",
                "BHARAT-SEATS": "BHARATSEAT",
                "STATE-BANK-OF-INDIA": "SBIN",
                "HDFC-BANK": "HDFCBANK",
                "ICICI-BANK": "ICICIBANK",
                "AXIS-BANK": "AXISBANK",
                "BANDHAN-BANK": "BANDHANBNK",
                "TATA-CONSULTANCY-SERVICES": "TCS",
                "LNT-MINDTREE": "LTIM",
                "COAL-INDIA": "COALINDIA",
                "HINDALCO-INDUSTRIES": "HINDALCO",
                "JSW-STEEL": "JSWSTEEL",
                "TATA-STEEL": "TATASTEEL",
                "DR-REDDYS-LABORATORIES": "DRREDDY",
                "SUN-PHARMA": "SUNPHARMA"
            }
            if symbol in mapping:
                symbol = mapping[symbol]
            else:
                symbol = symbol.replace("-", "")[:10]

            name = company_dir.name.replace("-", " ").title()

            # Ensure company exists in companies table
            cur.execute("""
                INSERT INTO companies (symbol, name, sector, total_shares)
                VALUES (%s, %s, %s, %s)
                ON CONFLICT (symbol) DO UPDATE SET sector = EXCLUDED.sector
            """, (symbol, name, sector, 1000000000))
            
            print(f"📥 Loading data for {symbol}...")

            # 1. Prices
            prices_csv = company_dir / "prices.csv"
            if prices_csv.exists():
                with open(prices_csv, "r", encoding="utf-8-sig") as f:
                    reader = csv.DictReader(f)
                    data = []
                    for row in reader:
                        date_val = row.get("Date")
                        if not date_val: continue
                        clean = lambda x: str(x).replace(",", "").strip() if x else "0"
                        trades = clean(row.get("No. of Trades", "0"))
                        if trades == "-" or not trades: trades = "0"
                        turnover = clean(row.get("Total Turnover (Rs.)", "0"))
                        if turnover == "-" or not turnover: turnover = "0"
                        shares = clean(row.get("No.of Shares", "0"))
                        if shares == "-" or not shares: shares = "0"
                        
                        data.append((
                            date_val, clean(row.get("Open Price", "0")), clean(row.get("High Price", "0")),
                            clean(row.get("Low Price", "0")), clean(row.get("Close Price", "0")),
                            shares, trades, turnover
                        ))
                        
                if data:
                    cur.execute("CREATE TEMP TABLE raw_prices (date DATE, open_price NUMERIC, high_price NUMERIC, low_price NUMERIC, close_price NUMERIC, no_of_shares NUMERIC, no_of_trades NUMERIC, total_turnover NUMERIC)")
                    execute_values(cur, "INSERT INTO raw_prices VALUES %s", data)
                    
                    cur.execute("""
                        INSERT INTO stock_prices (company_id, date, open_price, high_price, low_price, close_price, no_of_shares, no_of_trades, total_turnover)
                        SELECT c.id, r.date, r.open_price, r.high_price, r.low_price, r.close_price, r.no_of_shares, r.no_of_trades, r.total_turnover
                        FROM raw_prices r
                        JOIN companies c ON c.symbol = %s
                        ON CONFLICT (company_id, date) DO UPDATE SET
                            open_price = EXCLUDED.open_price,
                            high_price = EXCLUDED.high_price,
                            low_price = EXCLUDED.low_price,
                            close_price = EXCLUDED.close_price,
                            no_of_shares = EXCLUDED.no_of_shares,
                            no_of_trades = EXCLUDED.no_of_trades,
                            total_turnover = EXCLUDED.total_turnover
                    """, (symbol,))
                    cur.execute("DROP TABLE raw_prices")

            # 2. Profit & Loss
            pl_cols = ["year","sales","expenses","operating_profit","opm_percent","other_income","interest","depreciation","profit_before_tax","tax_percent","net_profit","eps","dividend_payout"]
            process_table(cur, company_dir / "profit_loss.csv", "profit_loss_temp", "profit_loss", symbol, pl_cols)

            # 3. Cash Flows
            cf_cols = ["year","cash_from_operating_activity","cash_from_investing_activity","cash_from_financing_activity","net_cash_flow"]
            process_table(cur, company_dir / "cash_flows.csv", "cash_flows_temp", "cash_flows", symbol, cf_cols)

            # 4. Balance Sheets
            bs_cols = ["year","equity_capital","reserves","borrowings","other_liabilities","total_liabilities","fixed_assets","cwip","investments","other_assets","total_assets"]
            process_table(cur, company_dir / "balance_sheets.csv", "balance_sheets_temp", "balance_sheets", symbol, bs_cols)

            # 5. Ratios
            r_cols = ["year","roe","debt_equity","opm","intrinsic_value","debtor_days","inventory_days","days_payable","cash_conversion_cycle","working_capital_days","roce_percent"]
            process_table(cur, company_dir / "ratios.csv", "ratios_temp", "ratios", symbol, r_cols)

    conn.commit()
    cur.close()
    conn.close()
    print("✅ All data successfully imported into the DB.")

if __name__ == "__main__":
    main()
