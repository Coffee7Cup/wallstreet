<script>
	import { onMount } from 'svelte';
	import gsap from 'gsap';
	import api from '$lib/api';
	import { user, admin, token } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { Shield, User, AlertCircle, TrendingUp, ArrowUpRight } from 'lucide-svelte';

	let isAdmin = $state(false);
	let username = $state('');
	let email = $state('');
	let error = $state(null);
	let loading = $state(false);

	onMount(() => {
		const tl = gsap.timeline();

		tl.from('.login-card', {
			x: -50,
			opacity: 0,
			duration: 0.8,
			ease: 'power3.out'
		});

		tl.from(
			'.preview-card',
			{
				x: 50,
				opacity: 0,
				duration: 0.8,
				ease: 'power3.out'
			},
			'-=0.6'
		);

		tl.fromTo(
			'.chart-line',
			{ strokeDasharray: 1000, strokeDashoffset: 1000 },
			{ strokeDashoffset: 0, duration: 2, ease: 'power2.inOut' },
			'-=0.4'
		);
	});

	async function handleLogin() {
		error = '';
		loading = true;
		try {
			const endpoint = isAdmin ? '/admin/login' : '/users/login';
			const response = await api.post(endpoint, {
				username: username,
				email: email
			});

			token.set(response.data.token);
			if (isAdmin) {
				admin.set(response.data.admin);
				goto('/admin');
			} else {
				user.set(response.data.user);
				goto('/dashboard');
			}
		} catch (err) {
			error = err.response?.data?.error || 'Login failed. Please check your credentials.';
		} finally {
			loading = false;
		}
	}

	function toggleAdmin() {
		isAdmin = !isAdmin;
		username = '';
		email = '';
		error = '';
		gsap.fromTo(
			'.login-header',
			{ opacity: 0, x: isAdmin ? 10 : -10 },
			{ opacity: 1, x: 0, duration: 0.3 }
		);
	}
</script>

<div class="flex min-h-[80vh] items-center justify-center p-4">
	<div class="grid w-full max-w-5xl gap-8 lg:grid-cols-2 lg:items-center">
		<!-- Left Side: Login Form -->
		<div class="login-card notion-card bg-(--bg-primary) p-8 shadow-xl lg:p-12">
			<div class="login-header mb-8 text-center lg:text-left">
				<div
					class="mb-6 inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-(--bg-hover) shadow-inner"
				>
					{#if isAdmin}
						<Shield size={32} class="text-[#EB5757]" />
					{:else}
						<User size={32} class="text-[#27AE60]" />
					{/if}
				</div>
				<h1 class="mb-2 text-3xl font-bold tracking-tight text-(--text-primary)">
					{isAdmin ? 'Admin Portal' : 'Market Access'}
				</h1>
				<p class="text-sm opacity-60">Welcome back. Enter your details to trade.</p>
			</div>

			{#if error}
				<div
					class="border-opacity-20 bg-opacity-10 mb-6 flex items-center gap-2 rounded-lg border border-[#EB5757] bg-[#EB5757] p-4 text-sm text-black"
				>
					<AlertCircle size={18} />
					{error}
				</div>
			{/if}

			<form
				onsubmit={(e) => {
					e.preventDefault();
					handleLogin();
				}}
				class="space-y-5"
			>
				<div>
					<label for="username" class="mb-2 block text-sm font-medium opacity-70">Username</label>
					<input
						type="text"
						id="username"
						bind:value={username}
						placeholder="e.g. jdoe_trader"
						class="focus:ring-opacity-10 w-full rounded-lg border border-(--border-color) bg-transparent px-4 py-3 transition-all focus:border-[#27AE60] focus:ring-4 focus:ring-[#27AE60] focus:outline-none"
					/>
				</div>

				<div>
					<label for="email" class="mb-2 block text-sm font-medium opacity-70"
						>Email (Optional for users)</label
					>
					<input
						type="text"
						id="email"
						bind:value={email}
						placeholder="your@email.com"
						class="focus:ring-opacity-10 w-full rounded-lg border border-(--border-color) bg-transparent px-4 py-3 transition-all focus:border-[#27AE60] focus:ring-4 focus:ring-[#27AE60] focus:outline-none"
					/>
				</div>

				<button
					type="submit"
					disabled={loading && !(username || email)}
					class="w-full rounded-lg py-3.5 text-sm font-bold tracking-wide transition-all active:scale-[0.98] disabled:opacity-50 {isAdmin
						? 'btn-red shadow-lg shadow-red-500/20'
						: 'btn-green shadow-lg shadow-green-500/20'}"
				>
					{loading ? 'Authenticating...' : 'Sign In'}
				</button>
			</form>

			<div class="mt-10 border-t border-(--border-color) pt-8 text-center">
				<button
					onclick={toggleAdmin}
					class="text-xs font-medium tracking-widest uppercase opacity-40 transition-all hover:text-(--text-primary) hover:opacity-100"
				>
					Switch to {isAdmin ? 'User' : 'Admin'} Mode
				</button>
			</div>
		</div>

		<!-- Right Side: Decorative Trading Interface -->
		<div class="preview-card hidden space-y-6 lg:block">
			<div
				class="notion-card rounded-2xl bg-linear-to-br from-(--bg-primary) to-(--bg-hover) p-6 shadow-2xl backdrop-blur-sm"
			>
				<div class="mb-6 flex items-center justify-between">
					<div>
						<h3 class="text-xs font-bold tracking-widest uppercase opacity-40">Preview Market</h3>
						<div class="flex items-center gap-2">
							<span class="text-2xl font-black">NVDA</span>
							<span class="rounded bg-[#27AE60]/10 px-2 py-0.5 text-xs font-bold text-[#27AE60]"
								>+4.2%</span
							>
						</div>
					</div>
					<div class="text-right">
						<p class="font-mono text-2xl font-bold">₹742.18</p>
						<p class="text-xs opacity-40">Real-time Simulation</p>
					</div>
				</div>

				<!-- Dummy Graph (SVG) -->
				<div class="relative mb-8 h-48 w-full overflow-hidden rounded-xl bg-(--bg-hover)/30">
					<svg viewBox="0 0 400 200" class="h-full w-full">
						<defs>
							<linearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">
								<stop offset="0%" stop-color="#27AE60" stop-opacity="0.2" />
								<stop offset="100%" stop-color="#27AE60" stop-opacity="0" />
							</linearGradient>
						</defs>
						<!-- Area -->
						<path
							d="M0 180 Q 50 160, 100 170 T 200 120 T 300 140 T 400 80 V 200 H 0 Z"
							fill="url(#chartGradient)"
						/>
						<!-- Line -->
						<path
							class="chart-line"
							d="M0 180 Q 50 160, 100 170 T 200 120 T 300 140 T 400 80"
							fill="none"
							stroke="#27AE60"
							stroke-width="3"
							stroke-linecap="round"
						/>
					</svg>
					<div
						class="pointer-events-none absolute inset-0 flex items-center justify-center opacity-5"
					>
						<TrendingUp size={120} />
					</div>
				</div>

				<!-- Dummy Buttons -->
				<div class="grid grid-cols-2 gap-4">
					<div
						class="group relative overflow-hidden rounded-xl border border-(--border-color) bg-(--bg-primary) p-4 transition-all hover:border-[#27AE60]/50 hover:shadow-lg"
					>
						<div class="flex items-center justify-between">
							<span class="text-xs font-bold uppercase opacity-40">Buy</span>
							<ArrowUpRight size={14} class="text-[#27AE60]" />
						</div>
						<p class="mt-1 text-lg font-bold">Market Order</p>
						<div
							class="absolute -right-2 -bottom-2 h-12 w-12 rounded-full bg-[#27AE60]/5 transition-transform group-hover:scale-150"
						></div>
					</div>
					<div
						class="group relative overflow-hidden rounded-xl border border-(--border-color) bg-(--bg-primary) p-4 transition-all hover:border-[#EB5757]/50 hover:shadow-lg"
					>
						<div class="flex items-center justify-between">
							<span class="text-xs font-bold uppercase opacity-40">Sell</span>
							<ArrowUpRight size={14} class="rotate-90 text-[#EB5757]" />
						</div>
						<p class="mt-1 text-lg font-bold">Limit Order</p>
						<div
							class="absolute -right-2 -bottom-2 h-12 w-12 rounded-full bg-[#EB5757]/5 transition-transform group-hover:scale-150"
						></div>
					</div>
				</div>
			</div>

			<!-- Decorative Info Cards -->
			<div class="grid grid-cols-3 gap-4">
				{#each ['Active Trades', 'Pending', 'Closed'] as label, i}
					<div
						class="rounded-xl border border-(--border-color) p-3 text-center transition-all hover:bg-(--bg-hover)"
					>
						<p class="text-[10px] font-bold tracking-widest uppercase opacity-40">{label}</p>
						<p class="text-sm font-bold">{[24, 5, 142][i]}</p>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>

<style>
	:global(.btn-green) {
		background-color: #27ae60;
		color: white;
	}
	:global(.btn-red) {
		background-color: #eb5757;
		color: white;
	}
	:global(.btn-green:hover) {
		background-color: #219150;
	}
	:global(.btn-red:hover) {
		background-color: #d14d4d;
	}
</style>
