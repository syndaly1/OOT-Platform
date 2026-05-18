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

  const seenTitles = new Set();

  result = result.filter((item) => {
    const normalizedTitle = item.title.toLowerCase().trim();

    if (seenTitles.has(normalizedTitle)) {
      return false;
    }

    seenTitles.add(normalizedTitle);
    return true;
  });

  return result;
}, [baseItems, filters]);
}