# TODO

The short list below contains the next concrete improvements.

- [X] Fix session-cache expiration so an expired token cannot remain valid in Redis.
- [ ] Add input validation and map expected failures to `400`, `404`, `409`, and `500` responses.
- [ ] Handle pgx errors correctly, including `pgx.ErrNoRows`.
- [ ] Add tests for authentication, token expiration, link ownership, duplicate values, and redirects.
- [ ] Add k6 tests for registration, login, link creation, redirects, and authenticated API load.
- [x] Replace browser tokens in `localStorage` with secure `HttpOnly`, `SameSite=Lax` cookies.
- [ ] Add an API-key center for creating, listing, revoking, and rotating long-lived keys.
- [x] Make `.env.example` use the same Redis configuration as the application: `REDIS_ADDR`.
- [x] Add logout, including current-session revocation.
- [ ] Add account deletion, including revocation of all sessions.
- [ ] Decide whether to implement redirect telemetry or remove its unused code and database table.
