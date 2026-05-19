# Streaming Platform API Documentation

This project uses one shared API Gateway as the REST entry point for the frontend.

The API Gateway forwards requests to internal gRPC microservices:

* Auth Service
* Content Service
* Playback Service

The gRPC contracts are defined in the `proto` directory. REST routes are implemented in `api-gateway/main.go`.

---

## Authentication System Endpoints

### Auth Routes

#### Register Account

**REST:** `POST /api/auth/register`
**gRPC:** `AuthService.Register`

#### Verify Email

**REST:** `POST /api/auth/verify-email`
**gRPC:** `AuthService.VerifyEmail`

#### Login

**REST:** `POST /api/auth/login`
**gRPC:** `AuthService.Login`

#### Refresh Token

**REST:** `POST /api/auth/refresh`
**gRPC:** `AuthService.RefreshToken`

#### Logout

**REST:** `POST /api/auth/logout`
**gRPC:** `AuthService.Logout`

#### Forgot Password

**REST:** `POST /api/auth/forgot-password`
**gRPC:** `AuthService.ForgotPassword`

#### Reset Password

**REST:** `POST /api/auth/reset-password`
**gRPC:** `AuthService.ResetPassword`

---

### Account Routes

#### Get Account

**REST:** `GET /api/auth/account/{accountId}`
**gRPC:** `AuthService.GetAccount`

#### Update Account

**REST:** `PUT /api/auth/account/{accountId}`
**gRPC:** `AuthService.UpdateAccount`

#### Delete Account

**REST:** `DELETE /api/auth/account/{accountId}`
**gRPC:** `AuthService.DeleteAccount`

---

### Profile Routes

#### Create Profile

**REST:** `POST /api/profiles`
**gRPC:** `AuthService.CreateProfile`

#### Get Profiles by Account

**REST:** `GET /api/profiles?accountId={accountId}`
**gRPC:** `AuthService.GetProfiles`

#### Get Profile by ID

**REST:** `GET /api/profiles/{profileId}`
**gRPC:** `AuthService.GetProfileById`

#### Update Profile

**REST:** `PUT /api/profiles/{profileId}`
**gRPC:** `AuthService.UpdateProfile`

#### Delete Profile

**REST:** `DELETE /api/profiles/{profileId}`
**gRPC:** `AuthService.DeleteProfile`

---

## Content Service Endpoints

### Content Routes

#### Create Content

**REST:** `POST /api/content`
**gRPC:** `ContentService.CreateContent`

#### List Content

**REST:** `GET /api/content`
**gRPC:** `ContentService.ListContent`

Example:

```bash
curl "http://localhost:8080/api/content?limit=5"
```

#### Get Content by ID

**REST:** `GET /api/content/{contentId}`
**gRPC:** `ContentService.GetContent`

Example:

```bash
curl "http://localhost:8080/api/content/10000000-0000-0000-0000-000000000001"
```

#### Update Content

**REST:** `PUT /api/content/{contentId}`
**gRPC:** `ContentService.UpdateContent`

#### Delete Content

**REST:** `DELETE /api/content/{contentId}`
**gRPC:** `ContentService.DeleteContent`

#### Search Content

**REST:** `GET /api/content/search?q={query}`
**gRPC:** `ContentService.SearchContent`

#### Get Content by Genre

**REST:** `GET /api/content/genre/{genre}`
**gRPC:** `ContentService.GetContentByGenre`

#### Get Trending Content

**REST:** `GET /api/content/trending`
**gRPC:** `ContentService.GetTrendingContent`

#### Get Featured Content

**REST:** `GET /api/content/featured`
**gRPC:** `ContentService.GetFeaturedContent`

---

### Episode Routes

#### Create Episode for Content

**REST:** `POST /api/contents/{contentId}/episodes`
**gRPC:** `ContentService.CreateEpisode`

#### List Episodes by Content

**REST:** `GET /api/contents/{contentId}/episodes`
**gRPC:** `ContentService.ListEpisodesByContent`

#### Get Episode by ID

**REST:** `GET /api/episodes/{episodeId}`
**gRPC:** `ContentService.GetEpisode`

#### Update Episode

**REST:** `PUT /api/episodes/{episodeId}`
**gRPC:** `ContentService.UpdateEpisode`

#### Delete Episode

**REST:** `DELETE /api/episodes/{episodeId}`
**gRPC:** `ContentService.DeleteEpisode`

#### Set Skip Intro Metadata

**REST:** `PUT /api/episodes/{episodeId}/skip-intro`
**gRPC:** `ContentService.SetSkipIntroMetadata`

---

## Playback Service Endpoints

The Playback proto file contains more gRPC methods, but the API Gateway currently exposes the following REST routes.

#### Start Playback

**REST:** `POST /api/playback/start`
**gRPC:** `PlaybackService.StartPlayback`

#### Update Progress

**REST:** `PUT /api/playback/progress`
**gRPC:** `PlaybackService.UpdateProgress`

#### Get Continue Watching

**REST:** `GET /api/playback/continue-watching/{profileId}`
**gRPC:** `PlaybackService.GetContinueWatching`

Example:

```bash
curl "http://localhost:8080/api/playback/continue-watching/profile-test-1"
```

#### Get Watch History

**REST:** `GET /api/playback/history/{profileId}`
**gRPC:** `PlaybackService.GetWatchHistory`

#### Remove from Continue Watching

**REST:** `DELETE /api/playback/continue-watching/{profileId}/{contentId}`
**gRPC:** `PlaybackService.RemoveFromContinueWatching`

---

## gRPC Contracts

The full internal gRPC methods are defined in:

* `proto/auth`
* `proto/content`
* `proto/playback`

Useful command to list gRPC services and methods:

```bash
grep -R "service \|rpc " -n proto
```

Useful command to list API Gateway REST routes:

```bash
grep -n "mux.HandleFunc\|mux.Handle" api-gateway/main.go
```

---

## NATS Message Queue

NATS is used for asynchronous domain events.

### Auth Events

* `user.registered`
* `password.reset.requested`

### Content Events

* `content.created`
* `content.updated`
* `content.deleted`

### Playback Events

* Progress update events
* Completion events

Useful commands:

```bash
docker compose ps nats
docker compose logs nats --tail=40
grep -R "Publish\|Subscribe\|NATS\|nats" -n services api-gateway
```

NATS monitoring page:

```text
http://localhost:8222
```

---

## Database Migrations

Each service owns its own database. Migrations are stored under service migration folders.

Useful command:

```bash
find services -path "*migrations*" -type f | sort
```

Check applied content migrations:

```bash
docker exec -it content-postgres psql -U postgres -d content_db -c "SELECT version, dirty FROM schema_migrations;"
```

Check seeded content count:

```bash
docker exec -it content-postgres psql -U postgres -d content_db -c "SELECT COUNT(*) FROM contents;"
```

Check sample seeded content:

```bash
docker exec -it content-postgres psql -U postgres -d content_db -c "SELECT id, title, poster_url FROM contents LIMIT 5;"
```

---

## Redis Cache

Redis is used as a cache / temporary storage layer.

Useful commands:

```bash
docker compose ps redis
docker exec -it redis redis-cli ping
grep -R "redis\|Redis\|cache\|Cache" -n services
```

Expected Redis ping result:

```text
PONG
```

---

## Grafana Observability

Grafana is available at:

```text
http://localhost:3000
```

Default credentials:

```text
admin / admin
```

The observability stack includes:

* Prometheus for metrics
* cAdvisor for container metrics
* Loki for logs
* Promtail for Docker log collection
* Tempo for traces
* OpenTelemetry tracing in API Gateway

---

## Prometheus Metrics Query

Grafana → Explore → Prometheus:

```promql
rate(container_cpu_usage_seconds_total[1m])
```

If `/metrics` is enabled in API Gateway, also show:

```promql
promhttp_metric_handler_requests_total
```

or:

```promql
go_goroutines
```

---

## Loki Logs Query

Grafana → Explore → Loki:

```logql
{job="docker"}
```

To generate fresh logs:

```bash
curl "http://localhost:8080/api/content?limit=5"
docker compose logs api-gateway --tail=20
docker compose logs content-service --tail=20
```

---

## Tempo Tracing Query

Grafana → Explore → Tempo:

```traceql
{resource.service.name="api-gateway"}
```

To generate fresh traces:

```bash
curl "http://localhost:8080/api/content?limit=5"
curl "http://localhost:8080/api/content/10000000-0000-0000-0000-000000000001"
curl "http://localhost:8080/api/playback/continue-watching/profile-test-1"
```

The API Gateway creates HTTP spans for incoming requests and gRPC client spans for calls to Auth, Content, and Playback services.

---

## Git Commands

After replacing the README content:

```bash
cd /Users/ingkarpahr/Desktop/streaming-platform

git add README.md
git status
```
