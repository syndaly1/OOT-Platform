import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import DramaCard from "../components/DramaCard";
import FilterChips from "../components/FilterChips";
import {
  allItems,
  collectionDataMap,
  collectionLabels,
} from "../data/catalog";

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
  all: {
    type: "all",
    subtitle: "all",
    region: "all",
    status: "all",
    sort: "popular",
  },
};

export default function BrowsePage() {
  const [searchParams, setSearchParams] = useSearchParams();

  const collection = searchParams.get("collection") || "all";
  const defaults = collectionDefaults[collection] || collectionDefaults.all;

  const filters = {
    type: searchParams.get("type") || defaults.type,
    subtitle: searchParams.get("subtitle") || defaults.subtitle,
    region: searchParams.get("region") || defaults.region,
    status: searchParams.get("status") || defaults.status,
    sort: searchParams.get("sort") || defaults.sort,
  };

  const updateParam = (key, value) => {
    const next = new URLSearchParams(searchParams);
    next.set("collection", collection);
    next.set(key, value);
    setSearchParams(next);
  };

  const baseItems = useMemo(() => {
    if (collection === "all") return allItems;
    return collectionDataMap[collection] || allItems;
  }, [collection]);

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
      result.sort((a, b) => a.title.localeCompare(b.title));
    }

    if (filters.sort === "last-update") {
      result.sort((a, b) => b.id.localeCompare(a.id));
    }

    return result;
  }, [baseItems, filters]);

  return (
    <main className="page page--browse">
      <div className="browse-header">
        <h1>{collectionLabels[collection] || "Browse Library"}</h1>
        <p>{filteredItems.length} items found</p>
      </div>

      <div className="browse-filters">
        <FilterChips
          options={filterGroups.type}
          selectedValue={filters.type}
          onSelect={(value) => updateParam("type", value)}
        />

        <FilterChips
          options={filterGroups.subtitle}
          selectedValue={filters.subtitle}
          onSelect={(value) => updateParam("subtitle", value)}
        />

        <FilterChips
          options={filterGroups.region}
          selectedValue={filters.region}
          onSelect={(value) => updateParam("region", value)}
        />

        <FilterChips
          options={filterGroups.status}
          selectedValue={filters.status}
          onSelect={(value) => updateParam("status", value)}
        />

        <FilterChips
          options={filterGroups.sort}
          selectedValue={filters.sort}
          onSelect={(value) => updateParam("sort", value)}
        />
      </div>

      <div className="browse-grid">
        {filteredItems.map((item) => (
          <DramaCard key={item.id} item={item} to={`/title/${item.id}`} />
        ))}
      </div>
    </main>
  );
}