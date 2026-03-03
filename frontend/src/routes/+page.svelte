<script>
	import { onMount } from 'svelte';
	import gsap from 'gsap';
	import { ArrowRight, TrendingUp, BarChart3, Zap, Users, Globe, Sparkles } from 'lucide-svelte';

	onMount(() => {
		const tl = gsap.timeline();
		
		// Hero animation
		tl.from('.hero-badge', {
			y: -20,
			opacity: 0,
			duration: 0.6,
			ease: 'power2.out'
		})
		.from('.hero-title', {
			y: 40,
			opacity: 0,
			duration: 0.8,
			ease: 'power3.out'
		}, '-=0.3')
		.from('.hero-subtitle', {
			y: 30,
			opacity: 0,
			duration: 0.7,
			ease: 'power2.out'
		}, '-=0.4')
		.from('.hero-cta', {
			y: 20,
			opacity: 0,
			duration: 0.6,
			ease: 'power2.out'
		}, '-=0.3')
		.from('.stats-card', {
			scale: 0.95,
			opacity: 0,
			stagger: 0.1,
			duration: 0.5,
			ease: 'back.out(1.2)'
		}, '-=0.2')
		.from('.feature-item', {
			x: -30,
			opacity: 0,
			stagger: 0.15,
			duration: 0.6,
			ease: 'power2.out'
		}, '-=0.3');

		// Floating animation for badge
		gsap.to('.hero-badge', {
			y: -5,
			duration: 2,
			ease: 'sine.inOut',
			repeat: -1,
			yoyo: true
		});
	});
</script>

<div class="flex flex-col items-center px-4 py-16 md:py-24">
	<!-- Hero Section -->
	<div class="hero-content mb-20 max-w-4xl text-center">
		<div class="hero-badge mb-6 inline-flex items-center gap-2 rounded-full border border-(--border-color) bg-(--bg-hover) px-4 py-2 text-sm font-medium">
			<Sparkles size={16} class="text-[#27AE60]" />
			Zero Risk. Real Experience.
		</div>
		
		<h1 class="hero-title mb-6 text-5xl font-bold leading-tight tracking-tight md:text-7xl">
			Master Trading<br />
			<span class="bg-gradient-to-r from-[#27AE60] to-[#2ecc71] bg-clip-text text-transparent">
				Without The Risk
			</span>
		</h1>
		
		<p class="hero-subtitle mx-auto mb-10 max-w-2xl text-lg leading-relaxed opacity-70 md:text-xl">
			Step into a live trading simulation powered by real market dynamics. Build your portfolio, 
			react to breaking news, and sharpen your instincts—all in a safe environment.
		</p>
		
		<div class="hero-cta flex flex-wrap justify-center gap-4">
			<a
				href="/login"
				class="btn-green group flex items-center gap-2 px-8 py-3.5 text-base font-semibold shadow-lg shadow-[#27AE60]/20 transition-all hover:shadow-xl hover:shadow-[#27AE60]/30"
			>
				Start Trading
				<ArrowRight size={20} class="transition-transform group-hover:translate-x-1" />
			</a>
			<a
			
				href="#how-it-works"
				class="notion-card flex items-center gap-2 px-8 py-3.5 text-base font-medium"
			>
				How It Works
			</a>
		</div>
	</div>

	<!-- Stats Cards -->
	<div class="mb-24 grid w-full max-w-5xl grid-cols-1 gap-4 md:grid-cols-3">
		<div class="stats-card notion-card text-center">
			<div class="mb-2 text-3xl font-bold text-[#27AE60]">Real-Time</div>
			<p class="text-sm opacity-70">Live market simulation</p>
		</div>
		<div class="stats-card notion-card text-center">
			<div class="mb-2 text-3xl font-bold text-[#F2C94C]">$100K</div>
			<p class="text-sm opacity-70">Starting capital</p>
		</div>
		<div class="stats-card notion-card text-center">
			<div class="mb-2 text-3xl font-bold text-[#EB5757]">24/7</div>
			<p class="text-sm opacity-70">Trading available</p>
		</div>
	</div>

	<!-- How It Works -->
	<div id="how-it-works" class="w-full max-w-5xl">
		<div class="mb-12 text-center">
			<h2 class="mb-3 text-3xl font-bold md:text-4xl">How It Works</h2>
			<p class="text-lg opacity-70">Three simple steps to start your trading journey</p>
		</div>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
			<div class="feature-item notion-card group relative overflow-hidden">
				<div class="absolute right-4 top-4 text-6xl font-bold opacity-5 transition-opacity group-hover:opacity-10">
					01
				</div>
				<div class="mb-4 flex h-14 w-14 items-center justify-center rounded-xl bg-[#27AE60] bg-opacity-10">
					<Users size={28} class="text-[#27AE60]" />
				</div>
				<h3 class="mb-2 text-xl font-bold">Create Account</h3>
				<p class="leading-relaxed opacity-70">
					Sign up instantly and receive your starting capital. No credit card required.
				</p>
			</div>

			<div class="feature-item notion-card group relative overflow-hidden">
				<div class="absolute right-4 top-4 text-6xl font-bold opacity-5 transition-opacity group-hover:opacity-10">
					02
				</div>
				<div class="mb-4 flex h-14 w-14 items-center justify-center rounded-xl bg-[#F2C94C] bg-opacity-10">
					<BarChart3 size={28} class="text-[#F2C94C]" />
				</div>
				<h3 class="mb-2 text-xl font-bold">Analyze Markets</h3>
				<p class="leading-relaxed opacity-70">
					Study real-time price charts, track trends, and stay updated with breaking news.
				</p>
			</div>

			<div class="feature-item notion-card group relative overflow-hidden">
				<div class="absolute right-4 top-4 text-6xl font-bold opacity-5 transition-opacity group-hover:opacity-10">
					03
				</div>
				<div class="mb-4 flex h-14 w-14 items-center justify-center rounded-xl bg-[#EB5757] bg-opacity-10">
					<TrendingUp size={28} class="text-[#EB5757]" />
				</div>
				<h3 class="mb-2 text-xl font-bold">Execute Trades</h3>
				<p class="leading-relaxed opacity-70">
					Buy and sell stocks in real-time. Watch your portfolio grow as you refine your strategy.
				</p>
			</div>
		</div>
	</div>

	<!-- Features Grid -->
	<div class="mt-24 w-full max-w-5xl">
		<div class="mb-12 text-center">
			<h2 class="mb-3 text-3xl font-bold md:text-4xl">Why WallStreet Sim?</h2>
			<p class="text-lg opacity-70">Everything you need to become a confident trader</p>
		</div>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
			<div class="feature-item notion-card">
				<div class="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-[#27AE60] bg-opacity-10">
					<Globe size={24} class="text-[#27AE60]" />
				</div>
				<h3 class="mb-2 text-lg font-bold">Live Market Dynamics</h3>
				<p class="text-sm leading-relaxed opacity-70">
					Experience authentic market behavior with real-time price movements, volatility, and market events.
				</p>
			</div>

			<div class="feature-item notion-card">
				<div class="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-[#EB5757] bg-opacity-10">
					<Zap size={24} class="text-[#EB5757]" />
				</div>
				<h3 class="mb-2 text-lg font-bold">Instant Execution</h3>
				<p class="text-sm leading-relaxed opacity-70">
					No delays, no complications. Execute trades instantly and see results in real-time.
				</p>
			</div>

			<div class="feature-item notion-card">
				<div class="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-[#F2C94C] bg-opacity-10">
					<BarChart3 size={24} class="text-[#F2C94C]" />
				</div>
				<h3 class="mb-2 text-lg font-bold">Portfolio Analytics</h3>
				<p class="text-sm leading-relaxed opacity-70">
					Track your performance with detailed analytics, profit/loss charts, and key financial metrics.
				</p>
			</div>

			<div class="feature-item notion-card">
				<div class="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-[#9B59B6] bg-opacity-10">
					<TrendingUp size={24} class="text-[#9B59B6]" />
				</div>
				<h3 class="mb-2 text-lg font-bold">Learning by Doing</h3>
				<p class="text-sm leading-relaxed opacity-70">
					The best way to learn is by doing. Make mistakes, learn from them, and improve—risk-free.
				</p>
			</div>
		</div>
	</div>

	<!-- CTA Section -->
	<div class="mt-24 w-full max-w-4xl text-center">
		<div class="notion-card bg-gradient-to-br from-[#27AE60]/5 to-[#2ecc71]/5 px-8 py-16">
			<h2 class="mb-4 text-3xl font-bold md:text-4xl">Ready to Start Trading?</h2>
			<p class="mb-8 text-lg opacity-70">
				Join thousands of traders mastering the markets in our simulation
			</p>
			<a
				href="/login"
				class="btn-green inline-flex items-center gap-2 px-10 py-4 text-lg font-semibold shadow-xl shadow-[#27AE60]/20"
			>
				Get Started Free
				<ArrowRight size={22} />
			</a>
		</div>
	</div>
</div>