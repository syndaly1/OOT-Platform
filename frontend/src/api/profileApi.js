import { getCurrentUser } from "./authApi";

const SELECTED_PROFILE_KEY = "dorabox:selectedProfileId";
const PROFILES_KEY = "dorabox:profiles";

function notifyProfileChanged() {
  window.dispatchEvent(new Event("profile-changed"));
}

function getProfilesMap() {
  const raw = localStorage.getItem(PROFILES_KEY);

  if (!raw) {
    return {};
  }

  try {
    return JSON.parse(raw);
  } catch {
    return {};
  }
}

function saveProfilesMap(map) {
  localStorage.setItem(PROFILES_KEY, JSON.stringify(map));
  notifyProfileChanged();
}

function getAccountId() {
  const user = getCurrentUser();
  return user?.accountId || user?.id || "guest-account";
}

function getDefaultProfile() {
  const user = getCurrentUser();

  return {
    id: "profile-test-1",
    profileId: "profile-test-1",
    accountId: getAccountId(),
    name: user?.fullName || user?.username || "Main",
    type: "adult",
    language: "en",
    avatar: "👤",
  };
}

function normalizeProfile(profile) {
  return {
    id: profile.id || profile.profileId || crypto.randomUUID(),
    profileId: profile.profileId || profile.id || crypto.randomUUID(),
    accountId: profile.accountId || getAccountId(),
    name: profile.name || "Main",
    type: String(profile.type || "adult").toLowerCase(),
    language: profile.language || "en",
    avatar: profile.avatar || "👤",
  };
}

export async function getProfiles() {
  const user = getCurrentUser();
  const accountId = getAccountId();

  if (!user) {
    return [getDefaultProfile()];
  }

  if (user.accountType !== "family") {
    return [getDefaultProfile()];
  }

  const map = getProfilesMap();
  const profiles = map[accountId] || [];

  if (!profiles.length) {
    const defaultProfile = getDefaultProfile();
    map[accountId] = [defaultProfile];
    saveProfilesMap(map);
    return [defaultProfile];
  }

  return profiles.map(normalizeProfile);
}

export async function createProfile(payload) {
  const user = getCurrentUser();

  if (!user) {
    throw new Error("Please sign in before creating profiles.");
  }

  if (user.accountType !== "family") {
    throw new Error("Only Family accounts can create multiple profiles.");
  }

  const accountId = getAccountId();
  const map = getProfilesMap();

  const profile = normalizeProfile({
    id: crypto.randomUUID(),
    profileId: crypto.randomUUID(),
    accountId,
    name: payload.name,
    type: payload.type || "adult",
    language: payload.language || "en",
    avatar: payload.avatar || "👤",
  });

  map[accountId] = [...(map[accountId] || []), profile];
  saveProfilesMap(map);

  return profile;
}

export async function updateProfile(profileId, payload) {
  const accountId = getAccountId();
  const map = getProfilesMap();
  const profiles = map[accountId] || [];

  map[accountId] = profiles.map((profile) =>
    profile.id === profileId || profile.profileId === profileId
      ? normalizeProfile({ ...profile, ...payload })
      : profile
  );

  saveProfilesMap(map);

  return map[accountId].find(
    (profile) => profile.id === profileId || profile.profileId === profileId
  );
}

export async function deleteProfile(profileId) {
  const accountId = getAccountId();
  const map = getProfilesMap();
  const profiles = map[accountId] || [];

  map[accountId] = profiles.filter(
    (profile) => profile.id !== profileId && profile.profileId !== profileId
  );

  saveProfilesMap(map);

  if (getSelectedProfileId() === profileId) {
    localStorage.removeItem(SELECTED_PROFILE_KEY);
  }

  notifyProfileChanged();

  return { message: "Profile deleted." };
}

export function selectProfile(profileId) {
  localStorage.setItem(SELECTED_PROFILE_KEY, profileId || "profile-test-1");
  notifyProfileChanged();
}

export function getSelectedProfileId() {
  return localStorage.getItem(SELECTED_PROFILE_KEY) || "profile-test-1";
}

export function getSelectedProfile() {
  const selectedId = getSelectedProfileId();
  const user = getCurrentUser();

  if (!user || user.accountType !== "family") {
    return getDefaultProfile();
  }

  const accountId = getAccountId();
  const map = getProfilesMap();
  const profiles = map[accountId] || [];

  return (
    profiles.find(
      (profile) => profile.id === selectedId || profile.profileId === selectedId
    ) ||
    profiles[0] ||
    getDefaultProfile()
  );
}