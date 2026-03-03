<script>
	import { onMount, untrack } from 'svelte';
	import { browser } from '$app/environment';
	import api from '$lib/api';

	let { data = [], type = 'candle_solid', tool = 'none', clearCount = 0, id = '' } = $props();

	let container;
	let d3Module;

	// === Persistent state across redraws ===
	let zoomTransform = null; // Survives WebSocket data updates!
	let zoomBehavior = null;
	let isRestoringZoom = false; // Guard against infinite loop

	// === Drawing state ===
	let drawings = $state([]);
	let activePoints = []; // points being placed for current shape
	let hoverPos = null; // current mouse in pixel coords within inner g

	// Store current scales for mouse coord conversion (updated on each draw)
	let _xScale = null;
	let _yScale = null;
	let _margin = { top: 20, right: 55, bottom: 30, left: 20 };
	let _innerW = 0;
	let _innerH = 0;

	const TOOL_COLORS = {
		line: '#F2C94C',
		rect: '#2D9CDB',
		triangle: '#9B59B6'
	};

	let lastLoadedId = null;
	let lastClearCount = 0;

	async function loadDrawings(companyId) {
		if (!browser || !companyId) return;
		try {
			const res = await api.get(`/market/drawings/${companyId}`);
			if (res.data && res.data.drawings) {
				drawings = res.data.drawings;
			} else {
				drawings = [];
			}
		} catch (e) {
			console.error('Failed to load drawings:', e);
			drawings = [];
		}
	}

	let saveTimeout;
	function saveDrawings(companyId, currentDrawings) {
		if (!browser || !companyId) return;
		if (saveTimeout) clearTimeout(saveTimeout);
		saveTimeout = setTimeout(async () => {
			try {
				await api.post(`/market/drawings/${companyId}`, { drawings: currentDrawings });
			} catch (e) {
				console.error('Failed to save drawings:', e);
			}
		}, 1000); // 1s debounce
	}

	// === Single consolidated effect — redraws whenever any relevant prop changes ===
	$effect(() => {
		// Touch all reactive props so Svelte tracks them
		const _data = data;
		const _type = type;
		const _tool = tool;
		const _clear = clearCount;
		const _id = id;

		if (_id !== lastLoadedId) {
			untrack(() => {
				loadDrawings(_id).then(() => {
					draw();
				});
			});
			lastLoadedId = _id;
		}

		if (_clear > lastClearCount) {
			drawings = [];
			activePoints = [];
			saveDrawings(_id, []);
			lastClearCount = _clear;
		}

		draw();
	});

	// Save drawings whenever they are modified (handled inside click listener for reliability,
	// but this effect captures external modifications if any)
	$effect(() => {
		if (id && drawings) {
			saveDrawings(id, drawings);
		}
	});

	// ─── Main draw function ───────────────────────────────────────────────────
	function draw() {
		if (!browser || !container || !d3Module || data.length === 0) return;
		const d3 = d3Module;

		const W = container.clientWidth;
		const H = container.clientHeight;
		if (W < 10 || H < 10) return;

		const m = _margin;
		const iW = Math.max(1, W - m.left - m.right);
		const iH = Math.max(1, H - m.top - m.bottom);
		_innerW = iW;
		_innerH = iH;

		// --- Build base scales ---
		const timestamps = data.map((d) => d.timestamp);
		// Add half-bar padding on each side so first/last candles aren't clipped
		const tMin = Math.min(...timestamps);
		const tMax = Math.max(...timestamps);
		const tSpan = tMax - tMin || 1;
		// Half-bar pad so first/last candle centers are at the edges, not clipped
		const tPad = data.length > 1 ? (tSpan / (data.length - 1)) * 0.6 : tSpan * 0.05;
		const xFull = d3
			.scaleTime()
			.domain([tMin - tPad, tMax + tPad])
			.range([0, iW]);

		// --- Apply zoom transform to x ---
		const xScale = zoomTransform ? zoomTransform.rescaleX(xFull) : xFull;

		// --- Y scale: auto-fit to VISIBLE candles only (natural zoom like TradingView) ---
		const [visStartRaw, visEndRaw] = xScale.domain();
		const visStart = visStartRaw instanceof Date ? visStartRaw.getTime() : Number(visStartRaw);
		const visEnd = visEndRaw instanceof Date ? visEndRaw.getTime() : Number(visEndRaw);
		const visibleData = data.filter(
			(d) => d.timestamp >= visStart - tPad && d.timestamp <= visEnd + tPad
		);
		const yData = visibleData.length > 0 ? visibleData : data;
		const yMin = Math.min(...yData.map((d) => d.low));
		const yMax = Math.max(...yData.map((d) => d.high));
		const yPad = (yMax - yMin) * 0.06 || 1;
		const yScale = d3
			.scaleLinear()
			.domain([yMin - yPad, yMax + yPad])
			.range([iH, 0]);

		_xScale = xScale;
		_yScale = yScale;

		// --- Rebuild SVG ---
		d3.select(container).selectAll('svg').remove();
		d3.select(container).selectAll('.d3-tooltip').remove();

		const tooltipEl = d3
			.select(container)
			.append('div')
			.attr('class', 'd3-tooltip')
			.style('position', 'absolute')
			.style('background', 'rgba(0,0,0,0.85)')
			.style('color', 'white')
			.style('padding', '7px 11px')
			.style('border-radius', '8px')
			.style('pointer-events', 'none')
			.style('font-size', '11px')
			.style('display', 'none')
			.style('z-index', '20')
			.style('border', '1px solid rgba(255,255,255,0.1)')
			.style('line-height', '1.7');

		const svg = d3
			.select(container)
			.append('svg')
			.attr('width', W)
			.attr('height', H)
			.style('display', 'block');

		// Clip path so candles/drawings don't leak outside the chart area
		const defs = svg.append('defs');
		defs
			.append('clipPath')
			.attr('id', 'chart-clip')
			.append('rect')
			.attr('width', iW)
			.attr('height', iH);

		const g = svg.append('g').attr('transform', `translate(${m.left},${m.top})`);

		// --- Grid ---
		g.append('g')
			.call(d3.axisLeft(yScale).ticks(5).tickSize(-iW).tickFormat(''))
			.call((gg) => gg.select('.domain').remove())
			.call((gg) => gg.selectAll('.tick line').attr('stroke', 'rgba(255,255,255,0.05)'));

		// --- Axes ---
		g.append('g')
			.attr('transform', `translate(0,${iH})`)
			.call(
				d3
					.axisBottom(xScale)
					.ticks(Math.floor(iW / 90))
					.tickSizeOuter(0)
			)
			.call((gg) => gg.select('.domain').attr('stroke', 'rgba(255,255,255,0.1)'))
			.call((gg) => gg.selectAll('.tick line').attr('stroke', 'rgba(255,255,255,0.1)'))
			.call((gg) =>
				gg.selectAll('.tick text').attr('fill', 'rgba(255,255,255,0.4)').style('font-size', '10px')
			);

		g.append('g')
			.attr('transform', `translate(${iW},0)`)
			.call(d3.axisRight(yScale).ticks(6))
			.call((gg) => gg.select('.domain').attr('stroke', 'rgba(255,255,255,0.1)'))
			.call((gg) => gg.selectAll('.tick line').attr('stroke', 'rgba(255,255,255,0.1)'))
			.call((gg) =>
				gg.selectAll('.tick text').attr('fill', 'rgba(255,255,255,0.4)').style('font-size', '10px')
			);

		// --- Data (candles/area) with clip ---
		const dataG = g.append('g').attr('clip-path', 'url(#chart-clip)');

		if (type === 'area') {
			const gradId = 'area-grad-' + Date.now();
			defs
				.append('linearGradient')
				.attr('id', gradId)
				.attr('x1', '0%')
				.attr('y1', '0%')
				.attr('x2', '0%')
				.attr('y2', '100%')
				.selectAll('stop')
				.data([
					{ offset: '0%', color: '#27AE60', opacity: 0.28 },
					{ offset: '100%', color: '#27AE60', opacity: 0 }
				])
				.enter()
				.append('stop')
				.attr('offset', (d) => d.offset)
				.attr('stop-color', (d) => d.color)
				.attr('stop-opacity', (d) => d.opacity);

			const area = d3
				.area()
				.x((d) => xScale(d.timestamp))
				.y0(iH)
				.y1((d) => yScale(d.close))
				.curve(d3.curveMonotoneX);
			const line = d3
				.line()
				.x((d) => xScale(d.timestamp))
				.y((d) => yScale(d.close))
				.curve(d3.curveMonotoneX);

			dataG.append('path').datum(data).attr('fill', `url(#${gradId})`).attr('d', area);
			dataG
				.append('path')
				.datum(data)
				.attr('fill', 'none')
				.attr('stroke', '#27AE60')
				.attr('stroke-width', 2)
				.attr('d', line);
		} else {
			// Candle width derived from real pixel spacing so it grows when zoomed
			const spacingPx =
				data.length > 1 ? Math.abs(xScale(data[1].timestamp) - xScale(data[0].timestamp)) : iW;
			const cw = Math.max(1, Math.min(40, spacingPx * 0.7));
			const candles = dataG
				.selectAll('.candle')
				.data(data)
				.enter()
				.append('g')
				.attr('class', 'candle')
				.attr('transform', (d) => `translate(${xScale(d.timestamp)},0)`);

			candles
				.append('line')
				.attr('y1', (d) => yScale(d.low))
				.attr('y2', (d) => yScale(d.high))
				.attr('stroke', (d) => (d.close >= d.open ? '#27AE60' : '#EB5757'))
				.attr('stroke-width', 1);

			candles
				.append('rect')
				.attr('x', -cw / 2)
				.attr('y', (d) => yScale(Math.max(d.open, d.close)))
				.attr('width', cw)
				.attr('height', (d) => Math.max(1, Math.abs(yScale(d.open) - yScale(d.close))))
				.attr('fill', (d) => (d.close >= d.open ? '#27AE60' : '#EB5757'))
				.attr('rx', 1);
		}

		// --- Drawings layer (clipped, redrawn with correct scales) ---
		const drawG = g.append('g').attr('clip-path', 'url(#chart-clip)');
		renderDrawings(drawG, d3, xScale, yScale);

		// --- Crosshair group (hidden by default) ---
		const crossG = g.append('g').style('display', 'none');
		crossG
			.append('line')
			.attr('class', 'cross-v')
			.attr('y1', 0)
			.attr('y2', iH)
			.attr('stroke', 'rgba(255,255,255,0.25)')
			.attr('stroke-width', 1)
			.attr('stroke-dasharray', '4,3');
		crossG
			.append('line')
			.attr('class', 'cross-h')
			.attr('x1', 0)
			.attr('x2', iW)
			.attr('stroke', 'rgba(255,255,255,0.25)')
			.attr('stroke-width', 1)
			.attr('stroke-dasharray', '4,3');

		// --- Interaction overlay ---
		const overlay = g
			.append('rect')
			.attr('width', iW)
			.attr('height', iH)
			.attr('fill', 'transparent')
			.style('cursor', tool !== 'none' ? 'crosshair' : 'default');

		if (tool !== 'none') {
			// Drawing mode — clicks place points
			overlay.on('click', (event) => {
				const [px, py] = d3.pointer(event);
				const dataX = _xScale.invert(px).getTime();
				const dataY = _yScale.invert(py);
				activePoints.push([dataX, dataY]);

				const needed = { line: 2, rect: 2, triangle: 3 };
				if (activePoints.length >= needed[tool]) {
					drawings = [
						...drawings,
						{ type: tool, points: [...activePoints], color: TOOL_COLORS[tool] }
					];
					activePoints = [];
					hoverPos = null;
				}
				draw();
			});

			overlay.on('mousemove', (event) => {
				const [px, py] = d3.pointer(event);
				hoverPos = [px, py];
				// Redraw just the drawing layer efficiently
				drawG.selectAll('*').remove();
				renderDrawings(drawG, d3, _xScale, _yScale);
			});

			overlay.on('mouseleave', () => {
				hoverPos = null;
				drawG.selectAll('*').remove();
				renderDrawings(drawG, d3, _xScale, _yScale);
			});
		} else {
			// Pan/zoom mode
			if (!zoomBehavior) {
				zoomBehavior = d3
					.zoom()
					.scaleExtent([1, 100]) // 1 = full data view; can't zoom out further
					.on('zoom', (event) => {
						if (isRestoringZoom) return;
						const t = event.transform;
						// Snap back to "all data" when fully zoomed out
						zoomTransform = t.k <= 1 ? null : t;
						draw();
					});
			}

			overlay.call(zoomBehavior);
			if (zoomTransform) {
				// Restore previous zoom silently (without triggering another draw)
				isRestoringZoom = true;
				overlay.call(zoomBehavior.transform, zoomTransform);
				isRestoringZoom = false;
			}

			// Crosshair + tooltip in pan mode
			overlay
				.on('mousemove', (event) => {
					const [mx] = d3.pointer(event);
					const date = _xScale.invert(mx);
					const bisect = d3.bisector((d) => d.timestamp).center;
					const idx = Math.max(0, Math.min(data.length - 1, bisect(data, date.getTime())));
					const d = data[idx];
					if (!d) return;

					crossG.style('display', null);
					crossG
						.select('.cross-v')
						.attr('x1', _xScale(d.timestamp))
						.attr('x2', _xScale(d.timestamp));
					crossG.select('.cross-h').attr('y1', _yScale(d.close)).attr('y2', _yScale(d.close));

					tooltipEl
						.style('display', 'block')
						.style('left', `${event.offsetX + 14}px`)
						.style('top', `${event.offsetY - 8}px`)
						.html(
							`<div style="font-weight:700;font-size:10px;opacity:.5;margin-bottom:2px">${new Date(d.timestamp).toLocaleDateString()}</div>
						<div>O <b>${d.open.toFixed(2)}</b> · H <b>${d.high.toFixed(2)}</b></div>
						<div>L <b>${d.low.toFixed(2)}</b> · C <b style="color:${d.close >= d.open ? '#27AE60' : '#EB5757'}">${d.close.toFixed(2)}</b></div>`
						);
				})
				.on('mouseleave', () => {
					crossG.style('display', 'none');
					tooltipEl.style('display', 'none');
				});
		}
	}

	// ─── Render all drawings + in-progress shape ──────────────────────────────
	function renderDrawings(g, d3, xScale, yScale) {
		for (const dr of drawings) {
			renderShape(g, d3, dr.type, dr.points, dr.color, xScale, yScale);
		}

		// In-progress shape preview
		if (activePoints.length > 0 && hoverPos) {
			const previewPoints = [
				...activePoints,
				[xScale.invert(hoverPos[0]).getTime(), yScale.invert(hoverPos[1])]
			];
			renderShape(g, d3, tool, previewPoints, TOOL_COLORS[tool] || '#fff', xScale, yScale, true);
		}
	}

	function renderShape(g, d3, type, points, color, xScale, yScale, preview = false) {
		if (points.length < 1) return;
		const px = points.map(([tx, ty]) => [xScale(new Date(tx)), yScale(ty)]);

		const sharedAttrs = (el) =>
			el
				.attr('stroke', color)
				.attr('stroke-width', preview ? 1.2 : 1.8)
				.attr('stroke-dasharray', preview ? '6,3' : null)
				.attr('fill', 'none');

		if (type === 'line') {
			if (px.length < 2) {
				// Show just a dot at first point
				g.append('circle')
					.attr('cx', px[0][0])
					.attr('cy', px[0][1])
					.attr('r', 3)
					.attr('fill', color);
				return;
			}
			// Extend the line across the entire chart (like TradingView)
			const slope = (px[1][1] - px[0][1]) / (px[1][0] - px[0][0]);
			const x0 = 0,
				x1 = _innerW;
			const y0 = px[0][1] + slope * (x0 - px[0][0]);
			const y1 = px[0][1] + slope * (x1 - px[0][0]);

			sharedAttrs(g.append('line')).attr('x1', x0).attr('y1', y0).attr('x2', x1).attr('y2', y1);
			// Endpoint handles
			if (!preview) {
				for (const [x, y] of px) {
					g.append('circle')
						.attr('cx', x)
						.attr('cy', y)
						.attr('r', 4)
						.attr('fill', color)
						.attr('stroke', 'white')
						.attr('stroke-width', 1);
				}
			}
		} else if (type === 'rect') {
			if (px.length < 2) {
				g.append('circle')
					.attr('cx', px[0][0])
					.attr('cy', px[0][1])
					.attr('r', 3)
					.attr('fill', color);
				return;
			}
			const x = Math.min(px[0][0], px[1][0]);
			const y = Math.min(px[0][1], px[1][1]);
			const w = Math.abs(px[1][0] - px[0][0]);
			const h = Math.abs(px[1][1] - px[0][1]);
			g.append('rect')
				.attr('x', x)
				.attr('y', y)
				.attr('width', w)
				.attr('height', h)
				.attr('fill', color + '18')
				.attr('stroke', color)
				.attr('stroke-width', preview ? 1.2 : 1.8)
				.attr('stroke-dasharray', preview ? '6,3' : null)
				.attr('rx', 2);

			if (!preview) {
				for (const [x, y] of px) {
					g.append('circle')
						.attr('cx', x)
						.attr('cy', y)
						.attr('r', 4)
						.attr('fill', color)
						.attr('stroke', 'white')
						.attr('stroke-width', 1);
				}
			}
		} else if (type === 'triangle') {
			if (px.length < 2) {
				g.append('circle')
					.attr('cx', px[0][0])
					.attr('cy', px[0][1])
					.attr('r', 3)
					.attr('fill', color);
				return;
			}
			const pts = preview ? px : [...px, px[0]]; // close the triangle
			const lineGen = d3
				.line()
				.x((d) => d[0])
				.y((d) => d[1]);
			g.append('path')
				.datum(pts)
				.attr('d', lineGen)
				.attr('fill', preview ? 'none' : color + '18')
				.attr('stroke', color)
				.attr('stroke-width', preview ? 1.2 : 1.8)
				.attr('stroke-dasharray', preview ? '6,3' : null)
				.attr('stroke-linejoin', 'round');

			if (!preview) {
				for (const [x, y] of px) {
					g.append('circle')
						.attr('cx', x)
						.attr('cy', y)
						.attr('r', 4)
						.attr('fill', color)
						.attr('stroke', 'white')
						.attr('stroke-width', 1);
				}
			}
		}
	}

	onMount(async () => {
		d3Module = await import('d3');
		draw();

		const ro = new ResizeObserver(() => draw());
		ro.observe(container);
		return () => ro.disconnect();
	});
</script>

<div class="relative h-full w-full" bind:this={container} style="min-height: 100px;"></div>
