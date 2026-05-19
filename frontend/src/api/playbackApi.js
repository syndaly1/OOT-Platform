import { getSelectedProfileId } from "./profileApi";
import { getContentById } from "./contentApi";

const API_BASE_URL = "/api";

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
  });

  const data = await response.json().catch(() => ({}));

  if (!response.ok) {
    throw new Error(data.error || data.message || `Playback API error: ${response.status}`);
  }

  return data;
}

function normalizePlaybackItem(entry) {
  const contentId = entry.contentId || entry.content_id;
  const episodeId = entry.episodeId || entry.episode_id;

  return {
    sessionId: entry.sessionId || entry.session_id,
    profileId: entry.profileId || entry.profile_id,
    contentId,
    episodeId,
    currentPositionSec:
      entry.currentPositionSec ||
      entry.current_position_sec ||
      entry.currentPosition ||
      0,
    durationSec: entry.durationSec || entry.duration_sec || 0,
    watchedPercent: entry.watchedPercent || entry.watched_percent || 0,
    completed: Boolean(entry.completed),
    status: entry.status || "PLAYING",
    updatedAt: entry.updatedAt || entry.updated_at || new Date().toISOString(),
  };
}

export async function startPlayback({
  profileId,
  contentId,
  episodeId = "episode-1",
  durationSec = 3600,
  device = "web",
}) {
  return request("/playback/start", {
    method: "POST",
    body: JSON.stringify({
      profileId,
      contentId,
      episodeId,
      duration_sec: durationSec,
      device,
    }),
  });
}

export async function updateProgress({
  profileId,
  contentId,
  episodeId = "episode-1",
  sessionId = "session-created",
  currentPositionSec,
  durationSec,
}) {
  let finalDurationSec = durationSec;

  if (!finalDurationSec) {
    try {
      const item = await getContentById(contentId);
      finalDurationSec = item.durationSec || 3600;
    } catch {
      finalDurationSec = 3600;
    }
  }

  const data = await request("/playback/progress", {
    method: "PUT",
    body: JSON.stringify({
      sessionId,
      profileId,
      contentId,
      episodeId,
      current_position_sec: currentPositionSec,
      duration_sec: finalDurationSec,
    }),
  });

  return {
    profileId,
    contentId,
    episodeId,
    sessionId,
    currentPositionSec,
    durationSec: finalDurationSec,
    watchedPercent: data.watchedPercent,
    completed: data.watchedPercent >= 90,
    updatedAt: new Date().toISOString(),
    ...data,
  };
}

export async function getProgress({ profileId, contentId, episodeId = "episode-1" }) {
  const items = await getContinueWatching(profileId);

  return (
    items.find(
      (entry) =>
        entry.profileId === profileId &&
        entry.contentId === contentId &&
        entry.episodeId === episodeId
    ) || {
      profileId,
      contentId,
      episodeId,
      currentPositionSec: 0,
      completed: false,
    }
  );
}

export async function getContinueWatching(profileId = getSelectedProfileId()) {
  if (!profileId) {
    return [];
  }

  const data = await request(`/playback/continue-watching/${profileId}`);
  const items = data.items || data.sessions || [];

  const normalized = items.map(normalizePlaybackItem);

  const withContent = await Promise.all(
    normalized.map(async (entry) => {
      try {
        const item = await getContentById(entry.contentId);
        return {
          ...entry,
          item,
        };
      } catch {
        return entry;
      }
    })
  );

  return withContent;
}

export async function markCompleted({
  profileId,
  contentId,
  episodeId = "episode-1",
  sessionId = "session-created",
  durationSec = 3600,
}) {
  return updateProgress({
    profileId,
    contentId,
    episodeId,
    sessionId,
    currentPositionSec: durationSec,
    durationSec,
  });
}

export async function clearProgress({
  profileId,
  contentId,
  episodeId = "episode-1",
}) {
  return request(
    `/playback/continue-watching/${profileId}/${contentId}?episodeId=${encodeURIComponent(
      episodeId
    )}`,
    {
      method: "DELETE",
    }
  );
}