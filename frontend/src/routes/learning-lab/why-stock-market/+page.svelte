<script>
	import { onMount } from 'svelte';
	import { TrendingUp, Flame, Shield, Clock } from 'lucide-svelte';

	// Inflation chart - interactive SVG
	let inflationCanvas;
	let selectedYear = $state(null);

	const data = [
		{ year: 2015, inflation: 5.9, investment: 100 },
		{ year: 2016, inflation: 4.9, investment: 112 },
		{ year: 2017, inflation: 3.6, investment: 128 },
		{ year: 2018, inflation: 3.4, investment: 148 },
		{ year: 2019, inflation: 4.8, investment: 171 },
		{ year: 2020, inflation: 6.2, investment: 196 },
		{ year: 2021, inflation: 5.1, investment: 229 },
		{ year: 2022, inflation: 6.7, investment: 268 },
		{ year: 2023, inflation: 5.4, investment: 310 },
		{ year: 2024, inflation: 4.6, investment: 361 }
	];

	// chart dimensions
	const W = 520,
		H = 200,
		PAD = { top: 20, right: 20, bottom: 36, left: 44 };
	const cW = W - PAD.left - PAD.right;
	const cH = H - PAD.top - PAD.bottom;

	const minInv = 0,
		maxInv = 400;
	const xScale = (i) => PAD.left + (i / (data.length - 1)) * cW;
	const yScale = (v) => PAD.top + cH - ((v - minInv) / (maxInv - minInv)) * cH;

	// Build SVG path
	let linePath = $state('');
	let areaPath = $state('');
	$effect(() => {
		linePath = data
			.map((d, i) => `${i === 0 ? 'M' : 'L'}${xScale(i)},${yScale(d.investment)}`)
			.join(' ');
		areaPath = `${linePath} L${xScale(data.length - 1)},${PAD.top + cH} L${xScale(0)},${PAD.top + cH} Z`;
	});

	const reasons = [
		{
			icon: Flame,
			color: '#EB5757',
			title: 'Beat Inflation',
			desc: 'Money kept in savings loses value every year due to inflation (~5-7% in India). Equity markets have historically returned 12–15% CAGR, far outpacing inflation.'
		},
		{
			icon: TrendingUp,
			color: '#27AE60',
			title: 'Wealth Creation',
			desc: 'Compounding turns small regular investments into substantial wealth. ₹10,000/month invested for 20 years at 12% CAGR grows to over ₹1 crore.'
		},
		{
			icon: Shield,
			color: '#2D9CDB',
			title: 'Financial Security',
			desc: 'Building an investment portfolio provides a safety net for emergencies, retirement, and life goals like education or buying a home.'
		},
		{
			icon: Clock,
			color: '#9B59B6',
			title: 'Passive Income',
			desc: 'Dividend-paying stocks and REITs generate regular income without active work. Companies like TCS, Infosys and HDFC Bank pay consistent dividends.'
		}
	];

	// Compounding visualizer
	const compoundingData = [
		{
			years: 5,
			value: 4.1,
			label: '5 Years',
			text: 'Short-term saving period. Gains are mostly from your own contributions.'
		},
		{
			years: 10,
			value: 11.6,
			label: '10 Years',
			text: 'Compounding starts to kick in. Interest begins to earn interest.'
		},
		{
			years: 20,
			value: 49.9,
			label: '20 Years',
			text: 'The "hockey stick" curve. Gains outpace contributions significantly.'
		},
		{
			years: 30,
			value: 176.5,
			label: '30 Years',
			text: 'Massive wealth creation. Small regular amounts turn into a multi-crore corpus.'
		}
	];
</script>

<svelte:head>
	<title>Why Stock Market? — Learning Lab</title>
	<meta
		name="description"
		content="Understand why people invest in stock markets — inflation, wealth creation, and financial freedom."
	/>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(39,174,96,0.12); color: #27AE60">Module 1</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			Why Invest in the Stock Market?
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			Most people keep money in savings accounts, but inflation silently erodes its value. The stock
			market offers a powerful way to grow wealth and achieve financial freedom through ownership in
			thriving businesses.
		</p>
	</header>

	<!-- Interactive Inflation Chart -->
	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-3xl font-extrabold">
			The Inflation Problem — Visualized
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Inflation is the rate at which prices for goods and services rise. If inflation is 6%, your
			₹100 today will only buy ₹94 worth of goods next year. Equity markets have historically
			returned 12–15% CAGR, effectively shielding your wealth. <b
				>Hover over the chart to explore the growth of a ₹100 investment since 2015.</b
			>
		</p>

		<div class="relative rounded-2xl border border-(--border-color) bg-(--bg-hover) p-5">
			<div class="mb-3 flex items-center gap-2 text-[0.75rem] opacity-60">
				<span class="h-2.5 w-2.5 shrink-0 rounded-full" style="background: #27AE60"></span> Nifty 50 Investment
				Value (₹ indexed to 100)
			</div>
			<svg
				viewBox="0 0 {W} {H}"
				class="w-full overflow-visible"
				role="img"
				aria-label="Nifty 50 growth chart"
			>
				<!-- Area fill -->
				<defs>
					<linearGradient id="greenGrad" x1="0" y1="0" x2="0" y2="1">
						<stop offset="0%" stop-color="#27AE60" stop-opacity="0.25" />
						<stop offset="100%" stop-color="#27AE60" stop-opacity="0.02" />
					</linearGradient>
				</defs>
				<path d={areaPath} fill="url(#greenGrad)" />
				<path
					d={linePath}
					fill="none"
					stroke="#27AE60"
					stroke-width="2.5"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>

				<!-- Y-axis labels -->
				{#each [0, 100, 200, 300, 400] as v}
					<text
						x={PAD.left - 6}
						y={yScale(v) + 4}
						text-anchor="end"
						class="fill-current text-[10px] opacity-45 select-none"
					>
						₹{v}
					</text>
					<line
						x1={PAD.left}
						y1={yScale(v)}
						x2={W - PAD.right}
						y2={yScale(v)}
						stroke="currentColor"
						stroke-opacity="0.06"
						stroke-dasharray="4 4"
					/>
				{/each}

				<!-- Data points + hover -->
				{#each data as d, i}
					<g
						class="cursor-pointer"
						onmouseenter={() => (selectedYear = d)}
						onmouseleave={() => (selectedYear = null)}
					>
						<circle
							cx={xScale(i)}
							cy={yScale(d.investment)}
							r={selectedYear?.year === d.year ? 7 : 5}
							fill="#27AE60"
							stroke="var(--bg-primary)"
							stroke-width="2.5"
							class="transition-[r] duration-150"
						/>
						<text
							x={xScale(i)}
							y={H - 8}
							text-anchor="middle"
							class="fill-current text-[10px] opacity-45 select-none"
						>
							{d.year}
						</text>
					</g>
				{/each}
			</svg>

			{#if selectedYear}
				<div
					class="mt-3 inline-flex items-center gap-6 rounded-[10px] border border-(--border-color) bg-(--bg-primary) p-3 px-4 text-[0.8rem]"
				>
					<div>
						<span class="block text-[10px] font-bold uppercase opacity-50">Year</span>
						<strong class="text-base">{selectedYear.year}</strong>
					</div>
					<div>
						<span class="block text-[10px] font-bold uppercase opacity-50">Value</span>
						<b style="color:#27AE60" class="text-base">₹{selectedYear.investment}</b>
					</div>
					<div>
						<span class="block text-[10px] font-bold uppercase opacity-50">Inflation</span>
						<b style="color:#EB5757" class="text-base">{selectedYear.inflation}%</b>
					</div>
				</div>
			{:else}
				<div class="flex h-[54px] items-center justify-center text-xs italic opacity-30">
					Hover over dots to see details
				</div>
			{/if}
		</div>
	</section>

	<!-- Why invest -->
	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-3xl font-extrabold">
			Four Powerful Reasons to Invest
		</h2>
		<div class="grid grid-cols-[repeat(auto-fill,minmax(240px,1fr))] gap-4">
			{#each reasons as r}
				<div
					role="button"
					tabindex="0"
					class="group rounded-[14px] border border-(--border-color) bg-(--bg-primary) p-5 transition-all duration-200"
					style="--rc: {r.color}"
					onmouseenter={(e) => {
						e.currentTarget.style.borderColor = r.color;
						e.currentTarget.style.boxShadow = `0 4px 16px color-mix(in srgb, ${r.color} 12%, transparent)`;
					}}
					onmouseleave={(e) => {
						e.currentTarget.style.borderColor = '';
						e.currentTarget.style.boxShadow = '';
					}}
				>
					<div
						class="mb-3 flex h-10 w-10 items-center justify-center rounded-[10px] text-white"
						style="background: color-mix(in srgb, {r.color} 14%, transparent); color: {r.color}"
					>
						<r.icon size={22} />
					</div>
					<h3 class="mb-1 text-2xl font-bold">{r.title}</h3>
					<p class="text-lg leading-relaxed opacity-60">{r.desc}</p>
				</div>
			{/each}
		</div>
	</section>

	<!-- Compounding visualizer -->
	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-3xl font-extrabold">The Power of Compounding</h2>
		<p class="mb-8 max-w-[800px] text-lg leading-relaxed opacity-80">
			Albert Einstein called compounding the "8th wonder of the world." Below is the growth of a
			₹5,000 monthly investment at a 12% annual return. Notice how the gap between 20 and 30 years
			is much larger than the gap between 5 and 10 years.
		</p>

		<div class="flex flex-col gap-6">
			{#each compoundingData as bar}
				<div class="flex flex-col">
					<div class="flex items-center gap-4">
						<div class="w-[70px] shrink-0 text-[0.75rem] font-semibold opacity-60">
							{bar.label}
						</div>
						<div class="h-8 flex-1 overflow-hidden rounded-full bg-(--bg-hover)">
							<div
								class="flex h-full min-w-[60px] items-center rounded-full pl-3 text-[0.75rem] font-bold text-white transition-[width] duration-800 ease-[cubic-bezier(0.22,1,0.36,1)]"
								style="width: {(bar.value / 176.5) * 100}%; background: #27AE60"
							>
								<span>₹{bar.value} Lakhs</span>
							</div>
						</div>
					</div>
					<p class="text-md mt-1 ml-[86px] opacity-50">{bar.text}</p>
				</div>
			{/each}
		</div>
	</section>

	<section class="mb-12 border-t border-(--border-color) pt-8">
		<h2 class="mb-5 flex items-center gap-3 text-3xl font-extrabold">Active vs. Passive Growth</h2>
		<p class="mb-8 max-w-[800px] text-lg leading-relaxed opacity-80">
			Investing isn't just about picking specific stocks. Broad market participation allows you to
			benefit from the collective success of a nation's economy.
		</p>
		<div class="rounded-2xl border border-(--border-color) bg-(--bg-hover) p-6">
			<ul class="space-y-4 text-lg text-sm leading-relaxed opacity-80">
				<li>
					• <b>Business Ownership:</b> When you buy a share of Reliance or HDFC, you become a partial
					owner. As they profit and grow, you profit.
				</li>
				<li>
					• <b>Capital Gains:</b> The increase in the share price over time as the company's valuation
					increases.
				</li>
				<li>
					• <b>Dividend Income:</b> Direct cash payments made by companies to their shareholders from
					their earnings.
				</li>
			</ul>
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
