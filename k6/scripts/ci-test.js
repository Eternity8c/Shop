import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Минимальная конфигурация для CI
export const options = {
  stages: [
    { duration: '10s', target: 2 },   
    { duration: '20s', target: 5 },   
    { duration: '10s', target: 0 },   
  ],
  thresholds: {
    business_errors: ['rate<0.1'],
    http_req_duration: ['p(95)<2000'],  
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Метрики
const businessErrors = new Rate('business_errors');

// Задержка перед запусков сервисов
export function setup() {
  const maxRetries = 10;
  for (let i = 0; i < maxRetries; i++) {
    try {
      const health = http.get(BASE_URL);
      if (health && typeof health.status !== 'undefined' && health.status > 0) {
        console.log(`Service reachable (status=${health.status})`);
        return;
      }
    } catch (e) {
      console.log('Health check request threw:', String(e));
    }
    console.log(`Service not reachable yet, retry ${i + 1}/${maxRetries}`);
    sleep(2);
  }
  console.error('Service did not become reachable after retries');
}
export default function () {
  // Простой тест регистрации
  const payload = JSON.stringify({
    email: `ci_test_${Date.now()}_${__VU}_${Math.random()}@test.com`,
    password: 'Test123!@#',
    full_name: 'CI Test User',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
    // mark these requests as "expected" so thresholds can ignore acceptable 4xx responses
    tags: {
      expected_response: 'true',
    },
  };

  const res = http.post(`${BASE_URL}/register`, payload, params);
  if (!res || typeof res.status === 'undefined' || res.status === 0) {
    console.error('Request failed (no HTTP response or status 0)', String(res));
    businessErrors.add(true);
    sleep(0.5);
    return;
  }

  check(res, {
    'status is 201 or 409': (r) => r.status === 201 || r.status === 409,
    'response time < 2s': (r) => r.timings && r.timings.duration < 2000,
  });

  const ok = res.status === 201 || res.status === 409;
  if (!ok) {
    try {
      console.error('Unexpected response', { status: res.status, body: res.body });
    } catch (e) {
      console.error('Failed to stringify response body', String(e));
    }
  }
  businessErrors.add(!ok);

  sleep(0.5); // Пауза между запросами
}
