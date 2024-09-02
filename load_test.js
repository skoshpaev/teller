import http from 'k6/http';
import { check, sleep } from 'k6';

const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.H1KnRyz0-_3OVZJJH-AYAwIMqR-9n5Uz9r97omtnGTc';
const JWT_SECRET = 'H1g$eCr3t!2S#cUr3T@256-bSecr3tIt';

export let options = {
    scenarios: {
        publish_scenario: {
            executor: 'constant-arrival-rate',
            rate: 1000, // Number of requests per second
            timeUnit: '1s', // Time unit
            duration: '1m', // Test duration
            preAllocatedVUs: 1000, // Number of pre-allocated VUs
            maxVUs: 1000, // Maximum number of VUs
            exec: 'publishScenario',
        },
        subscribe_scenario: {
            executor: 'constant-arrival-rate',
            rate: 1000, // Number of requests per second
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 1000,
            maxVUs: 1000,
            exec: 'subscribeScenario',
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.01'], // Less than 1% of requests should fail
    },
};


export function publishScenario() {
    let pubRes = http.post('http://localhost:8080/publish', JSON.stringify({
        channel: 'test-channel',
        message: { key: 'value' },
    }), {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${TOKEN}`,
        },
        timeout: '120s',
    });   
    
    check(pubRes, {
        'Publication status is 200': (r) => r.status === 200,
    });
}

export function subscribeScenario() {
    let res = http.get('http://localhost:8080/subscribe?channel=test-channel&test=true', {
        headers: { 'Authorization': `Bearer ${TOKEN}` },
        timeout: '120s',
    });
    
    check(res, {
        'Connected and test data received': (r) => r.status === 200 && r.json().test === true
    });
}