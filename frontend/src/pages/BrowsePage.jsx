import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import DramaCard from "../components/DramaCard";
import FilterChips from "../components/FilterChips";
import {
  getAllContent,
  getFeaturedContent,
  getTrendingContent,
} from "../api/contentApi";

const filterGroups = {
  type: [
    { label: "All", value: "all" },
    { label: "TV Series", value: "tv-series" },
    { label: "Movie", value: "movie" },
  ],
  subtitle: [
    { label: "All Subtitles", value: "all" },
    { label: "English", value: "english" },
    { label: "Russian", value: "russian" },
  ],
  region: [
    { label: "All Regions", value: "all" },
    { label: "South Korea", value: "south-korea" },
    { label: "Japan", value: "japan" },
    { label: "China", value: "china" },
  ],
  status: [
    { label: "All", value: "all" },
    { label: "Ongoing", value: "ongoing" },
    { label: "Completed", value: "completed" },
    { label: "Upcoming", value: "upcoming" },
  ],
  sort: [
    { label: "Popular", value: "popular" },
    { label: "Last Update", value: "last-update" },
  ],
};

const collectionDefaults = {
  "latest-update": {
    type: "all",
    subtitle: "all",
    region: "all",
    status: "all",
    sort: "last-update",
  },
  "top-k-drama": {
    type: "all",
    subtitle: "all",
    region: "south-korea",
    status: "all",
    sort: "popular",
  },
  "top-c-drama": {
    type: "all",
    subtitle: "all",
    region: "china",
    status: "all",
    sort: "popular",
  },
  "top-j-drama": {
    type: "all",
    subtitle: "all",
    region: "japan",
    status: "all",
    sort: "popular",
  },
  upcoming: {
    type: "all",
    subtitle: "all",
    region: "all",
    status: "upcoming",
    sort: "last-update",
  },
  featured: {
    type: "all",
    subtitle: "all",
    region: "all",
    status: "all",
    sort: "popular",
  },
  trending: {
    type: "all",
    subtitle: "all",
    region: "all",
    status: "all",
    sort: "popular",
  },
  all: {
    type: "all",
    subtitle: "all",
    region: "all",
    status: "all",
    sort: "popular",
  },
};

const collectionLabels = {
  "latest-update": "Latest Update",
  "top-k-drama": "Top K-Drama",
  "top-c-drama": "Top C-Drama",
  "top-j-drama": "Top J-Drama",
  upcoming: "Upcoming",
  featured: "Featured",
  trending: "Trending",
  all: "All Dramas",
};

function normalizeType(type) {
  if (!type) return "tv-series";
  if (type === "tv_series") return "tv-series";
  return type;
}

function normalizeRegion(region) {
  if (!region) return "all";

  const lower = String(region).toLowerCase();

  if (lower.includes("korea")) return "south-korea";
  if (lower.includes("china")) return "china";
  if (lower.includes("japan")) return "japan";

  return lower;
}

function inferCollection(item) {
  if (item.collection) return item.collection;
  if (item.status === "upcoming") return "upcoming";
  return "all";
}

function enrichItem(item) {
  return {
    ...item,
    type: normalizeType(item.type),
    region: normalizeRegion(item.region),
    subtitle: item.subtitle || "english",
    status: item.status || "ongoing",
    collection: inferCollection(item),
  };
}

export default function BrowsePage() {
  const [searchParams, setSearchParams] = useSearchParams();

  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);

  const collection = searchParams.get("collection") || "all";
  const defaults = collectionDefaults[collection] || collectionDefaults.all;

  const filters = {
    type: searchParams.get("type") || defaults.type,
    subtitle: searchParams.get("subtitle") || defaults.subtitle,
    region: searchParams.get("region") || defaults.region,
    status: searchParams.get("status") || defaults.status,
    sort: searchParams.get("sort") || defaults.sort,
  };

  useEffect(() => {
    async function loadItems() {
      try {
        let data;

        if (collection === "featured") {
          data = await getFeaturedContent(50);
        } else if (collection === "trending") {
          data = await getTrendingContent(50);
        } else {
          data = await getAllContent();
        }

        setItems(data.map(enrichItem));
      } catch (error) {
        console.error("Failed to load browse content:", error);
        setItems([]);
      } finally {
        setLoading(false);
      }
    }

    loadItems();
  }, [collection]);

  const updateParam = (key, value) => {
    const next = new URLSearchParams(searchParams);
    next.set("collection", collection);
    next.set(key, value);
    setSearchParams(next);
  };

  const baseItems = useMemo(() => {
    if (collection === "all") {
      return items;
    }

    if (collection === "latest-update") {
      return items;
    }

    if (collection === "top-k-drama") {
      return items.filter((item) => item.region === "south-korea");
    }

    if (collection === "top-c-drama") {
      return items.filter((item) => item.region === "china");
    }

    if (collection === "top-j-drama") {
      return items.filter((item) => item.region === "japan");
    }

    if (collection === "upcoming") {
      return items.filter((item) => item.status === "upcoming");
    }

    if (collection === "featured") {
      return items.filter((item) => item.isFeatured);
    }

    if (collection === "trending") {
      return items.filter((item) => item.isTrending);
    }

    return items;
  }, [collection, items]);

  const filteredItems = useMemo(() => {
    let result = [...baseItems];

    if (filters.type !== "all") {
      result = result.filter((item) => item.type === filters.type);
    }

    if (filters.subtitle !== "all") {
      result = result.filter((item) => item.subtitle === filters.subtitle);
    }

    if (filters.region !== "all") {
      result = result.filter((item) => item.region === filters.region);
    }

    if (filters.status !== "all") {
      result = result.filter((item) => item.status === filters.status);
    }

    if (filters.sort === "popular") {
      result.sort((a, b) => Number(b.rating || 0) - Number(a.rating || 0));
    }

    if (filters.sort === "last-update") {
      result.sort((a, b) =>
        String(b.updatedAt || "").localeCompare(String(a.updatedAt || ""))
      );
    }

    const seenTitles = new Set();

    return result.filter((item) => {
      const normalizedTitle = item.title.toLowerCase().trim();

      if (seenTitles.has(normalizedTitle)) {
        return false;
      }

      seenTitles.add(normalizedTitle);
      return true;
    });
  }, [baseItems, filters]);

  return (
    <main className="page browse-page">
      <section className="browse-hero">
        <p className="eyebrow">Browse</p>
        <h1>{collectionLabels[collection] || "All Dramas"}</h1>
        <p>
          Filter dramas by type, subtitles, region, release status, and sorting.
        </p>
      </section>

      <section className="browse-filters">
        <FilterChips
  label="Type"
  options={filterGroups.type}
  value={filters.type}
  onChange={(value) => updateParam("type", value)}
/>

<FilterChips
  label="Subtitles"
  options={filterGroups.subtitle}
  value={filters.subtitle}
  onChange={(value) => updateParam("subtitle", value)}
/>

<FilterChips
  label="Region"
  options={filterGroups.region}
  value={filters.region}
  onChange={(value) => updateParam("region", value)}
/>

<FilterChips
  label="Status"
  options={filterGroups.status}
  value={filters.status}
  onChange={(value) => updateParam("status", value)}
/>

<FilterChips
  label="Sort"
  options={filterGroups.sort}
  value={filters.sort}
  onChange={(value) => updateParam("sort", value)}
/>
      </section>

      <section className="browse-results">
        {loading ? (
          <div className="continue-empty">
            <h2>Loading titles...</h2>
            <p>The backend is thinking. Tragic, but necessary.</p>
          </div>
        ) : filteredItems.length ? (
          <div className="browse-grid">
            {filteredItems.map((item) => (
              <DramaCard key={item.id} item={item} to={`/title/${item.id}`} />
            ))}
          </div>
        ) : (
          <div className="continue-empty">
            <h2>No titles found</h2>
            <p>
              Try changing filters or add more content through the backend.
            </p>
          </div>
        )}
      </section>
    </main>
  );
}