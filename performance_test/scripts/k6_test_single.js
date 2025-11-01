import http from 'k6/http';
import { check } from 'k6';

// Single test configuration for verification
export const options = {
  vus: 1,
  iterations: 1,
};

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'ory_at_b3qecKKOoLeW-3CTVARpvO8dJpkTaKvm_3F1jY2fsf4.3Lz12Ub06Pgn5sUtPT68ltDuLa5c6deJ_KqRzDjUM2Q';
const CSRF_TOKEN = __ENV.CSRF_TOKEN || 'csrf_token_be481debe9e1ebcf14d99f6f631d9a520ca6701ba0f3e4398508af30ebb1f509=9znDvETbFEvXcMNkmypOZFop5yibQ94nZCxAMmDVlj8=';

// Load test file
const testFile = open('../test_files/test_1mb.txt', 'b');

export default function () {
  console.log('Testing file upload to:', `${BASE_URL}/api/files/encrypt`);
  
  const fileName = 'test_1mb.txt';
  const boundary = '----WebKitFormBoundaryABC123';
  
  // Build the multipart/form-data payload
  const formData = [];
  
  // Add boundary and file headers
  formData.push(`--${boundary}\r\n`);
  formData.push(`Content-Disposition: form-data; name="file"; filename="${fileName}"\r\n`);
  formData.push(`Content-Type: text/plain\r\n\r\n`);
  
  // Add file content
  formData.push(testFile);
  
  // Add closing boundary
  formData.push(`\r\n--${boundary}--\r\n`);
  
  // Join all parts
  const body = formData.join('');

  const headers = {
    'Authorization': `Bearer ${AUTH_TOKEN}`,
    'Cookie': CSRF_TOKEN,
    'Content-Type': `multipart/form-data; boundary=${boundary}`,
  };

  console.log('Sending request with Content-Type:', headers['Content-Type']);

  const res = http.post(`${BASE_URL}/api/files/encrypt`, body, {
    headers: headers,
    timeout: '120s',
  });

  console.log('Response status:', res.status);
  console.log('Response body:', res.body);

  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'response has body': (r) => r.body && r.body.length > 0,
    'no server error': (r) => r.status < 500,
  });

  if (success) {
    console.log('✓ Test PASSED');
  } else {
    console.error('✗ Test FAILED');
  }
}
