import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom Metrics
const errorRate = new Rate('errors');
const encryptionDuration = new Trend('encryption_duration');
const serverErrors = new Counter('server_errors');
const clientErrors = new Counter('client_errors');

// Stress Test Configuration
// Purpose: Find the breaking point of the system
export const options = {
  stages: [
    { duration: '2m', target: 10 },   // Warm up to 10 users
    { duration: '3m', target: 20 },   // Ramp up to 20 users
    { duration: '3m', target: 40 },   // Ramp up to 40 users
    { duration: '3m', target: 60 },   // Ramp up to 60 users
    { duration: '3m', target: 80 },   // Ramp up to 80 users
    { duration: '3m', target: 100 },  // Ramp up to 100 users - Peak load
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '3m', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<15000'], // 95% of requests within 15s
    http_req_failed: ['rate<0.1'],      // Less than 10% failure rate
    errors: ['rate<0.1'],
    encryption_duration: ['p(99)<20000'], // 99% within 20s
  },
};

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'ory_at_b3qecKKOoLeW-3CTVARpvO8dJpkTaKvm_3F1jY2fsf4.3Lz12Ub06Pgn5sUtPT68ltDuLa5c6deJ_KqRzDjUM2Q';
const CSRF_TOKEN = __ENV.CSRF_TOKEN || 'csrf_token_be481debe9e1ebcf14d99f6f631d9a520ca6701ba0f3e4398508af30ebb1f509=9znDvETbFEvXcMNkmypOZFop5yibQ94nZCxAMmDVlj8=';

// Load test files from the test_files directory
// Using different sizes for stress testing
const testFiles = [
  { content: open('../test_files/test_1mb.txt', 'b'), name: 'test_1mb.txt' },
  { content: open('../test_files/test_3mb.txt', 'b'), name: 'test_3mb.txt' },
  { content: open('../test_files/test_5mb.txt', 'b'), name: 'test_5mb.txt' },
  { content: open('../test_files/test_5mb.txt', 'b'), name: 'test_5mb_copy.txt' },
  { content: open('../test_files/test_3mb.txt', 'b'), name: 'test_3mb_copy.txt' },
];

export default function () {
  // Randomly select from pre-generated files
  const selectedFile = testFiles[Math.floor(Math.random() * testFiles.length)];

  // Build multipart/form-data payload manually
  const boundary = '----WebKitFormBoundaryABC123';
  const formData = [];
  
  formData.push(`--${boundary}\r\n`);
  formData.push(`Content-Disposition: form-data; name="file"; filename="${selectedFile.name}"\r\n`);
  formData.push(`Content-Type: text/plain\r\n\r\n`);
  formData.push(selectedFile.content);
  formData.push(`\r\n--${boundary}--\r\n`);
  
  const body = formData.join('');

  const headers = {
    'Authorization': `Bearer ${AUTH_TOKEN}`,
    'Cookie': CSRF_TOKEN,
    'Content-Type': `multipart/form-data; boundary=${boundary}`,
  };

  const startTime = Date.now();
  const res = http.post(`${BASE_URL}/api/files/encrypt`, body, {
    headers: headers,
    timeout: '120s',
  });
  const duration = Date.now() - startTime;

  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'status is not 5xx': (r) => r.status < 500,
    'response has body': (r) => r.body && r.body.length > 0,
    'response time < 15s': (r) => r.timings.duration < 15000,
  });

  // Record metrics
  errorRate.add(!success);
  encryptionDuration.add(duration);
  
  if (res.status >= 500) {
    serverErrors.add(1);
    console.error(`[VU ${__VU}] Server Error: Status ${res.status} - ${res.body.substring(0, 100)}`);
  } else if (res.status >= 400) {
    clientErrors.add(1);
    console.error(`[VU ${__VU}] Client Error: Status ${res.status}`);
  }

  if (!success && res.status === 200) {
    console.error(`[VU ${__VU}] Validation failed for ${selectedFile.name}`);
  }

  // Minimal sleep to increase stress
  sleep(Math.random() * 0.5 + 0.1); // Sleep between 0.1-0.6 seconds
}

export function handleSummary(data) {
  console.log('\n=== Stress Test Summary ===');
  console.log(`Total requests: ${data.metrics.http_reqs.values.count}`);
  console.log(`Failed requests: ${(data.metrics.http_req_failed.values.rate * 100).toFixed(2)}%`);
  console.log(`Error rate: ${(data.metrics.errors.values.rate * 100).toFixed(2)}%`);
  console.log(`Server errors (5xx): ${data.metrics.server_errors ? data.metrics.server_errors.values.count : 0}`);
  console.log(`Client errors (4xx): ${data.metrics.client_errors ? data.metrics.client_errors.values.count : 0}`);
  console.log(`\nResponse Times:`);
  console.log(`  Avg: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms`);
  console.log(`  Min: ${data.metrics.http_req_duration.values.min.toFixed(2)}ms`);
  console.log(`  Max: ${data.metrics.http_req_duration.values.max.toFixed(2)}ms`);
  console.log(`  P50: ${data.metrics.http_req_duration.values.med.toFixed(2)}ms`);
  if (data.metrics.http_req_duration.values['p(95)']) {
    console.log(`  P95: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms`);
  }
  if (data.metrics.http_req_duration.values['p(99)']) {
    console.log(`  P99: ${data.metrics.http_req_duration.values['p(99)'].toFixed(2)}ms`);
  }
  
  return {
    'stdout': JSON.stringify(data, null, 2),
    'results/k6_stress_test_results.json': JSON.stringify(data),
  };
}
