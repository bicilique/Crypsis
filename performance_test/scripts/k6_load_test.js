import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';


// Custom Metrics
const errorRate = new Rate('errors');
const encryptionDuration = new Trend('encryption_duration');

// Load Test Configuration
// Purpose: Include baseline, sustained load, and sudden bursts
export const options = {
  discardResponseBodies: true,
  stages: [
    { duration: '1m', target: 10 },  // Warm-up: ramp up to 10 VUs
    { duration: '2m', target: 10 },  // Hold baseline
    { duration: '15s', target: 50 }, // Burst #1: spike to 50 VUs instantly
    { duration: '30s', target: 50 }, // Sustain burst
    { duration: '30s', target: 10 }, // Cool down
    { duration: '1m', target: 10 },  // Stabilize again
    { duration: '10s', target: 80 }, // Burst #2: extreme spike
    { duration: '30s', target: 80 }, // Sustain burst
    { duration: '1m', target: 0 },    // Ramp down gracefully
  ],
  thresholds: {
    http_req_duration: ['p(95)<10000'],  // 95% of requests <10s
    http_req_failed: ['rate<0.05'],      // <5% failure rate
    errors: ['rate<0.05'],
    encryption_duration: ['p(95)<10000'],
  },
};

// Environment Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'ory_at_TjGvfiK8DGbf4XiogKRGxMQrNlcyEpvgU_yDOYklZXg.3VZ4PuBzqeyWWuGAeCZn2bztmt9Q4PxcCkfGLVxMcXY';


// Test files
const testFiles = {
  small: { path: open('../test_files/test_1mb.txt', 'b'), name: 'test_1mb.txt', contentType: 'text/plain' },
  medium: { path: open('../test_files/test_3mb.txt', 'b'), name: 'test_3mb.txt', contentType: 'text/plain' },
  large: { path: open('../test_files/test_5mb.txt', 'b'), name: 'test_5mb.txt', contentType: 'text/plain' },
  xlarge: { path: open('../test_files/10mb.pdf', 'b'), name: '10mb.pdf', contentType: 'application/pdf' },
};

export default function () {
  // Random file to simulate variety
  // const fileTypes = ['small', 'medium', 'large', 'xlarge'];
  const fileTypes = ['xlarge'];
  const selectedType = fileTypes[Math.floor(Math.random() * fileTypes.length)];
  const selectedFile = testFiles[selectedType];

  // Multipart payload - use k6 http.file to send binary file correctly
  const payload = {
    file: http.file(selectedFile.path, selectedFile.name, selectedFile.contentType),
  };

  const params = {
    headers: {
      Authorization: `Bearer ${AUTH_TOKEN}`,
    },
    timeout: '120s',
  };

  const startTime = Date.now();
  const res = http.post(`${BASE_URL}/api/files/encrypt`, payload, params);

  // Prefer server-side timing when available (res.timings.duration). Fallback to wall-clock.
  const duration = res && res.timings && typeof res.timings.duration === 'number'
    ? res.timings.duration
    : Date.now() - startTime;

  const success = check(res, {
    'status is 200': (r) => r && r.status === 200,
    'response time < 10s': (r) => r && r.timings && r.timings.duration < 10000,
    'no server error': (r) => r && r.status < 500,
  });

  errorRate.add(!success);
  encryptionDuration.add(duration);

  if (!success) {
    const status = res ? res.status : 'no-response';
    console.error(`[VU ${__VU}] ❌ Failed: ${selectedFile.name} - Status ${status} - Duration ${duration}ms`);
  }

  // Simulate realistic user think time
  sleep(Math.random() * 1.5 + 0.5);
}

