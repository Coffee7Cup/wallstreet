<script>
	import { Landmark, Briefcase, TrendingUp, Search } from 'lucide-svelte';

	const exchanges = [
		{
			name: 'NSE',
			full: 'National Stock Exchange',
			year: 1992,
			tech: 'Electronic',
			color: '#1B5E20'
		},
		{ name: 'BSE', full: 'Bombay Stock Exchange', year: 1875, tech: 'Oldest', color: '#0D47A1' }
	];

	const indices = [
		{ name: 'Nifty 50', focus: '50 Large Cap Stocks', parent: 'NSE' },
		{ name: 'Sensex', focus: '30 Large Cap Stocks', parent: 'BSE' }
	];
</script>

<svelte:head>
	<title>Indian Stock Market — Learning Lab</title>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(242,153,74,0.12); color: #F2994A">Module 3</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			The Indian Stock Market
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			India has a massive, highly digitalized financial ecosystem. It is home to some of the world's
			oldest and most liquid stock exchanges.
		</p>
	</header>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			National vs. Bombay Stock Exchange
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Most stocks in India are "dual-listed," meaning they trade on both the NSE and BSE. While they
			are competitors, they both follow the same SEBI guidelines and market hours.
		</p>
		<div class="grid gap-6 sm:grid-cols-2">
			{#each exchanges as ex}
				<div
					role="button"
					tabindex="0"
					class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-8 transition-all"
					style="--item-color: {ex.color}"
					onmouseenter={(e) => (e.currentTarget.style.borderColor = ex.color)}
					onmouseleave={(e) => (e.currentTarget.style.borderColor = '')}
				>
					<div class="mb-4 flex items-center justify-between">
						<span
							class="rounded-lg px-4 py-1.5 text-xs font-black text-white"
							style="background: {ex.color}">{ex.name}</span
						>
						<span class="text-xs font-bold opacity-40">Est. {ex.year}</span>
					</div>
					<h3 class="mb-2 text-xl font-bold">{ex.full}</h3>
					<p class="text-sm leading-relaxed opacity-70">
						{ex.name === 'BSE'
							? 'Located on Dalal Street, it is the oldest exchange in Asia. It is famous for its Sensex index.'
							: 'The largest exchange in India by volume. It revolutionized the market with fully automated screen-based trading.'}
					</p>
				</div>
			{/each}
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			Indices: The Market's Pulse
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			An index is a statistical measure of the changes in a portfolio of stocks representing a
			portion of the overall market. When people say "the market is up," they usually mean these
			indices are up.
		</p>

		<div class="my-6 rounded-2xl border border-(--border-color) bg-(--bg-hover) p-6">
			<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
				{#each indices as idx}
					<div
						class="flex items-start gap-5 rounded-2xl border border-(--border-color) bg-(--bg-primary) p-6 transition-shadow hover:shadow-md"
					>
						<div class="rounded-xl bg-green-500/10 p-3 text-green-600">
							<TrendingUp size={28} />
						</div>
						<div>
							<h4 class="mb-1 text-lg font-bold">{idx.name} ({idx.parent})</h4>
							<p class="mb-3 text-xs font-bold tracking-wider uppercase opacity-50">
								Tracks {idx.focus}
							</p>
							<p class="text-sm leading-relaxed opacity-70">
								{idx.name === 'Nifty 50'
									? 'Includes 50 of the largest and most liquid Indian companies across 12 sectors.'
									: 'The first real-time stock index in India, tracking 30 blue-chip companies.'}
							</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">Market Regulation: SEBI</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			In the early 90s, the Indian market was prone to scams. SEBI was empowered to make sure
			investors are protected.
		</p>
		<div
			class="relative flex flex-col gap-6 overflow-hidden rounded-3xl border border-blue-500/20 bg-blue-500/5 p-8 md:flex-row"
		>
			<div class="absolute -right-8 -bottom-8 text-blue-500 opacity-5">
				<Landmark size={200} />
			</div>
			<div class="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl bg-blue-500/10">
				<Landmark size={32} class="text-blue-500" />
			</div>
			<div class="relative z-10">
				<h3 class="mb-3 text-xl font-bold">Securities and Exchange Board of India (SEBI)</h3>
				<p class="max-w-2xl text-base leading-relaxed opacity-80">
					SEBI acts as the watchdog of the Indian capital markets. Its primary roles include:
				</p>
				<ul
					class="mt-4 grid grid-cols-1 gap-x-8 gap-y-3 text-sm font-medium opacity-70 sm:grid-cols-2"
				>
					<li>• <b>Protecting</b> the interests of retail investors.</li>
					<li>• <b>Regulating</b> stock exchanges and brokers.</li>
					<li>• <b>Preventing</b> insider trading and market manipulation.</li>
					<li>• <b>Establishing</b> fair rules for companies going public.</li>
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
