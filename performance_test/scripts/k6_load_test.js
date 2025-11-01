import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom Metrics
const errorRate = new Rate('errors');
const encryptionDuration = new Trend('encryption_duration');

// Load Test Configuration
// Purpose: Test system performance under expected load
export const options = {
  stages: [
    { duration: '2m', target: 10 },  // Ramp up to 10 users over 2 minutes
    { duration: '5m', target: 10 },  // Stay at 10 users for 5 minutes
    { duration: '2m', target: 20 },  // Ramp up to 20 users
    { duration: '5m', target: 20 },  // Stay at 20 users for 5 minutes
    { duration: '2m', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<10000'], // 95% of requests should complete within 10s
    http_req_failed: ['rate<0.05'],     // Less than 5% of requests should fail
    errors: ['rate<0.05'],
    encryption_duration: ['p(95)<10000'],
  },
};

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'ory_at_RlIjxvgPt_PKyjNRttMK2N5p4IRCLmHUosPAALsTxmY.JMwklzikGc_vosv5wH0CD9B-rjklKrlzpaBg9e94ruc';
const CSRF_TOKEN = __ENV.CSRF_TOKEN || 'csrf_token_be481debe9e1ebcf14d99f6f631d9a520ca6701ba0f3e4398508af30ebb1f509=9znDvETbFEvXcMNkmypOZFop5yibQ94nZCxAMmDVlj8=';

// Load test files from the test_files directory
// Make sure to run this from the performance_test directory
const testFiles = {
  small: { path: open('../test_files/test_1mb.txt', 'b'), name: 'test_1mb.txt' },
  medium: { path: open('../test_files/test_3mb.txt', 'b'), name: 'test_3mb.txt' },
  large: { path: open('../test_files/test_5mb.txt', 'b'), name: 'test_5mb.txt' },
};

export default function () {
  // Randomly select file size to simulate real-world usage
  const fileTypes = ['small', 'medium', 'large'];
  const selectedType = fileTypes[Math.floor(Math.random() * fileTypes.length)];
  const selectedFile = testFiles[selectedType];

  // Build multipart/form-data payload manually
  const boundary = '----WebKitFormBoundaryABC123';
  const formData = [];
  
  formData.push(`--${boundary}\r\n`);
  formData.push(`Content-Disposition: form-data; name="file"; filename="${selectedFile.name}"\r\n`);
  formData.push(`Content-Type: text/plain\r\n\r\n`);
  formData.push(selectedFile.path);
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
    'response has body': (r) => r.body.length > 0,
    'response time < 10s': (r) => r.timings.duration < 10000,
    'no server error': (r) => r.status < 500,
  });

  // Record custom metrics
  errorRate.add(!success);
  encryptionDuration.add(duration);

  if (!success) {
    console.error(`[VU ${__VU}] Request failed: ${selectedFile.name} - Status ${res.status} - Body: ${res.body}`);
  }

  // Variable sleep time to simulate real user behavior
  sleep(Math.random() * 2 + 1); // Sleep between 1-3 seconds
}

export function handleSummary(data) {
  // Ensure the results directory exists or warn the user
  const resultsDir = 'results';
  try {
    // This will throw if the directory does not exist
    open(`${resultsDir}/.keep`, 'r');
  } catch (e) {
    console.warn(`WARNING: The '${resultsDir}' directory does not exist. Please create it before running the test to save summary results.`);
  }

  // console.log('\n=== Load Test Summary ===');
  // console.log(`Total requests: ${data.metrics.http_reqs.values.count}`);
  // console.log(`Failed requests: ${data.metrics.http_req_failed.values.rate * 100}%`);
  // console.log(`Avg response time: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms`);
  // console.log(`P95 response time: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms`);
  
  return {
    'stdout': JSON.stringify(data, null, 2),
    'results/k6_load_test_results.json': JSON.stringify(data),
  };
}
