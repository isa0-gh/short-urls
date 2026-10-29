# k6 Read-Heavy Load Test Report

> **Status:** Passed
>
> **Command:** `k6 run k6/read-heavy.js`
>
> **Execution:** Local

## Executive Summary

The read-heavy workload completed **2,640 iterations** across **10 VUs** with no interrupted iterations. Every functional check passed and no HTTP requests failed.

## Workload

| Setting | Value |
| --- | ---: |
| Scenario | `read_heavy` |
| VUs | 10 |
| Active duration | 1m |
| Maximum duration | 1m30s, including graceful stop |
| Completed iterations | 2,640 |
| Interrupted iterations | 0 |

## Thresholds

| Metric | Threshold | Result | Status |
| --- | --- | ---: | :---: |
| `http_req_duration` | `p(95)<750ms` | `p(95)=6.86ms` | **Passed** |
| `http_req_failed` | `rate<0.01` | `0.00%` | **Passed** |

## Functional Checks

All checks passed: **13,230 of 13,230**.

- Registration succeeds
- Login succeeds
- Link creation succeeds
- Authenticated link listing succeeds
- Redirect succeeds

## HTTP Metrics

| Metric | Value |
| --- | ---: |
| Requests | 13,230 |
| Request rate | 219.83 requests/s |
| Failed requests | 0.00% |
| Average duration | 4.85ms |
| Minimum duration | 1.31ms |
| Median duration | 2.47ms |
| Maximum duration | 1.21s |
| p(90) duration | 5.07ms |
| p(95) duration | 6.86ms |

## Execution Metrics

| Metric | Value |
| --- | ---: |
| Average iteration duration | 227.71ms |
| Minimum iteration duration | 210.2ms |
| Median iteration duration | 215.48ms |
| Maximum iteration duration | 2.61s |
| p(90) iteration duration | 228.49ms |
| p(95) iteration duration | 236.54ms |
| VUs | 10 |

## Network Metrics

| Direction | Volume | Rate |
| --- | ---: | ---: |
| Received | 3.0 MB | 50 kB/s |
| Sent | 1.9 MB | 31 kB/s |

## Assessment

The read-heavy workload passed cleanly. The API sustained approximately **220 requests per second** with a **6.86ms p95 latency** and zero failed requests.

---

# k6 Balanced Load Test Report

> **Status:** Passed
>
> **Command:** `k6 run k6/balanced.js`
>
> **Execution:** Local

## Executive Summary

The balanced workload completed **1,439 iterations** across up to **10 VUs** with no interrupted iterations. Every functional check passed, and no HTTP requests failed.

## Workload

| Setting | Value |
| --- | ---: |
| Scenario | `balanced` |
| Maximum VUs | 10 |
| Active duration | 1m30s |
| Maximum duration | 2m, including graceful stop |
| Completed iterations | 1,439 |
| Interrupted iterations | 0 |

## Thresholds

| Metric | Threshold | Result | Status |
| --- | --- | ---: | :---: |
| `http_req_duration` | `p(95)<1000ms` | `p(95)=14.57ms` | **Passed** |
| `http_req_failed` | `rate<0.01` | `0.00%` | **Passed** |

## Functional Checks

All checks passed: **4,347 of 4,347**.

- Registration succeeds
- Login succeeds
- Link creation succeeds
- Authenticated link listing succeeds
- Redirect succeeds

## HTTP Metrics

| Metric | Value |
| --- | ---: |
| Requests | 4,347 |
| Request rate | 48.16 requests/s |
| Failed requests | 0.00% |
| Average duration | 8.37ms |
| Minimum duration | 1.4ms |
| Median duration | 6.32ms |
| Maximum duration | 352.54ms |
| p(90) duration | 11.98ms |
| p(95) duration | 14.57ms |

## Execution Metrics

| Metric | Value |
| --- | ---: |
| Average iteration duration | 527.88ms |
| Minimum iteration duration | 513.18ms |
| Median iteration duration | 520.15ms |
| Maximum iteration duration | 1.23s |
| p(90) iteration duration | 532.05ms |
| p(95) iteration duration | 543.69ms |
| VUs | 1-10 |

## Network Metrics

| Direction | Volume | Rate |
| --- | ---: | ---: |
| Received | 11 MB | 118 kB/s |
| Sent | 795 kB | 8.8 kB/s |

## Assessment

The balanced workload passed cleanly. The API handled approximately **48 requests per second** with a **14.57ms p95 latency** and zero failed requests.
