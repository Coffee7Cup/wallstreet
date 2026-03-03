<script>
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	let { data = [] } = $props();
	let container;
	let d3Module;

	function initChart() {
		if (!browser || !container || !d3Module || data.length < 2) return;
		const d3 = d3Module;

		const w = container.clientWidth;
		const h = container.clientHeight;
		if (w < 10 || h < 10) return;

		d3.select(container).selectAll('*').remove();

		const isUp = data[data.length - 1] >= data[0];
		const margin = { top: 2, right: 2, bottom: 2, left: 2 };
		const iw = Math.max(1, w - margin.left - margin.right);
		const ih = Math.max(1, h - margin.top - margin.bottom);

		const x = d3
			.scaleLinear()
			.domain([0, data.length - 1])
			.range([0, iw]);
		const yMin = Math.min(...data);
		const yMax = Math.max(...data);
		const y = d3
			.scaleLinear()
			.domain([yMin, yMax === yMin ? yMin + 1 : yMax])
			.range([ih, 0]);

		const line = d3
			.line()
			.x((d, i) => x(i))
			.y((d) => y(d))
			.curve(d3.curveBasis);

		const svg = d3
			.select(container)
			.append('svg')
			.attr('width', w)
			.attr('height', h)
			.style('display', 'block');

		svg
			.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`)
			.append('path')
			.datum(data)
			.attr('fill', 'none')
			.attr('stroke', isUp ? '#27AE60' : '#EB5757')
			.attr('stroke-width', 1.5)
			.attr('d', line);
	}

	$effect(() => {
		if (data) initChart();
	});

	onMount(async () => {
		d3Module = await import('d3');
		initChart();
		const ro = new ResizeObserver(() => initChart());
		ro.observe(container);
		return () => ro.disconnect();
	});
</script>

<div class="h-full w-full" bind:this={container} style="min-height: 20px;"></div>
