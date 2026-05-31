<script>
	import { onMount } from 'svelte';
	import api from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { marketState } from '$lib/stores/market';
	import { Wallet, Briefcase, User as UserIcon, Calendar } from 'lucide-svelte';
	import gsap from 'gsap';

	let portfolio = $derived($marketState.userPortfolio || []);
	let loading = $state(true);

	// Calculate portfolio value reactively using real-time prices and holdings from marketState
	let portfolioValue = $derived.by(() => {
		return portfolio.reduce((acc, curr) => {
			const livePrice =
				$marketState.stocks.find((s) => s.company_id === curr.company_id)?.close_price ||
				curr.current_price;
			return acc + curr.quantity * livePrice;
		}, 0);
	});

	onMount(async () => {
		if (!$user) return;

		try {
			const [profileRes, portRes] = await Promise.all([
				api.get('/users/profile'),
				api.get('/trade/portfolio')
			]);

			user.update((u) => ({ ...u, ...profileRes.data.user }));

			// Initialize marketState with current data
			marketState.update((s) => ({
				...s,
				userPortfolio: portRes.data.portfolio || [],
				userBalance: profileRes.data.user.cash_balance
			}));

			loading = false;

			gsap.from('.profile-section', {
				x: -30,
				opacity: 0,
				duration: 0.8,
				ease: 'power2.out'
			});
			gsap.from('.stat-card', {
				y: 20,
				opacity: 0,
				stagger: 0.1,
				duration: 0.6,
				ease: 'power2.out'
			});
		} catch (err) {
			console.error('Profile load failed', err);
		}
	});
</script>

<div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
	<!-- User Info Sidebar -->
	<div class="space-y-6">
		<div class="profile-section notion-card bg-(--bg-primary)">
			<div class="mb-6 flex items-center gap-4">
				<div
					class="bg-opacity-20 flex h-16 w-16 items-center justify-center rounded-full bg-[#27AE60] text-[#27AE60]"
				>
					<UserIcon size={32} />
				</div>
				<div>
					<h2 class="text-xl font-bold">{$user?.username}</h2>
					<p class="text-sm opacity-60">{$user?.email}</p>
					<p class="text-sm opacity-60">{$user?.id}</p>
				</div>
			</div>

			<div class="space-y-4">
				<div class="flex items-center justify-between rounded bg-(--bg-hover) p-3">
					<div class="flex items-center gap-2 opacity-70">
						<Wallet size={18} />
						<span class="text-sm">Cash Balance</span>
					</div>
					<span class="font-bold text-[#27AE60]">
						₹{$marketState.userBalance?.toLocaleString() || '0'}
					</span>
				</div>

				<div class="flex items-center justify-between rounded bg-(--bg-hover) p-3">
					<div class="flex items-center gap-2 opacity-70">
						<Briefcase size={18} />
						<span class="text-sm">Portfolio Value</span>
					</div>
					<span class="font-bold">₹{portfolioValue.toLocaleString()}</span>
				</div>
			</div>

			<div class="mt-8 border-t border-(--border-color) pt-6">
				<div class="mb-4 flex items-center gap-2 text-xs tracking-widest uppercase opacity-60">
					<Calendar size={14} /> Account Details
				</div>
				<div class="space-y-2 text-sm">
					<div class="flex justify-between">
						<span class="opacity-60">Status</span>
						<span class="font-medium text-[#27AE60]">Active</span>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="space-y-8 lg:col-span-2">
		<!-- Portfolio List -->
		<div>
			<h3 class="mb-4 flex items-center gap-2 text-lg font-bold">
				<Briefcase size={20} /> Your Holdings
			</h3>

			{#if loading}
				<div class="space-y-4">
					<div class="h-16 animate-pulse rounded bg-(--bg-hover)"></div>
					<div class="h-16 animate-pulse rounded bg-(--bg-hover)"></div>
				</div>
			{:else if portfolio.length === 0}
				<div class="notion-card py-12 text-center opacity-60">
					<p>No stocks in your portfolio yet.</p>
					<a href="/dashboard" class="mt-2 text-[#27AE60] hover:underline">Start Trading</a>
				</div>
			{:else}
				<div class="space-y-3">
					{#each portfolio as item (item.company_id)}
						<div
							class="portfolio-item notion-card flex items-center justify-between transition-all hover:translate-x-1 hover:border-[#27AE60]"
						>
							<div class="flex items-center gap-4">
								<div
									class="flex h-10 w-10 min-w-[40px] items-center justify-center rounded-xl bg-(--bg-hover) text-[10px] font-black text-[#27AE60]"
								>
									{item.company_symbol}
								</div>
								<div>
									<h4 class="font-bold">{item.company_name}</h4>
									<p class="text-xs opacity-60">Holding {item.quantity} shares</p>
								</div>
							</div>
							<div class="text-right">
								<p class="font-bold">₹{(item.quantity * item.current_price).toLocaleString()}</p>
								<p class="text-[10px] opacity-40">@ ₹{item.current_price.toLocaleString()}</p>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
