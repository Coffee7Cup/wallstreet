<script>
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	let { history = { cpu: [], connections: [], memory: [], labels: [] } } = $props();
	let container;
	let d3Module;

	const margin = { top: 10, right: 10, bottom: 20, left: 30 };

	function draw() {
		if (!browser || !container || !d3Module || history.cpu.length === 0) return;
		const d3 = d3Module;

		const width = container.clientWidth;
		const height = container.clientHeight;
		if (width < 50 || height < 50) return;

		const iw = Math.max(1, width - margin.left - margin.right);
		const ih = Math.max(1, height - margin.top - margin.bottom);

		d3.select(container).selectAll('svg').remove();

		const svg = d3
			.select(container)
			.append('svg')
			.attr('width', width)
			.attr('height', height)
			.style('display', 'block');
		const g = svg.append('g').attr('transform', `translate(${margin.left},${margin.top})`);

		const x = d3
			.scaleLinear()
			.domain([0, history.cpu.length - 1])
			.range([0, iw]);
		const allVals = [...history.cpu, ...history.connections, ...history.memory];
		const maxVal = Math.max(...allVals.filter((v) => !isNaN(v)), 1);
		const y = d3
			.scaleLinear()
			.domain([0, maxVal * 1.1])
			.range([ih, 0]);

		// Grid
		g.append('g')
			.call(d3.axisLeft(y).ticks(5).tickSize(-iw).tickFormat(''))
			.call((gg) => gg.select('.domain').remove())
			.call((gg) => gg.selectAll('.tick line').attr('stroke', 'rgba(255,255,255,0.05)'));

		// Lines
		const line = d3
			.line()
			.x((d, i) => x(i))
			.y((d) => y(d))
			.curve(d3.curveMonotoneX);
		const datasets = [
			{ data: history.cpu, color: '#27AE60' },
			{ data: history.connections, color: '#2D9CDB' },
			{ data: history.memory, color: '#9B59B6' }
		];

		datasets.forEach((set) => {
			if (set.data.length > 0) {
				g.append('path')
					.datum(set.data)
					.attr('fill', 'none')
					.attr('stroke', set.color)
					.attr('stroke-width', 2)
					.attr('d', line);
			}
		});

		g.append('g')
			.attr('transform', `translate(0,${ih})`)
			.call(d3.axisBottom(x).ticks(0))
			.call((gg) => gg.select('.domain').attr('stroke', 'rgba(255,255,255,0.1)'));

		g.append('g')
			.call(d3.axisLeft(y).ticks(5))
			.call((gg) => gg.select('.domain').attr('stroke', 'rgba(255,255,255,0.1)'))
			.call((gg) =>
				gg.selectAll('.tick text').attr('fill', 'rgba(255,255,255,0.4)').style('font-size', '10px')
			);
	}

	$effect(() => {
		if (history) draw();
	});

	onMount(async () => {
		d3Module = await import('d3');
		draw();
		const ro = new ResizeObserver(() => draw());
		ro.observe(container);
		return () => ro.disconnect();
	});
</script>

<div class="h-full w-full" bind:this={container} style="min-height: 50px;"></div>
