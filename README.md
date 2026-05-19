# Authentication system endpoints
- POST
/api/auth/register
gRPC: Register
- POST
/api/auth/verify-email
gRPC: VerifyEmail
- POST
/api/auth/login
gRPC: Login
- POST
/api/auth/refresh
gRPC: RefreshToken
- POST
/api/auth/logout
gRPC: Logout
- POST
/api/auth/forgot-password
gRPC: ForgotPassword
- POST
/api/auth/reset-password

- GET
/api/auth/account/{accountId}
gRPC: GetAccount
- PUT
/api/auth/account/{accountId}
gRPC: UpdateAccount
- DELETE
/api/auth/account/{accountId}

- POST
/api/profiles
gRPC: CreateProfile
- GET
/api/profiles?account_id=…
gRPC: GetProfiles
- GET
/api/profiles/{profileId}
gRPC: GetProfileById
- PUT
/api/profiles/{profileId}
gRPC: UpdateProfile
- DELETE
/api/profiles/{profileId}

# Content Service system endpoints

- POST
/api/content
gRPC: CreateContent
- GET
/api/content/{id}
gRPC: GetContent
- GET
/api/content
gRPC: ListContent
- PUT
/api/content/{id}
gRPC: UpdateContent
- DELETE
/api/content/{id}
gRPC: DeleteContent
- GET
/api/content/search?q={query}
gRPC: SearchContent
- GET
/api/content/genre/{genre}
gRPC: GetContentByGenre
- GET
/api/content/trending
gRPC: GetTrendingContent
- GET
/api/content/featured
gRPC: GetFeaturedContent
- POST
/api/episodes
gRPC: CreateEpisode
- GET
/api/episodes/{id}
gRPC: GetEpisode
- GET
/api/content/{contentId}/episodes
gRPC: ListEpisodesByContent
- PUT
/api/episodes/{id}
gRPC: UpdateEpisode
- DELETE
/api/episodes/{id}
gRPC: DeleteEpisode
- PATCH
/api/episodes/{id}/skip-intro
gRPC: SetSkipIntroMetadata

# Playback Service system endpoints

- POST
/api/playback/start
gRPC: StartPlayback

- PUT
/api/playback/progress
gRPC: UpdateProgress

- GET
/api/playback/progress/:profileId/:contentId/:episodeId
gRPC: GetProgress

- GET
/api/playback/resume/:profileId/:contentId/:episodeId
gRPC: GetResumeData

- GET
/api/playback/continue-watching/:profileId
gRPC: GetContinueWatching

- DELETE
/api/playback/continue-watching/:profileId/:contentId/:episodeId
gRPC: RemoveFromContinueWatching

- PUT
/api/playback/completed
gRPC: MarkCompleted

- POST
/api/playback/end-session
gRPC: EndPlaybackSession

- GET
/api/playback/session/:sessionId
gRPC: GetPlaybackSession

- GET
/api/playback/history/:profileId
gRPC: GetWatchHistory

- DELETE
/api/playback/history/:profileId
gRPC: ClearWatchHistory

- GET
/api/playback/stats/:profileId
gRPC: GetProfileStats

- POST
/api/playback/skip-intro
gRPC: SkipIntro

- GET
/api/playback/last-watched/:profileId
gRPC: GetLastWatchedContent

- DELETE
/api/playback/progress/:profileId/:contentId/:episodeId
gRPC: DeleteProgress

# Grafana query commands

## Prometheous:

1. rate(container_cpu_usage_seconds_total[1m])

## Loki:

1. {job="docker"}
2. curl "http://localhost:8080/api/content?limit=5"
docker compose logs api-gateway --tail=20
docker compose logs content-service --tail=20

## Tempo:

1. {resource.service.name="api-gateway"}