import pandas as pd
import re
import os

# Configuration
COMPANIES = ['relience-industries', 'tcs', 'itc', 'infosys', 'airtel']
DATA_DIR = '../' # Relative to scripts directory

# DB-ready column names for fundamentals
FUNDAMENTAL_COLS = [
    'year',
    'sales',
    'operating_profit',
    'net_profit',
    'eps_in_rs',
    'equity_capital',
    'reserves',
    'borrowings',
    'total_assets'
]

# Growth rates for Intrinsic Value calculation
GROWTH_RATES = {
    'relience-industries': 7,
    'tcs': 10,
    'infosys': 10,
    'itc': 6,
    'airtel': 8
}

def normalize_column_name(col: str) -> str:
    """Sanitize column names for database compatibility."""
    col = col.lower()
    col = re.sub(r'[+%]', '', col)
    col = re.sub(r'\s+', '_', col)
    col = re.sub(r'_+', '_', col)
    return col.strip('_')

def clean_numeric_value(val):
    """Clean numeric strings (remove commas, percentages)."""
    if pd.isna(val):
        return val
    # Remove commas, percentages, and other non-numeric chars except decimal and minus
    cleaned = re.sub(r'[^0-9.\-]', '', str(val))
    return cleaned if cleaned else None

def preprocess_fundamentals():
    print("--- Preprocessing Fundamentals ---")
    for company in COMPANIES:
        filepath = os.path.join(DATA_DIR, company, 'fundamentals.csv')
        if not os.path.exists(filepath):
            print(f"File not found: {filepath}")
            continue
            
        print(f"Processing: {company}")
        df = pd.read_csv(filepath)
        
        # 1. Normalize headers (DB-ready)
        print(f"  > Normalizing columns...")
        df.columns = [normalize_column_name(c) for c in df.columns]
        
        # 2. Drop TTM row if exists
        if 'year' in df.columns:
            df = df[df['year'] != 'TTM']
            # Convert year to datetime (e.g., 'Mar 2024' -> '2024-03-01')
            df['year'] = pd.to_datetime(df['year'], format='mixed', errors='coerce')
        
        # 3. Clean numeric columns and keep only relevant ones
        print(f"  > Cleaning numeric data and dropping irrelevant columns...")
        # Ensure 'eps_in_rs' is the correct name for 'eps'
        if 'eps_in_rs' not in df.columns and 'eps' in df.columns:
             df.rename(columns={'eps': 'eps_in_rs'}, inplace=True)
             
        cols_to_keep = [c for c in FUNDAMENTAL_COLS if c in df.columns]
        df = df[cols_to_keep]
        
        for col in df.columns:
            if col != 'year':
                df[col] = df[col].apply(clean_numeric_value)
                df[col] = pd.to_numeric(df[col], errors='coerce')
        
        # 4. Drop rows with too many NaNs
        df.dropna(subset=['year'], inplace=True)
        
        # 5. Save cleaned fundamentals
        df.to_csv(filepath, index=False)
        print(f"  > Saved cleaned fundamentals to {filepath}")
        
        # ---------- RATIOS ----------
        print(f"  > Generating ratios...")
        ratios = pd.DataFrame()
        ratios['year'] = df['year']
        
        # Calculations
        equity = df['equity_capital'] + df['reserves']
        ratios['roe'] = df['net_profit'] / equity
        ratios['debt_equity'] = df['borrowings'] / equity
        ratios['opm'] = df['operating_profit'] / df['sales']
        
        # Intrinsic Value: P/E formula (simplified)
        g = GROWTH_RATES[company]
        ratios['intrinsic_value'] = df['eps_in_rs'] * (8.5 + 2 * g)
        
        # Save ratios
        ratio_path = os.path.join(DATA_DIR, company, 'ratios.csv')
        ratios.to_csv(ratio_path, index=False)
        print(f"  > Saved ratios to {ratio_path}")

if __name__ == "__main__":
    preprocess_fundamentals()
    print("\nPre-processing complete!")
