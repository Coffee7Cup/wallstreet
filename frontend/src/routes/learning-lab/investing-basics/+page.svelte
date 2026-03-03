<script>
	import { Calculator, BarChart2, PieChart, Activity, TrendingUp } from 'lucide-svelte';

	const ratios = [
		{
			name: 'P/E Ratio',
			full: 'Price to Earnings',
			formula: 'Share Price / EPS',
			use: 'Tells if stock is cheap or expensive relative to profit.'
		},
		{
			name: 'P/B Ratio',
			full: 'Price to Book',
			formula: 'Price / Book Value',
			use: 'How much you pay for company net assets.'
		},
		{
			name: 'ROE',
			full: 'Return on Equity',
			formula: 'Net Profit / Equity',
			use: "Efficiency of generating profit from shareholders' money."
		},
		{
			name: 'Dividend Yield',
			full: 'Yield',
			formula: 'Div per share / Price',
			use: 'Cash return on investment.'
		}
	];
</script>

<svelte:head>
	<title>Investing Basics — Learning Lab</title>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(39,174,96,0.12); color: #27AE60">Module 6</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			Investing Basics & Ratios
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			Fundamental analysis is the art of evaluating a stock by measuring its "Intrinsic Value."
			Financial ratios act as shorthand indicators of a company's health.
		</p>
	</header>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			The Four Pillars of Financial Analysis
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			While there are hundreds of ratios, beginners should master these four first. They help you
			compare companies of different sizes on an equal footing.
		</p>
		<div class="grid gap-6 sm:grid-cols-2">
			{#each ratios as r}
				<div
					class="group relative cursor-default overflow-hidden rounded-3xl border border-(--border-color) bg-(--bg-primary) p-8 transition-all hover:scale-[1.02]"
				>
					<div
						class="absolute -top-4 -right-4 text-[#27AE60] opacity-[0.03] transition-transform duration-500 group-hover:scale-125"
					>
						<Calculator size={140} />
					</div>
					<div class="mb-4 flex items-center gap-2 text-[#27AE60]">
						<Calculator size={20} />
						<span class="text-xs font-black tracking-[0.2em] uppercase">{r.name}</span>
					</div>
					<h3 class="mb-2 text-xl font-bold">{r.full}</h3>
					<div
						class="mb-5 inline-block rounded-full bg-green-500/10 px-4 py-1.5 font-mono text-xs font-bold tracking-tighter text-green-700 dark:text-green-400"
					>
						{r.formula}
					</div>
					<p class="text-base leading-relaxed opacity-70">{r.use}</p>
				</div>
			{/each}
		</div>
	</section>

	<section class="mb-12 border-t border-(--border-color) pt-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			The Two Major Investing Philosophies
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Most experienced investors usually lean toward one of these two camps, though many use a
			hybrid approach called GARP (Growth at a Reasonable Price).
		</p>
		<div class="mt-4 flex flex-col gap-8 md:flex-row">
			<div
				class="group relative flex-1 overflow-hidden rounded-3xl border border-green-500/20 bg-green-500/3 p-8"
			>
				<div
					class="absolute -top-8 -right-8 text-green-600 opacity-5 transition-transform duration-700 group-hover:rotate-12"
				>
					<TrendingUp size={160} />
				</div>
				<h3 class="mb-4 text-2xl font-bold text-green-600">Growth Investing</h3>
				<p class="text-base leading-relaxed opacity-80">
					Growth investors prioritize companies that are growing revenues and profits much faster
					than the rest of the market.
				</p>
				<ul class="mt-6 space-y-3 text-sm opacity-70">
					<li>• <b>Focus:</b> Future potential over current profits.</li>
					<li>• <b>Ratios:</b> High P/E is common (investors pay for future).</li>
					<li>• <b>Example:</b> Tech firms, new-age startups (Zomato, Tesla).</li>
				</ul>
			</div>
			<div
				class="group relative flex-1 overflow-hidden rounded-3xl border border-orange-500/20 bg-orange-500/3 p-8"
			>
				<div
					class="absolute -top-8 -right-8 text-orange-600 opacity-5 transition-transform duration-700 group-hover:-rotate-12"
				>
					<PieChart size={160} />
				</div>
				<h3 class="mb-4 text-2xl font-bold text-orange-600">Value Investing</h3>
				<p class="text-base leading-relaxed opacity-80">
					Value investors look for "hidden gems" or stable bargains. They want to buy ₹1 worth of
					assets for ₹0.70.
				</p>
				<ul class="mt-6 space-y-3 text-sm opacity-70">
					<li>• <b>Focus:</b> Strong fundamentals at a discount.</li>
					<li>• <b>Ratios:</b> Low P/E and Low P/B.</li>
					<li>• <b>Example:</b> Banks, PSU firms, mature utilities.</li>
				</ul>
			</div>
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
