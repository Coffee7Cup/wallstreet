<script>
	import { onMount } from 'svelte';
	import api from '$lib/api';
	import { Newspaper, AlertCircle, Loader, Search } from 'lucide-svelte';

	let news = $state([]);
	let sectors = $state([]);
	let selectedSector = $state('');
	let err = $state('');
	let loading = $state(true);
	let loadingMore = $state(false);
	let searchQuery = $state('');
	let fixedTick = $state(null);
	let offset = $state(0);
	let limit = 20;
	let hasMore = $state(true);

	async function fetchSectors() {
		try {
			const res = await api.get('/market/sectors');
			sectors = res.data.sectors || [];
		} catch (e) {
			console.error('Failed to fetch sectors:', e);
		}
	}

	async function fetchNews(isLoadMore = false) {
		if (isLoadMore) {
			loadingMore = true;
		} else {
			loading = true;
			offset = 0;
		}

		try {
			const endpoint = selectedSector ? '/market/news/sector' : '/market/news';
			const res = await api.get(endpoint, {
				params: {
					search: searchQuery,
					tick: fixedTick,
					limit: limit,
					offset: offset,
					sector: selectedSector
				}
			});

			const newNews = res.data.news || [];
			if (isLoadMore) {
				news = [...news, ...newNews];
			} else {
				news = newNews;
			}

			if (fixedTick === null && res.data.tick !== undefined) {
				fixedTick = res.data.tick;
			}

			hasMore = newNews.length === limit;
		} catch (e) {
			err = e.message;
		} finally {
			loading = false;
			loadingMore = false;
		}
	}

	onMount(() => {
		fetchSectors();
		fetchNews();
	});

	function handleSectorChange(sector) {
		selectedSector = sector;
		fetchNews();
	}

	let debounceTimer;
	function handleSearch() {
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			fetchNews();
		}, 1000);
	}

	function manualSearch() {
		clearTimeout(debounceTimer);
		fetchNews();
	}

	function loadMore() {
		offset += limit;
		fetchNews(true);
	}
</script>

<div class="mx-auto px-20">
	<div class="mb-8 flex items-center justify-between">
		<div class="flex items-center gap-3">
			<Newspaper size={32} />
			<h1 class="text-3xl font-bold">News</h1>
		</div>
		{#if fixedTick !== null}
			<div class="text-sm opacity-50">
				Fixed at Tick: {fixedTick}
			</div>
		{/if}
	</div>

	{#if loading && offset === 0}
		<div class="flex items-center justify-center py-12">
			<Loader size={24} class="animate-spin text-gray-400" />
		</div>
	{:else if err}
		<div class="notion-card flex items-start gap-3 bg-red-50 dark:bg-red-950/20">
			<AlertCircle size={20} class="mt-0.5 text-[#EB5757]" />
			<div>
				<p class="font-medium text-[#EB5757]">Error loading news</p>
				<p class="text-sm text-gray-600 dark:text-gray-400">{err}</p>
			</div>
		</div>
	{:else}
		<!-- Filters Row -->
		<div class="mb-8 flex flex-wrap items-center gap-4">
			<!-- Search Bar -->
			<div class="flex max-w-lg flex-1 gap-2">
				<div class="relative flex-1">
					<Search size={18} class="absolute top-1/2 left-3 -translate-y-1/2 opacity-30" />
					<input
						type="text"
						bind:value={searchQuery}
						oninput={handleSearch}
						placeholder="Search news..."
						class="w-full rounded-lg border border-(--border-color) bg-(--bg-primary)
							py-2.5 pr-4 pl-10 text-sm
							transition-colors focus:ring-2 focus:ring-gray-300 focus:outline-none dark:focus:ring-gray-700"
					/>
				</div>
				<button
					onclick={manualSearch}
					class="rounded-md border border-(--border-color) bg-(--bg-primary) px-3 py-2 text-sm font-medium transition-colors hover:bg-(--bg-hover)"
				>
					Search
				</button>
			</div>

			<div class="flex items-center gap-3">
				<span class="text-sm font-bold uppercase opacity-40">Filter by Sector:</span>
				<div class="relative">
					<select
						value={selectedSector}
						onchange={(e) => handleSectorChange(e.target.value)}
						class="cursor-pointer appearance-none rounded-xl border border-(--border-color) bg-(--bg-hover) px-6 py-2.5 pr-10 text-sm font-black text-[#27AE60] transition-all hover:bg-(--bg-primary) focus:ring-4 focus:ring-[#27AE60]/10 focus:outline-none"
					>
						<option value="">All Sectors</option>
						{#each sectors as sector}
							<option value={sector}>{sector}</option>
						{/each}
					</select>
					<div class="pointer-events-none absolute top-1/2 right-4 -translate-y-1/2 opacity-50">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="lucide lucide-chevron-down"><path d="m6 9 6 6 6-6" /></svg
						>
					</div>
				</div>
			</div>
		</div>

		{#if news.length === 0}
			<div class="notion-card py-12 text-center">
				<Newspaper size={48} class="mx-auto mb-3 text-gray-300 dark:text-gray-700" />
				<p class="text-gray-500 dark:text-gray-400">
					{searchQuery ? `No news found for "${searchQuery}"` : 'No news available at the moment'}
				</p>
			</div>
		{:else}
			<div class="space-y-1">
				{#each news as n}
					<div class="notion-card">
						<div class="mb-2 flex items-center justify-between">
							<h2 class="text-2xl font-bold">{n.title}</h2>
							<div class="flex items-center gap-2">
								{#if n.company_symbol}
									<span
										class="rounded bg-(--bg-hover) px-2 py-1 text-[10px] font-black tracking-widest text-[#27AE60] uppercase"
									>
										{n.company_symbol}
									</span>
								{/if}
								<span class="text-xs font-bold opacity-50"
									>{n.company_name || `Company #${n.company_id}`}</span
								>
							</div>
						</div>
						<p class="text-lg leading-relaxed text-gray-600 dark:text-gray-300">
							{n.content}
						</p>
						{#if n.release_date}
							<p class="mt-3 text-xs text-gray-400">
								{new Date(n.release_date).toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'long',
									day: 'numeric'
								})}
							</p>
						{/if}
					</div>
				{/each}
			</div>

			{#if hasMore}
				<div class="mt-8 flex justify-center pb-12">
					<button
						onclick={loadMore}
						disabled={loadingMore}
						class="notion-card bg-(--bg-primary) px-8 py-2 text-sm font-medium transition-colors hover:bg-(--bg-hover) disabled:opacity-50"
					>
						{#if loadingMore}
							<Loader size={16} class="mr-2 inline animate-spin" /> Loading...
						{:else}
							Load More
						{/if}
					</button>
				</div>
			{/if}
		{/if}
	{/if}
</div>
