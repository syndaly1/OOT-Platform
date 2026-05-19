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
    throw new Error(
      data.error || data.message || `Content API error: ${response.status}`
    );
  }

  return data;
}

function normalizeType(type) {
  if (!type) return "tv-series";
  if (type === "tv_series") return "tv-series";
  return type;
}

function inferRegion(item) {
  const text = `${item.genre || ""} ${item.posterUrl || item.poster_url || ""} ${
    item.title || ""
  }`.toLowerCase();

  if (
    text.includes("south-korea") ||
    text.includes("kdrama") ||
    text.includes("k-drama") ||
    text.includes("korea")
  ) {
    return "south-korea";
  }

  if (
    text.includes("china") ||
    text.includes("cdrama") ||
    text.includes("c-drama")
  ) {
    return "china";
  }

  if (
    text.includes("japan") ||
    text.includes("jdrama") ||
    text.includes("j-drama")
  ) {
    return "japan";
  }

  return "all";
}

function inferStatus(item) {
  const text = `${item.genre || ""} ${item.title || ""}`.toLowerCase();

  if (text.includes("upcoming")) return "upcoming";
  if (text.includes("ongoing")) return "ongoing";
  if (text.includes("completed")) return "completed";

  return "completed";
}

function inferSubtitle(item) {
  const text = `${item.genre || ""}`.toLowerCase();

  if (text.includes("russian")) return "russian";
  return "english";
}

function inferAgeRating(item) {
  const text = `${item.genre || ""}`.toLowerCase();

  if (text.includes("18+")) return "18+";
  if (text.includes("16+")) return "16+";
  return "13+";
}

function inferCollection(item) {
  const text = `${item.genre || ""}`.toLowerCase();

  if (text.includes("latest-update")) return "latest-update";
  if (text.includes("top-k-drama")) return "top-k-drama";
  if (text.includes("top-c-drama")) return "top-c-drama";
  if (text.includes("top-j-drama")) return "top-j-drama";
  if (text.includes("upcoming")) return "upcoming";

  return item.collection || "all";
}

function inferEpisodeLabel(item) {
  const status = item.status || inferStatus(item);

  if (status === "upcoming") {
    return "Coming Soon";
  }

  return item.episodeLabel || item.episode_label || "EP 01";
}

function inferVideoUrl() {
  return "/videos/demo.mp4";
}

function normalizeContent(item) {
  const posterUrl = item.posterUrl || item.poster_url || item.imageUrl;
  const releaseYear = item.releaseYear || item.release_year || item.year;
  const status = item.status || inferStatus(item);
  const collection = item.collection || inferCollection(item);

  return {
    id: item.id,
    title: item.title,
    description: item.description,

    imageUrl: posterUrl,
    posterUrl,
    

    type: normalizeType(item.type),
    genre: item.genre,

    region: item.region || inferRegion(item),
    subtitle: item.subtitle || inferSubtitle(item),
    status,
    collection,

    year: releaseYear,
    releaseYear,

    rating: item.rating,
    isFeatured: item.isFeatured ?? item.is_featured ?? false,
    isTrending: item.isTrending ?? item.is_trending ?? false,

    createdAt: item.createdAt || item.created_at,
    updatedAt: item.updatedAt || item.updated_at,

    durationSec: item.durationSec || item.duration_sec || 3600,
    introStartSec: item.introStartSec || item.intro_start_sec || 8,
    introEndSec: item.introEndSec || item.intro_end_sec || 72,
    videoUrl: item.videoUrl || item.video_url || inferVideoUrl(item),

    ageRating: item.ageRating || item.age_rating || inferAgeRating(item),
    imageLabel: item.imageLabel || item.image_label || item.title || "poster",
    episodeLabel: inferEpisodeLabel({ ...item, status }),
  };
}



function normalizeContentList(data) {
  const contents = data.contents || data.content || data.items || [];
  return Array.isArray(contents) ? contents.map(normalizeContent) : [];
}

export async function getAllContent() {
  const data = await request("/content?limit=100");
  return normalizeContentList(data);
}

export async function getContentById(id) {
  const data = await request(`/content/${id}`);
  return normalizeContent(data.content || data);
}

export async function searchContent(query) {
  const data = await request(`/content/search?q=${encodeURIComponent(query)}`);
  return normalizeContentList(data);
}

export async function getContentByGenre(genre) {
  const data = await request(`/content/genre/${encodeURIComponent(genre)}`);
  return normalizeContentList(data);
}

export async function getTrendingContent(limit = 50) {
  const data = await request(`/content/trending?limit=${limit}`);
  return normalizeContentList(data);
}

export async function getFeaturedContent(limit = 50) {
  const data = await request(`/content/featured?limit=${limit}`);
  return normalizeContentList(data);
}

export async function getEpisodesByContent(contentId) {
  const data = await request(`/contents/${contentId}/episodes`);
  return data.episodes || [];
}

export async function createContent(payload) {
  const data = await request("/content", {
    method: "POST",
    body: JSON.stringify({
      title: payload.title,
      description: payload.description,
      poster_url: payload.posterUrl || payload.poster_url || payload.imageUrl,
      type: payload.type || "tv_series",
      genre: payload.genre || "drama",
      release_year: payload.releaseYear || payload.release_year || payload.year,
      rating: payload.rating || 0,
      is_featured: Boolean(payload.isFeatured || payload.is_featured),
      is_trending: Boolean(payload.isTrending || payload.is_trending),
    }),
  });

  return normalizeContent(data.content || data);
}

export async function getCollection(collection) {
  if (!collection || collection === "all") {
    return getAllContent();
  }

  if (collection === "trending") {
    return getTrendingContent(50);
  }

  if (collection === "featured") {
    return getFeaturedContent(50);
  }

  const all = await getAllContent();

  if (collection === "latest-update") {
    return all.filter((item) => item.collection === "latest-update");
  }

  if (collection === "top-k-drama") {
    return all.filter((item) => item.collection === "top-k-drama");
  }

  if (collection === "top-c-drama") {
    return all.filter((item) => item.collection === "top-c-drama");
  }

  if (collection === "top-j-drama") {
    return all.filter((item) => item.collection === "top-j-drama");
  }

  if (collection === "upcoming") {
    return all.filter((item) => item.collection === "upcoming");
  }

  return all.filter((item) => item.genre?.includes(collection));
}