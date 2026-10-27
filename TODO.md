# TODO

The short list below contains the next concrete improvements.

- [ ] Fix session-cache expiration so an expired token cannot remain valid in Redis.
- [ ] Add input validation and map expected failures to `400`, `404`, `409`, and `500` responses.
- [ ] Handle pgx errors correctly, including `pgx.ErrNoRows`.
- [ ] Add tests for authentication, token expiration, link ownership, duplicate values, and redirects.
- [ ] Add k6 tests for registration, login, link creation, redirects, and authenticated API load.
- [ ] Improve credential management: replace browser tokens in `localStorage` with secure `HttpOnly`, `SameSite` cookies, and add an API-key center for creating, listing, revoking, and rotating long-lived keys.
- [ ] Make `.env.example` use the same Redis configuration as the application: `REDIS_ADDR`.
- [ ] Add logout and account deletion, including session revocation.
- [ ] Decide whether to implement redirect telemetry or remove its unused code and database table.
