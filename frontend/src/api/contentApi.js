const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

async function request(path) {
  const response = await fetch(`${API_BASE_URL}${path}`);

  if (!response.ok) {
    throw new Error(`Content API error: ${response.status}`);
  }

  return response.json();
}

export async function getAllContent() {
  return request("/content");
}

export async function getContentById(id) {
  return request(`/content/${id}`);
}

export async function searchContent(query) {
  return request(`/content/search?q=${encodeURIComponent(query)}`);
}

export async function getContentByGenre(genre) {
  return request(`/content/genre/${encodeURIComponent(genre)}`);
}

export async function getTrendingContent() {
  return request("/content/trending");
}

export async function getFeaturedContent() {
  return request("/content/featured");
}

export async function getEpisodesByContent(contentId) {
  return request(`/content/${contentId}/episodes`);
}

export async function getCollection(collection) {
  if (!collection || collection === "all") {
    return getAllContent();
  }

  if (collection === "trending") {
    return getTrendingContent();
  }

  if (collection === "featured") {
    return getFeaturedContent();
  }

  return getContentByGenre(collection);
}
