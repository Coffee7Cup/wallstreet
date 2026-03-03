<script>
	import {
		Settings2,
		Clock,
		MousePointer2,
		Zap,
		ArrowRight,
		ShieldCheck,
		Target
	} from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	const orders = [
		{
			type: 'Market Order',
			desc: 'Buy or sell immediately at the current market price.',
			icon: Zap,
			color: '#F2C94C'
		},
		{
			type: 'Limit Order',
			desc: 'Buy or sell at a specific price or better.',
			icon: Target,
			color: '#2D9CDB'
		}
	];
</script>

<svelte:head>
	<title>Trading Mechanics — Learning Lab</title>
</svelte:head>

<div class="animate-[fadeIn_0.5s_ease-out]">
	<header class="mb-10">
		<span
			class="mb-4 inline-block rounded-full px-3 py-1 text-[0.65rem] font-extrabold tracking-widest uppercase"
			style="background: rgba(45,156,219,0.12); color: #2D9CDB">Module 8</span
		>
		<h1
			class="mb-5 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl lg:text-[3.5rem]"
		>
			Trading Mechanics
		</h1>
		<p class="max-w-[750px] text-xl leading-relaxed opacity-70">
			How do orders actually execute? Understanding the pipes and clocks of the market helps you
			trade smarter and avoid costly execution errors.
		</p>
	</header>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">Mastering Order Types</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			Choosing the wrong order type can lead to buying a stock for more than you intended. Use <b
				>Market Orders</b
			>
			for speed and <b>Limit Orders</b> for precision.
		</p>
		<div class="grid gap-6 sm:grid-cols-2">
			{#each orders as order}
				<div
					class="group relative overflow-hidden rounded-3xl border border-(--border-color) bg-(--bg-primary) p-8 transition-all hover:border-(--item-color)"
					style="--item-color: {order.color}"
				>
					<div
						class="absolute -right-8 -bottom-8 opacity-5 transition-opacity duration-700 group-hover:scale-110 group-hover:opacity-10"
					>
						<order.icon size={180} />
					</div>
					<h3 class="mb-4 flex items-center gap-3 text-2xl font-bold" style="color: {order.color}">
						<order.icon size={24} />
						{order.type}
					</h3>
					<p class="relative z-10 text-base leading-relaxed opacity-70">{order.desc}</p>
				</div>
			{/each}
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-5 flex items-center gap-3 text-2xl font-extrabold">The Trading Clock (IST)</h2>
		<p class="mb-8 max-w-[800px] text-[1.1rem] leading-relaxed opacity-80">
			The Indian equity market follows a disciplined schedule. Trading only happens during specific
			windows on weekdays.
		</p>
		<div class="my-6 rounded-2xl border border-(--border-color) bg-blue-500/2! p-6">
			<div class="grid grid-cols-1 gap-8 p-4 md:grid-cols-3">
				<div class="flex flex-col items-center gap-4 text-center">
					<div
						class="flex h-16 w-16 items-center justify-center rounded-3xl bg-orange-100 text-orange-600 shadow-sm dark:bg-orange-900/20"
					>
						<Clock size={32} />
					</div>
					<div>
						<h4 class="text-lg font-bold">Pre-Market</h4>
						<p class="mb-2 text-xs font-black opacity-30">9:00 AM - 9:15 AM</p>
						<p class="text-sm opacity-60">
							Equilibrium price discovery through order matching. No execution.
						</p>
					</div>
				</div>
				<div class="relative flex flex-col items-center gap-4 text-center">
					<div class="absolute top-8 -left-4 hidden h-px w-8 bg-(--border-color) md:block"></div>
					<div class="absolute top-8 -right-4 hidden h-px w-8 bg-(--border-color) md:block"></div>
					<div
						class="flex h-16 w-16 items-center justify-center rounded-3xl bg-green-100 text-green-600 shadow-sm dark:bg-green-900/20"
					>
						<Zap size={32} />
					</div>
					<div>
						<h4 class="text-lg font-bold">Market Hours</h4>
						<p class="mb-2 text-xs font-black opacity-30">9:15 AM - 3:30 PM</p>
						<p class="text-sm opacity-60">
							The main continuous session. All retail and institutional trading happens here.
						</p>
					</div>
				</div>
				<div class="flex flex-col items-center gap-4 text-center">
					<div
						class="flex h-16 w-16 items-center justify-center rounded-3xl bg-blue-100 text-blue-600 shadow-sm dark:bg-blue-900/20"
					>
						<Clock size={32} />
					</div>
					<div>
						<h4 class="text-lg font-bold">Post-Market</h4>
						<p class="mb-2 text-xs font-black opacity-30">3:40 PM - 4:00 PM</p>
						<p class="text-sm opacity-60">
							Closing trades at the weighted average price of the last 30 minutes.
						</p>
					</div>
				</div>
			</div>
		</div>
	</section>

	<section class="mb-12 border-t border-(--border-color) pt-12">
		<div
			class="flex flex-col items-center justify-between gap-8 rounded-3xl border-2 border-green-500/20 bg-green-500/5 p-10 md:flex-row"
		>
			<div class="max-w-xl">
				<h3 class="mb-4 text-3xl font-black text-green-700 dark:text-green-500">
					Well done!
				</h3>
				<p class="text-lg leading-relaxed opacity-70">
					You now have the fundamental knowledge needed to navigate the stock market with
					confidence. From inflation to order types, you've covered the essentials.
				</p>
				<button
					class="mt-8 rounded-xl bg-green-600 px-8 py-3 font-bold text-white shadow-lg transition-all hover:bg-green-700 hover:shadow-green-500/20"
					onclick={() => goto(resolve('/dashboard'))}
				>
					Explore live Dashboard
				</button>
			</div>
			<div
				class="flex h-24 w-24 animate-pulse items-center justify-center rounded-full bg-green-500 text-white shadow-xl"
			>
				<ShieldCheck size={48} />
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
