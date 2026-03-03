import { writable, get } from 'svelte/store';
import { token } from './auth';
import { PUBLIC_HOST } from '$env/static/public'

export const marketState = writable({
    tick: 0,
    isActive: false,
    isPaused: false,
    stocks: [],
    news: [],
    priceHistory: {}, // Symbol -> Array of prices
    connected: false,
    simulationEnded: false,
    error: null,
    lastError: null, // For transient errors like trade failures
    connectionStatus: 'disconnected', // 'connected', 'connecting', 'disconnected', 'reconnecting'
    userPortfolio: [], // Real-time portfolio entries
    userBalance: 0    // Real-time cash balance
});

let ws;
let reconnectTimeout;
let retryCount = 0;
const MAX_RETRY_DELAY = 10000;

export function connectMarketWS() {

    if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return;

    if (reconnectTimeout) clearTimeout(reconnectTimeout);

    const t = get(token);
    if (!t) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const url = `${protocol}//${PUBLIC_HOST}:3000/api/v1/trade/ws?token=${t}`;

    const currentStatus = get(marketState).connectionStatus;
    marketState.update(s => ({
        ...s,
        connectionStatus: currentStatus === 'disconnected' ? 'connecting' : 'reconnecting'
    }));

    ws = new WebSocket(url);

    ws.onopen = () => {
        console.log('Market WS connected');
        marketState.update(s => ({ ...s, connected: true, connectionStatus: 'connected', error: null }));
        retryCount = 0;
    };

    ws.onmessage = (event) => {
        // Reset retryCount and ensure status is connected on successful message
        retryCount = 0;
        if (get(marketState).connectionStatus !== 'connected') {
            marketState.update(s => ({ ...s, connectionStatus: 'connected' }));
        }

        const msg = JSON.parse(event.data);
        if (msg.type === 'INITIAL_STATE' || msg.type === 'UPDATE') {
            marketState.update(state => {
                const newHistory = { ...state.priceHistory };

                if (msg.stocks) {
                    msg.stocks.forEach(s => {
                        const id = s.company_id;
                        if (!newHistory[id]) newHistory[id] = [];
                        newHistory[id] = [...newHistory[id], {
                            price: s.close_price,
                            timestamp: s.date || new Date().toISOString()
                        }].slice(-30); // Keep last 30
                    });
                }

                const updatedNews = msg.news && msg.news.length > 0
                    ? msg.news
                    : state.news;

                return {
                    ...state,
                    tick: msg.tick,
                    isActive: msg.is_active,
                    isPaused: msg.is_paused,
                    stocks: msg.stocks || state.stocks,
                    news: updatedNews,
                    priceHistory: newHistory
                };
            });
        } else if (msg.type === 'TRADE_UPDATE') {
            console.log('Real-time trade update received');
            marketState.update(state => ({
                ...state,
                userPortfolio: msg.portfolio || state.userPortfolio,
                userBalance: msg.balance !== undefined ? msg.balance : state.userBalance
            }));
        } else if (msg.type === 'SIMULATION_ENDED') {
            console.log('Simulation has ended');
            marketState.update(s => ({ ...s, isActive: false, simulationEnded: true, connected: false }));
            if (ws) ws.close();
        } else if (msg.type === 'ERROR') {
            console.error('Market error received:', msg.error);
            marketState.update(s => ({ ...s, lastError: msg.error }));
            setTimeout(() => {
                marketState.update(s => ({ ...s, lastError: null }));
            }, 5000);
        }
    };

    ws.onclose = () => {
        console.log('Market WS closed');
        marketState.update(s => ({ ...s, connected: false, connectionStatus: 'disconnected' }));

        if (get(marketState).simulationEnded) {
            console.log('Simulation ended, stopping reconnection attempts.');
            return;
        }

        const delay = Math.min(1000 * Math.pow(2, retryCount), MAX_RETRY_DELAY);
        console.log(`Retrying in ${delay}ms... (Attempt ${retryCount + 1})`);
        retryCount++;
        reconnectTimeout = setTimeout(connectMarketWS, delay);
    };

    ws.onerror = (err) => {
        console.error('Market WS error:', err);
        // Error will often be followed by onclose, which handles retry
    };
}

export function sendTrade(companyId, tradeType, quantity, userId) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        throw new Error('WebSocket not connected');
    }

    ws.send(JSON.stringify({
        user_id: userId,
        company_id: companyId,
        trade_type: tradeType,
        quantity: quantity
    }));
}

export function closeMarketWS() {
    if (ws) ws.close();
}
