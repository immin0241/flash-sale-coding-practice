import http from 'k6/http';
import encoding from 'k6/encoding';
import { check } from 'k6';
import { Counter, Rate } from 'k6/metrics';

const baseURL = __ENV.BASE_URL || 'http://localhost:8000';
const productID = Number(__ENV.PRODUCT_ID || 1);
const virtualUsers = Number(__ENV.VUS || 100);
const iterationsPerUser = Number(__ENV.ITERATIONS || 10);
const password = __ENV.PASSWORD || 'password';

const successfulSales = new Counter('successful_sales');
const soldOutResponses = new Counter('sold_out_responses');
const duplicateResponses = new Counter('duplicate_responses');
const unexpectedResponses = new Counter('unexpected_responses');
const requestTimeouts = new Counter('request_timeouts');
const validResponseRate = new Rate('valid_response_rate');

export const options = {
    scenarios: {
        flash_sale: {
            executor: 'per-vu-iterations',
            vus: virtualUsers,
            iterations: iterationsPerUser,
            maxDuration: __ENV.MAX_DURATION || '2m',
        },
    },
    thresholds: {
        successful_sales: [`count==${Number(__ENV.EXPECTED_SALES || 100)}`],
        unexpected_responses: ['count==0'],
        request_timeouts: ['count==0'],
        valid_response_rate: ['rate==1'],
    },
};

export default function () {
    const username = __ENV.USERNAME || `user${__VU}`;
    const response = http.post(
        `${baseURL}/api/transactions/buy`,
        JSON.stringify({ product_id: productID, amount: 1 }),
        {
            headers: {
                Authorization: `Basic ${encoding.b64encode(`${username}:${password}`)}`,
                'Content-Type': 'application/json',
            },
            tags: { endpoint: 'buy' },
            timeout: __ENV.REQUEST_TIMEOUT || '10s',
        },
    );

    const isSuccess = response.status === 200;
    const isSoldOut = response.status === 410;
    const isDuplicate = response.status === 400;
    const isTimeout = response.status === 0;

    if (isSuccess) {
        successfulSales.add(1);
    } else if (isSoldOut) {
        soldOutResponses.add(1);
    } else if (isDuplicate) {
        duplicateResponses.add(1);
    } else {
        unexpectedResponses.add(1, { status: String(response.status) });
    }

    if (isTimeout) {
        requestTimeouts.add(1);
    }

    validResponseRate.add(isSuccess || isSoldOut || isDuplicate);
    check(response, {
        'response is success, duplicate, or explicit sold-out error': () => isSuccess || isDuplicate || isSoldOut,
        'response is not a timeout': () => !isTimeout,
        'response is not a server error': () => response.status < 500,
    });
}

export function handleSummary(data) {
    const duration = data.metrics.http_req_duration?.values || {};
    const totalRequests = virtualUsers * iterationsPerUser;
    const summary = {
        scenario: 'flash sale: 10 requests per user',
        virtual_users: virtualUsers,
        iterations_per_user: iterationsPerUser,
        total_requests: totalRequests,
        expected_sales: Number(__ENV.EXPECTED_SALES || 100),
        successful_sales: data.metrics.successful_sales?.values.count || 0,
        sold_out_responses: data.metrics.sold_out_responses?.values.count || 0,
        duplicate_responses: data.metrics.duplicate_responses?.values.count || 0,
        unexpected_responses: data.metrics.unexpected_responses?.values.count || 0,
        request_timeouts: data.metrics.request_timeouts?.values.count || 0,
        http_req_duration_ms: {
            p95: duration['p(95)'] || null,
            p99: duration['p(99)'] || null,
        },
    };

    return {
        stdout: `${JSON.stringify(summary, null, 2)}\n`,
        'load-test/results/flash-sale-summary.json': JSON.stringify(summary, null, 2),
    };
}
