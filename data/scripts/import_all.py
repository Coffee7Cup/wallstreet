import subprocess
import sys
from pathlib import Path

# List of scripts to run in order
SCRIPTS = [
    "add_users_admins.py",
    "preprocess_data.py",
    "add_companies.py",
    "add_profit_loss.py",
    "add_balance_sheets.py",
    "add_cash_flows.py",
    "add_ratios.py",
    "add_prices.py",
    "add_news.py",
]

def run_script(script_name):
    print(f"\n{'='*50}")
    print(f"Running {script_name}...")
    print(f"{'='*50}\n")
    
    try:
        # Use sys.executable to ensure we use the same python interpreter
        result = subprocess.run([sys.executable, script_name], check=True, text=True)
        return True
    except subprocess.CalledProcessError as e:
        print(f"\n❌ Error running {script_name}: {e}")
        return False
    except FileNotFoundError:
        print(f"\n❌ Script not found: {script_name}")
        return False

def main():
    script_dir = Path(__file__).parent
    
    print("🚀 Starting full data ingestion process...")
    
    for script in SCRIPTS:
        success = run_script(script_dir / script)
        if not success:
            print(f"\n🛑 Ingestion stopped due to error in {script}")
            sys.exit(1)
            
    print("\n✨ All scripts executed successfully. Data ingestion complete!")

if __name__ == "__main__":
    main()
