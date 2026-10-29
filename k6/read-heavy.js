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
        read_heavy: {
            executor: 'constant-vus',
            vus: Number(__ENV.VUS || 10),
            duration: __ENV.DURATION || '1m',
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<750'],
    },
};

let shortCode;

export default function () {
    if (!shortCode) {
        registerAndLogin();
        shortCode = createLink();
    }

    listLinks();
    listLinks();
    redirect(shortCode);
    redirect(shortCode);
    redirect(shortCode);
    sleep(0.2);
}
