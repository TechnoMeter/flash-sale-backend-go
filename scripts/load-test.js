import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate } from 'k6/metrics';  

// Custom metrics
const stockDepleted = new Counter('stock_depleted');
const realErrors = new Rate('real_errors');

export const options = {
  stages: [
    { duration: '30s', target: 200 },
    { duration: '1m', target: 500 },
    { duration: '2m', target: 1000 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(99)<200'],
    real_errors: ['rate<0.01'],   // only 5xx are errors
  },
};

export default function () {
  const payload = JSON.stringify({
    product_id: 1,
    user_id: `user-${__VU}-${__ITER}`,
  });
  const params = {
    headers: { 'Content-Type': 'application/json' },
  };
  const res = http.post('http://localhost:8080/reserve', payload, params);

  // Check acceptable status codes (200 or 429)
  check(res, {
    'status is 200 or 429': (r) => r.status === 200 || r.status === 429,
  });

  // Track stock depletion separately
  if (res.status === 429) {
    stockDepleted.add(1);
  }

  // Count only 5xx as real errors
  if (res.status >= 500) {
    realErrors.add(1);
  } else {
    realErrors.add(0);
  }

  sleep(0.1);
}