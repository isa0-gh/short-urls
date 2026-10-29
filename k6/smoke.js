import { check } from 'k6';
import http from 'k6/http';
import {
    API_URL,
    registerAndLogin,
    createLink,
    listLinks,
    redirect,
} from './helpers.js';

export const options = {
    vus: 1,
    iterations: 1,
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<1000'],
    },
};

export default function () {
    registerAndLogin();
    const shortCode = createLink();
    listLinks();
    redirect(shortCode);

    const unauthenticatedResponse = http.get(`${API_URL}/links`);
    check(unauthenticatedResponse, {
        'session protects authenticated routes': (response) => response.status === 401,
    });
}
