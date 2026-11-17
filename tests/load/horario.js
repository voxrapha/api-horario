import http from "k6/http";
import { check } from "k6";

const targetRPM = Number(__ENV.TARGET_RPM || 600);
const duration = __ENV.DURATION || "1m";
const baseURL = (__ENV.BASE_URL || "http://localhost:8080").replace(/\/$/, "");

export const options = {
  scenarios: {
    horario: {
      executor: "constant-arrival-rate",
      rate: targetRPM,
      timeUnit: "1m",
      duration,
      preAllocatedVUs: 5,
      maxVUs: 25,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"],
    dropped_iterations: ["count==0"],
  },
};

export default function () {
  const response = http.get(
    `${baseURL}/api/v1/horario?timezone=America%2FSao_Paulo`,
  );

  const validStatus = response.status === 200;
  const body = validStatus ? response.json() : {};

  check(response, {
    "status 200": () => validStatus,
    "content type JSON": (res) =>
      res.headers["Content-Type"]?.includes("application/json") ?? false,
    "payload contains timezone and date/time": () =>
      body.timezone === "America/Sao_Paulo" &&
      typeof body.data === "string" &&
      typeof body.hora === "string",
  });
}
