import http from "k6/http";
import ws from "k6/ws";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    default: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "5s", target: 10 },
        { duration: "5s", target: 100 },
        { duration: "10s", target: 500 },
        { duration: "10s", target: 1000 },
        { duration: "5s", target: 10 },
        { duration: "10s", target: 1000 },

      ],
      gracefulStop: "30s",
    },
  },
};

const BASE_URL = "http://localhost:3000/api/v1";
const WS_URL = "ws://localhost:3000/api/v1/trade/ws";

export default function () {
  // 1. Login to get JWT
  const loginRes = http.post(
    `${BASE_URL}/users/login`,
    JSON.stringify({ username: "jash" }),
    { headers: { "Content-Type": "application/json" } }
  );

  const isLoginOk = check(loginRes, {
    "login status is 200": (r) => r.status === 200,
  });

  if (!isLoginOk) return;

  const loginData = loginRes.json();
  const token = loginData.token;
  const userId = loginData.user.id;
  const authHeaders = { headers: { Authorization: `Bearer ${token}` } };

  // 2. Fetch valid company IDs and sectors once to avoid errors
  const companiesRes = http.get(`${BASE_URL}/market/companies`, authHeaders);
  const companyIds = (companiesRes.json().companies || []).map(c => c.id);

  const sectorsRes = http.get(`${BASE_URL}/market/sectors`, authHeaders);
  const sectors = sectorsRes.json().sectors || [];

  if (companyIds.length === 0) {
    console.error("No companies found, test might fail");
    return;
  }

  // 3. Continuous Background HTTP Requests (Simulating active browsing)
  const doBackgroundRequests = () => {
    const companyId = companyIds[Math.floor(Math.random() * companyIds.length)];
    const sector = sectors.length > 0 ? sectors[Math.floor(Math.random() * sectors.length)] : null;

    // User & Trade Routes
    http.get(`${BASE_URL}/users/profile`, authHeaders);
    http.get(`${BASE_URL}/trade/portfolio`, authHeaders);
    http.get(`${BASE_URL}/trade/trades`, authHeaders);

    // Market Routes (General)
    http.get(`${BASE_URL}/market/companies`, authHeaders);
    http.get(`${BASE_URL}/market/sectors`, authHeaders);
    http.get(`${BASE_URL}/market/news?limit=10`, authHeaders);
    if (sector) {
      http.get(`${BASE_URL}/market/news/sector?sector=${encodeURIComponent(sector)}&limit=5`, authHeaders);
    }

    // Market Routes (Company Specific)
    http.get(`${BASE_URL}/market/companies/${companyId}`, authHeaders);
    http.get(`${BASE_URL}/market/companies/${companyId}/peers`, authHeaders);
    http.get(`${BASE_URL}/market/profit-loss/${companyId}`, authHeaders);
    http.get(`${BASE_URL}/market/balance-sheets/${companyId}`, authHeaders);
    http.get(`${BASE_URL}/market/cash-flows/${companyId}`, authHeaders);
    http.get(`${BASE_URL}/market/ratios/${companyId}`, authHeaders);
    http.get(`${BASE_URL}/market/stocks-history/${companyId}`, authHeaders);
    http.get(`${BASE_URL}/market/stocks-history`, authHeaders); // All history
    http.get(`${BASE_URL}/market/drawings/${companyId}`, authHeaders);

    // Occasional drawing save (POST)
    if (Math.random() > 0.9) {
      http.post(
        `${BASE_URL}/market/drawings/${companyId}`,
        JSON.stringify({ drawings: [] }),
        { headers: { ...authHeaders.headers, "Content-Type": "application/json" } }
      );
    }
  };

  // 4. WebSocket Connection & Trades
  const url = `${WS_URL}?token=${token}`;

  ws.connect(url, {}, function (socket) {
    socket.on("open", function () {
      // Send 10 trades per second (100ms interval)
      socket.setInterval(function () {
        const targetId = companyIds[Math.floor(Math.random() * companyIds.length)];
        const qty = Math.floor(Math.random() * 10) + 1;

        socket.send(JSON.stringify({
          user_id: userId,
          company_id: targetId,
          trade_type: Math.random() > 0.5 ? "BUY" : "SELL",
          quantity: qty,
        }));

        // Occasionally do an HTTP request to spice up the 10req/s goal
        if (Math.random() > 0.8) {
          doBackgroundRequests();
        }
      }, 100);
    });

    socket.on("message", (msg) => {
      // Logic for handling messages
    });

    // Close after session duration
    socket.setTimeout(() => socket.close(), 45000); // 45 second session
  });

  sleep(Math.random() * 2 + 1); // Randomized wait between VU iterations
}
