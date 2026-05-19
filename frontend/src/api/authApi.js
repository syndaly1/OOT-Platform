const AUTH_USER_KEY = "dorabox:user";
const AUTH_TOKEN_KEY = "dorabox:accessToken";
const AUTH_REFRESH_TOKEN_KEY = "dorabox:refreshToken";
const AUTH_ACCOUNTS_KEY = "dorabox:accounts";

function notifyAuthChanged() {
  window.dispatchEvent(new Event("auth-changed"));
}

function getAccounts() {
  const raw = localStorage.getItem(AUTH_ACCOUNTS_KEY);

  if (!raw) {
    return [];
  }

  try {
    return JSON.parse(raw);
  } catch {
    return [];
  }
}

function saveAccounts(accounts) {
  localStorage.setItem(AUTH_ACCOUNTS_KEY, JSON.stringify(accounts));
}

function normalizeEmail(email) {
  return String(email || "").trim().toLowerCase();
}

function makeUserFromAccount(account) {
  return {
    id: account.id,
    accountId: account.id,
    username: account.username || account.fullName || account.email,
    fullName: account.fullName || account.username || account.email,
    email: account.email,
    accountType: account.accountType || "basic",
  };
}

function saveAuthSession(user) {
  localStorage.setItem(AUTH_USER_KEY, JSON.stringify(user));
  localStorage.setItem(AUTH_TOKEN_KEY, "frontend-demo-session");
  localStorage.setItem(AUTH_REFRESH_TOKEN_KEY, "frontend-demo-refresh");
  notifyAuthChanged();

  return {
    user,
    accessToken: "frontend-demo-session",
    refreshToken: "frontend-demo-refresh",
    message: "Authenticated successfully.",
  };
}

export async function registerUser(payload) {
  const email = normalizeEmail(payload.email);
  const accounts = getAccounts();

  const existing = accounts.find((account) => normalizeEmail(account.email) === email);

  if (existing) {
    throw new Error("This email is already registered. Please sign in instead.");
  }

  const account = {
    id: crypto.randomUUID(),
    fullName: payload.fullName || payload.username || email,
    username: payload.username || payload.fullName || email,
    email,
    password: payload.password,
    accountType: payload.accountType || "basic",
    createdAt: new Date().toISOString(),
  };

  saveAccounts([...accounts, account]);

  const user = makeUserFromAccount(account);
  return saveAuthSession(user);
}

export async function loginUser(payload) {
  const email = normalizeEmail(payload.email);
  const accounts = getAccounts();

  const account = accounts.find((item) => normalizeEmail(item.email) === email);

  if (!account) {
    throw new Error("Account not found. Please register first.");
  }

  if (account.password !== payload.password) {
    throw new Error("Incorrect password.");
  }

  const user = makeUserFromAccount(account);
  return saveAuthSession(user);
}

export function getCurrentUser() {
  const raw = localStorage.getItem(AUTH_USER_KEY);

  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

export function getAccessToken() {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function getRefreshToken() {
  return localStorage.getItem(AUTH_REFRESH_TOKEN_KEY);
}

export function isAuthenticated() {
  return Boolean(localStorage.getItem(AUTH_TOKEN_KEY));
}

export async function refreshAuthToken() {
  const user = getCurrentUser();

  if (!user) {
    throw new Error("No active session.");
  }

  return saveAuthSession(user);
}

export async function logoutUser() {
  localStorage.removeItem(AUTH_USER_KEY);
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(AUTH_REFRESH_TOKEN_KEY);
  localStorage.removeItem("dorabox:selectedProfileId");

  notifyAuthChanged();
}