<script>
	import { BarChart3, Users, Building, Percent } from 'lucide-svelte';

	let companyValue = $state(1000); // 1000 Crores
	let publicStake = $state(25); // 25%

	let marketCap = $derived(companyValue);
	let publicValue = $derived((companyValue * publicStake) / 100);
</script>

<svelte:head>
	<title>IPO & Shares — Learning Lab</title>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(155,89,182,0.12); color: #9B59B6">Module 4</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			IPO & Understanding Shares
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			A share is the smallest unit of ownership in a company. When you buy a share, you are
			literally buying a small piece of that business.
		</p>
	</header>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			Interactive Market Cap Visualizer
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Market Capitalization (Market Cap) is the total market value of a company's outstanding
			shares. Use the sliders to see how valuation and the percentage sold to the public (IPO)
			create the "Issue Size."
		</p>

		<div class="my-6 rounded-2xl border border-(--border-color) bg-purple-500/3! p-6">
			<div class="space-y-8 p-4">
				<div>
					<div class="mb-3 flex items-end justify-between">
						<span class="text-sm font-bold tracking-wide uppercase opacity-60"
							>Total Company Valuation</span
						>
						<span class="text-2xl font-black text-[#9B59B6]">₹{companyValue} Cr</span>
					</div>
					<input
						type="range"
						bind:value={companyValue}
						min="100"
						max="5000"
						step="100"
						class="h-2 w-full cursor-pointer appearance-none rounded-lg bg-(--border-color) accent-[#9B59B6]"
					/>
					<div class="mt-2 flex justify-between text-[10px] font-bold opacity-40">
						<span>₹100 CR</span>
						<span>₹5,000 CR</span>
					</div>
				</div>

				<div>
					<div class="mb-3 flex items-end justify-between">
						<span class="text-sm font-bold tracking-wide uppercase opacity-60"
							>Public Stake to be Sold (IPO)</span
						>
						<span class="text-2xl font-black text-[#27AE60]">{publicStake}%</span>
					</div>
					<input
						type="range"
						bind:value={publicStake}
						min="10"
						max="49"
						class="h-2 w-full cursor-pointer appearance-none rounded-lg bg-(--border-color) accent-[#27AE60]"
					/>
					<div class="mt-2 flex justify-between text-[10px] font-bold opacity-40">
						<span>10% (MIN)</span>
						<span>49% (MAX)</span>
					</div>
				</div>

				<div class="grid grid-cols-1 gap-6 border-t border-(--border-color) pt-8 sm:grid-cols-2">
					<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-6 shadow-sm">
						<div class="mb-2 text-xs font-bold tracking-widest uppercase opacity-40">
							Market Cap
						</div>
						<div class="text-3xl font-black text-[#9B59B6]">
							₹{marketCap} <span class="text-sm opacity-50">CR</span>
						</div>
						<p class="mt-2 text-xs opacity-60">Total value of all shares combined.</p>
					</div>
					<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-6 shadow-sm">
						<div class="mb-2 text-xs font-bold tracking-widest uppercase opacity-40">
							IPO Issue Size
						</div>
						<div class="text-3xl font-black text-[#27AE60]">
							₹{publicValue} <span class="text-sm opacity-50">CR</span>
						</div>
						<p class="mt-2 text-xs opacity-60">Total money being raised from the public.</p>
					</div>
				</div>
			</div>
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">
			Why do companies go public?
		</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Going public isn't just about getting famous; it's a strategic move for a business. Here's why
			they do it:
		</p>
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-6">
				<h4 class="mb-2 text-lg flex items-center gap-2 font-bold text-blue-500">
					<Building size={18} /> Business Expansion
				</h4>
				<p class="text-md opacity-70">
					Building new factories, entering new markets, or acquiring other companies requires
					massive capital.
				</p>
			</div>
			<div class="rounded-2xl border border-(--border-color) bg-(--bg-primary) p-6">
				<h4 class="mb-2 text-lg flex items-center gap-2 font-bold text-orange-500">
					<Users size={18} /> Exit for Early Investors
				</h4>
				<p class="text-md opacity-70">
					Allows early venture capital investors and founders to sell their stake to the general
					public.
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
