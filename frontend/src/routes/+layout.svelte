<script>
	import '../app.css';
	import { theme } from '$lib/stores/theme';
	import { user, admin, logout, token } from '$lib/stores/auth';
	import {
		Moon,
		Sun,
		LogOut,
		User,
		LayoutDashboard,
		Settings,
		TrendingUp,
		X,
		Paperclip,
		Newspaper,
		BookOpen
	} from 'lucide-svelte';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { cubicIn, cubicOut } from 'svelte/easing';
	import { fly } from 'svelte/transition';
	import gsap from 'gsap';
	import { untrack } from 'svelte';

	let areYouSure = $state({
		show: false,
		message: null,
		color: null,
		callback: null
	});

	function toggleTheme() {
		theme.update((t) => (t === 'light' ? 'dark' : 'light'));
	}

	import { connectMarketWS, marketState } from '$lib/stores/market';
	import { fade, slide } from 'svelte/transition';

	onMount(() => {
		// Simple entrance animation for the whole layout
		gsap.from('nav', { y: -20, opacity: 0, duration: 0.8, ease: 'power2.out' });

		// Initial check for auth
		checkAuth();

		// Initialize global market data stream
		connectMarketWS();
	});

	// Reactive check handles token removal/expiry during session
	$effect(() => {
		if (!$token && page.url.pathname !== '/' && page.url.pathname !== '/login') {
			goto('/login');
		}
		// Prevent logged-in users from accessing login page
		if ($token && page.url.pathname === '/login') {
			if ($admin) {
				goto('/admin');
			} else if ($user) {
				goto('/dashboard');
			}
		}
	});

	function checkAuth() {
		const path = page.url.pathname;
		if (!$token && path !== '/' && path !== '/login') {
			goto('/login');
		}
	}

	let notifications = $state([]);
	let shownNewsIds = $state(new Set());

	function addNotification(type, title, content, color = 'red') {
		const id = Math.random().toString(36).substring(2, 9);
		const newNotification = { id, type, title, content, color };
		notifications = [newNotification, ...notifications];

		setTimeout(() => {
			notifications = notifications.filter((n) => n.id !== id);
		}, 5000);
	}

	function removeNotification(id) {
		notifications = notifications.filter((n) => n.id !== id);
	}

	$effect(() => {
		const latestNews = $marketState.news;
		if (latestNews.length > 0) {
			const newest = latestNews[latestNews.length - 1];
			const newsId = newest.id || newest.title;

			if (!shownNewsIds.has(newsId)) {
				shownNewsIds.add(newsId);
				untrack(() => {
					addNotification('news', newest.title, newest.content, 'red');
				});
			}
		}
	});

	$effect(() => {
		if ($marketState.lastError) {
			const errorMsg = $marketState.lastError;
			untrack(() => {
				addNotification('error', 'Trade Error', errorMsg, 'red');
			});
		}
	});

	let { children } = $props();
</script>

<div class="flex min-h-screen flex-col transition-colors duration-300">
	<nav class="sticky top-0 z-50 border-b border-(--border-color) bg-(--bg-primary) px-4 py-3">
		<div class="mx-auto flex max-w-screen-2xl items-center justify-between">
			<div class="flex items-center gap-6">
				<a
					href="https://coffee7cup.github.io/cie-1/"
					target="_blank"
					class="text-xl font-bold tracking-tight">WallStreet - Sjcknl</a
				>
				<div class="hidden items-center gap-4 md:flex">
					{#if $user || $admin}
						<a
							href="/dashboard"
							class="flex items-center gap-1.5 rounded px-2 py-1 hover:bg-(--bg-hover) {page.url
								.pathname === '/dashboard'
								? 'font-bold'
								: ''}"
						>
							<LayoutDashboard size={18} /> Dashboard
						</a>
					{/if}
					{#if $user}
						<a
							href="/trades"
							class="flex items-center gap-1.5 rounded px-2 py-1 hover:bg-(--bg-hover) {page.url
								.pathname === '/trades'
								? 'font-bold'
								: ''}"
						>
							<TrendingUp size={18} /> Trades
						</a>
						<a
							href="/profile"
							class="flex items-center gap-1.5 rounded px-2 py-1 hover:bg-(--bg-hover) {page.url
								.pathname === '/profile'
								? 'font-bold'
								: ''}"
						>
							<User size={18} /> Profile
						</a>
					{/if}
					{#if $admin}
						<a
							href="/admin"
							class="flex items-center gap-1.5 rounded px-2 py-1 hover:bg-(--bg-hover) {page.url
								.pathname === '/admin'
								? 'font-bold'
								: ''}"
						>
							<Settings size={18} /> Admin
						</a>
						<a
							href="/admin/top-traders"
							class="flex items-center gap-1.5 rounded px-2 py-1 hover:bg-(--bg-hover) {page.url
								.pathname === '/admin/top-traders'
								? 'font-bold'
								: ''}"
						>
							<TrendingUp size={18} /> Top Traders
						</a>
					{/if}

					{#if $admin || $user}
						<a
							href="/news"
							class="flex items-center gap-1.5 rounded px-2 py-1 hover:bg-(--bg-hover) {page.url
								.pathname === '/news'
								? 'font-bold'
								: ''}"
						>
							<Newspaper size={18} /> News
						</a>
						<a
							href="/learning-lab"
							class="flex items-center gap-1.5 rounded px-2 py-1 hover:bg-(--bg-hover) {page.url
								.pathname === '/learning-lab'
								? 'font-bold'
								: ''}"
						>
							<BookOpen size={18} /> Learn
						</a>
					{/if}
				</div>
			</div>

			<div class="flex items-center gap-4">
				<button onclick={toggleTheme} class="rounded p-2 transition-colors hover:bg-(--bg-hover)">
					{#if $theme === 'light'}
						<Moon size={20} />
					{:else}
						<Sun size={20} />
					{/if}
				</button>

				{#if $user || $admin}
					<button
						onclick={() => {
							areYouSure.show = true;
							areYouSure.color = '#EB5757';
							areYouSure.message = 'Logout';
							areYouSure.callback = () => logout();
						}}
						class="flex items-center gap-1.5 rounded bg-(--bg-hover) px-3 py-1.5 transition-all hover:text-[#EB5757]"
					>
						<LogOut size={18} /> <span class="hidden sm:inline">Logout</span>
					</button>
				{:else}
					<a href="/login" class="btn-green rounded px-4 py-1.5 font-medium">Login</a>
				{/if}
			</div>
		</div>
	</nav>

	{#if !$marketState.connected && $user && !$marketState.isActive}
		<div transition:slide class="bg-[#EB5757] px-4 py-2 text-center text-xs font-bold text-white">
			<span class="flex items-center justify-center gap-2">
				<div class="h-2 w-2 animate-pulse rounded-full bg-white"></div>
				{$marketState.isActive
					? 'Disconnected from Market Server. Reconnecting...'
					: 'Wallstreet has ended'}
			</span>
		</div>
	{/if}

	<div class="fixed top-20 right-6 z-100 flex flex-col gap-3">
		{#each notifications as n (n.id)}
			<div
				in:fly={{ x: 300, duration: 400, easing: cubicOut }}
				out:fly={{ x: 300, duration: 300, easing: cubicIn }}
				class="relative flex max-w-sm cursor-pointer gap-3 rounded-xl border border-red-500/50
					   bg-(--bg-primary) p-4 shadow-2xl backdrop-blur-md transition-all hover:bg-(--bg-hover)"
			>
				<div class="mt-1 shrink-0 text-red-500">
					{#if n.type === 'news'}
						<Newspaper size={20} />
					{:else}
						<div class="mt-1 h-3 w-3 animate-pulse rounded-full bg-red-500"></div>
					{/if}
				</div>
				<div class="flex flex-col gap-1 pr-6">
					<h2 class="text-lg leading-none font-bold text-red-500">{n.title}</h2>
					<p class="text-xs leading-relaxed opacity-60">{n.content}</p>
				</div>
				<button
					onclick={() => removeNotification(n.id)}
					class="absolute top-3 right-3 rounded-md p-1 text-(--text-muted) hover:bg-(--bg-hover) hover:text-(--text-primary)"
				>
					<X size={14} />
				</button>
			</div>
		{/each}
	</div>

	{#if areYouSure.show}
		<div
			class="fixed inset-0 z-50 flex items-center justify-center
		       bg-black/60 backdrop-blur-sm"
		>
			<div
				class="flex w-[90%] max-w-md
			       flex-col
			       gap-4 rounded-xl
			       border border-(--border-color) bg-(--bg-primary) p-6 shadow-2xl"
			>
				<h2 class="text-center text-2xl font-semibold">Are you sure?</h2>

				<p class="text-center text-sm text-gray-400">
					{areYouSure.message}
				</p>

				<div class="mt-4 flex gap-3">
					<button
						class="flex-1 rounded-lg
					       py-2 font-medium text-white
					       transition hover:brightness-110"
						style="background-color: {areYouSure.color};"
						onclick={() => {
							areYouSure.show = false;
							areYouSure.callback();
						}}
					>
						Proceed
					</button>

					<button
						class="flex-1 rounded-lg bg-[#2F2F2F]
					       py-2 font-medium text-gray-300
					       transition hover:bg-[#3A3A3A]"
						onclick={() => {
							areYouSure.show = false;
						}}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	{/if}

	<main class="mx-auto w-full max-w-screen-2xl grow p-4 md:p-8">
		{@render children?.()}
	</main>

	<!-- Connection Status Indicator -->
	<div class="fixed right-6 bottom-6 z-50 flex flex-col items-end gap-2">
		<div
			class="flex items-center gap-2 rounded-full border border-(--border-color) bg-(--bg-primary)/80 px-3 py-1 text-[10px] font-black tracking-widest uppercase shadow-xl backdrop-blur-md"
			class:text-green-500={$marketState.connectionStatus === 'connected'}
			class:text-yellow-500={$marketState.connectionStatus === 'connecting' ||
				$marketState.connectionStatus === 'reconnecting'}
			class:text-red-500={$marketState.connectionStatus === 'disconnected'}
		>
			<span class="relative flex h-2 w-2">
				{#if $marketState.connectionStatus === 'connected'}
					<span
						class="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75"
					></span>
					<span class="relative inline-flex h-2 w-2 rounded-full bg-green-500"></span>
				{:else if $marketState.connectionStatus === 'disconnected'}
					<span class="relative inline-flex h-2 w-2 rounded-full bg-red-500"></span>
				{:else}
					<span
						class="absolute inline-flex h-full w-full animate-ping rounded-full bg-yellow-400 opacity-75"
					></span>
					<span class="relative inline-flex h-2 w-2 rounded-full bg-yellow-500"></span>
				{/if}
			</span>
			{$marketState.connectionStatus}
		</div>
	</div>

	{#if $marketState.simulationEnded}
		<div
			class="fixed inset-0 z-100 flex items-center justify-center
		       bg-black/40 backdrop-blur-md"
			transition:fade
		>
			<div
				class="flex w-[90%] max-w-lg
			       flex-col items-center
			       gap-6 rounded-2xl
			       border border-(--border-color) bg-(--bg-primary) p-10 text-center shadow-2xl"
				in:fly={{ y: 20, duration: 600, delay: 200 }}
			>
				<div
					class="flex h-16 w-16 items-center justify-center rounded-full bg-[#27AE60]/10 text-[#27AE60]"
				>
					<TrendingUp size={32} />
				</div>

				<div class="space-y-2">
					<h2 class="text-3xl font-bold tracking-tight">Simulation Ended</h2>
					<p class="text-lg opacity-60">
						The market has closed. All historical data has been processed.
					</p>
				</div>

				<div class="h-px w-full bg-(--border-color) opacity-50"></div>

				<p class="text-sm opacity-50">
					Thank you for participating in WallStreet. You can still view your trade history and
					portfolio.
				</p>

				<div class="mt-4 flex gap-4">
					<button
						onclick={() => marketState.update((s) => ({ ...s, simulationEnded: false }))}
						class="rounded-lg border border-(--border-color) bg-(--bg-hover) px-6 py-3 font-medium
						       transition-all hover:bg-(--bg-primary)"
					>
						Dismiss
					</button>
					<a
						href="/profile"
						onclick={() => marketState.update((s) => ({ ...s, simulationEnded: false }))}
						class="rounded-lg bg-(--text-primary) px-8 py-3 font-bold text-(--bg-primary)
						       transition-all hover:scale-[1.02] active:scale-[0.98]"
					>
						Go to Profile
					</a>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	:global(html) {
		scrollbar-gutter: stable;
	}
</style>
