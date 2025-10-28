import http from 'k6/http';
import { check, sleep } from 'k6';

// Smoke Test Configuration
// Purpose: Verify the system works with minimal load
export const options = {
  stages: [
    { duration: '1m', target: 1 }, // 1 user for 1 minute
  ],
  thresholds: {
    http_req_duration: ['p(95)<5000'], // 95% of requests should complete within 5s
    http_req_failed: ['rate<0.01'],    // Less than 1% of requests should fail
  },
};

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'ory_at_b3qecKKOoLeW-3CTVARpvO8dJpkTaKvm_3F1jY2fsf4.3Lz12Ub06Pgn5sUtPT68ltDuLa5c6deJ_KqRzDjUM2Q';
const CSRF_TOKEN = __ENV.CSRF_TOKEN || 'csrf_token_be481debe9e1ebcf14d99f6f631d9a520ca6701ba0f3e4398508af30ebb1f509=9znDvETbFEvXcMNkmypOZFop5yibQ94nZCxAMmDVlj8=';

// Load test file from the test_files directory
const testFileContent = open('../test_files/test_5mb.txt', 'b');
const testFileName = 'test_5mb.txt';

export default function () {
  // Build multipart/form-data payload manually
  const boundary = '----WebKitFormBoundaryABC123';
  const formData = [];
  
  formData.push(`--${boundary}\r\n`);
  formData.push(`Content-Disposition: form-data; name="file"; filename="${testFileName}"\r\n`);
  formData.push(`Content-Type: text/plain\r\n\r\n`);
  formData.push(testFileContent);
  formData.push(`\r\n--${boundary}--\r\n`);
  
  const body = formData.join('');

  const headers = {
    'Authorization': `Bearer ${AUTH_TOKEN}`,
    'Cookie': CSRF_TOKEN,
    'Content-Type': `multipart/form-data; boundary=${boundary}`,
  };

  console.log(`Smoke Test - Sending ${testFileName} (${(testFileContent.length / 1024 / 1024).toFixed(2)} MB)`);

  const res = http.post(`${BASE_URL}/api/files/encrypt`, body, {
    headers: headers,
    timeout: '60s',
  });

  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'response has body': (r) => r.body.length > 0,
    'response time < 5s': (r) => r.timings.duration < 5000,
  });

  if (!success) {
    console.error(`Request failed: Status ${res.status}, Body: ${res.body.substring(0, 200)}`);
  } else {
    console.log(`✓ Request successful - Duration: ${res.timings.duration.toFixed(2)}ms`);
  }

  sleep(1);
}

export function handleSummary(data) {
  return {
    'stdout': JSON.stringify(data, null, 2),
    'results/k6_smoke_test_results.json': JSON.stringify(data),
  };
}
