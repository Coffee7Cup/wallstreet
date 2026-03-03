<script>
	import { onMount, onDestroy } from 'svelte';
	import api from '$lib/api';
	import { token } from '$lib/stores/auth';
	import {
		Play,
		Square,
		Pause,
		RotateCcw,
		Activity,
		Users,
		Clock,
		AlertCircle,
		Terminal,
		RefreshCw,
		Cpu,
		HardDrive,
		TrendingUp,
		TrendingDown
	} from 'lucide-svelte';
	import gsap from 'gsap';
	import { PUBLIC_HOST } from '$env/static/public';
	import MonitorChart from '$lib/components/MonitorChart.svelte';

	let stats = $state(null);
	let loading = $state(true);
	let actionLoading = $state(false);
	let startTick = $state(0);
	let logs = $state('');
	let logsLoading = $state(false);
	let ws = $state(null);
	let monitoringHistory = $state({ cpu: [], connections: [], memory: [], labels: [] });
	const MAX_HISTORY = 40;

	let areYouSure = $state({
		show: false,
		message: null,
		color: null,
		callback: null
	});

	// Trades state
	let trades = $state([]);
	let tradesLoading = $state(false);
	let userIdFilter = $state('');
	let tradesError = $state(null);

	function connectWS() {
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const url = `${protocol}//${PUBLIC_HOST}:3000/api/v1/admin/monitor/ws?token=${$token}`;

		ws = new WebSocket(url);

		ws.onmessage = (event) => {
			stats = JSON.parse(event.data);
			loading = false;

			const cpuVal = parseFloat(stats.cpu_usage);
			const connVal = stats.active_connections;
			const memVal = parseFloat(stats.ram_usage); // Extract number from "123.45 MB"
			const timeLabel = new Date().toLocaleTimeString();

			monitoringHistory = {
				cpu: [...monitoringHistory.cpu.slice(-MAX_HISTORY + 1), cpuVal],
				connections: [...monitoringHistory.connections.slice(-MAX_HISTORY + 1), connVal],
				memory: [...monitoringHistory.memory.slice(-MAX_HISTORY + 1), memVal],
				labels: [...monitoringHistory.labels.slice(-MAX_HISTORY + 1), timeLabel]
			};
		};

		ws.onclose = () => {
			console.log('Admin WS closed, retrying...');
			setTimeout(connectWS, 3000);
		};

		ws.onerror = (err) => {
			console.error('Admin WS error', err);
		};
	}

	async function fetchLogs() {
		logsLoading = true;
		try {
			const res = await api.get('/admin/stats/logs');
			logs = res.data;
		} catch (err) {
			console.error('Failed to fetch logs', err);
		} finally {
			logsLoading = false;
		}
	}

	onMount(() => {
		connectWS();
		fetchLogs();

		gsap.from('.admin-card', {
			y: 20,
			opacity: 0,
			stagger: 0.1,
			duration: 0.6,
			ease: 'power2.out'
		});
	});

	onDestroy(() => {
		if (ws) ws.close();
	});

	async function handleReset() {
		actionLoading = true;
		try {
			await api.post('/admin/simulation/reset');
		} catch (err) {
			alert(`Simulation reset failed: ` + (err.response?.data?.error || err.message));
		} finally {
			actionLoading = false;
		}
	}

	async function handleEngine(action) {
		actionLoading = true;
		try {
			if (action === 'start') {
				await api.post('/admin/engine/start', { start_tick: startTick });
			} else {
				await api.post(`/admin/engine/${action}`);
			}
		} catch (err) {
			alert(`Engine ${action} failed: ` + (err.response?.data?.error || err.message));
		} finally {
			actionLoading = false;
		}
	}

	async function fetchTrades() {
		if (!userIdFilter || userIdFilter.trim() === '') {
			tradesError = 'Please enter a user ID';
			return;
		}

		tradesLoading = true;
		tradesError = null;
		try {
			const res = await api.get(`/admin/trade/trades/${userIdFilter}`);
			trades = res.data || [];
		} catch (err) {
			console.error('Failed to fetch trades', err);
			tradesError = err.response?.data?.error || 'Failed to load trades';
			trades = [];
		} finally {
			tradesLoading = false;
		}
	}

	function formatDate(dateString) {
		return new Date(dateString).toLocaleString();
	}
</script>

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

<div class="space-y-8">
	<header class="flex items-center justify-between">
		<div>
			<h1 class="mb-1 text-3xl font-bold">Admin Control Center</h1>
			<p class="opacity-60">Monitor and manage the WallStreet simulation engine</p>
		</div>
		<div
			class="flex items-center gap-2 rounded border border-(--border-color) bg-(--bg-hover) px-3 py-1 text-xs font-bold tracking-wider uppercase"
		>
			<span class="h-2 w-2 rounded-full {ws?.readyState === 1 ? 'bg-[#27AE60]' : 'bg-[#EB5757]'}"
			></span>
			{ws?.readyState === 1 ? 'Live' : 'Disconnected'}
		</div>
	</header>

	<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
		<!-- Connection Stats -->
		<div class="admin-card notion-card bg-(--bg-primary) p-6">
			<div class="mb-4 flex items-center gap-3 opacity-70">
				<Users size={20} />
				<span class="text-sm font-bold tracking-wider uppercase">Traders</span>
			</div>
			<p class="text-4xl font-bold">{stats?.active_connections || 0}</p>
			<p class="mt-2 text-sm opacity-60">Connected Hub</p>
		</div>

		<!-- Current Tick -->
		<div class="admin-card notion-card bg-(--bg-primary) p-6">
			<div class="mb-4 flex items-center gap-3 opacity-70">
				<Clock size={20} />
				<span class="text-sm font-bold tracking-wider uppercase">Tick</span>
			</div>
			<p class="text-4xl font-bold">{stats?.simulation_tick || 0}</p>
			<p class="mt-2 text-sm opacity-60">Simulation Step</p>
		</div>

		<!-- CPU Usage -->
		<div class="admin-card notion-card bg-(--bg-primary) p-6 transition-all hover:border-[#27AE60]">
			<div class="mb-4 flex items-center gap-3 opacity-70">
				<Cpu size={20} />
				<span class="text-sm font-bold tracking-wider uppercase">CPU</span>
			</div>
			<p class="text-4xl font-bold">{stats?.cpu_usage || '0.00%'}</p>
			<p class="mt-2 text-sm opacity-60">Server Load</p>
		</div>

		<!-- RAM Usage -->
		<div class="admin-card notion-card bg-(--bg-primary) p-6 transition-all hover:border-[#F2C94C]">
			<div class="mb-4 flex items-center gap-3 opacity-70">
				<HardDrive size={20} />
				<span class="text-sm font-bold tracking-wider uppercase">RAM</span>
			</div>
			<p class="truncate text-lg font-bold">{stats?.ram_usage || '0 / 0 GB'}</p>
			<p class="mt-2 text-sm opacity-60">Memory Usage</p>
		</div>
	</div>

	<!-- Monitoring Graph -->
	<div class="admin-card notion-card overflow-hidden bg-(--bg-primary) p-0">
		<div class="border-b border-(--border-color) p-6">
			<h3 class="flex items-center gap-2 text-xl font-bold">
				<Activity size={20} class="text-[#27AE60]" /> System Performance
			</h3>
		</div>
		<div class="h-64 w-full p-6">
			<div class="h-full w-full">
				<MonitorChart history={monitoringHistory} />
			</div>
			<div class="mt-4 flex gap-6 text-xs font-bold tracking-wider uppercase">
				<div class="flex items-center gap-2">
					<div class="h-3 w-3 rounded-full bg-[#27AE60]"></div>
					<span class="opacity-60">CPU Usage (%)</span>
				</div>
				<div class="flex items-center gap-2">
					<div class="h-3 w-3 rounded-full bg-[#2D9CDB]"></div>
					<span class="opacity-60">Connections</span>
				</div>
				<div class="flex items-center gap-2">
					<div class="h-3 w-3 rounded-full bg-[#9B59B6]"></div>
					<span class="opacity-60">Memory (GB)</span>
				</div>
			</div>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
		<!-- Engine Controls -->
		<div class="admin-card notion-card bg-(--bg-primary) p-8">
			<h3
				class="mb-6 flex items-center gap-2 border-b border-(--border-color) pb-4 text-xl font-bold"
			>
				<Activity size={20} class="text-[#EB5757]" /> Engine Status:
				<span class={stats?.is_active ? 'text-[#27AE60]' : 'text-[#EB5757]'}>
					{stats?.is_active ? 'WORKING' : 'OFFLINE'}
				</span>
			</h3>

			<div class="flex flex-wrap gap-4">
				{#if !stats?.is_active}
					<div class="mb-4 flex w-full items-center gap-3 md:mb-0 md:w-auto">
						<input
							type="number"
							bind:value={startTick}
							placeholder="Start Tick"
							class="notion-card w-32 rounded border-(--border-color) bg-transparent px-4 py-2"
						/>
						<button
							onclick={() => {
								areYouSure.show = true;
								areYouSure.color = '#27AE60';
								areYouSure.message = 'Start Engine';
								areYouSure.callback = () => handleEngine('start');
							}}
							disabled={actionLoading}
							class="btn-green flex items-center gap-2 rounded px-4 py-2 font-bold"
						>
							<Play size={18} /> Start Engine
						</button>
					</div>
				{:else}
					<button
						onclick={() => {
							areYouSure.show = true;
							areYouSure.color = '#EB5757';
							areYouSure.message = 'Stop Engine';
							areYouSure.callback = () => handleEngine('stop');
						}}
						disabled={actionLoading}
						class="btn-red flex items-center gap-2 rounded px-4 py-2 font-bold"
					>
						<Square size={18} /> Stop
					</button>

					<button
						onclick={() => {
							areYouSure.show = true;
							areYouSure.color = stats?.is_paused ? '#27AE60' : '#EB5757';
							areYouSure.message = stats?.is_paused ? 'Resume Engine' : 'Pause Engine';
							areYouSure.callback = () => handleEngine(stats?.is_paused ? 'resume' : 'pause');
						}}
						disabled={actionLoading}
						class="btn-yellow flex items-center gap-2 rounded bg-[#F2C94C] px-4 py-2 font-bold text-black"
					>
						{#if stats?.is_paused}
							<RotateCcw size={18} /> Resume
						{:else}
							<Pause size={18} /> Pause
						{/if}
					</button>
				{/if}

				<button
					onclick={() => {
						areYouSure.show = true;
						areYouSure.color = '#EB5757';
						areYouSure.message = 'Reset Engine';
						areYouSure.callback = () => handleReset();
					}}
					disabled={actionLoading}
					class="btn-red flex items-center gap-2 rounded px-4 py-2 font-bold"
				>
					<RefreshCw size={18} /> Reset
				</button>
			</div>

			<div
				class="bg-opacity-5 mt-8 flex items-start gap-3 rounded bg-[#EB5757] p-4 text-sm text-black"
			>
				<AlertCircle size={20} class="shrink-0" />
				<p>
					<strong>Caution:</strong> Actions here directly affect all connected traders. Resetting the
					engine might cause data loss for current round.
				</p>
			</div>
		</div>

		<!-- System Logs -->
		<div class="admin-card notion-card flex flex-col bg-(--bg-primary) p-8">
			<div class="mb-6 flex items-center justify-between border-b border-(--border-color) pb-4">
				<h3 class="flex items-center gap-2 text-xl font-bold">
					<Terminal size={20} /> Engine Logs
				</h3>
				<button
					onclick={fetchLogs}
					disabled={logsLoading}
					class="rounded p-2 transition-all hover:bg-(--bg-hover)"
					class:animate-spin={logsLoading}
				>
					<RefreshCw size={18} />
				</button>
			</div>

			<div
				class="relative max-h-[400px] grow overflow-y-auto rounded-lg bg-black p-4 font-mono text-xs whitespace-pre-wrap"
				id="log-container"
			>
				{#if logs}
					{#each logs.split('\n').filter((l) => l.trim()) as line}
						{@const entry = (() => {
							try {
								return JSON.parse(line);
							} catch (e) {
								return null;
							}
						})()}
						{#if entry}
							<div class="mb-1 border-b border-white/5 pb-1 last:border-0">
								<span class="opacity-40">[{entry.ts}]</span>
								<span
									class="font-bold uppercase {entry.level === 'error'
										? 'text-black'
										: entry.level === 'warn'
											? 'text-[#F2C94C]'
											: 'text-[#27AE60]'}"
								>
									{entry.level}
								</span>
								<span class="text-white/80">{entry.msg}</span>
								{#if entry.error}
									<p class="mt-1 text-black opacity-80">Error: {entry.error}</p>
								{/if}
								{#if entry.category}
									<span class="ml-2 rounded bg-white/10 px-1 text-[10px] opacity-40"
										>{entry.category}</span
									>
								{/if}
							</div>
						{:else}
							<div class="text-[#27AE60]">{line}</div>
						{/if}
					{/each}
				{:else}
					<span class="opacity-40">-- no logs available --</span>
				{/if}
			</div>
		</div>
	</div>

	<!-- User Trades Viewer -->
	<div class="admin-card notion-card bg-(--bg-primary) p-8">
		<div class="mb-6 border-b border-(--border-color) pb-4">
			<h3 class="mb-4 flex items-center gap-2 text-xl font-bold">
				<TrendingUp size={20} class="text-[#2D9CDB]" /> User Trades Viewer
			</h3>
			<div class="flex gap-3">
				<input
					type="number"
					bind:value={userIdFilter}
					placeholder="Enter User ID"
					class="notion-card flex-1 rounded border-(--border-color) bg-transparent px-4 py-2"
					onkeydown={(e) => e.key === 'Enter' && fetchTrades()}
				/>
				<button
					onclick={fetchTrades}
					disabled={tradesLoading}
					class="btn-green flex items-center gap-2 rounded px-6 py-2 font-bold"
				>
					{#if tradesLoading}
						<RefreshCw size={18} class="animate-spin" />
					{:else}
						<TrendingUp size={18} />
					{/if}
					Load Trades
				</button>
			</div>
		</div>

		{#if tradesError}
			<div class="rounded bg-[#EB5757]/10 p-4 text-center text-[#EB5757]">
				{tradesError}
			</div>
		{:else if trades.length === 0 && !tradesLoading}
			<div class="rounded bg-(--bg-hover) p-8 text-center opacity-60">
				Enter a user ID to view their trade history
			</div>
		{:else if tradesLoading}
			<div class="flex h-32 items-center justify-center">
				<div
					class="h-8 w-8 animate-spin rounded-full border-4 border-(--border-color) border-t-[#27AE60]"
				></div>
			</div>
		{:else}
			<div class="overflow-x-auto rounded-lg border border-(--border-color)">
				<table class="w-full">
					<thead class="border-b border-(--border-color) bg-(--bg-hover)">
						<tr>
							<th class="px-4 py-3 text-left text-xs font-bold tracking-wider uppercase opacity-60">
								ID
							</th>
							<th class="px-4 py-3 text-left text-xs font-bold tracking-wider uppercase opacity-60">
								Company
							</th>
							<th class="px-4 py-3 text-left text-xs font-bold tracking-wider uppercase opacity-60">
								Type
							</th>
							<th class="px-4 py-3 text-left text-xs font-bold tracking-wider uppercase opacity-60">
								Quantity
							</th>
							<th class="px-4 py-3 text-left text-xs font-bold tracking-wider uppercase opacity-60">
								Date
							</th>
						</tr>
					</thead>
					<tbody>
						{#each trades as trade}
							<tr class="border-b border-(--border-color) transition-colors hover:bg-(--bg-hover)">
								<td class="px-4 py-3 font-mono text-sm opacity-60">#{trade.id}</td>
								<td class="px-4 py-3 font-medium">Company #{trade.company_id}</td>
								<td class="px-4 py-3">
									<span
										class="inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs font-bold uppercase {trade.trade_type ===
										'BUY'
											? 'bg-[#27AE60]/10 text-[#27AE60]'
											: 'bg-[#EB5757]/10 text-[#EB5757]'}"
									>
										{#if trade.trade_type === 'BUY'}
											<TrendingUp size={12} />
										{:else}
											<TrendingDown size={12} />
										{/if}
										{trade.trade_type}
									</span>
								</td>
								<td class="px-4 py-3 font-bold">{trade.quantity}</td>
								<td class="px-4 py-3 text-sm opacity-60">{formatDate(trade.timestamp)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			<div class="mt-4 text-center text-sm opacity-60">
				Showing <span class="font-bold">{trades.length}</span> trades for User #{userIdFilter}
			</div>
		{/if}
	</div>
</div>
