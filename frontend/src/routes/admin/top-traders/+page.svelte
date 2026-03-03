<script>
	import { onMount, onDestroy } from 'svelte';
	import api from '$lib/api';
	import { TrendingUp, Trophy, Users, RefreshCw, ArrowLeft } from 'lucide-svelte';
	import gsap from 'gsap';

	let topTraders = $state([]);
	let loading = $state(true);
	let error = $state(null);

	async function fetchTopTraders() {
		loading = true;
		error = null;
		try {
			const res = await api.get('/admin/stats/top-traders');
			topTraders = res.data || [];
		} catch (err) {
			console.error('Failed to fetch top traders', err);
			error = err.response?.data?.error || 'Failed to load top traders';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchTopTraders();

		gsap.from('.trader-row', {
			opacity: 0,
			x: -20,
			stagger: 0.05,
			duration: 0.5,
			ease: 'power2.out'
		});

		const interval = setInterval(() => {
			fetchTopTraders();
		}, 10000);

		onDestroy(() => clearInterval(interval));
	});

	function formatCurrency(value) {
		return new Intl.NumberFormat('en-IN', {
			style: 'currency',
			currency: 'INR',
			maximumFractionDigits: 0
		}).format(value);
	}
</script>

<div class="space-y-8">
	<header class="flex items-center justify-between">
		<div class="flex items-center gap-4">
			<a
				href="/admin"
				class="group flex h-10 w-10 items-center justify-center rounded-full border border-(--border-color) bg-(--bg-primary) transition-all hover:bg-(--bg-hover)"
			>
				<ArrowLeft size={20} class="transition-transform group-hover:-translate-x-1" />
			</a>
			<div>
				<h1 class="mb-1 text-3xl font-bold">Top Traders</h1>
				<p class="opacity-60">Ranking of all participants by total wealth</p>
			</div>
		</div>
		<button
			onclick={fetchTopTraders}
			disabled={loading}
			class="btn-green flex items-center gap-2 rounded px-4 py-2 font-bold transition-all hover:scale-105 active:scale-95"
		>
			<RefreshCw size={18} class={loading ? 'animate-spin' : ''} />
			Refresh
		</button>
	</header>

	{#if error}
		<div
			class="rounded-xl border border-[#EB5757]/20 bg-[#EB5757]/10 p-8 text-center text-[#EB5757]"
		>
			<p class="text-lg font-bold">Error loading rankings</p>
			<p class="mt-1 opacity-80">{error}</p>
			<button onclick={fetchTopTraders} class="mt-4 font-bold underline">Try again</button>
		</div>
	{:else if loading && topTraders.length === 0}
		<div class="flex h-64 flex-col items-center justify-center gap-4">
			<div
				class="h-12 w-12 animate-spin rounded-full border-4 border-(--border-color) border-t-[#27AE60]"
			></div>
			<p class="font-medium opacity-60">Calculating rankings...</p>
		</div>
	{:else if topTraders.length === 0}
		<div class="rounded-xl border border-(--border-color) bg-(--bg-primary) p-12 text-center">
			<Users size={48} class="mx-auto mb-4 opacity-20" />
			<p class="text-xl font-medium opacity-60">No participants registered yet.</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
			{#each topTraders.slice(0, 3) as trader, i}
				<div
					class="notion-card relative overflow-hidden bg-(--bg-primary) p-8 text-center transition-all hover:border-[#F2C94C]/90"
				>
					<div class="absolute -top-4 -right-4 opacity-50">
						<Trophy
							size={100}
							class={i === 0 ? 'text-[#F2C94C]' : i === 1 ? 'text-[#BDBDBD]' : 'text-[#CD7F32]'}
						/>
					</div>
					<div class="mb-4 flex items-center justify-center">
						<div
							class="flex h-16 w-16 items-center justify-center rounded-full border-2 border-(--border-color) bg-(--bg-hover) text-2xl font-bold shadow-inner"
						>
							#{trader.rank}
						</div>
					</div>
					<h3 class="mb-1 truncate text-2xl font-bold">{trader.username}</h3>
					<p class="text-3xl font-black text-[#27AE60]">{formatCurrency(trader.total_value)}</p>
					<div
						class="mt-4 flex justify-center gap-4 text-xs font-bold tracking-wider uppercase opacity-60"
					>
						<div>Cash: {formatCurrency(trader.cash_balance)}</div>
						<div>Stocks: {formatCurrency(trader.stock_value)}</div>
					</div>
				</div>
			{/each}
		</div>

		<div class="notion-card overflow-hidden bg-(--bg-primary) p-0">
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead class="border-b border-(--border-color) bg-(--bg-hover)">
						<tr>
							<th class="px-6 py-4 text-left text-xs font-bold tracking-wider uppercase opacity-60"
								>Rank</th
							>
							<th class="px-6 py-4 text-left text-xs font-bold tracking-wider uppercase opacity-60"
								>Username</th
							>
							<th class="px-6 py-4 text-right text-xs font-bold tracking-wider uppercase opacity-60"
								>Cash Balance</th
							>
							<th class="px-6 py-4 text-right text-xs font-bold tracking-wider uppercase opacity-60"
								>Stock Value</th
							>
							<th class="px-6 py-4 text-right text-xs font-bold tracking-wider uppercase opacity-60"
								>Total Value</th
							>
						</tr>
					</thead>
					<tbody class="divide-y divide-(--border-color)">
						{#each topTraders as trader}
							<tr class="trader-row group transition-colors hover:bg-(--bg-hover)">
								<td class="px-6 py-4">
									<span
										class="inline-flex h-8 w-8 items-center justify-center rounded-full border border-(--border-color) bg-(--bg-hover) text-sm font-bold transition-colors group-hover:border-[#27AE60]/50"
									>
										{trader.rank}
									</span>
								</td>
								<td class="px-6 py-4 text-lg font-bold">{trader.username}</td>
								<td class="px-6 py-4 text-right font-mono opacity-80"
									>{formatCurrency(trader.cash_balance)}</td
								>
								<td class="px-6 py-4 text-right font-mono opacity-80"
									>{formatCurrency(trader.stock_value)}</td
								>
								<td class="px-6 py-4 text-right font-mono text-lg font-bold text-[#27AE60]"
									>{formatCurrency(trader.total_value)}</td
								>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>
