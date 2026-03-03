<script>
	import { BookOpen, Repeat, Building, ArrowRightLeft } from 'lucide-svelte';

	const visualizationSteps = [
		{ icon: Building, label: 'Company', color: '#2D9CDB', text: 'Needs money to grow' },
		{ icon: BookOpen, label: 'Stock Market', color: '#27AE60', text: 'The platform for exchange' },
		{ icon: Repeat, label: 'Investors', color: '#F2994A', text: 'People who provide capital' }
	];
</script>

<svelte:head>
	<title>What is Stock Market? — Learning Lab</title>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(45,156,219,0.12); color: #2D9CDB">Module 2</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			What is a Stock Market?
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			A stock market (or equity market) is a secure, regulated environment where market participants
			can trade shares and other financial instruments with confidence.
		</p>
	</header>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">The Market Ecosystem</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Think of a stock market as a high-tech auction house. It brings together people who have
			capital (investors) and people who need capital (companies). Instead of physical products,
			they trade "certificates of ownership" or shares.
		</p>
		<div class="my-6 rounded-2xl border border-(--border-color) bg-(--bg-hover) p-6">
			<div class="relative flex items-center justify-between p-8 px-4">
				{#each visualizationSteps as step, i}
					<div
						class="group z-1 flex w-[120px] flex-col items-center gap-3 text-center"
						style="--step-color: {step.color}"
					>
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full border-2 border-(--border-color) bg-(--bg-primary) transition-all duration-300 group-hover:scale-110 group-hover:border-(--step-color) group-hover:shadow-[0_0_15px_color-mix(in_srgb,var(--step-color)_30%,transparent)]"
						>
							<step.icon size={24} color={step.color} />
						</div>
						<div class="text-[0.75rem] font-bold">{step.label}</div>
						<p class="text-[0.7rem] font-medium opacity-70">{step.text}</p>
					</div>
					{#if i < visualizationSteps.length - 1}
						<div
							class="cycle-arrow relative -top-[15px] -mx-2.5 h-0.5 flex-1 bg-(--border-color)"
						></div>
					{/if}
				{/each}
			</div>
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-3xl font-extrabold">
			Primary vs Secondary Market
		</h2>
		<p class="mb-8 max-w-[800px] text-lg leading-relaxed opacity-80">
			Markets are divided into two main stages. One for "new" shares and one for "used" (already
			existing) shares.
		</p>
		<div class="mt-6 grid gap-6 md:grid-cols-2">
			<div class="rounded-2xl border border-(--border-color) bg-blue-500/5 p-8">
				<h3 class="mb-4 flex items-center gap-2 text-2xl font-bold text-blue-500">
					<Building size={22} /> Primary Market
				</h3>
				<p class="text-xl leading-relaxed opacity-80">
					This is where securities are created. It's the first time a company issues its shares to
					the public. The transaction happens directly between the <b>Company</b> and the
					<b>Investor</b>.
				</p>
				<ul class="mt-6 space-y-3 text-lg opacity-70">
					<li>• <b>Purpose:</b> Raising fresh capital for expansion or debt reduction.</li>
					<li>• <b>Transaction:</b> Directly from the company to the public.</li>
					<li>• <b>Example:</b> IPOs (Initial Public Offerings) and FPOs.</li>
				</ul>
			</div>
			<div class="rounded-2xl border border-(--border-color) bg-orange-500/5 p-8">
				<h3 class="mb-4 flex items-center gap-2 text-2xl font-bold text-orange-500">
					<ArrowRightLeft size={22} /> Secondary Market
				</h3>
				<p class="text-xl leading-relaxed opacity-80">
					This is what most people mean by "trading". Here, investors buy and sell shares among
					themselves on exchanges like the NSE or BSE. The <b>Company</b> is not part of this trade.
				</p>
				<ul class="mt-6 space-y-3 text-lg opacity-70">
					<li>• <b>Purpose:</b> Providing "liquidity" (the ability to exit/enter positions).</li>
					<li>• <b>Transaction:</b> Between two different investors.</li>
					<li>• <b>Example:</b> Buying Reliance today on your mobile app.</li>
				</ul>
			</div>
		</div>
	</section>

	<section class="mb-12 border-t border-(--border-color) pt-8">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			The Role of Stock Exchanges
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Exchanges act as neutral facilitators. They don't set the price—they just provide the rules,
			the technology, and the security to ensure that every trade is fair and settled correctly.
		</p>
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
			<div class="rounded-xl border border-(--border-color) bg-(--bg-primary) p-4">
				<h4 class="mb-1 text-lg font-bold">Standardization</h4>
				<p class="text-md opacity-60">
					Ensures all shares of a company are identical and tradeable.
				</p>
			</div>
			<div class="rounded-xl border border-(--border-color) bg-(--bg-primary) p-4">
				<h4 class="mb-1 text-lg font-bold">Price Transparency</h4>
				<p class="text-md opacity-60">Ensures everyone sees the same price at the same time.</p>
			</div>
			<div class="rounded-xl border border-(--border-color) bg-(--bg-primary) p-4">
				<h4 class="mb-1 text-lg font-bold">Security</h4>
				<p class="text-md opacity-60">
					Guarantees that you receive your shares if you pay for them.
				</p>
			</div>
		</div>
	</section>
</div>

<style>
	.cycle-arrow::after {
		content: '';
		position: absolute;
		right: 0;
		top: -3px;
		width: 8px;
		height: 8px;
		border-top: 2px solid var(--border-color);
		border-right: 2px solid var(--border-color);
		transform: rotate(45deg);
	}

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
