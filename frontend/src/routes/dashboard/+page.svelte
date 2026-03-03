<script>
	import { onMount, tick } from 'svelte';
	import api from '$lib/api';
	import { token } from '$lib/stores/auth';
	import { marketState } from '$lib/stores/market';
	import {
		Search,
		TrendingUp,
		TrendingDown,
		Wallet,
		Briefcase,
		Clock,
		ArrowRight
	} from 'lucide-svelte';
	import gsap from 'gsap';
	import Sparkline from '$lib/components/Sparkline.svelte';

	let companies = $state([]);
	let searchQuery = $state('');
	let loading = $state(true);
	let portfolio = $state([]);

	let filteredCompanies = $derived(
		companies
			.map((c) => {
				const stockData = portfolio.find((p) => p.company_id === c.id);
				const marketData = $marketState.stocks.find((s) => s.company_id === c.id);
				const priceHistory = $marketState.priceHistory[c.id] || [];

				return {
					...c,
					price: marketData?.close_price || 0,
					owned: stockData?.quantity || 0,
					history: priceHistory.map((h) => h.price)
				};
			})
			.filter(
				(c) =>
					c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
					c.symbol.toLowerCase().includes(searchQuery.toLowerCase())
			)
	);

	$effect(() => {
		console.log('Dashboard: Rendered', filteredCompanies.length, 'companies after filter');
	});

	async function refreshUserData() {
		if (!$token) return;
		try {
			const portRes = await api.get('/trade/portfolio');
			portfolio = portRes.data.portfolio || [];
		} catch (err) {
			console.error('Failed to refresh user data', err);
		}
	}

	onMount(async () => {
		try {
			const [res, historyRes] = await Promise.all([
				api.get('/market/companies'),
				api.get('/market/stocks-history')
			]);
			console.log('Dashboard: Fetched', res.data.companies.length, 'companies');
			companies = res.data.companies;
			console.log('Dashboard: Companies state set', companies);

			// Populate marketState with historical data
			marketState.update((state) => {
				const newHistory = { ...state.priceHistory };
				historyRes.data.history.forEach((h) => {
					if (!newHistory[h.company_id]) newHistory[h.company_id] = [];
					newHistory[h.company_id].push({
						price: h.close_price,
						timestamp: h.date || h.active_time // date from DB
					});
				});
				return { ...state, priceHistory: newHistory };
			});

			loading = false;

			refreshUserData();

			// Wait for DOM update then animate
			await tick();
			gsap.to('.company-card', {
				opacity: 1,
				y: 0,
				startAt: { y: 20, opacity: 0 },
				stagger: 0.05,
				duration: 0.6,
				ease: 'power2.out'
			});
		} catch (err) {
			console.error('Dashboard load failed', err);
		}
	});

	// Refresh user data when market updates (to get updated portfolio after trades)
	$effect(() => {
		if ($marketState.tick) {
			refreshUserData();
		}
	});
</script>

<div class="space-y-8">
	<header class="flex flex-col justify-between gap-4 md:flex-row md:items-center">
		<div>
			<h1 class="text-3xl font-black tracking-tight">Market Overview</h1>
			<p class="text-sm opacity-50">Real-time stock updates • Tick {$marketState.tick}</p>
		</div>

		<div class="relative w-full md:w-96">
			<Search class="absolute top-1/2 left-4 -translate-y-1/2 opacity-30" size={18} />
			<input
				type="text"
				placeholder="Search companies or symbols..."
				bind:value={searchQuery}
				class="w-full rounded-2xl border border-(--border-color) bg-(--bg-primary) py-3 pr-4 pl-12 transition-all focus:border-[#27AE60] focus:ring-4 focus:ring-[#27AE60]/10 focus:outline-none"
			/>
		</div>
	</header>

	{#if loading}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each Array(6) as _}
				<div class="h-48 animate-pulse rounded-3xl bg-(--bg-hover)"></div>
			{/each}
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each filteredCompanies as company (company.id)}
				<a
					href="/dashboard/{company.id}"
					class="company-card group relative overflow-hidden rounded-3xl border border-(--border-color) bg-(--bg-primary) p-6 opacity-0 transition-all hover:-translate-y-1 hover:border-[#27AE60]/50 hover:shadow-2xl hover:shadow-[#27AE60]/5"
				>
					<div class="mb-4 flex items-start justify-between">
						<div>
							<h3 class="font-black group-hover:text-[#27AE60]">{company.symbol}</h3>
							<p class="line-clamp-1 text-sm opacity-40">{company.name}</p>
						</div>
						<div class="text-right">
							<p class="text-xl font-black">₹{company.price.toFixed(2)}</p>
							{#if company.history.length > 1}
								{@const change = company.price - company.history[0]}
								<p
									class="flex items-center justify-end gap-1 text-xs font-bold {change >= 0
										? 'text-[#27AE60]'
										: 'text-[#EB5757]'}"
								>
									{change >= 0 ? '+' : ''}{change.toFixed(2)}
									{#if change >= 0}
										<TrendingUp size={12} />
									{:else}
										<TrendingDown size={12} />
									{/if}
								</p>
							{/if}
						</div>
					</div>

					<div class="flex items-end justify-between">
						<div class="space-y-1">
							<p class="flex items-center gap-1.5 text-xs font-bold opacity-30">
								<Briefcase size={12} />
								{company.owned} Shares Owned
							</p>
							<p class="text-xs opacity-20">{company.sector}</p>
						</div>
						<Sparkline data={company.history} />
					</div>

					<div
						class="absolute bottom-0 left-0 h-1 w-0 bg-[#27AE60] transition-all duration-500 group-hover:w-full"
					></div>
				</a>
			{/each}
		</div>
	{/if}

	{#if !loading && filteredCompanies.length === 0}
		<div class="flex flex-col items-center justify-center py-20 opacity-30">
			<Search size={48} class="mb-4" />
			<p class="text-xl font-bold">No companies found match your search</p>
		</div>
	{/if}
</div>
