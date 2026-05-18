import { getItemById } from "../data/catalog";
import { getSelectedProfileId } from "./profileApi";

const PROGRESS_KEY = "dorabox:progress";

function getProgressMap() {
  const raw = localStorage.getItem(PROGRESS_KEY);
  return raw ? JSON.parse(raw) : {};
}

function saveProgressMap(map) {
  localStorage.setItem(PROGRESS_KEY, JSON.stringify(map));
}

function makeProgressKey(profileId, contentId) {
  return `${profileId}:${contentId}`;
}

export async function updateProgress({ profileId, contentId, currentPositionSec }) {
  const progressMap = getProgressMap();
  const item = getItemById(contentId);

  if (!item) {
    throw new Error("Content not found.");
  }

  const watchedPercent = currentPositionSec / item.durationSec;
  const completed = watchedPercent >= 0.9;

  progressMap[makeProgressKey(profileId, contentId)] = {
    profileId,
    contentId,
    currentPositionSec,
    durationSec: item.durationSec,
    completed,
    updatedAt: new Date().toISOString(),
  };

  saveProgressMap(progressMap);

  return progressMap[makeProgressKey(profileId, contentId)];
}

export async function getProgress({ profileId, contentId }) {
  const progressMap = getProgressMap();

  return (
    progressMap[makeProgressKey(profileId, contentId)] || {
      profileId,
      contentId,
      currentPositionSec: 0,
      completed: false,
    }
  );
}

export async function getContinueWatching(profileId = getSelectedProfileId()) {
  const progressMap = getProgressMap();

  return Object.values(progressMap)
    .filter((entry) => entry.profileId === profileId)
    .filter((entry) => !entry.completed)
    .filter((entry) => entry.currentPositionSec > 10)
    .map((entry) => {
      const item = getItemById(entry.contentId);

      return {
        ...entry,
        item,
      };
    })
    .filter((entry) => Boolean(entry.item))
    .sort((a, b) => new Date(b.updatedAt) - new Date(a.updatedAt));
}

export async function markCompleted({ profileId, contentId }) {
  const item = getItemById(contentId);

  if (!item) {
    throw new Error("Content not found.");
  }

  return updateProgress({
    profileId,
    contentId,
    currentPositionSec: item.durationSec,
  });
}

export async function clearProgress({ profileId, contentId }) {
  const progressMap = getProgressMap();
  delete progressMap[makeProgressKey(profileId, contentId)];
  saveProgressMap(progressMap);
}