<script>
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import {
		TrendingUp,
		BookOpen,
		Building2,
		BarChart3,
		DollarSign,
		Calculator,
		List,
		Settings2,
		ArrowRight
	} from 'lucide-svelte';

	const topics = [
		{
			href: '/learning-lab/why-stock-market',
			label: 'Why Stock Market?',
			icon: TrendingUp,
			color: '#27AE60',
			desc: 'Inflation, wealth creation, and why people invest.'
		},
		{
			href: '/learning-lab/what-is-stock-market',
			label: 'What is a Stock Market?',
			icon: BookOpen,
			color: '#2D9CDB',
			desc: 'Purpose, primary vs secondary markets.'
		},
		{
			href: '/learning-lab/indian-stock-market',
			label: 'Indian Stock Market',
			icon: Building2,
			color: '#F2994A',
			desc: 'NSE, BSE, Sensex, Nifty 50 and indices.'
		},
		{
			href: '/learning-lab/ipo-and-shares',
			label: 'IPO & Shares',
			icon: BarChart3,
			color: '#9B59B6',
			desc: 'How companies go public, face value, market cap.'
		},
		{
			href: '/learning-lab/price-discovery',
			label: 'Price Discovery',
			icon: DollarSign,
			color: '#EB5757',
			desc: 'Demand, supply, news and market sentiment.'
		},
		{
			href: '/learning-lab/stock-market-terms',
			label: 'Stock Market Terms',
			icon: List,
			color: '#F2C94C',
			desc: 'Price, equity, volume, margin, support, resistance.'
		},
		{
			href: '/learning-lab/investing-basics',
			label: 'Investing Basics',
			icon: Calculator,
			color: '#27AE60',
			desc: 'P/E, P/B, ROE, ROCE, Dividend Yield, Growth vs Value.'
		},
		{
			href: '/learning-lab/trading-mechanics',
			label: 'Trading Mechanics',
			icon: Settings2,
			color: '#2D9CDB',
			desc: 'Market vs limit orders, trading hours, open/close.'
		}
	];
</script>

<svelte:head>
	<title>Learning — Stock Market — Wallstreet </title>
	<meta
		name="description"
		content="Master the stock market with interactive lessons on investing, trading, IPOs, and financial ratios."
	/>
</svelte:head>

<div class="max-w-[800px]">
	<header class="mb-10">
		<div
			class="mb-4 inline-block rounded-full bg-[#27ae60]/12 px-3 py-1 text-[0.65rem] font-extrabold tracking-widest text-[#27ae60] uppercase"
		>
			Resource
		</div>
		<h1 class="mb-4 text-4xl leading-[1.1] font-black tracking-tighter md:text-5xl">
			"WallStreet" Stock Market <br /><span class="text-[#27ae60]">Learn</span>
		</h1>
		<p class="max-w-[520px] text-lg leading-relaxed opacity-60">
			From understanding why markets exist to executing your first trade — learn everything step by
			step with interactive charts and real explanations.
		</p>
	</header>

	<section class="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-3.5">
		{#each topics as topic, i}
			<a
				href={topic.href}
				class="group flex items-center gap-4 rounded-[14px] border border-(--border-color) bg-(--bg-primary) p-4 text-inherit no-underline transition-all duration-200"
				style="--c: {topic.color}; animation: fadeUp 0.4s ease both; animation-delay: {i *
					0.05}s; --hover-shadow: color-mix(in srgb, {topic.color} 18%, transparent); --hover-bg: color-mix(in srgb, {topic.color} 4%, var(--bg-primary));"
				onmouseenter={(e) => {
					e.currentTarget.style.borderColor = topic.color;
					e.currentTarget.style.boxShadow = `0 4px 20px ${e.currentTarget.style.getPropertyValue('--hover-shadow')}`;
					e.currentTarget.style.transform = 'translateY(-2px)';
					e.currentTarget.style.background = e.currentTarget.style.getPropertyValue('--hover-bg');
				}}
				onmouseleave={(e) => {
					e.currentTarget.style.borderColor = '';
					e.currentTarget.style.boxShadow = '';
					e.currentTarget.style.transform = '';
					e.currentTarget.style.background = '';
				}}
			>
				<div
					class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl"
					style="background: color-mix(in srgb, {topic.color} 14%, transparent); color: {topic.color}"
				>
					<topic.icon size={24} />
				</div>
				<div class="min-w-0 flex-1">
					<h2 class="mb-1 text-[0.9rem] font-bold">{topic.label}</h2>
					<p class="text-[0.75rem] leading-tight opacity-50">{topic.desc}</p>
				</div>
				<ArrowRight
					size={16}
					class="shrink-0 opacity-30 transition-all duration-200 group-hover:translate-x-1 group-hover:opacity-80"
				/>
			</a>
		{/each}
	</section>
</div>

<style>
	@keyframes fadeUp {
		from {
			opacity: 0;
			transform: translateY(12px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
