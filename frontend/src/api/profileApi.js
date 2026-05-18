const PROFILES_KEY = "dorabox:profiles";
const SELECTED_PROFILE_KEY = "dorabox:selectedProfileId";

const defaultProfiles = [
  {
    id: "profile-1",
    name: "Main",
    type: "adult",
    avatar: "🌙",
  },
  {
    id: "profile-2",
    name: "Family",
    type: "family",
    avatar: "🍿",
  },
  {
    id: "profile-3",
    name: "Kid",
    type: "kid",
    avatar: "🧸",
  },
];

function ensureProfiles() {
  const saved = localStorage.getItem(PROFILES_KEY);

  if (!saved) {
    localStorage.setItem(PROFILES_KEY, JSON.stringify(defaultProfiles));
    return defaultProfiles;
  }

  return JSON.parse(saved);
}

export async function getProfiles() {
  return ensureProfiles();
}

export async function createProfile(payload) {
  const profiles = ensureProfiles();

  const newProfile = {
    id: crypto.randomUUID(),
    name: payload.name,
    type: payload.type,
    avatar: payload.avatar || "👤",
  };

  const next = [...profiles, newProfile];
  localStorage.setItem(PROFILES_KEY, JSON.stringify(next));

  return newProfile;
}

export function selectProfile(profileId) {
  localStorage.setItem(SELECTED_PROFILE_KEY, profileId);
}

export function getSelectedProfileId() {
  return localStorage.getItem(SELECTED_PROFILE_KEY) || "profile-1";
}

export function getSelectedProfile() {
  const profiles = ensureProfiles();
  const selectedId = getSelectedProfileId();

  return profiles.find((profile) => profile.id === selectedId) || profiles[0];
}