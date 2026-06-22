# Test Server

Mock Litmos REST API for `baton-litmos`, used by CI when no real-tenant credentials are available. Replicates the upstream API's auth flow, endpoints, pagination contract, and error shapes.

## Auth

| Real API | Test server |
|---|---|
| Custom `apikey` header on every request | Same — validates `apikey` header |
| API key from customer's Litmos account | Hardcoded: `test-api-key` |
| `source` query param required | Accepted but not validated |

## Endpoints

| Path | Method | Description |
|---|---|---|
| `/v1.svc/users` | GET | List users (paginated) |
| `/v1.svc/teams` | GET | List teams (paginated) |
| `/v1.svc/teams/{id}/users` | GET | List team members (paginated) |
| `/v1.svc/courses` | GET | List courses (paginated) |
| `/v1.svc/courses/{id}` | GET | Get course by ID |
| `/v1.svc/courses/{id}/users` | GET | List course enrollments (paginated) |
| `/v1.svc/courses/{id}/modules` | GET | List course modules (paginated) |
| `/v1.svc/users/{id}/courses` | POST | Assign course to user |
| `/v1.svc/users/{id}/courses/{courseId}` | DELETE | Remove course from user |

All responses are `application/xml`. Pagination uses `?start=N&limit=M` (default limit 500).

## Seed data

- **5 users**: alice (admin, active), bob (active), carol (disabled — exercises `STATUS_DISABLED`), dave (no enrollments — exercises empty-grants path), eve (active)
- **3 teams**: Engineering (alice+bob), Product (bob+carol), Operations (eve) — overlapping membership
- **3 courses**: Introduction to Security (active), Compliance Training (active), Legacy Systems Overview (inactive)
- **3 modules per course**
- **Course enrollments**: alice completed intro+compliance; bob in-progress intro+deprecated; carol in-progress intro; dave has none; eve completed compliance

## Running locally

```bash
# Start the test server (from project root)
go run ./cmd/test-server/

# In a separate terminal, point the connector at it
export LITMOS_BASE_URL=http://localhost:8080
export BATON_API_KEY=test-api-key
export BATON_SOURCE=test-source
./baton-litmos
```

## Curl examples

```bash
# List users
curl -H 'apikey: test-api-key' 'http://localhost:8080/v1.svc/users?source=test&limit=500'

# List teams
curl -H 'apikey: test-api-key' 'http://localhost:8080/v1.svc/teams?source=test'

# Get course by ID
curl -H 'apikey: test-api-key' 'http://localhost:8080/v1.svc/courses/course-intro?source=test'

# List course users
curl -H 'apikey: test-api-key' 'http://localhost:8080/v1.svc/courses/course-intro/users?source=test'

# Assign a course to a user
curl -X POST -H 'apikey: test-api-key' -H 'Content-Type: application/xml' \
  'http://localhost:8080/v1.svc/users/user-4/courses?source=test' \
  -d '<Courses><Course><Id>course-intro</Id></Course></Courses>'

# Remove a course from a user
curl -X DELETE -H 'apikey: test-api-key' \
  'http://localhost:8080/v1.svc/users/user-4/courses/course-intro?source=test'

# Rejected: missing apikey
curl 'http://localhost:8080/v1.svc/users?source=test'
# → 401 Unauthorized
```
