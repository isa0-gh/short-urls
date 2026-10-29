import http from 'k6/http';
import { check } from 'k6';

export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
export const API_URL = `${BASE_URL}/api/v1`;

export function uniqueSuffix() {
    return `${Date.now()}-${__VU}-${__ITER}`;
}

export function registerAndLogin() {
    const suffix = uniqueSuffix();
    const user = {
        username: `k6-user-${suffix}`,
        email: `k6-${suffix}@example.com`,
        password: __ENV.TEST_PASSWORD || 'k6-password-123',
    };

    const registerResponse = http.post(
        `${API_URL}/register`,
        JSON.stringify(user),
        { headers: { 'Content-Type': 'application/json' } },
    );
    check(registerResponse, {
        'registration succeeds': (response) => response.status === 201,
    });

    const loginResponse = http.post(
        `${API_URL}/login`,
        JSON.stringify({ email: user.email, password: user.password }),
        { headers: { 'Content-Type': 'application/json' } },
    );
    check(loginResponse, {
        'login succeeds': (response) => response.status === 201,
    });

    return user;
}

export function createLink() {
    const shortCode = `k6-${uniqueSuffix()}`;
    const response = http.post(
        `${API_URL}/create`,
        JSON.stringify({
            short_code: shortCode,
            destination_url: __ENV.DESTINATION_URL || 'https://example.com/k6',
        }),
        { headers: { 'Content-Type': 'application/json' } },
    );
    check(response, {
        'link creation succeeds': (result) => result.status === 201,
    });
    return shortCode;
}

export function listLinks() {
    const response = http.get(`${API_URL}/links`);
    check(response, {
        'authenticated link list succeeds': (result) => result.status === 200,
    });
}

export function redirect(shortCode) {
    const response = http.get(`${BASE_URL}/s/${shortCode}`, {
        redirects: false,
    });
    check(response, {
        'redirect succeeds': (result) => result.status === 302,
    });
}
