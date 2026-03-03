<script>
	import { page } from '$app/state';
	import {
		TrendingUp,
		BookOpen,
		Building2,
		BarChart3,
		DollarSign,
		LineChart,
		Calculator,
		List,
		Settings2,
		Menu,
		X,
		ChevronRight
	} from 'lucide-svelte';

	let { children } = $props();
	let sidebarOpen = $state(false);

	const navItems = [
		{
			href: '/learning-lab/why-stock-market',
			label: 'Why Stock Market?',
			icon: TrendingUp,
			color: '#27AE60'
		},
		{
			href: '/learning-lab/what-is-stock-market',
			label: 'What is a Stock Market?',
			icon: BookOpen,
			color: '#2D9CDB'
		},
		{
			href: '/learning-lab/indian-stock-market',
			label: 'Indian Stock Market',
			icon: Building2,
			color: '#F2994A'
		},
		{
			href: '/learning-lab/ipo-and-shares',
			label: 'IPO & Shares',
			icon: BarChart3,
			color: '#9B59B6'
		},
		{
			href: '/learning-lab/price-discovery',
			label: 'Price Discovery',
			icon: DollarSign,
			color: '#EB5757'
		},
		{
			href: '/learning-lab/stock-market-terms',
			label: 'Stock Market Terms',
			icon: List,
			color: '#F2C94C'
		},
		{
			href: '/learning-lab/investing-basics',
			label: 'Investing Basics',
			icon: Calculator,
			color: '#27AE60'
		},
		{
			href: '/learning-lab/trading-mechanics',
			label: 'Trading Mechanics',
			icon: Settings2,
			color: '#2D9CDB'
		}
	];

	function isActive(href) {
		return page.url.pathname === href;
	}

	function closeSidebar() {
		sidebarOpen = false;
	}
</script>

<div class="relative flex min-h-[calc(100vh-57px)]">
	<!-- Mobile overlay -->
	{#if sidebarOpen}
		<button
			class="fixed inset-0 z-39 cursor-default border-none bg-black/50 md:hidden"
			onclick={closeSidebar}
			aria-label="Close sidebar"
		></button>
	{/if}

	<!-- Sidebar -->
	<aside
		class="fixed inset-y-0 left-0 z-40 flex h-[calc(100vh-57px)] w-[260px] min-w-[260px] flex-col overflow-y-auto border-r border-(--border-color) bg-(--bg-primary) transition-transform duration-300 ease-in-out md:sticky md:top-[57px] {sidebarOpen
			? 'translate-x-0'
			: '-translate-x-full md:translate-x-0'} shadow-2xl md:shadow-none"
	>
		<div class="flex items-center justify-between border-b border-(--border-color) p-5">
			<div class="flex items-center gap-2 text-[0.95rem] font-bold text-[#27ae60]">
				<BookOpen size={20} />
				<span>Learn</span>
			</div>
			<button
				class="flex rounded-md border-none bg-transparent p-1 text-inherit transition-colors hover:bg-(--bg-hover) md:hidden"
				onclick={closeSidebar}
				aria-label="Close"
			>
				<X size={18} />
			</button>
		</div>

		<nav class="flex flex-1 flex-col gap-0.5 p-3">
			<p class="px-2.5 pt-2 pb-1 text-[0.65rem] font-bold tracking-widest uppercase opacity-40">
				Topics
			</p>
			{#each navItems as item}
				{@const active = isActive(item.href)}
				<a
					href={item.href}
					class="relative flex items-center gap-2.5 rounded-xl px-3 py-2.5 text-sm font-medium text-inherit no-underline transition-colors duration-150 hover:bg-(--bg-hover) {active
						? 'font-semibold'
						: ''}"
					style="--item-color: {item.color}; {active
						? `background: color-mix(in srgb, ${item.color} 12%, transparent); color: ${item.color};`
						: ''}"
					onclick={closeSidebar}
				>
					<span class="flex shrink-0 items-center" style="color: {active ? item.color : 'inherit'}">
						<item.icon size={18} />
					</span>
					<span class="flex-1 leading-tight">{item.label}</span>
					{#if active}
						<ChevronRight size={14} class="opacity-60" />
					{/if}
				</a>
			{/each}
		</nav>

		<div class="border-t border-(--border-color) p-4">
			<div class="text-xs font-semibold opacity-50">
				<p>WallStreet — CIE</p>
				<p class="text-[0.65rem] font-normal opacity-70">Stock Market Education</p>
			</div>
		</div>
	</aside>

	<!-- Main Content -->
	<div class="flex min-w-0 flex-1 flex-col">
		<!-- Mobile top bar -->
		<div
			class="sticky top-[57px] z-30 flex items-center gap-3 border-b border-(--border-color) bg-(--bg-primary) p-3 md:hidden"
		>
			<button
				class="cursor-pointer rounded-lg border-none bg-transparent p-1.5 text-inherit transition-colors hover:bg-(--bg-hover)"
				onclick={() => (sidebarOpen = true)}
				aria-label="Open sidebar"
			>
				<Menu size={20} />
			</button>
			<span class="text-[0.95rem] font-bold">
				{navItems.find((n) => isActive(n.href))?.label ?? 'Learning Lab'}
			</span>
		</div>

		<main class="max-w-[900px] p-4 md:p-8">
			{@render children()}
		</main>
	</div>
</div>
