---
name: MSUBot OpenAPI Spec
overview: Design an OpenAPI 3.0 specification for a modernized MSUBot Go backend hosted on Render, using PostgreSQL as the database and custom JWT auth (replacing Firebase/Firestore).
todos:
  - id: openapi-spec
    content: Write api/openapi.yaml — full OpenAPI 3.0 spec with all schemas and the 17 endpoints above
    status: completed
isProject: false
---

# MSUBot OpenAPI Specification

## Output File

`api/openapi.yaml`

## Architecture Assumptions

- **Hosting**: Render (Go service + managed Postgres)
- **Auth**: Email/password → server-issued JWT (`Authorization: Bearer <token>`)
- **Admin/cron protection**: `X-Admin-Key` shared secret header
- **Data source for `GET /courses`**: Postgres index (populated by the nightly update-index job) — seat counts reflect last index run, not real-time
- **Data source for `GET /courses/{crn}`**: Live scrape of MSU Atlas — always returns current occupancy
- **Term format**: 6-digit string (e.g., `202530` = Spring 2025, `202550` = Summer, `202570` = Fall)

---

## Schemas

### `Section`

Derived from [old-service/server/models/models.go](old-service/server/models/models.go):

- `crn` string
- `term` string
- `deptAbbr` string
- `deptName` string
- `courseNumber` string
- `courseName` string
- `courseType` string
- `sectionNumber` string
- `credits` string
- `instructor` string
- `time` string
- `location` string
- `totalSeats` integer
- `takenSeats` integer
- `availableSeats` integer
- `seatsLive` boolean — `true` when the seat counts came from a live Atlas scrape, `false` when from the Postgres index

### `User`

- `id` string (UUID)
- `email` string (login email)
- `notificationPreferences` `NotificationPreferences`
- `createdAt` string (ISO 8601)

### `NotificationPreferences`

Designed to be extended with future channels (SMS, push, etc.) without a schema change:

- `email` string — address alerts are sent to (may differ from login email)
- `emailEnabled` boolean

### `CourseSuggestion`

Lightweight object for autocomplete dropdowns:

- `deptAbbr` string
- `deptName` string
- `courseNumber` string
- `courseName` string

### `OccupancyResult`

Live seat data keyed by CRN, returned by the batch occupancy endpoint:

- `crn` string
- `totalSeats` integer
- `takenSeats` integer
- `availableSeats` integer
- `seatsLive` boolean (always `true` from this endpoint)

### `TrackedSection`

- `crn` string
- `term` string
- `section` `Section` (denormalized metadata from Postgres; seat counts reflect last cron check)
- `addedAt` string (ISO 8601)
- `lastCheckedAt` string (ISO 8601) — timestamp of the last cron-triggered live scrape for this CRN

### `AuthResponse`

- `token` string (JWT)
- `expiresIn` integer (seconds)
- `user` `User`

### `Stats`

- `term` string — the term these stats are scoped to (or `"all"` if no term filter applied)
- `totalSections` integer — sections in the Postgres index for the term
- `totalDepartments` integer — departments in the index
- `totalUsers` integer — registered users
- `totalTrackedSections` integer — unique CRNs currently being watched across all users
- `alertsSentThisTerm` integer — total open-seat notifications sent for the term (requires a notifications log table)
- `seatAvailability` object:
  - `openSections` integer — sections with `availableSeats > 0`
  - `fullSections` integer — sections with `availableSeats == 0`
  - `percentOpen` number — `openSections / totalSections * 100`
- `mostTrackedCourses` `MostTrackedCourse[]` — top N courses by watcher count
- `indexLastUpdatedAt` string (ISO 8601)

### `MostTrackedCourse`

- `deptAbbr` string
- `courseNumber` string
- `courseName` string
- `trackedCount` integer — number of users watching any section of this course

### `ErrorResponse`

- `code` string
- `message` string

---

## Endpoints

### Public

```
GET /stats
```

Returns aggregate analytics about the service. Scoped to a term if provided; otherwise uses the current/upcoming term.

- **Query params**:
  - `term` string — e.g., `202570`; defaults to current term
  - `limit` integer — number of entries in `mostTrackedCourses`, default `5`
- **Response** `200`: `Stats`
- **Notes**: `alertsSentThisTerm` requires a `notifications_log` table (records each alert sent with timestamp and term). Most contested courses (watcher count vs. available seats) is a future addition requiring periodic seat-count snapshots.

```
GET /courses
```

Query sections from the Postgres index. Seat counts reflect the last nightly index run; `seatsLive` will be `false`. Accepts either structured filters or a free-text query via `q`.

- **Query params**:
  - `term` string — e.g., `202570`. Required unless `q` is provided; defaults to current/upcoming terms when omitted with `q`.
  - `q` string — unified free-text search across `courseName`, `deptName`, `deptAbbr`, `courseNumber`, and `instructor`
  - `dept` string — department abbreviation (e.g., `CS`)
  - `courseNumber` string — e.g., `101`
  - `instructor` string — partial match
  - `availableOnly` boolean — filter to sections where `availableSeats > 0` (based on last index)
  - `courseType` string — `Lecture`, `Lab`, `Recitation`, etc.
- **Response** `200`: `Section[]`

```
GET /courses/suggestions
```

Autocomplete endpoint for the search bar. Returns lightweight `CourseSuggestion` objects (no seat data), suitable for dropdown rendering. Always queries from the Postgres index.

- **Query params**:
  - `q` string **required** — partial match against `courseName`, `deptAbbr`, `courseNumber`
  - `term` string — restrict suggestions to a term; omit to search all indexed terms
  - `limit` integer — max results, default `10`
- **Response** `200`: `CourseSuggestion[]`

```
GET /courses/{crn}
```

Live-scrapes MSU Atlas for a single section by CRN. The backend looks up the section's `term`, `dept`, and `courseNumber` from Postgres, then makes a live Atlas request to return current seat occupancy. `seatsLive` will be `true`.

- **Path param**: `crn` string
- **Response** `200`: `Section`
- **Response** `404`: section not found in index
- **Response** `502`: Atlas upstream unavailable

```
POST /courses/occupancy
```

Batch live-scrape for a list of CRNs. Intended for the table view — fetches current seat counts for all visible rows in a single request. The backend groups CRNs by dept+course+term (resolved from Postgres) and fans out Atlas requests concurrently, then returns results keyed by CRN.

- **Body**: `{ crns: string[] }`
- **Response** `200`: `OccupancyResult[]`
- **Response** `207`: partial success — some CRNs could not be fetched (each result includes an optional `error` field)

---

### Auth

```
POST /auth/register
```

- **Body**: `{ email, password, name? }`
- **Response** `201`: `AuthResponse`

```
POST /auth/login
```

- **Body**: `{ email, password }`
- **Response** `200`: `AuthResponse`

```
DELETE /auth/session
```

Invalidates the current JWT (requires the backend to maintain a token blocklist in Postgres).

- **Auth**: Bearer JWT
- **Response** `204`

---

### Authed User (Bearer JWT)

```
GET /users/me
```

- **Response** `200`: `User`

```
PATCH /users/me
```

Update user profile and notification preferences.

- **Body**: `{ notificationPreferences?: NotificationPreferences }`
- **Response** `200`: `User`

```
DELETE /users/me
```

Permanently deletes the account and all tracked sections.

- **Response** `204`

```
GET /users/me/sections
```

- **Response** `200`: `TrackedSection[]`

```
POST /users/me/sections
```

- **Body**: `{ crn, term }`
- **Response** `201`: `TrackedSection`

```
DELETE /users/me/sections/{crn}
```

- **Path param**: `crn`
- **Response** `204`

---

### Admin / Cron (`X-Admin-Key` header)

```
POST /admin/scrape-sections
```

Triggers a check of all tracked sections against MSU Atlas and sends notifications for any that have open seats.

- **Response** `202`

```
POST /admin/update-index
```

Triggers the nightly scrape of all MSU departments and course sections, writing results to Postgres.

- **Response** `202`

---

## Security Schemes

```yaml
securitySchemes:
  BearerAuth:
    type: http
    scheme: bearer
    bearerFormat: JWT
  AdminKey:
    type: apiKey
    in: header
    name: X-Admin-Key
```

---

## Endpoint Summary

```
[public]
GET  /stats
GET  /courses
GET  /courses/suggestions
GET  /courses/{crn}
POST /courses/occupancy

[auth]
POST   /auth/register
POST   /auth/login
DELETE /auth/session

[authed user — BearerAuth]
GET    /users/me
PATCH  /users/me
DELETE /users/me
GET    /users/me/sections
POST   /users/me/sections
DELETE /users/me/sections/{crn}

[admin — AdminKey]
POST /admin/scrape-sections
POST /admin/update-index
```

