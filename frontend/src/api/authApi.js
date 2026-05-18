const BASE_URL = "/api";

const AUTH_USER_KEY = "dorabox:user";
const AUTH_TOKEN_KEY = "dorabox:accessToken";
const REFRESH_TOKEN_KEY = "dorabox:refreshToken";
const ACCOUNT_ID_KEY = "dorabox:accountId";

async function request(path, options = {}) {
  const token = localStorage.getItem(AUTH_TOKEN_KEY);

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

function normalizeUser(data, fallback = {}) {
  return {
    id: data.accountId || data.account_id || fallback.id || fallback.accountId,
    accountId: data.accountId || data.account_id || fallback.accountId || fallback.id,
    email: data.email || fallback.email || "",
    username: data.username || fallback.username || fallback.fullName || "User",
    fullName: data.username || fallback.fullName || fallback.username || "User",
    isVerified: data.isVerified ?? data.is_verified ?? fallback.isVerified ?? false,
    accountType: fallback.accountType || "standard",
    createdAt: data.createdAt || data.created_at || fallback.createdAt || "",
  };
}

export async function registerUser(payload) {
  const username =
    payload.username ||
    payload.fullName ||
    payload.name ||
    payload.email?.split("@")[0] ||
    "user";

  const data = await request("/auth/register", {
    method: "POST",
    body: JSON.stringify({
      email: payload.email,
      username,
      password: payload.password,
    }),
  });

  const user = normalizeUser(
    {
      accountId: data.accountId,
      email: payload.email,
      username,
      isVerified: false,
    },
    payload
  );

  localStorage.setItem(AUTH_USER_KEY, JSON.stringify(user));
  localStorage.setItem(ACCOUNT_ID_KEY, user.accountId);

  return {
    user,
    accountId: user.accountId,
    message: data.message || "Registration completed.",
  };
}

export async function loginUser(payload) {
  const data = await request("/auth/login", {
    method: "POST",
    body: JSON.stringify({
      email: payload.email,
      password: payload.password,
    }),
  });

  localStorage.setItem(AUTH_TOKEN_KEY, data.accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, data.refreshToken);
  localStorage.setItem(ACCOUNT_ID_KEY, data.accountId);

  let user = normalizeUser(
    {
      accountId: data.accountId,
      email: payload.email,
    },
    payload
  );

  try {
    const account = await request(`/auth/account/${data.accountId}`);
    user = normalizeUser(account, {
      accountId: data.accountId,
      email: payload.email,
    });
  } catch {
    // Login still succeeds even if account fetch fails.
    // Frontend can survive. Barely. Like all of us.
  }

  localStorage.setItem(AUTH_USER_KEY, JSON.stringify(user));

  return {
    user,
    accessToken: data.accessToken,
    refreshToken: data.refreshToken,
    accountId: data.accountId,
  };
}

export async function refreshAccessToken() {
  const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);

  if (!refreshToken) {
    throw new Error("No refresh token found");
  }

  const data = await request("/auth/refresh", {
    method: "POST",
    body: JSON.stringify({
      refreshToken,
    }),
  });

  localStorage.setItem(AUTH_TOKEN_KEY, data.accessToken);

  if (data.refreshToken) {
    localStorage.setItem(REFRESH_TOKEN_KEY, data.refreshToken);
  }

  return data;
}

export function getCurrentUser() {
  const raw = localStorage.getItem(AUTH_USER_KEY);

  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw);
  } catch {
    localStorage.removeItem(AUTH_USER_KEY);
    return null;
  }
}

export function getCurrentAccountId() {
  const user = getCurrentUser();
  return (
    localStorage.getItem(ACCOUNT_ID_KEY) ||
    user?.accountId ||
    user?.id ||
    null
  );
}

export function getAccessToken() {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function isAuthenticated() {
  return Boolean(localStorage.getItem(AUTH_TOKEN_KEY));
}

export async function logoutUser() {
  const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);

  if (refreshToken) {
    try {
      await request("/auth/logout", {
        method: "POST",
        body: JSON.stringify({
          refreshToken,
        }),
      });
    } catch {
      // Local logout should still happen even if backend logout fails.
    }
  }

  localStorage.removeItem(AUTH_USER_KEY);
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(ACCOUNT_ID_KEY);
  localStorage.removeItem("dorabox:selectedProfileId");
}