import http from "k6/http";
import ws from "k6/ws";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    default: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "5s", target: 50 },
        { duration: "3s", target: 500 },
        { duration: "3s", target: 700 },
        { duration: "30s", target: 1000 },
      ],
      gracefulStop: "30s",
    },
  },
};

const BASE_URL = "http://localhost:3000";
const WS_URL = "ws://localhost:3000/api/v1/trade/ws";

export default function () {
  // 1. Login to get JWT
  const loginRes = http.post(
    `${BASE_URL}/api/v1/users/login`,
    JSON.stringify({
      username: "jash",
    }),
    {
      headers: { "Content-Type": "application/json" },
    },
  );

  const isLoginOk = check(loginRes, {
    "logged in status is 200": (r) => r.status === 200,
    "has body": (r) => r.body !== null,
  });

  if (!isLoginOk) {
    console.error(
      `Login failed with status ${loginRes.status}: ${loginRes.body}`,
    );
    return;
  }

  let loginData;
  try {
    loginData = loginRes.json();
  } catch (e) {
    console.error(`Failed to parse login JSON: ${e}`);
    return;
  }

  const token = loginData.token;
  const user = loginData.user;

  if (!token || !user) {
    console.error("Login response missing token or user data");
    return;
  }

  const userId = user.id;

  // 2. Fetch companies to have a target company ID
  const companiesRes = http.get(`${BASE_URL}/api/v1/users/companies`);
  if (companiesRes.status !== 200) {
    console.error(`Failed to fetch companies: ${companiesRes.status}`);
    return;
  }
  const companies = companiesRes.json().companies;
  const targetCompanyId =
    companies && companies.length > 0 ? companies[0].id : 1;

  // 3. Connect via WebSocket
  const url = `${WS_URL}?token=${token}`;

  // In k6, ws.connect is used for WebSocket load testing
  const res = ws.connect(url, {}, function (socket) {
    socket.on("open", function () {
      console.log("WebSocket connected");

      // Periodically send BUY requests
      socket.setInterval(function () {
        const tradeReq = {
          user_id: userId,
          company_id: targetCompanyId,
          trade_type: "BUY",
          quantity: 1,
        };
        socket.send(JSON.stringify(tradeReq));
        console.log(`Sent BUY request for company ${targetCompanyId}`);
      }, 1000); // Send every 1 second
    });

    socket.on("message", function (message) {
      const msg = JSON.parse(message);
      if (msg.Type === "ERROR") {
        console.error(`Trade Error: ${msg.Error}`);
      } else if (msg.Type === "INITIAL_STATE") {
        console.log("Received initial state");
      }
    });

    socket.on("close", () => console.log("WebSocket closed"));
    socket.on("error", (e) => console.error("WebSocket error: ", e.error()));

    // Keep the connection open for some time
    socket.setTimeout(function () {
      socket.close();
    }, 20000); // Stay connected for 20 seconds
  });

  check(res, { "status is 101": (r) => r && r.status === 101 });

  sleep(1);
}
