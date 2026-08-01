<script>
	import { onMount, onDestroy, tick, untrack } from 'svelte';
	import { page } from '$app/state';
	import api from '$lib/api';
	import { user, admin, token } from '$lib/stores/auth';
	import { marketState, sendTrade } from '$lib/stores/market';
	import {
		TrendingUp,
		TrendingDown,
		Clock,
		Info,
		Wallet,
		ShoppingCart,
		ArrowLeft,
		DollarSign,
		Maximize2,
		Minimize2,
		Layers,
		Newspaper
	} from 'lucide-svelte';
	import gsap from 'gsap';
	import { browser } from '$app/environment';
	import StockChart from '$lib/components/StockChart.svelte';
	import PeerComparison from '$lib/components/PeerComparison.svelte';

	let companyId = $derived(page.params.id);
	let company = $state(null);
	let details = $state({ pl: [], bs: [], cf: [], ratios: [] });
	let loading = $state(true);
	let history = $state([]);
	let quantity = $state(1);
	let tradeLoading = $state(false);
	let tradeError = $state('');
	let tradeSuccess = $state('');
	let userData = $state({ balance: 0, owned: 0 });
	let chartType = $state('candle_solid'); // 'area' or 'candle_solid'
	let selectedYear = $state(null);
	let news = $state([]);
	let newsLoading = $state(false);
	let drawingTool = $state('none'); // 'none' | 'line' | 'rect' | 'triangle'
	let clearCount = $state(0);

	// Derived values from market store
	let currentPrice = $derived(
		$marketState.stocks.find((s) => s.company_id == companyId)?.close_price || company?.price || 0
	);

	let totalPrice = $derived((currentPrice * quantity).toFixed(2));
	let isUp = $derived(history.length > 1 ? currentPrice >= history[0].close : true);

	// Derived metrics calculated on the fly for EACH year
	let yearlyMetrics = $derived.by(() => {
		if (!company || !details.pl.length) return [];

		const totalShares = company.total_shares || 0;

		return details.pl.map((pl) => {
			const yearDate = new Date(pl.year);
			const yearInt = yearDate.getFullYear();

			// Find price at end of that financial year (March 31 or nearest)
			// For simplicity, we'll find the closest price to the 'year' timestamp in history
			const targetTs = yearDate.getTime();
			const closestPricePoint = history.reduce(
				(prev, curr) => {
					return Math.abs(curr.timestamp - targetTs) < Math.abs(prev.timestamp - targetTs)
						? curr
						: prev;
				},
				history[0] || { close: currentPrice }
			);

			const yearPrice = closestPricePoint.close;
			const marketCap = yearPrice * totalShares;

			const bs = details.bs.find((b) => new Date(b.year).getFullYear() === yearInt);
			const ratio = details.ratios.find((r) => new Date(r.year).getFullYear() === yearInt);

			let peRatio = 0;
			let pbRatio = 0;
			let divYield = 0;
			let roe = 0;
			let debtEquity = 0;
			let opm = 0;
			let bookValue = 0;

			if (bs) {
				const totalEquity = bs.equity_capital + bs.reserves;
				bookValue = (totalEquity * 10000000) / totalShares;
				if (bookValue > 0) {
					pbRatio = yearPrice / bookValue;
				}
				if (totalEquity > 0) {
					debtEquity = bs.borrowings / totalEquity;
					roe = pl.net_profit / totalEquity;
				}
			}

			if (pl.sales > 0) {
				opm = pl.operating_profit / pl.sales;
			}

			if (pl.net_profit > 0) {
				peRatio = marketCap / (pl.net_profit * 10000000);
				const totalDividend = pl.net_profit * (pl.dividend_payout / 100) * 10000000;
				divYield = (totalDividend / marketCap) * 100;
			}

			return {
				year: yearInt,
				market_cap: marketCap,
				pe_ratio: peRatio,
				pb_ratio: pbRatio,
				dividend_yield: divYield,
				roe,
				debt_equity: debtEquity,
				opm,
				book_value: bookValue
			};
		});
	});

	let latestMetrics = $derived(yearlyMetrics[0] || null);

	// Debug log for WebSocket connection status
	$effect(() => {
		// console.log('Market WS Connected:', $marketState.connected, 'Tick:', $marketState.tick);
	});

	async function loadInitialData() {
		loading = true; // Ensure loading is reset if called again
		try {
			const [compRes, plRes, bsRes, cfRes, ratioRes, histRes, portRes] = await Promise.all([
				api.get(`/market/companies/${companyId}`),
				api.get(`/market/profit-loss/${companyId}`),
				api.get(`/market/balance-sheets/${companyId}`),
				api.get(`/market/cash-flows/${companyId}`),
				api.get(`/market/ratios/${companyId}`),
				api.get(`/market/stocks-history/${companyId}`),
				api.get('/trade/portfolio')
			]);

			company = compRes.data.company;
			details.pl = plRes.data.profit_loss || [];
			details.bs = bsRes.data.balance_sheets || [];
			details.cf = cfRes.data.cash_flows || [];
			details.ratios = ratioRes.data.ratios || [];
			history = (histRes.data.history || []).map((h) => ({
				open: Number(h.open_price),
				high: Number(h.high_price),
				low: Number(h.low_price),
				close: Number(h.close_price),
				timestamp: new Date(h.date).getTime()
			}));

			if (details.pl.length > 0) {
				selectedYear = new Date(details.pl[0].year).getFullYear();
			}

			const portfolioItem = portRes.data.portfolio?.find((p) => p.company_id == companyId);
			userData.owned = portfolioItem?.quantity || 0;

			const profileRes = await api.get('/users/profile');
			userData.balance = profileRes.data.user.cash_balance;

			marketState.update((s) => ({
				...s,
				userPortfolio: portRes.data.portfolio || [],
				userBalance: profileRes.data.user.cash_balance
			}));

			loading = false;
			console.log('Data loaded, initializing chart and animations');

			// Use tick to ensure DOM is ready for animation
			await tick();

			gsap.to('.fade-in', {
				opacity: 1,
				y: 0,
				startAt: { opacity: 0, y: 10 },
				stagger: 0.1,
				duration: 0.5,
				ease: 'power2.out'
			});
		} catch (err) {
			console.error('Failed to load company detail', err);
			// Log specific error details if available
			if (err.response) {
				console.error('API Error Response:', err.response.status, err.response.data);
			}
		} finally {
			loading = false;
		}
	}

	// Reactively update chart when new market data arrives
	$effect(() => {
		if ($marketState.tick && company) {
			const newPoint = $marketState.stocks.find((s) => s.company_id == companyId);
			if (newPoint) {
				untrack(() => {
					const dateObj = new Date(newPoint.date || new Date().toISOString());
					const newEntry = {
						timestamp: dateObj.getTime(),
						open: Number(newPoint.open_price),
						high: Number(newPoint.high_price),
						low: Number(newPoint.low_price),
						close: Number(newPoint.close_price)
					};

					// Update our history array
					history = [...history, newEntry];
				});
			}
		}
	});

	async function fetchNews() {
		newsLoading = true;
		try {
			const res = await api.get('/market/news', {
				params: { company_id: companyId, limit: 10 }
			});
			news = res.data.news || [];
		} catch (err) {
			console.error('Failed to fetch news', err);
		} finally {
			newsLoading = false;
		}
	}

	onMount(async () => {
		loadInitialData();
		fetchNews();
	});

	// Derived year specific data for some displays if needed
	let currentPL = $derived(details.pl.find((p) => new Date(p.year).getFullYear() === selectedYear));
	let currentRatio = $derived(
		details.ratios.find((r) => new Date(r.year).getFullYear() === selectedYear)
	);

	// Reactively update userData when marketState updates
	$effect(() => {
		if ($marketState.userBalance !== undefined && $marketState.userBalance !== null) {
			userData.balance = $marketState.userBalance;
		}
		if ($marketState.userPortfolio && Array.isArray($marketState.userPortfolio)) {
			const portfolioItem = $marketState.userPortfolio.find((p) => p.company_id == companyId);
			userData.owned = portfolioItem?.quantity || 0;
		}
	});

	async function handleTrade(type) {
		if (!$user) {
			tradeError = 'User session not found. Please log in again.';
			return;
		}

		tradeLoading = true;
		tradeError = '';
		tradeSuccess = '';
		try {
			// Send trade via WebSocket
			sendTrade(parseInt(companyId), type, parseInt(quantity), $user.id);

			// Optimistic success message
			tradeSuccess = `${type} order for ${quantity} shares placed.`;

			// Refresh portfolio & balance from API as well to guarantee immediate UI update
			setTimeout(async () => {
				try {
					const [portRes, profileRes] = await Promise.all([
						api.get('/trade/portfolio'),
						api.get('/users/profile')
					]);
					if (portRes.data.portfolio && profileRes.data.user) {
						marketState.update((s) => ({
							...s,
							userPortfolio: portRes.data.portfolio,
							userBalance: profileRes.data.user.cash_balance
						}));
					}
				} catch (e) {
					console.error('Error refreshing portfolio after trade', e);
				}
			}, 300);

			// Clear success message after 3s
			setTimeout(() => {
				tradeSuccess = '';
			}, 3000);
		} catch (err) {
			tradeError = err.message || err.response?.data?.error || 'Trade failed';
		} finally {
			tradeLoading = false;
		}
	}

	onDestroy(() => {});
</script>

<div class="space-y-8">
	<a
		href="/dashboard"
		class="group flex items-center gap-2 text-sm opacity-50 transition-all hover:opacity-100"
	>
		<ArrowLeft size={16} class="transition-transform group-hover:-translate-x-1" /> Back to Market
	</a>

	{#if loading}
		<div class="flex h-96 items-center justify-center">
			<div
				class="h-8 w-8 animate-spin rounded-full border-2 border-[#27AE60] border-t-transparent"
			></div>
		</div>
	{:else if !company}
		<div class="flex h-96 flex-col items-center justify-center gap-4 text-center opacity-50">
			<Info size={48} />
			<div>
				<h2 class="text-xl font-bold">Company not found</h2>
				<p class="text-sm">Unable to load company details. Please try again.</p>
			</div>
			<a
				href="/dashboard"
				class="rounded bg-(--bg-hover) px-4 py-2 font-bold text-white transition-colors hover:bg-(--bg-hover)/80"
			>
				Back to Dashboard
			</a>
		</div>
	{:else}
		<header class="fade-in flex flex-col justify-between gap-6 lg:flex-row lg:items-end">
			<div>
				<div class="mb-2 flex items-center gap-3">
					<span
						class="rounded bg-(--bg-hover) px-2 py-1 text-xs font-black tracking-widest text-[#27AE60]"
					>
						{company.symbol}
					</span>
					<h1 class="text-4xl font-black">{company.name}</h1>
				</div>
				<p class="text-lg opacity-40">{company.sector} • {company.industry}</p>
			</div>
			<div class="text-right">
				<p class="text-xs font-black tracking-widest uppercase opacity-40">Current Price</p>
				<p class="text-5xl font-black text-[#27AE60]">₹{currentPrice.toFixed(2)}</p>
			</div>
		</header>

		<div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
			<div class="fade-in space-y-8 lg:col-span-2">
				<section
					class="relative h-[450px] w-full rounded-3xl border border-(--border-color) bg-(--bg-primary) p-6 shadow-xl"
				>
					<div class="mb-4 flex flex-wrap items-center justify-between gap-4">
						<div class="flex gap-2 rounded-xl bg-(--bg-hover)/50 p-1">
							<button
								onclick={() => (chartType = 'area')}
								class="rounded-lg px-4 py-1.5 text-xs font-black transition-all {chartType ===
								'area'
									? 'bg-[#27AE60] text-white shadow-lg'
									: 'opacity-40 hover:opacity-100'}"
							>
								AREA
							</button>
							<button
								onclick={() => (chartType = 'candle_solid')}
								class="rounded-lg px-4 py-1.5 text-xs font-black transition-all {chartType ===
								'candle_solid'
									? 'bg-[#27AE60] text-white shadow-lg'
									: 'opacity-40 hover:opacity-100'}"
							>
								CANDLE
							</button>
						</div>

						<div class="flex items-center gap-1 rounded-xl bg-(--bg-hover)/50 p-1">
							<!-- Drawing tools -->
							{#each [['line', '╱ LINE', '#F2C94C'], ['rect', '▭ RECT', '#2D9CDB'], ['triangle', '△ TRI', '#9B59B6']] as [toolId, label, color]}
								<button
									onclick={() => (drawingTool = drawingTool === toolId ? 'none' : toolId)}
									class="rounded-lg px-3 py-1.5 text-xs font-black transition-all
										{drawingTool === toolId ? 'text-white shadow-lg' : 'opacity-40 hover:opacity-100'}"
									style={drawingTool === toolId ? `background:${color}` : ''}
								>
									{label}
								</button>
							{/each}
							<div class="mx-1 h-4 w-px bg-white/10"></div>
							<button
								onclick={() => {
									clearCount++;
									drawingTool = 'none';
								}}
								class="rounded-lg px-3 py-1.5 text-xs font-black text-[#EB5757] opacity-60 transition-all hover:opacity-100"
								title="Clear all drawings"
							>
								✕ CLEAR
							</button>
						</div>
					</div>

					{#if drawingTool !== 'none'}
						<div
							class="absolute top-16 right-6 z-10 rounded-xl border border-white/10 bg-black/60 px-3 py-2 text-[10px] font-bold tracking-wide backdrop-blur-sm"
						>
							{#if drawingTool === 'line'}Click 2 points to draw a line
							{:else if drawingTool === 'rect'}Click 2 corners to draw a rectangle
							{:else if drawingTool === 'triangle'}Click 3 points to draw a triangle{/if}
						</div>
					{/if}

					<div class="h-[340px] w-full">
						<StockChart
							data={history}
							type={chartType}
							tool={drawingTool}
							{clearCount}
							id={companyId}
						/>
					</div>
				</section>
			</div>

			<aside class="fade-in space-y-6">
				<section class="rounded-3xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-xl">
					<h2 class="mb-6 flex items-center gap-2 text-xl font-bold">
						<ShoppingCart size={24} /> Trade
					</h2>

					<div class="mb-6 space-y-4 rounded-2xl bg-(--bg-hover)/30 p-4">
						<div class="flex items-center justify-between text-sm">
							<span class="flex items-center gap-1.5 opacity-50"><Wallet size={14} /> Balance</span>
							<span class="font-bold">₹{userData.balance.toLocaleString()}</span>
						</div>
						<div class="flex items-center justify-between text-sm">
							<span class="flex items-center gap-1.5 opacity-50"
								><ShoppingCart size={14} /> Position</span
							>
							<span class="font-bold">{userData.owned} shares</span>
						</div>
					</div>

					<div class="space-y-6">
						<div>
							<label for="qty" class="mb-2 block text-xs font-bold uppercase opacity-40"
								>Quantity</label
							>
							<input
								type="number"
								id="qty"
								bind:value={quantity}
								min="1"
								class="w-full rounded-2xl border border-(--border-color) bg-transparent px-5 py-4 text-2xl font-black focus:border-[#27AE60] focus:ring-4 focus:ring-[#27AE60]/10 focus:outline-none"
							/>
						</div>

						<div class="space-y-2">
							<div class="flex justify-between text-sm opacity-60">
								<span>Estimated Total</span>
								<span>₹{totalPrice}</span>
							</div>
							<div class="grid grid-cols-2 gap-4">
								<button
									onclick={() => handleTrade('BUY')}
									disabled={tradeLoading || userData.balance < currentPrice * quantity}
									class="rounded-2xl bg-[#27AE60] py-4 font-black text-white shadow-lg shadow-green-500/20 transition-all hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50"
								>
									BUY
								</button>
								<button
									onclick={() => handleTrade('SELL')}
									disabled={tradeLoading || userData.owned < quantity}
									class="rounded-2xl bg-[#EB5757] py-4 font-black text-white shadow-lg shadow-red-500/20 transition-all hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50"
								>
									SELL
								</button>
							</div>
						</div>

						{#if tradeError}
							<div
								class="animate-shake mt-4 rounded-xl bg-[#EB5757] p-4 text-sm font-bold text-white shadow-lg shadow-red-500/20"
							>
								{tradeError}
							</div>
						{/if}

						{#if tradeSuccess}
							<div
								class="mt-4 rounded-xl bg-[#27AE60] p-4 text-sm font-bold text-white shadow-lg shadow-green-500/20"
							>
								{tradeSuccess}
							</div>
						{/if}
					</div>
				</section>
			</aside>
			<div class="fade-in flex w-full flex-col gap-8 lg:col-span-3">
				<!-- Ratios and Cards -->
				<section class="rounded-3xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-xl">
					<div class="mb-8 flex items-center justify-between">
						<h2 class="flex items-center gap-2 text-xl font-bold">
							<TrendingUp size={24} /> Key Ratios
						</h2>
						{#if latestMetrics}
							<span class="text-xs font-bold opacity-40">FY {latestMetrics.year}</span>
						{/if}
					</div>

					{#if latestMetrics}
						<div class="mb-10 grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
							<div
								class="group flex flex-col items-center justify-center rounded-2xl border border-white/5 bg-white/5 py-4 transition-all hover:bg-white/10"
							>
								<span class="text-[10px] font-black tracking-widest text-[#27AE60] uppercase"
									>P/E Ratio</span
								>
								<span class="text-xl font-black">{latestMetrics.pe_ratio.toFixed(2)}</span>
							</div>
							<div
								class="group flex flex-col items-center justify-center rounded-2xl border border-white/5 bg-white/5 py-4 transition-all hover:bg-white/10"
							>
								<span class="text-[10px] font-black tracking-widest text-[#3498DB] uppercase"
									>P/B Ratio</span
								>
								<span class="text-xl font-black">{latestMetrics.pb_ratio.toFixed(2)}</span>
							</div>
							<div
								class="group flex flex-col items-center justify-center rounded-2xl border border-white/5 bg-white/5 py-4 transition-all hover:bg-white/10"
							>
								<span class="text-[10px] font-black tracking-widest text-[#9B59B6] uppercase"
									>ROE</span
								>
								<span class="text-xl font-black">{(latestMetrics.roe * 100).toFixed(2)}%</span>
							</div>
							<div
								class="group flex flex-col items-center justify-center rounded-2xl border border-white/5 bg-white/5 py-4 transition-all hover:bg-white/10"
							>
								<span class="text-[10px] font-black tracking-widest text-[#F2C94C] uppercase"
									>Debt/Equity</span
								>
								<span class="text-xl font-black">{latestMetrics.debt_equity.toFixed(2)}</span>
							</div>
							<div
								class="group flex flex-col items-center justify-center rounded-2xl border border-white/5 bg-white/5 py-4 transition-all hover:bg-white/10"
							>
								<span class="text-[10px] font-black tracking-widest text-[#E67E22] uppercase"
									>OPM</span
								>
								<span class="text-xl font-black">{(latestMetrics.opm * 100).toFixed(2)}%</span>
							</div>
							<div
								class="group flex flex-col items-center justify-center rounded-2xl border border-white/5 bg-white/5 py-4 transition-all hover:bg-white/10"
							>
								<span class="text-[10px] font-black tracking-widest text-[#27AE60] uppercase"
									>Div. Yield</span
								>
								<span class="text-xl font-black">{latestMetrics.dividend_yield.toFixed(2)}%</span>
							</div>
						</div>
					{/if}

					<!-- <div class="overflow-x-auto">
						<table class="w-full text-left text-sm">
							<thead>
								<tr class="border-b border-(--border-color)">
									<th class="py-3 font-bold opacity-40">Year</th>
									<th class="py-3 text-right font-bold opacity-40">P/E</th>
									<th class="py-3 text-right font-bold opacity-40">P/B</th>
									<th class="py-3 text-right font-bold opacity-40">ROE</th>
									<th class="py-3 text-right font-bold opacity-40">D/E</th>
									<th class="py-3 text-right font-bold opacity-40">OPM</th>
									<th class="py-3 text-right font-bold opacity-40">Div. Yield</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-(--border-color)/50">
								{#each yearlyMetrics as m}
									<tr>
										<td class="py-3 font-bold">{m.year}</td>
										<td class="py-3 text-right font-mono">{m.pe_ratio.toFixed(2)}</td>
										<td class="py-3 text-right font-mono">{m.pb_ratio.toFixed(2)}</td>
										<td class="py-3 text-right font-mono">{(m.roe * 100).toFixed(2)}%</td>
										<td class="py-3 text-right font-mono">{m.debt_equity.toFixed(2)}</td>
										<td class="py-3 text-right font-mono">{(m.opm * 100).toFixed(2)}%</td>
										<td class="py-3 text-right font-mono">{m.dividend_yield.toFixed(2)}%</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div> -->
				</section>

				<section class="rounded-3xl border border-(--border-color) bg-(--bg-primary) p-6">
					<h2 class="mb-6 flex items-center gap-2 text-lg font-bold">
						<Info size={20} /> Fundamentals
					</h2>
					<div class="overflow-x-auto">
						<table class="w-full text-left text-sm">
							<thead>
								<tr class="border-b border-(--border-color)">
									<th class="py-3 font-bold opacity-40">Attribute</th>
									{#each details.pl as p}
										<th class="py-3 text-right font-bold opacity-40">
											{new Date(p.year).getFullYear()}
										</th>
									{/each}
								</tr>
							</thead>
							<tbody class="divide-y divide-(--border-color)/50">
								<!-- Profit & Loss -->
								<tr>
									<td class="py-4 font-bold text-[#27AE60]" colspan={details.pl.length + 1}
										>Profit & Loss</td
									>
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Sales</td>
									{#each details.pl as p}
										<td class="py-3 text-right font-mono">₹{(p.sales / 100).toFixed(2)}Cr</td>
									{/each}
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Operating Profit</td>
									{#each details.pl as p}
										<td class="py-3 text-right font-mono"
											>₹{(p.operating_profit / 100).toFixed(2)}Cr</td
										>
									{/each}
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Net Profit</td>
									{#each details.pl as p}
										<td class="py-3 text-right font-mono">₹{(p.net_profit / 100).toFixed(2)}Cr</td>
									{/each}
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">EPS</td>
									{#each details.pl as p}
										<td class="py-3 text-right font-mono">₹{p.eps.toFixed(2)}</td>
									{/each}
								</tr>
								<!-- Balance Sheet -->
								<tr>
									<td class="py-4 font-bold text-[#27AE60]" colspan={details.bs.length + 1}
										>Balance Sheet</td
									>
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Equity Capital</td>
									{#each details.bs as b}
										<td class="py-3 text-right font-mono"
											>₹{(b.equity_capital / 100).toFixed(2)}Cr</td
										>
									{/each}
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Reserves</td>
									{#each details.bs as b}
										<td class="py-3 text-right font-mono">₹{(b.reserves / 100).toFixed(2)}Cr</td>
									{/each}
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Borrowings</td>
									{#each details.bs as b}
										<td class="py-3 text-right font-mono">₹{(b.borrowings / 100).toFixed(2)}Cr</td>
									{/each}
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Total Assets</td>
									{#each details.bs as b}
										<td class="py-3 text-right font-mono">₹{(b.total_assets / 100).toFixed(2)}Cr</td
										>
									{/each}
								</tr>
								<!-- Cash Flows -->
								<tr>
									<td class="py-4 font-bold text-[#27AE60]" colspan={details.cf.length + 1}
										>Cash Flows</td
									>
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Operating Activity</td>
									{#each details.cf as c}
										<td class="py-3 text-right font-mono"
											>₹{(c.cash_from_operating_activity / 100).toFixed(2)}Cr</td
										>
									{/each}
								</tr>
								<tr>
									<td class="py-3 pl-4 opacity-40">Net Cash Flow</td>
									{#each details.cf as c}
										<td class="py-3 text-right font-mono"
											>₹{(c.net_cash_flow / 100).toFixed(2)}Cr</td
										>
									{/each}
								</tr>
							</tbody>
						</table>
					</div>
				</section>

				<!-- News & Reports Section -->
				<div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
					<section
						class="rounded-3xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-xl"
					>
						<h2 class="mb-6 flex items-center gap-2 text-xl font-bold">
							<Newspaper size={24} /> Latest News
						</h2>

						{#if newsLoading && news.length === 0}
							<div class="space-y-4">
								{#each Array(3) as _}
									<div class="h-20 w-full animate-pulse rounded-2xl bg-(--bg-hover)/20"></div>
								{/each}
							</div>
						{:else if news.length === 0}
							<p class="py-8 text-center opacity-40">No recent news available for this company.</p>
						{:else}
							<div class="custom-scrollbar max-h-[400px] space-y-4 overflow-y-auto pr-2">
								{#each news as n}
									<div
										class="group cursor-pointer rounded-2xl bg-(--bg-hover)/10 p-4 transition-all hover:bg-(--bg-hover)/20"
									>
										<p class="mb-1 text-[10px] font-black tracking-widest text-[#27AE60] uppercase">
											{new Date(n.release_date).toLocaleDateString('en-US', {
												month: 'short',
												day: 'numeric',
												year: 'numeric'
											})}
										</p>
										<h3 class="leading-snug font-bold transition-colors group-hover:text-[#27AE60]">
											{n.title}
										</h3>
									</div>
								{/each}
							</div>
						{/if}
					</section>

					<section
						class="rounded-3xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-xl"
					>
						<div class="mb-6 flex items-center justify-between">
							<h2 class="flex items-center gap-2 text-xl font-bold">
								<Layers size={24} /> Financial Trends
							</h2>
						</div>

						<div class="space-y-6">
							<div>
								<div class="mb-2 flex justify-between text-xs font-bold uppercase opacity-40">
									<span>Sales Growth (Crores)</span>
									<span>Trend</span>
								</div>
								<div class="flex h-12 items-end gap-1.5 px-2">
									{#each details.pl.slice().reverse() as p}
										<div
											class="group relative w-full rounded-t-md bg-[#27AE60]/20 transition-all hover:bg-[#27AE60]/40"
											style="height: {Math.max(
												20,
												(p.sales / Math.max(...details.pl.map((pl) => pl.sales))) * 100
											)}%"
										>
											<div
												class="pointer-events-none absolute -top-8 left-1/2 -translate-x-1/2 rounded bg-slate-900 px-2 py-1 text-[10px] font-bold text-white opacity-0 transition-opacity group-hover:opacity-100"
											>
												₹{(p.sales / 100).toFixed(0)}Cr
											</div>
										</div>
									{/each}
								</div>
							</div>

							<div>
								<div class="mb-2 flex justify-between text-xs font-bold uppercase opacity-40">
									<span>Net Profit Trend</span>
									<span>Trend</span>
								</div>
								<div class="flex h-12 items-end gap-1.5 px-2">
									{#each details.pl.slice().reverse() as p}
										<div
											class="group relative w-full rounded-t-md bg-[#3498DB]/20 transition-all hover:bg-[#3498DB]/40"
											style="height: {Math.max(
												20,
												(p.net_profit / Math.max(...details.pl.map((pl) => pl.net_profit))) * 100
											)}%"
										>
											<div
												class="pointer-events-none absolute -top-8 left-1/2 -translate-x-1/2 rounded bg-slate-900 px-2 py-1 text-[10px] font-bold text-white opacity-0 transition-opacity group-hover:opacity-100"
											>
												₹{(p.net_profit / 100).toFixed(0)}Cr
											</div>
										</div>
									{/each}
								</div>
							</div>
						</div>
					</section>
				</div>

				<!-- Peer Comparison Section -->
				<PeerComparison {companyId} />
			</div>
		</div>
	{/if}
</div>

<style>
	@keyframes shake {
		0%,
		100% {
			transform: translateX(0);
		}
		25% {
			transform: translateX(-4px);
		}
		75% {
			transform: translateX(4px);
		}
	}
	.animate-shake {
		animation: shake 0.2s ease-in-out 0s 2;
	}
</style>
