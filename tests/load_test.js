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
        { duration: "20s", target: 1000 },
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

  // 2. Real GET Requests (Simulating page loads/navigation)

  // Profile
  http.get(`${BASE_URL}/users/profile`, authHeaders);

  // Companies
  const companiesRes = http.get(`${BASE_URL}/market/companies`, authHeaders);
  const companies = companiesRes.json().companies || [];
  const targetCompanyId = companies.length > 0 ? companies[0].id : 1;

  // Portfolio
  http.get(`${BASE_URL}/trade/portfolio`, authHeaders);

  // Trade History
  http.get(`${BASE_URL}/trade/trades?limit=10&offset=0`, authHeaders);

  // News (Global)
  http.get(`${BASE_URL}/market/news?limit=10`, authHeaders);

  // Sectors
  const sectorsRes = http.get(`${BASE_URL}/market/sectors`, authHeaders);
  const sectors = sectorsRes.json().sectors || [];
  const targetSector = sectors.length > 0 ? sectors[0] : "";

  // Sector News
  if (targetSector) {
    http.get(`${BASE_URL}/market/news/sector?sector=${targetSector}&limit=5`, authHeaders);
  }

  // Stock Details (Simulation of visiting a specific dashboard)
  http.get(`${BASE_URL}/market/companies/${targetCompanyId}`, authHeaders);

  // 3. WebSocket Connection & Trades
  const url = `${WS_URL}?token=${token}`;

  ws.connect(url, {}, function (socket) {
    socket.on("open", function () {
      // Send a few buy requests during the session
      socket.setInterval(function () {
        const qty = Math.floor(Math.random() * 10) + 1;
        socket.send(JSON.stringify({
          user_id: userId,
          company_id: targetCompanyId,
          trade_type: "BUY",
          quantity: qty,
        }));
      }, 5000); // Every 5 seconds
    });

    socket.on("message", (msg) => {
      const data = JSON.parse(msg);
      if (data.type === "ERROR") {
        console.warn(`Trade Error: ${data.error}`);
      }
    });

    // Close after session duration
    socket.setTimeout(() => socket.close(), 15000);
  });

  sleep(Math.random() * 2 + 1); // Randomized wait between VU iterations
}
