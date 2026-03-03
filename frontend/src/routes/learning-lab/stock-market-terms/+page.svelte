<script>
	import { List, ShieldAlert, Target, BarChart, Layers, Wallet } from 'lucide-svelte';

	const terms = [
		{
			name: 'Price',
			desc: 'The current value at which a stock is bought or sold in the market.'
		},
		{
			name: 'Market Capitalization',
			desc: 'Total value of a company calculated as share price multiplied by total outstanding shares.'
		},
		{
			name: 'Liquidity',
			desc: 'How easily a stock can be bought or sold without significantly affecting its price.'
		},
		{
			name: 'Volatility',
			desc: 'The degree of price variation over time; higher volatility means higher risk.'
		},
		{
			name: 'Bid Price',
			desc: 'The highest price a buyer is willing to pay for a stock.'
		},
		{
			name: 'Ask Price',
			desc: 'The lowest price a seller is willing to accept for a stock.'
		},
		{
			name: 'Spread',
			desc: 'The difference between the bid price and the ask price.'
		},
		{
			name: 'Open Price',
			desc: 'The price at which a stock first trades when the market opens.'
		},
		{
			name: 'Close Price',
			desc: 'The final trading price of a stock when the market closes.'
		},
		{
			name: 'High',
			desc: 'The highest price a stock trades at during a given period.'
		},
		{
			name: 'Low',
			desc: 'The lowest price a stock trades at during a given period.'
		},
		{
			name: 'Dividend',
			desc: 'A portion of company profits distributed to shareholders.'
		},
		{
			name: 'Dividend Yield',
			desc: 'Dividend expressed as a percentage of the stock’s current price.'
		},
		{
			name: 'P/E Ratio',
			desc: 'Price-to-Earnings ratio showing how much investors pay for one unit of earnings.'
		},
		{
			name: 'P/B Ratio',
			desc: 'Price-to-Book ratio comparing market value to the company’s book value.'
		},
		{
			name: 'ROE',
			desc: 'Return on Equity measuring how efficiently a company uses shareholders’ funds.'
		},
		{
			name: 'ROCE',
			desc: 'Return on Capital Employed indicating profitability relative to capital used.'
		},
		{
			name: 'Bear Market',
			desc: 'A market condition where prices are falling or expected to fall.'
		},
		{
			name: 'Stop Loss',
			desc: 'A predefined price at which a position is automatically exited to limit losses.'
		},
		{
			name: 'Circuit Breaker',
			desc: 'A regulatory halt in trading triggered when prices move beyond predefined limits.'
		},
		{
			name: 'Candlestick',
			desc: 'A type of price chart that displays the high, low, open, and close prices of a security for a specific period.'
		}
	];
</script>

<svelte:head>
	<title>Stock Market Terms — Learning Lab</title>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(242,201,76,0.12); color: #F2C94C">Module 7</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			Stock Market Terms
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			To navigate the market effectively, you need to speak its language. Here are the core terms
			every investor must know to communicate with brokers and other traders.
		</p>
	</header>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">Essential Glossary</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			These terms form the foundation of most financial news and analysis reports. Understanding
			them helps you filter out the noise.
		</p>
		<div class="grid gap-6 sm:grid-cols-2">
			{#each terms as term}
				<div
					class="group rounded-2xl border border-(--border-color) bg-(--bg-primary) p-8 transition-all hover:border-[#F2C94C]"
				>
					<h3 class="mb-3 flex items-center gap-3 text-lg font-bold">
						<div
							class="h-3 w-3 rounded-full bg-[#F2C94C]/20 transition-colors group-hover:bg-[#F2C94C]"
						></div>
						{term.name}
					</h3>
					<p class="text-base leading-relaxed opacity-70">{term.desc}</p>
				</div>
			{/each}
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-3xl font-extrabold">Anatomy of a Candlestick</h2>
		<p class="mb-8 max-w-[800px] text-lg leading-relaxed opacity-80">
			A candlestick is a way of displaying information about an asset's price movement. Candlestick
			charts are the most popular way to analyze a stock because they pack a lot of data into a
			single "candle."
		</p>

		<div class="grid grid-cols-1 gap-8 md:grid-cols-2">
			<!-- Bullish Candle Viz -->
			<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-sm">
				<div class="mb-4 text-center">
					<span
						class="rounded-full bg-green-500/10 px-4 py-1.5 text-xs font-black text-green-600 uppercase"
					>
						Bullish (Price Up)
					</span>
				</div>
				<div class="flex items-center justify-center py-6">
					<svg width="200" height="300" viewBox="0 0 200 300" class="overflow-visible">
						<!-- Upper Wick -->
						<line x1="100" y1="20" x2="100" y2="80" stroke="#27AE60" stroke-width="3" />
						<text x="115" y="25" fill="#27AE60" class="text-[12px] font-bold">High</text>

						<!-- Body -->
						<rect x="60" y="80" width="80" height="140" fill="#27AE60" rx="4" />
						<text x="150" y="85" fill="currentColor" class="text-[12px] font-bold opacity-70"
							>Close</text
						>
						<text x="150" y="225" fill="currentColor" class="text-[12px] font-bold opacity-70"
							>Open</text
						>
						<text
							x="0"
							y="155"
							fill="currentColor"
							class="rotate-[-90deg] text-[12px] font-bold opacity-70">Real Body</text
						>

						<!-- Lower Wick -->
						<line x1="100" y1="220" x2="100" y2="280" stroke="#27AE60" stroke-width="3" />
						<text x="115" y="285" fill="#27AE60" class="text-[12px] font-bold">Low</text>
					</svg>
				</div>
				<p class="mt-4 text-center text-sm opacity-60">
					A green candle means the stock <b>closed higher</b> than it opened.
				</p>
			</div>

			<!-- Bearish Candle Viz -->
			<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-sm">
				<div class="mb-4 text-center">
					<span
						class="rounded-full bg-red-500/10 px-4 py-1.5 text-xs font-black text-red-600 uppercase"
					>
						Bearish (Price Down)
					</span>
				</div>
				<div class="flex items-center justify-center py-6">
					<svg width="200" height="300" viewBox="0 0 200 300" class="overflow-visible">
						<!-- Upper Wick -->
						<line x1="100" y1="20" x2="100" y2="80" stroke="#EB5757" stroke-width="3" />
						<text x="115" y="25" fill="#EB5757" class="text-[12px] font-bold">High</text>

						<!-- Body -->
						<rect x="60" y="80" width="80" height="140" fill="#EB5757" rx="4" />
						<text x="150" y="85" fill="currentColor" class="text-[12px] font-bold opacity-70"
							>Open</text
						>
						<text x="150" y="225" fill="currentColor" class="text-[12px] font-bold opacity-70"
							>Close</text
						>

						<!-- Lower Wick -->
						<line x1="100" y1="220" x2="100" y2="280" stroke="#EB5757" stroke-width="3" />
						<text x="115" y="285" fill="#EB5757" class="text-[12px] font-bold">Low</text>
					</svg>
				</div>
				<p class="mt-4 text-center text-sm opacity-60">
					A red candle means the stock <b>closed lower</b> than it opened.
				</p>
			</div>
		</div>

		<div class="mt-8 rounded-3xl border border-(--border-color) bg-(--bg-hover) p-8">
			<h4 class="mb-4 text-xl font-bold">Key Observations:</h4>
			<ul class="space-y-4 text-lg opacity-80">
				<li>• <b>Wicks (Shadows):</b> Represent the price range of the trading period.</li>
				<li>• <b>Body:</b> Represents the difference between the opening and closing price.</li>
				<li>• <b>OHLC:</b> Stands for <b>O</b>pen, <b>H</b>igh, <b>L</b>ow, and <b>C</b>lose.</li>
			</ul>
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-3xl font-extrabold">
			Support & Resistance Visualized
		</h2>
		<div
			class="relative my-6 overflow-hidden rounded-2xl border border-(--border-color) bg-slate-900 p-8 text-white"
		>
			<svg viewBox="0 0 500 200" class="w-full">
				<!-- Resistance line -->
				<line
					x1="0"
					y1="40"
					x2="500"
					y2="40"
					stroke="#EB5757"
					stroke-width="2"
					stroke-dasharray="8 4"
				/>
				<text x="10" y="32" fill="#EB5757" class="text-[10px] font-bold"
					>RESISTANCE (Selling Zone)</text
				>

				<!-- Support line -->
				<line
					x1="0"
					y1="160"
					x2="500"
					y2="160"
					stroke="#27AE60"
					stroke-width="2"
					stroke-dasharray="8 4"
				/>
				<text x="10" y="180" fill="#27AE60" class="text-[10px] font-bold"
					>SUPPORT (Buying Zone)</text
				>

				<!-- Bouncing price path -->
				<path
					d="M0,100 L50,155 L100,50 L150,140 L200,45 L250,158 L300,55 L350,145 L400,42 L450,155 L500,80"
					fill="none"
					stroke="#6366f1"
					stroke-width="3"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
			</svg>
			<p class="mt-4 text-center text-[10px] font-bold tracking-widest uppercase opacity-50">
				Price oscillations between critical zones
			</p>
		</div>
	</section>
</div>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
