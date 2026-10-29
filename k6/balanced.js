import { sleep } from 'k6';
import {
    registerAndLogin,
    createLink,
    listLinks,
    redirect,
} from './helpers.js';

export const options = {
    noCookiesReset: true,
    scenarios: {
        balanced: {
            executor: 'ramping-vus',
            startVUs: 1,
            stages: [
                { duration: '15s', target: Number(__ENV.VUS || 10) },
                { duration: __ENV.DURATION || '1m', target: Number(__ENV.VUS || 10) },
                { duration: '15s', target: 0 },
            ],
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<1000'],
    },
};

let shortCode;

export default function () {
    if (!shortCode) {
        registerAndLogin();
        shortCode = createLink();
    }

    createLink();
    listLinks();
    redirect(shortCode);
    sleep(0.5);
}
