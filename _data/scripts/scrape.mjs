import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { ScreenerScraperPro } from 'screener-scraper-pro';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const dataDir = path.resolve(__dirname, '..');

// Helper to clean and parse number values
function parseNumber(val) {
    if (val === undefined || val === null || val === '') return '';
    if (typeof val === 'number') return val;
    val = String(val).replace(/,/g, '').replace(/%/g, '').trim();
    if (val === '') return '';
    return Number(val);
}

// Helper to convert "Mar 2024" to "2024-03-01"
function parseDate(header) {
    if (header === "TTM") return new Date().toISOString().slice(0, 10);
    const parts = header.split(" ");
    if (parts.length !== 2) return header;
    const months = {
        "Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04", "May": "05", "Jun": "06",
        "Jul": "07", "Aug": "08", "Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12"
    };
    const mm = months[parts[0]];
    const yyyy = parts[1];
    return `${yyyy}-${mm}-01`;
}

async function getCompanyUrl(folderName) {
    // Attempt to search screener
    const query = folderName.replace(/-/g, ' ');
    try {
        const res = await fetch(`https://www.screener.in/api/company/search/?q=${encodeURIComponent(query)}`);
        const json = await res.json();
        if (json && json.length > 0) {
            return `https://www.screener.in${json[0].url}`;
        }
    } catch (err) {
        console.warn(`Search failed for ${query}, returning default URL`);
    }
    // Fallback based on typical mapping
    let sym = folderName.toUpperCase();
    if (folderName === 'mahindra-and-mahindra') sym = 'M&M';
    if (folderName === 'maruthi-suzuki') sym = 'MARUTI';
    if (folderName === 'tata-motors') sym = 'TATAMOTORS';
    if (folderName === 'bajaj-auto') sym = 'BAJAJ-AUTO';
    return `https://www.screener.in/company/${sym}/`;
}

async function scrapeAll() {
    // Read all sector directories
    const sectors = fs.readdirSync(dataDir, { withFileTypes: true })
        .filter(d => d.isDirectory() && d.name !== 'scripts')
        .map(d => d.name);

    for (const sector of sectors) {
        const sectorPath = path.join(dataDir, sector);
        const companies = fs.readdirSync(sectorPath, { withFileTypes: true })
            .filter(d => d.isDirectory())
            .map(d => d.name);

        for (const company of companies) {
            const companyPath = path.join(sectorPath, company);
            const url = await getCompanyUrl(company);
            console.log(`\nfetching ${company} from ${url}...`);

            try {
                const data = await ScreenerScraperPro(url);
                if (!data || !data.profitLoss) {
                    console.log(`Failed to fetch proper data for ${company}`);
                    continue;
                }

                // 1. Profit & Loss: year,sales,expenses,operating_profit,opm_percent,other_income,interest,depreciation,profit_before_tax,tax_percent,net_profit,eps,dividend_payout
                if (data.profitLoss && data.profitLoss.headers) {
                    let plCsv = `year,sales,expenses,operating_profit,opm_percent,other_income,interest,depreciation,profit_before_tax,tax_percent,net_profit,eps,dividend_payout\n`;
                    for (const h of data.profitLoss.headers) {
                        const date = parseDate(h);
                        const pd = data.profitLoss.data;
                        const row = [
                            date,
                            parseNumber(pd["Sales"]?.[h]),
                            parseNumber(pd["Expenses"]?.[h]),
                            parseNumber(pd["Operating Profit"]?.[h]),
                            parseNumber(pd["OPM %"]?.[h]),
                            parseNumber(pd["Other Income"]?.[h]),
                            parseNumber(pd["Interest"]?.[h]),
                            parseNumber(pd["Depreciation"]?.[h]),
                            parseNumber(pd["Profit before tax"]?.[h]),
                            parseNumber(pd["Tax %"]?.[h]),
                            parseNumber(pd["Net Profit"]?.[h]),
                            parseNumber(pd["EPS in Rs"]?.[h]),
                            parseNumber(pd["Dividend Payout %"]?.[h])
                        ];
                        plCsv += row.join(',') + '\n';
                    }
                    fs.writeFileSync(path.join(companyPath, 'profit_loss.csv'), plCsv);
                }

                // 2. Cash flows: year,cash_from_operating_activity,cash_from_investing_activity,cash_from_financing_activity,net_cash_flow
                if (data.cashFlow && data.cashFlow.headers) {
                    let cfCsv = `year,cash_from_operating_activity,cash_from_investing_activity,cash_from_financing_activity,net_cash_flow\n`;
                    for (const h of data.cashFlow.headers) {
                        const date = parseDate(h);
                        const pd = data.cashFlow.data;
                        const row = [
                            date,
                            parseNumber(pd["Cash from Operating Activity"]?.[h]),
                            parseNumber(pd["Cash from Investing Activity"]?.[h]),
                            parseNumber(pd["Cash from Financing Activity"]?.[h]),
                            parseNumber(pd["Net Cash Flow"]?.[h])
                        ];
                        cfCsv += row.join(',') + '\n';
                    }
                    fs.writeFileSync(path.join(companyPath, 'cash_flows.csv'), cfCsv);
                }

                // 3. Balance sheets: year,equity_capital,reserves,borrowings,other_liabilities,total_liabilities,fixed_assets,cwip,investments,other_assets,total_assets
                if (data.balanceSheet && data.balanceSheet.headers) {
                    let bsCsv = `year,equity_capital,reserves,borrowings,other_liabilities,total_liabilities,fixed_assets,cwip,investments,other_assets,total_assets\n`;
                    for (const h of data.balanceSheet.headers) {
                        const date = parseDate(h);
                        const pd = data.balanceSheet.data;
                        const row = [
                            date,
                            parseNumber(pd["Equity Capital"]?.[h]),
                            parseNumber(pd["Reserves"]?.[h]),
                            parseNumber(pd["Borrowings"]?.[h]),
                            parseNumber(pd["Other Liabilities"]?.[h]),
                            parseNumber(pd["Total Liabilities"]?.[h]),
                            parseNumber(pd["Fixed Assets"]?.[h]),
                            parseNumber(pd["CWIP"]?.[h]),
                            parseNumber(pd["Investments"]?.[h]),
                            parseNumber(pd["Other Assets"]?.[h]),
                            parseNumber(pd["Total Assets"]?.[h])
                        ];
                        bsCsv += row.join(',') + '\n';
                    }
                    fs.writeFileSync(path.join(companyPath, 'balance_sheets.csv'), bsCsv);
                }

                // 4. Ratios: year,roe,debt_equity,opm,intrinsic_value,debtor_days,inventory_days,days_payable,cash_conversion_cycle,working_capital_days,roce_percent
                // We'll combine from ratios & profit_loss (OPM)
                if (data.ratios && data.ratios.headers) {
                    let rCsv = `year,roe,debt_equity,opm,intrinsic_value,debtor_days,inventory_days,days_payable,cash_conversion_cycle,working_capital_days,roce_percent\n`;
                    for (const h of data.ratios.headers) {
                        const date = parseDate(h);
                        const rd = data.ratios.data;
                        const row = [
                            date,
                            0, // ROE default
                            0, // debt_equity 
                            parseNumber(data.profitLoss?.data?.["OPM %"]?.[h]) || 0, // opm
                            0, // intrinsic_value
                            parseNumber(rd["Debtor Days"]?.[h]),
                            parseNumber(rd["Inventory Days"]?.[h]),
                            parseNumber(rd["Days Payable"]?.[h]),
                            parseNumber(rd["Cash Conversion Cycle"]?.[h]),
                            parseNumber(rd["Working Capital Days"]?.[h]),
                            parseNumber(rd["ROCE %"]?.[h])
                        ];
                        rCsv += row.join(',') + '\n';
                    }
                    fs.writeFileSync(path.join(companyPath, 'ratios.csv'), rCsv);
                }

                console.log(`Saved CSVs for ${company}`);
            } catch (err) {
                console.error(`Error scraping ${company}:`, err);
            }
        }
    }
}

scrapeAll().catch(console.error);
