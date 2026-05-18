export const collectionLabels = {
  "latest-update": "Latest Update",
  "top-k-drama": "Top K-Drama",
  "top-c-drama": "Top C-Drama",
  "top-j-drama": "Top J-Drama",
  upcoming: "Upcoming",
  all: "Browse Library",
};

const placeholderGradients = [
  "linear-gradient(135deg, #16324f, #274c77)",
  "linear-gradient(135deg, #3a0ca3, #7209b7)",
  "linear-gradient(135deg, #1f4037, #99f2c8)",
  "linear-gradient(135deg, #4b134f, #c94b4b)",
  "linear-gradient(135deg, #232526, #414345)",
  "linear-gradient(135deg, #2c3e50, #4ca1af)",
  "linear-gradient(135deg, #42275a, #734b6d)",
  "linear-gradient(135deg, #134e5e, #71b280)",
];

export function getPlaceholderBackground(seed = "") {
  const total = seed.split("").reduce((sum, ch) => sum + ch.charCodeAt(0), 0);
  return placeholderGradients[total % placeholderGradients.length];
}

function buildItems(prefix, count, overrides = {}) {
  const typeCycle = ["tv-series", "tv-series", "movie"];
  const subtitleCycle = ["english", "russian"];
  const statusCycle = ["ongoing", "completed", "upcoming"];

  return Array.from({ length: count }, (_, index) => {
    const number = index + 1;
    const isUpcoming = overrides.status === "upcoming";

    return {
      id: `${prefix}-${number}`,
      title: `${prefix}-${number}`,
      imageLabel: `${prefix}-${number}`,
      episodeLabel: isUpcoming ? "Coming Soon" : `EP ${String(number).padStart(2, "0")}`,
      description:
        "A placeholder drama description for the streaming project. Replace this with real metadata from Content Service later.",
      year: 2026,
      durationSec: 2400 + number * 120,
      introStartSec: 8,
      introEndSec: 72,
      videoUrl:
        "https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4",
      type: overrides.type ?? typeCycle[index % typeCycle.length],
      subtitle: overrides.subtitle ?? subtitleCycle[index % subtitleCycle.length],
      region: overrides.region,
      status: overrides.status ?? statusCycle[index % statusCycle.length],
      collection: overrides.collection ?? prefix,
      ageRating: number % 3 === 0 ? "18+" : number % 2 === 0 ? "16+" : "13+",
    };
  });
}

export const heroSlides = [
  {
    id: "featured-1",
    title: "latest-update-1",
    imageLabel: "featured-1",
    subtitle: "Fresh releases and recently updated series.",
    linkedItemId: "latest-update-1",
    browseParams: {
      collection: "latest-update",
      sort: "last-update",
    },
  },
  {
    id: "featured-2",
    title: "kdrama-1",
    imageLabel: "featured-2",
    subtitle: "Popular Korean titles and fan favorites.",
    linkedItemId: "kdrama-1",
    browseParams: {
      collection: "top-k-drama",
      region: "south-korea",
      sort: "popular",
    },
  },
  {
    id: "featured-3",
    title: "cdrama-1",
    imageLabel: "featured-3",
    subtitle: "Explore the latest Chinese drama placeholders.",
    linkedItemId: "cdrama-1",
    browseParams: {
      collection: "top-c-drama",
      region: "china",
      sort: "popular",
    },
  },
  {
    id: "featured-4",
    title: "upcoming-1",
    imageLabel: "featured-4",
    subtitle: "See upcoming titles waiting to premiere.",
    linkedItemId: "upcoming-1",
    browseParams: {
      collection: "upcoming",
      status: "upcoming",
      sort: "last-update",
    },
  },
];

export const latestUpdates = buildItems("latest-update", 11, {
  collection: "latest-update",
  region: "south-korea",
});

export const topKDrama = buildItems("kdrama", 11, {
  collection: "top-k-drama",
  region: "south-korea",
});

export const topCDrama = buildItems("cdrama", 11, {
  collection: "top-c-drama",
  region: "china",
});

export const topJDrama = buildItems("jdrama", 11, {
  collection: "top-j-drama",
  region: "japan",
});

export const upcoming = buildItems("upcoming", 7, {
  collection: "upcoming",
  status: "upcoming",
  region: "south-korea",
});

export const collectionDataMap = {
  "latest-update": latestUpdates,
  "top-k-drama": topKDrama,
  "top-c-drama": topCDrama,
  "top-j-drama": topJDrama,
  upcoming,
};

export const allItems = [
  ...latestUpdates,
  ...topKDrama,
  ...topCDrama,
  ...topJDrama,
  ...upcoming,
];

export function getItemById(id) {
  return allItems.find((item) => item.id === id);
}