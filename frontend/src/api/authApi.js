const AUTH_USER_KEY = "dorabox:user";
const AUTH_TOKEN_KEY = "dorabox:accessToken";

function delay(ms = 350) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function registerUser(payload) {
  await delay();

  const user = {
    id: crypto.randomUUID(),
    fullName: payload.fullName,
    email: payload.email,
    accountType: payload.accountType,
  };

  localStorage.setItem(AUTH_USER_KEY, JSON.stringify(user));
  localStorage.setItem(AUTH_TOKEN_KEY, "mock-access-token");

  return {
    user,
    accessToken: "mock-access-token",
    message: "Registration completed.",
  };
}

export async function loginUser(payload) {
  await delay();

  const user = {
    id: "mock-user-1",
    fullName: "Demo User",
    email: payload.email,
    accountType: "family",
  };

  localStorage.setItem(AUTH_USER_KEY, JSON.stringify(user));
  localStorage.setItem(AUTH_TOKEN_KEY, "mock-access-token");

  return {
    user,
    accessToken: "mock-access-token",
  };
}

export function getCurrentUser() {
  const raw = localStorage.getItem(AUTH_USER_KEY);
  return raw ? JSON.parse(raw) : null;
}

export function isAuthenticated() {
  return Boolean(localStorage.getItem(AUTH_TOKEN_KEY));
}

export function logoutUser() {
  localStorage.removeItem(AUTH_USER_KEY);
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem("dorabox:selectedProfileId");
}