import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom Metrics
const errorRate = new Rate('errors');
const recoveryRate = new Rate('recovery');
const encryptionDuration = new Trend('encryption_duration');

// Spike Test Configuration
// Purpose: Test system behavior under sudden traffic spikes
export const options = {
  stages: [
    { duration: '1m', target: 5 },    // Normal load
    { duration: '10s', target: 100 }, // Sudden spike!
    { duration: '2m', target: 100 },  // Stay at spike
    { duration: '1m', target: 5 },    // Recovery
    { duration: '1m', target: 5 },    // Verify recovery
    { duration: '10s', target: 150 }, // Second spike even higher!
    { duration: '2m', target: 150 },  // Stay at higher spike
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<20000'], // More lenient during spikes
    http_req_failed: ['rate<0.15'],     // Allow higher failure during spike
    errors: ['rate<0.15'],
  },
};

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'ory_at_b3qecKKOoLeW-3CTVARpvO8dJpkTaKvm_3F1jY2fsf4.3Lz12Ub06Pgn5sUtPT68ltDuLa5c6deJ_KqRzDjUM2Q';
const CSRF_TOKEN = __ENV.CSRF_TOKEN || 'csrf_token_be481debe9e1ebcf14d99f6f631d9a520ca6701ba0f3e4398508af30ebb1f509=9znDvETbFEvXcMNkmypOZFop5yibQ94nZCxAMmDVlj8=';

// Load test file from the test_files directory
const testFile = {
  content: open('../test_files/test_5mb.txt', 'b'),
  name: 'test_5mb.txt'
};

export default function () {
  // Build multipart/form-data payload manually
  const boundary = '----WebKitFormBoundaryABC123';
  const formData = [];
  
  formData.push(`--${boundary}\r\n`);
  formData.push(`Content-Disposition: form-data; name="file"; filename="${testFile.name}"\r\n`);
  formData.push(`Content-Type: text/plain\r\n\r\n`);
  formData.push(testFile.content);
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
    'response has body': (r) => r.body && r.body.length > 0,
    'response time < 20s': (r) => r.timings.duration < 20000,
    'no server error': (r) => r.status < 500,
  });

  // Record metrics
  errorRate.add(!success);
  encryptionDuration.add(duration);
  
  // Check if system is recovering (low VU count with success)
  if (__VU <= 10 && success) {
    recoveryRate.add(1);
  } else if (__VU <= 10) {
    recoveryRate.add(0);
  }

  if (!success) {
    console.error(`[VU ${__VU}] Spike Test - Request failed: Status ${res.status}, Duration: ${duration}ms`);
  }

  // Very short sleep during spike
  sleep(0.1);
}

export function handleSummary(data) {
  console.log('\n=== Spike Test Summary ===');
  console.log(`Total requests: ${data.metrics.http_reqs.values.count}`);
  console.log(`Failed requests: ${(data.metrics.http_req_failed.values.rate * 100).toFixed(2)}%`);
  console.log(`Error rate: ${(data.metrics.errors.values.rate * 100).toFixed(2)}%`);
  
  if (data.metrics.recovery) {
    console.log(`Recovery rate: ${(data.metrics.recovery.values.rate * 100).toFixed(2)}%`);
  }
  
  console.log(`\nResponse Times:`);
  console.log(`  Avg: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms`);
  if (data.metrics.http_req_duration.values['p(95)']) {
    console.log(`  P95: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms`);
  }
  if (data.metrics.http_req_duration.values['p(99)']) {
    console.log(`  P99: ${data.metrics.http_req_duration.values['p(99)'].toFixed(2)}ms`);
  }
  console.log(`  Max: ${data.metrics.http_req_duration.values.max.toFixed(2)}ms`);
  
  return {
    'stdout': JSON.stringify(data, null, 2),
    'results/k6_spike_test_results.json': JSON.stringify(data),
  };
}
