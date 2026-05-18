import { getCurrentAccountId } from "./authApi";

const BASE_URL = "/api";
const SELECTED_PROFILE_KEY = "dorabox:selectedProfileId";

async function request(path, options = {}) {
  const token = localStorage.getItem("dorabox:accessToken");

  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options.headers || {}),
    },
  });

  const data = await res.json().catch(() => ({}));

  if (!res.ok) {
    throw new Error(data.error || data.message || "Request failed");
  }

  return data;
}

function normalizeProfile(profile) {
  return {
    id: profile.profileId || profile.profile_id || profile.id,
    profileId: profile.profileId || profile.profile_id || profile.id,
    accountId: profile.accountId || profile.account_id,
    name: profile.name || "Profile",
    type: normalizeProfileType(profile.type),
    language: profile.language || "en",
    avatar: profile.avatar || "👤",
    createdAt: profile.createdAt || profile.created_at || "",
  };
}

function normalizeProfileType(type) {
  if (typeof type === "number") {
    if (type === 1) return "adult";
    if (type === 2) return "teen";
    if (type === 3) return "kid";
  }

  const raw = String(type || "adult").toLowerCase();

  if (raw.includes("teen")) return "teen";
  if (raw.includes("kid")) return "kid";
  if (raw.includes("child")) return "kid";
  if (raw.includes("family")) return "adult";

  return "adult";
}

export async function getProfiles(accountIdArg) {
  const accountId = accountIdArg || getCurrentAccountId();

  if (!accountId) {
    return [];
  }

  const data = await request(`/profiles?accountId=${encodeURIComponent(accountId)}`);

  const profiles = data.profiles || [];

  return profiles.map(normalizeProfile);
}

export async function createProfile(payload) {
  const accountId = payload.accountId || payload.account_id || getCurrentAccountId();

  if (!accountId) {
    throw new Error("No accountId found. Please login first.");
  }

  const data = await request("/profiles", {
    method: "POST",
    body: JSON.stringify({
      accountId,
      name: payload.name,
      type: payload.type || "adult",
      language: payload.language || "en",
      avatar: payload.avatar || "👤",
    }),
  });

  const profileId = data.profileId || data.profile_id;

  return {
    id: profileId,
    profileId,
    accountId,
    name: payload.name,
    type: payload.type || "adult",
    language: payload.language || "en",
    avatar: payload.avatar || "👤",
  };
}

export async function updateProfile(profileId, payload) {
  const data = await request(`/profiles/${profileId}`, {
    method: "PUT",
    body: JSON.stringify({
      name: payload.name,
      language: payload.language || "en",
      avatar: payload.avatar || "👤",
    }),
  });

  return data;
}

export async function deleteProfile(profileId) {
  const data = await request(`/profiles/${profileId}`, {
    method: "DELETE",
  });

  const selectedId = getSelectedProfileId();

  if (selectedId === profileId) {
    localStorage.removeItem(SELECTED_PROFILE_KEY);
  }

  return data;
}

export function selectProfile(profileId) {
  localStorage.setItem(SELECTED_PROFILE_KEY, profileId);
}

export function getSelectedProfileId() {
  return localStorage.getItem(SELECTED_PROFILE_KEY);
}

export async function getSelectedProfile() {
  const profiles = await getProfiles();
  const selectedId = getSelectedProfileId();

  if (!profiles.length) {
    return null;
  }

  return profiles.find((profile) => profile.id === selectedId) || profiles[0];
}