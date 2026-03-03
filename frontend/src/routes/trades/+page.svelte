<script>
	import { onMount } from 'svelte';
	import api from '$lib/api';
	import { marketState } from '$lib/stores/market';
	import {
		TrendingUp,
		TrendingDown,
		Calendar,
		Hash,
		Loader,
		Search,
		DollarSign,
		Briefcase,
		IndianRupee
	} from 'lucide-svelte';
	import gsap from 'gsap';

	let trades = $state([]);
	let loading = $state(true);
	let loadingMore = $state(false);
	let error = $state(null);
	let offset = $state(0);
	let limit = 10;
	let searchQuery = $state('');
	let hasMore = $state(true);

	async function fetchTrades(isLoadMore = false) {
		if (isLoadMore) {
			loadingMore = true;
		} else {
			loading = true;
			offset = 0;
		}

		try {
			const res = await api.get('/trade/trades', {
				params: {
					limit: limit,
					offset: offset,
					search: searchQuery
				}
			});
			const newTrades = res.data.trades || [];
			if (isLoadMore) {
				trades = [...trades, ...newTrades];
			} else {
				trades = newTrades;
			}
			hasMore = newTrades.length === limit;
		} catch (err) {
			console.error('Failed to fetch trades', err);
			error = err.response?.data?.error || 'Failed to load trades';
		} finally {
			loading = false;
			loadingMore = false;
		}
	}

	onMount(async () => {
		const [tradesRes, portRes] = await Promise.all([fetchTrades(), api.get('/trade/portfolio')]);

		marketState.update((s) => ({
			...s,
			userPortfolio: portRes.data.portfolio || []
		}));

		gsap.from('.trade-row', {
			y: 20,
			opacity: 0,
			stagger: 0.05,
			duration: 0.6,
			ease: 'power2.out'
		});
	});

	let debounceTimer;
	function handleSearch() {
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			fetchTrades();
		}, 1000);
	}

	function manualSearch() {
		clearTimeout(debounceTimer);
		fetchTrades();
	}

	function loadMore() {
		offset += limit;
		fetchTrades(true);
	}

	function formatDate(dateString) {
		return new Date(dateString).toLocaleString();
	}

	function calculatePL(trade) {
		const currentPrice =
			$marketState.stocks.find((s) => s.company_id == trade.company_id)?.close_price || trade.price;
		const pl = (currentPrice - trade.price) * trade.quantity;
		const plPercentage = ((currentPrice - trade.price) / trade.price) * 100;
		return {
			value: trade.trade_type === 'BUY' ? pl : -pl,
			percentage: trade.trade_type === 'BUY' ? plPercentage : -plPercentage
		};
	}

	let totalPL = $derived(trades.reduce((sum, t) => sum + calculatePL(t).value, 0));
	let positionValue = $derived(
		$marketState.userPortfolio.reduce((sum, item) => {
			const currentPrice =
				$marketState.stocks.find((s) => s.company_id == item.company_id)?.close_price ||
				item.current_price ||
				0;
			return sum + currentPrice * item.quantity;
		}, 0)
	);
</script>

<div class="space-y-6">
	<div class="flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
		<header>
			<h1 class="mb-2 text-3xl font-bold">My Trades</h1>
			<p class="opacity-60">View your complete trading history and performance</p>
		</header>

		<div class="flex gap-4">
			<div class="notion-card flex items-center gap-3 bg-(--bg-primary) px-6 py-4">
				<div class="rounded-full bg-(--bg-hover) p-2 text-[#27AE60]">
					<Briefcase size={20} />
				</div>
				<div>
					<p class="text-[10px] font-black tracking-widest uppercase opacity-40">Portfolio Value</p>
					<p class="text-xl font-bold">
						₹{positionValue.toLocaleString(undefined, {
							minimumFractionDigits: 2,
							maximumFractionDigits: 2
						})}
					</p>
				</div>
			</div>
			<div class="notion-card flex items-center gap-3 bg-(--bg-primary) px-6 py-4">
				<div
					class="rounded-full bg-(--bg-hover) p-2 {totalPL >= 0
						? 'text-[#27AE60]'
						: 'text-[#EB5757]'}"
				>
					{#if totalPL >= 0}
						<TrendingUp size={20} />
					{:else}
						<TrendingDown size={20} />
					{/if}
				</div>
				<div>
					<p class="text-[10px] font-black tracking-widest uppercase opacity-40">Total P/L</p>
					<p class="text-xl font-bold {totalPL >= 0 ? 'text-[#27AE60]' : 'text-[#EB5757]'}">
						{totalPL >= 0 ? '+' : ''}₹{Math.abs(totalPL).toLocaleString(undefined, {
							minimumFractionDigits: 2,
							maximumFractionDigits: 2
						})}
					</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Search Bar -->
	<div class="mb-6 flex max-w-lg gap-2">
		<div class="relative flex-1">
			<input
				type="text"
				bind:value={searchQuery}
				oninput={handleSearch}
				placeholder="Search by company symbol or trade type..."
				class="w-full rounded-lg border border-(--border-color) bg-(--bg-primary)
					   py-2.5 pr-4 pl-10 text-sm
					   transition-colors focus:ring-2 focus:ring-gray-300 focus:outline-none dark:focus:ring-gray-700"
			/>
		</div>
		<button
			onclick={manualSearch}
			class="rounded-md border border-(--border-color) bg-(--bg-primary) px-3 py-2 text-sm font-medium transition-colors hover:bg-(--bg-hover)"
		>
			<Search size={18} class="inline" />
		</button>
	</div>

	{#if loading && offset === 0}
		<div class="flex h-64 items-center justify-center">
			<Loader size={32} class="animate-spin opacity-40" />
		</div>
	{:else if error}
		<div class="notion-card bg-(--bg-primary) p-8 text-center text-[#EB5757]">
			{error}
		</div>
	{:else if trades.length === 0}
		<div class="notion-card bg-(--bg-primary) p-12 text-center">
			<TrendingUp size={48} class="mx-auto mb-4 opacity-30" />
			<p class="text-lg opacity-60">No trades found. Start trading to see your history here!</p>
		</div>
	{:else}
		<div class="notion-card overflow-hidden bg-(--bg-primary)">
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead class="border-b border-(--border-color) bg-(--bg-hover)">
						<tr>
							<th class="px-6 py-4 text-left text-sm font-bold tracking-wider uppercase opacity-60">
								<Hash size={16} class="mr-2 inline" /> ID
							</th>
							<th class="px-6 py-4 text-left text-sm font-bold tracking-wider uppercase opacity-60">
								Company
							</th>
							<th class="px-6 py-4 text-left text-sm font-bold tracking-wider uppercase opacity-60">
								Type
							</th>
							<th
								class="px-6 py-4 text-right text-sm font-bold tracking-wider uppercase opacity-60"
							>
								Quantity
							</th>
							<th
								class="px-6 py-4 text-right text-sm font-bold tracking-wider uppercase opacity-60"
							>
								Price
							</th>
							<th
								class="px-6 py-4 text-right text-sm font-bold tracking-wider uppercase opacity-60"
							>
								P/L
							</th>
							<th class="px-6 py-4 text-left text-sm font-bold tracking-wider uppercase opacity-60">
								<Calendar size={16} class="mr-2 inline" /> Date
							</th>
						</tr>
					</thead>
					<tbody>
						{#each trades as trade}
							{@const pl = calculatePL(trade)}
							<tr
								class="trade-row border-b border-(--border-color) transition-colors hover:bg-(--bg-hover)"
							>
								<td class="px-6 py-4 font-mono text-sm opacity-60">#{trade.id}</td>
								<td class="px-6 py-4">
									<div class="flex flex-col">
										<span class="font-bold">{trade.company_name}</span>
										<span class="text-xs opacity-40">{trade.company_symbol}</span>
									</div>
								</td>
								<td class="px-6 py-4">
									<span
										class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-bold uppercase {trade.trade_type ===
										'BUY'
											? 'bg-[#27AE60]/10 text-[#27AE60]'
											: 'bg-[#EB5757]/10 text-[#EB5757]'}"
									>
										{#if trade.trade_type === 'BUY'}
											<TrendingUp size={14} />
										{:else}
											<TrendingDown size={14} />
										{/if}
										{trade.trade_type}
									</span>
								</td>
								<td class="px-6 py-4 text-right font-bold">{trade.quantity}</td>
								<td class="px-6 py-4 text-right font-mono">₹{trade.price.toFixed(2)}</td>
								<td class="px-6 py-4 text-right">
									<div class="flex flex-col items-end">
										<span class="font-bold {pl.value >= 0 ? 'text-[#27AE60]' : 'text-[#EB5757]'}">
											{pl.value >= 0 ? '+' : ''}₹{Math.abs(pl.value).toFixed(2)}
										</span>
										<span
											class="text-xs {pl.percentage >= 0
												? 'text-[#27AE60]'
												: 'text-[#EB5757]'} opacity-80"
										>
											{pl.percentage >= 0 ? '+' : ''}{pl.percentage.toFixed(2)}%
										</span>
									</div>
								</td>
								<td class="px-6 py-4 text-sm opacity-60"
									>{trade.date.toLocaleString().split('T')[0]}</td
								>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

		<div class="flex items-center justify-between">
			<div class="text-sm opacity-60">
				Displaying <span class="font-bold">{trades.length}</span> trades
			</div>

			{#if hasMore}
				<button
					onclick={loadMore}
					disabled={loadingMore}
					class="notion-card bg-(--bg-primary) px-8 py-2 text-sm font-medium transition-colors hover:bg-(--bg-hover) disabled:opacity-50"
				>
					{#if loadingMore}
						<Loader size={16} class="mr-2 inline animate-spin" /> Loading...
					{:else}
						Load More
					{/if}
				</button>
			{/if}
		</div>
	{/if}
</div>
