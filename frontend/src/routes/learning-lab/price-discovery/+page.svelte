<script>
	import { DollarSign, MessageCircle, Newspaper, TrendingUp, TrendingDown } from 'lucide-svelte';

	let demand = $state(50);
	let supply = $state(50);

	let price = $derived(100 + (demand - supply));
	let sentiment = $derived(demand > supply ? 'Bullish' : demand < supply ? 'Bearish' : 'Neutral');
</script>

<svelte:head>
	<title>Price Discovery — Learning Lab</title>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(235,87,87,0.12); color: #EB5757">Module 5</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			Price Discovery
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			Price discovery is the continuous process by which the market determines the "fair price" of
			an asset based on the tug-of-war between buyers and sellers.
		</p>
	</header>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			Demand vs. Supply: The Interactive Tug-Of-War
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Every second, thousands of buy and sell orders hit the exchange. When more people want to buy
			than sell (High Demand), the price goes up. When more people want to sell than buy (High
			Supply), the price falls. <b>Play with the sliders below to see this in action.</b>
		</p>

		<div class="my-6 rounded-2xl border border-(--border-color) bg-(--bg-hover) p-6">
			<div class="space-y-10 p-4">
				<div class="flex flex-col items-center gap-12 md:flex-row">
					<div class="w-full flex-1 space-y-8">
						<div>
							<div class="mb-3 flex items-end justify-between uppercase">
								<span class="text-sm font-black tracking-widest text-blue-500">Buyers (Demand)</span
								>
								<span class="text-2xl font-black text-blue-500">{demand}</span>
							</div>
							<input
								type="range"
								bind:value={demand}
								min="0"
								max="100"
								class="h-3 w-full cursor-pointer appearance-none rounded-lg bg-blue-100 accent-blue-500 dark:bg-blue-900/30"
							/>
						</div>
						<div>
							<div class="mb-3 flex items-end justify-between uppercase">
								<span class="text-sm font-black tracking-widest text-red-500">Sellers (Supply)</span
								>
								<span class="text-2xl font-black text-red-500">{supply}</span>
							</div>
							<input
								type="range"
								bind:value={supply}
								min="0"
								max="100"
								class="h-3 w-full cursor-pointer appearance-none rounded-lg bg-red-100 accent-red-500 dark:bg-red-900/30"
							/>
						</div>
					</div>

					<div
						class="relative flex h-48 w-48 flex-col items-center justify-center overflow-hidden rounded-3xl border-2 border-(--border-color) bg-(--bg-primary) shadow-xl transition-all duration-500"
						style="border-color: {sentiment === 'Bullish'
							? '#27AE60'
							: sentiment === 'Bearish'
								? '#EB5757'
								: '#ddd'}"
					>
						<div class="mb-1 text-xs font-bold tracking-[0.2em] opacity-30">EQUILIBRIUM</div>
						<div class="text-4xl font-black">₹{price}</div>
						<div
							class="mt-3 rounded-full px-4 py-1.5 text-[10px] font-black tracking-widest text-white uppercase transition-colors duration-500"
							style="background: {sentiment === 'Bullish'
								? '#27AE60'
								: sentiment === 'Bearish'
									? '#EB5757'
									: '#999'}"
						>
							{sentiment}
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			What invisible hands move the price?
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Demand and supply don't move randomly. They react to external stimuli. In the stock market,
			knowledge is literally money.
		</p>
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
			<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-sm">
				<div
					class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-100 text-blue-600"
				>
					<Newspaper size={28} />
				</div>
				<h3 class="mb-3 text-xl font-bold">Company News</h3>
				<p class="text-base leading-relaxed opacity-70">
					Positive news like high profits or a new invention increases buyer interest (Demand),
					while negative news like a lawsuit or poor sales increases seller interest (Supply).
				</p>
			</div>
			<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-8 shadow-sm">
				<div
					class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-orange-100 text-orange-600"
				>
					<TrendingUp size={28} />
				</div>
				<h3 class="mb-3 text-xl font-bold">Market Sentiment</h3>
				<p class="text-base leading-relaxed opacity-70">
					It's not just about facts; it's about <b>emotions</b>. Fear and Greed are the two primary
					drivers that create momentum, often moving prices further than logic would suggest.
				</p>
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
