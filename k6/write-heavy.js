import { sleep } from 'k6';
import { registerAndLogin, createLink, listLinks } from './helpers.js';

export const options = {
    scenarios: {
        write_heavy: {
            executor: 'ramping-vus',
            startVUs: 1,
            stages: [
                { duration: '15s', target: Number(__ENV.VUS || 10) },
                { duration: __ENV.DURATION || '45s', target: Number(__ENV.VUS || 10) },
                { duration: '15s', target: 0 },
            ],
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.02'],
        http_req_duration: ['p(95)<1000'],
    },
};

export default function () {
    registerAndLogin();
    createLink();
    createLink();
    createLink();
    listLinks();
    sleep(0.2);
}
