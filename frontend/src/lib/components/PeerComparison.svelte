<script>
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import api from '$lib/api';
	import { GitCompare, TrendingUp, ArrowUpDown } from 'lucide-svelte';
	let d3Module;

	let { companyId } = $props();

	let peers = $state([]);
	let loading = $state(true);
	let error = $state('');
	let sortKey = $state('market_cap');
	let sortDir = $state(-1); // -1 = desc, 1 = asc

	let radarContainer = $state();

	// Colours for the radar chart lines (one per peer)
	const COLORS = [
		'#27AE60',
		'#3498DB',
		'#9B59B6',
		'#E67E22',
		'#E74C3C',
		'#1ABC9C',
		'#F1C40F',
		'#2ECC71'
	];

	// Sorted peers derived from state
	let sortedPeers = $derived(
		[...peers].sort((a, b) => {
			let av = getValue(a, sortKey);
			let bv = getValue(b, sortKey);
			return (av - bv) * sortDir;
		})
	);

	function computeMetrics(peer) {
		const totalShares = peer.company?.total_shares || 0;
		const marketCap = (peer.price ?? 0) * totalShares;

		const latestPL = peer.latest_pl;
		const latestBS = peer.latest_bs;
		const latestRatio = peer.latest_ratio;

		let peRatio = 0;
		let pbRatio = 0;
		let divYield = 0;
		let bookValue = 0;
		let roe = latestRatio?.roe || 0;
		let debtEquity = latestRatio?.debt_equity || 0;
		let opm = latestPL?.opm_percent || latestRatio?.opm || 0;
		let roce = latestRatio?.roce_percent || 0;

		if (latestPL && latestPL.net_profit > 0) {
			peRatio = marketCap / (latestPL.net_profit * 10000000);
			const totalDividend = latestPL.net_profit * (latestPL.dividend_payout / 100) * 10000000;
			divYield = (totalDividend / marketCap) * 100;
		}

		if (latestBS && totalShares > 0) {
			const bookValueTotal = latestBS.equity_capital + latestBS.reserves;
			bookValue = (bookValueTotal * 10000000) / totalShares;
			if (bookValue > 0) {
				pbRatio = (peer.price ?? 0) / bookValue;
			}
		}

		return {
			market_cap: marketCap,
			pe_ratio: peRatio,
			pb_ratio: pbRatio,
			dividend_yield: divYield,
			roe,
			debt_equity: debtEquity,
			opm,
			roce,
			book_value: bookValue
		};
	}

	function getValue(peer, key) {
		const m = computeMetrics(peer);
		switch (key) {
			case 'price':
				return peer.price ?? 0;
			case 'market_cap':
				return m.market_cap;
			case 'roe':
				return m.roe;
			case 'debt_equity':
				return m.debt_equity;
			case 'opm':
				return m.opm;
			case 'roce':
				return m.roce;
			case 'pe':
				return m.pe_ratio;
			default:
				return 0;
		}
	}

	function toggleSort(key) {
		if (sortKey === key) {
			sortDir = sortDir * -1;
		} else {
			sortKey = key;
			sortDir = -1;
		}
	}

	function fmt(v, decimals = 2) {
		return (v ?? 0).toFixed(decimals);
	}

	function fmtCr(v) {
		return '₹' + ((v ?? 0) / 10000000).toFixed(2) + 'Cr';
	}

	function computePE(peer) {
		return computeMetrics(peer).pe_ratio;
	}

	async function loadPeers() {
		try {
			const res = await api.get(`/market/companies/${companyId}/peers`);
			peers = res.data.peers ?? [];
		} catch (e) {
			error = 'Failed to load peer data';
			console.error(e);
		} finally {
			loading = false;
		}
	}

	function buildRadar() {
		if (!browser || !radarContainer || !d3Module || peers.length === 0) return;
		const d3 = d3Module;

		const width = radarContainer.clientWidth;
		const height = radarContainer.clientHeight;
		if (width < 50 || height < 50) return;

		// Clear previous
		d3.select(radarContainer).selectAll('*').remove();

		const radius = Math.min(width, height) / 2 - 40;
		const center = { x: width / 2, y: height / 2 };

		const metrics = [
			{ label: 'ROE', get: (p) => (computeMetrics(p).roe || 0) * 100 },
			{ label: 'OPM', get: (p) => (computeMetrics(p).opm || 0) * 100 },
			{ label: 'Debt/Eq', get: (p) => computeMetrics(p).debt_equity || 0 },
			{ label: 'ROCE', get: (p) => computeMetrics(p).roce || 0 }
		];

		const angleSlice = (Math.PI * 2) / metrics.length;

		function normalise(values) {
			const filtered = values.map((v) => (Number.isNaN(Number(v)) ? 0 : Number(v)));
			const max = Math.max(...filtered, 0.001);
			return filtered.map((v) => (v / max) * 100);
		}

		const svg = d3.select(radarContainer).append('svg').attr('width', width).attr('height', height);

		const g = svg.append('g').attr('transform', `translate(${center.x},${center.y})`);

		// Draw background circles
		const levels = 5;
		for (let i = 1; i <= levels; i++) {
			const r = (radius / levels) * i;
			g.append('circle')
				.attr('r', r)
				.attr('fill', 'none')
				.attr('stroke', 'rgba(255,255,255,0.05)')
				.attr('stroke-width', 1);
		}

		// Draw axes
		metrics.forEach((m, i) => {
			const angle = angleSlice * i - Math.PI / 2;
			const x = Math.cos(angle) * radius;
			const y = Math.sin(angle) * radius;

			g.append('line')
				.attr('x1', 0)
				.attr('y1', 0)
				.attr('x2', x)
				.attr('y2', y)
				.attr('stroke', 'rgba(255,255,255,0.05)')
				.attr('stroke-width', 1);

			g.append('text')
				.attr('x', Math.cos(angle) * (radius + 20))
				.attr('y', Math.sin(angle) * (radius + 20))
				.attr('text-anchor', 'middle')
				.attr('dominant-baseline', 'middle')
				.attr('fill', 'rgba(255,255,255,0.4)')
				.style('font-size', '10px')
				.style('font-weight', 'bold')
				.text(m.label);
		});

		// Draw blobs
		const radarLine = d3
			.lineRadial()
			.radius((d) => (d / 100) * radius)
			.angle((d, i) => i * angleSlice)
			.curve(d3.curveLinearClosed);

		peers.slice(0, 8).forEach((peer, i) => {
			const raw = metrics.map((m) => m.get(peer));
			const normalisedData = normalise(raw);
			const isCurrent = peer.company?.id == companyId;

			const color = COLORS[i % COLORS.length];

			const blob = g.append('g').attr('class', 'blob');

			blob
				.append('path')
				.datum(normalisedData)
				.attr('d', radarLine)
				.attr('fill', color)
				.attr('fill-opacity', 0.1)
				.attr('stroke', color)
				.attr('stroke-width', isCurrent ? 2.5 : 1)
				.style('cursor', 'pointer');

			// Points
			normalisedData.forEach((d, j) => {
				const angle = angleSlice * j - Math.PI / 2;
				const r = (d / 100) * radius;
				blob
					.append('circle')
					.attr('cx', Math.cos(angle) * r)
					.attr('cy', Math.sin(angle) * r)
					.attr('r', isCurrent ? 4 : 2)
					.attr('fill', color);
			});
		});
	}

	// Rebuild radar whenever peers data is ready
	$effect(() => {
		if (peers.length > 0 && radarContainer && d3Module) {
			buildRadar();
		}
	});

	onMount(async () => {
		d3Module = await import('d3');
		loadPeers();
	});

	onDestroy(() => {});

	const colDefs = [
		{ key: 'company', label: 'Company', sortable: false },
		{ key: 'price', label: 'Price', sortable: true },
		{ key: 'market_cap', label: 'Mkt Cap', sortable: true },
		{ key: 'pe', label: 'P/E', sortable: true },
		{ key: 'roe', label: 'ROE', sortable: true },
		{ key: 'debt_equity', label: 'D/E', sortable: true },
		{ key: 'opm', label: 'OPM', sortable: true },
		{ key: 'roce', label: 'ROCE', sortable: true }
	];
</script>

<section class="rounded-3xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-xl">
	<div class="mb-8 flex items-center gap-3">
		<GitCompare size={24} class="text-[#27AE60]" />
		<h2 class="text-xl font-bold">Peer Comparison</h2>
		{#if !loading && peers.length > 0}
			<span class="rounded-full bg-(--bg-hover) px-3 py-0.5 text-xs font-black opacity-50">
				{peers[0]?.company?.sector ?? ''}
			</span>
		{/if}
	</div>

	{#if loading}
		<div class="flex h-40 items-center justify-center">
			<div
				class="h-8 w-8 animate-spin rounded-full border-2 border-[#27AE60] border-t-transparent"
			></div>
		</div>
	{:else if error}
		<p class="py-8 text-center opacity-40">{error}</p>
	{:else if peers.length === 0}
		<p class="py-8 text-center opacity-40">No peers found in this sector.</p>
	{:else}
		<!-- Radar Chart -->
		<div class="mb-8 flex justify-center">
			<div class="h-[320px] w-full max-w-[480px]" bind:this={radarContainer}></div>
		</div>

		<!-- Comparison Table -->
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead>
					<tr class="border-b border-(--border-color)">
						{#each colDefs as col}
							<th
								class="py-3 pr-4 {col.sortable ? 'cursor-pointer select-none' : ''}"
								onclick={() => col.sortable && toggleSort(col.key)}
							>
								<span
									class="flex items-center gap-1 text-[10px] font-black uppercase opacity-40
									{col.sortable ? 'hover:opacity-80' : ''}
									{sortKey === col.key ? 'text-[#27AE60] !opacity-100' : ''}"
								>
									{col.label}
									{#if col.sortable}
										<ArrowUpDown
											size={10}
											class={sortKey === col.key ? 'opacity-100' : 'opacity-30'}
										/>
									{/if}
								</span>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody class="divide-y divide-(--border-color)/30">
					{#each sortedPeers as peer}
						{@const isCurrent = peer.company?.id == companyId}
						{@const pe = computePE(peer)}
						{@const m = computeMetrics(peer)}
						<tr
							class="cursor-pointer transition-colors hover:bg-(--bg-hover)/30
								{isCurrent ? 'bg-[#27AE60]/10' : ''}"
							onclick={() => goto(`/dashboard/${peer.company?.id}`)}
						>
							<td class="py-3 pr-4">
								<div class="flex items-center gap-2">
									<span
										class="rounded bg-(--bg-hover) px-1.5 py-0.5 text-[10px] font-black
										{isCurrent ? 'bg-[#27AE60]/20 text-[#27AE60]' : ''}"
									>
										{peer.company?.symbol}
									</span>
									<span class="font-medium opacity-80 {isCurrent ? 'text-[#27AE60]' : ''}">
										{peer.company?.name}
									</span>
									{#if isCurrent}
										<span
											class="rounded-full bg-[#27AE60]/20 px-2 py-0.5 text-[9px] font-black text-[#27AE60] uppercase"
										>
											You
										</span>
									{/if}
								</div>
							</td>
							<td class="py-3 pr-4 text-right font-mono">
								₹{fmt(peer.price)}
							</td>
							<td class="py-3 pr-4 text-right font-mono opacity-70">
								{fmtCr(m.market_cap)}
							</td>
							<td class="py-3 pr-4 text-right font-mono opacity-70">
								{fmt(m.pe_ratio)}
							</td>
							<td
								class="py-3 pr-4 text-right font-mono
								{(m.roe ?? 0) > 0.15 ? 'text-[#27AE60]' : (m.roe ?? 0) < 0.05 ? 'text-[#EB5757]' : ''}"
							>
								{fmt((m.roe ?? 0) * 100)}%
							</td>
							<td
								class="py-3 pr-4 text-right font-mono
								{(m.debt_equity ?? 0) > 1 ? 'text-[#EB5757]' : 'text-[#27AE60]'}"
							>
								{fmt(m.debt_equity)}
							</td>
							<td class="py-3 pr-4 text-right font-mono opacity-70">
								{fmt((m.opm ?? 0) * 100)}%
							</td>
							<td class="py-3 text-right font-mono text-[#27AE60]">
								{fmt(m.roce)}%
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<p class="mt-4 text-[10px] opacity-30">
			Click any row to navigate to that company's dashboard. Colour coding: ROE &gt;15% green /
			&lt;5% red · D/E &gt;1 red.
		</p>
	{/if}
</section>
