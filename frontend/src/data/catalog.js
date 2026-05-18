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
      imageUrl: `/images/${prefix}-${number}.jpg`,
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
    title: "Night Has Come",
    imageLabel: "featured-1",
    imageUrl: "/images/featured-1.jpg",
    subtitle:
      "A Korean mystery thriller where a school trip turns into a deadly mafia game.",
    linkedItemId: "kdrama-1",
    browseParams: {
      collection: "top-k-drama",
      region: "south-korea",
      sort: "popular",
    },
  },
  {
    id: "featured-2",
    title: "The First Frost",
    imageLabel: "featured-2",
    imageUrl: "/images/featured-2.jpg",
    subtitle:
      "A Chinese romance about old feelings, distance, and a reunion that refuses to stay simple.",
    linkedItemId: "cdrama-2",
    browseParams: {
      collection: "top-c-drama",
      region: "china",
      sort: "popular",
    },
  },
  {
    id: "featured-3",
    title: "Alice in Borderland",
    imageLabel: "featured-3",
    imageUrl: "/images/featured-3.jpg",
    subtitle:
      "A Japanese survival thriller where every game can become someone’s last mistake.",
    linkedItemId: "jdrama-1",
    browseParams: {
      collection: "top-j-drama",
      region: "japan",
      sort: "popular",
    },
  },
  {
    id: "featured-4",
    title: "The Remarried Empress",
    imageLabel: "featured-4",
    imageUrl: "/images/featured-4.jpg",
    subtitle:
      "An upcoming royal drama about betrayal, power, dignity, and a second marriage.",
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
}).map((item, index) => {
  const titles = [
    "Night Has Come",
    "Hidden Love",
    "First Love",
    "Twinkling Watermelon",
    "The First Frost",
    "Alice in Borderland",
    "Lovely Runner",
    "When I Fly Towards You",
    "Silent",
    "My Demon",
    "Reset",
  ];

  const descriptions = [
    "A Korean mystery thriller where a school trip becomes a deadly mafia game. Students are forced to accuse, suspect, and survive each other while the rules become more brutal.",
    "A Chinese romance about a girl who quietly falls for her brother’s friend and later meets him again as an adult. Soft, warm, and painfully good at slow-burn feelings.",
    "A Japanese romance drama about first love, memory, and timing. Two people who once loved each other are separated by life, then quietly pulled back into each other’s orbit.",
    "A Korean time-slip youth drama about music, family, friendship, and growing up. A boy travels to the past and meets his parents as teenagers, because apparently family trauma needed a soundtrack.",
    "A Chinese romance about two people reconnecting after years apart, carrying old feelings, misunderstandings, and unfinished history that refuses to stay buried.",
    "A Japanese survival thriller where people are trapped in an empty Tokyo and forced to play deadly games. Intelligence, trust, and survival instincts get tested very rudely.",
    "A Korean fantasy romance about a devoted fan who gets a second chance to save the idol she loves. Time travel, grief, comedy, and emotional damage included, naturally.",
    "A Chinese youth romance about school life, awkward feelings, and first love slowly becoming impossible to ignore. Very soft, very sweet, and dangerously good at making people nostalgic.",
    "A Japanese emotional drama about grief, communication, and love after silence changes everything. Quiet, restrained, and absolutely built to ruin your mood in a pretty way.",
    "A Korean fantasy romance about a demon who loses his powers after becoming involved with a wealthy heiress. Contract marriage, danger, and emotional denial do most of the work.",
    "A Chinese time-loop thriller where two strangers relive the same bus explosion again and again. They must uncover the truth before the loop destroys them emotionally and literally.",
  ];

  const years = [
    2023,
    2023,
    2022,
    2023,
    2025,
    2020,
    2024,
    2023,
    2022,
    2023,
    2022,
  ];

  const regions = [
    "south-korea",
    "china",
    "japan",
    "south-korea",
    "china",
    "japan",
    "south-korea",
    "china",
    "japan",
    "south-korea",
    "china",
  ];

  const subtitles = [
    "english",
    "english",
    "russian",
    "english",
    "russian",
    "english",
    "english",
    "russian",
    "english",
    "russian",
    "english",
  ];

  const statuses = [
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
  ];

  const ageRatings = [
    "16+",
    "13+",
    "16+",
    "13+",
    "16+",
    "18+",
    "13+",
    "13+",
    "13+",
    "16+",
    "16+",
  ];

  return {
    ...item,
    title: titles[index] || item.title,
    description: descriptions[index] || item.description,
    year: years[index] || item.year,
    region: regions[index] || item.region,
    subtitle: subtitles[index] || item.subtitle,
    status: statuses[index] || item.status,
    ageRating: ageRatings[index] || item.ageRating,
    imageUrl: `/images/latest-update-${index + 1}.jpg`,
  };
});

export const topKDrama = buildItems("kdrama", 11, {
  collection: "top-k-drama",
  region: "south-korea",
}).map((item, index) => {
  const titles = [
    "Night Has Come",
    "Twinkling Watermelon",
    "True Beauty",
    "Lovely Runner",
    "Head Over Heels",
    "My Demon",
    "Good Boy",
    "Running Man",
    "Queen of Tears",
    "Business Proposal",
    "Twenty Five Twenty One",
  ];

  const descriptions = [
    "A Korean mystery thriller where a school trip turns into a deadly mafia game. The students must survive paranoia, accusations, and the fear that anyone around them could be hiding something.",
    "A time-slip youth drama about music, family, friendship, and the messy ache of growing up. A boy travels to the past and meets his parents as teenagers.",
    "A romantic school comedy about beauty standards, insecurity, and learning to be seen beyond appearances. The heroine hides behind makeup while trying to survive high school life.",
    "A romance fantasy about a devoted fan who gets a second chance to save the idol she loves. Time travel, grief, comedy, and emotional damage included, naturally.",
    "A supernatural youth romance about fate, curses, and a girl whose life becomes tied to dangerous secrets. Soft visuals, dramatic longing, and suspiciously pretty people everywhere.",
    "A fantasy romance about a cold demon who loses his powers after becoming involved with a wealthy heiress. Contract marriage, danger, and emotional denial do most of the heavy lifting.",
    "An action comedy drama about people chasing justice in ridiculous and dangerous situations. Stylish, chaotic, and built for characters with too much confidence.",
    "A long-running Korean variety show where celebrities complete missions, betray each other casually, and somehow turn games into emotional warfare.",
    "A melodrama romance about a married couple facing crisis, pride, illness, and the fragile mess of loving someone after resentment has built walls everywhere.",
    "A workplace romantic comedy about a fake blind date that spirals into a contract relationship with a CEO. Corporate nonsense, fake dating, and exactly the kind of chaos humans apparently enjoy.",
    "A nostalgic youth romance about ambition, friendship, fencing, first love, and the kind of growing up that hurts because nobody gets to stay seventeen forever.",
  ];

  const years = [
    2023,
    2023,
    2020,
    2024,
    2025,
    2023,
    2025,
    2010,
    2024,
    2022,
    2022,
  ];

  const ageRatings = [
    "16+",
    "13+",
    "13+",
    "13+",
    "16+",
    "16+",
    "16+",
    "13+",
    "16+",
    "13+",
    "13+",
  ];

  const statuses = [
    "completed",
    "completed",
    "completed",
    "completed",
    "ongoing",
    "completed",
    "ongoing",
    "ongoing",
    "completed",
    "completed",
    "completed",
  ];

  return {
    ...item,
    title: titles[index] || item.title,
    description: descriptions[index] || item.description,
    year: years[index] || item.year,
    ageRating: ageRatings[index] || item.ageRating,
    status: statuses[index] || item.status,
    imageUrl: `/images/kdrama-${index + 1}.jpg`,
  };
});

export const topCDrama = buildItems("cdrama", 11, {
  collection: "top-c-drama",
  region: "china",
}).map((item, index) => {
  const titles = [
    "Hidden Love",
    "The First Frost",
    "Love Between Fairy and Devil",
    "When I Fly Towards You",
    "Falling Into Your Smile",
    "You Are My Glory",
    "Reset",
    "Go Ahead",
    "The Long Ballad",
    "Meet Yourself",
    "Amidst a Snowstorm of Love",
  ];

  const descriptions = [
    "A soft Chinese romance about a girl who quietly falls for her brother’s friend and later meets him again as an adult. Gentle, warm, and extremely committed to emotional slow burn.",
    "A romance drama about two people reconnecting after years apart, carrying old feelings, misunderstandings, and the kind of unfinished history that refuses to stay buried.",
    "A fantasy romance about a feared demon lord and a fairy whose lives become magically connected. Their relationship begins with conflict and grows into loyalty, sacrifice, and very dramatic eye contact.",
    "A youth romance about a quiet girl and the bright, kind boy who slowly becomes important to her. School life, awkward affection, and first love doing its usual emotional vandalism.",
    "A romantic comedy set in the world of esports, where a talented female gamer joins a professional team and clashes with a cold team captain. Gaming, rivalry, romance, and men learning basic emotional literacy.",
    "A romance between a popular actress and an aerospace engineer who reconnect through an online game. It mixes celebrity life, ambition, science, and a very polished kind of longing.",
    "A time-loop thriller where two strangers repeatedly wake up on a bus before an explosion. They must uncover the truth, change the outcome, and survive the same nightmare again and again.",
    "A family drama about three unrelated children who grow up as siblings after being hurt by their own families. Warm, painful, and very good at reminding people that family can be built, not just inherited.",
    "A historical drama about revenge, war, politics, and a princess trying to survive after her world collapses. Strategy, loyalty, grief, and grand costumes doing most of the heavy lifting.",
    "A healing slice-of-life romance about a woman who leaves city life behind and moves to a quiet village. It focuses on grief, recovery, community, and learning how to breathe again.",
    "A winter romance about a former billiards genius and a professional player who meet abroad. Quiet attraction, snow-covered streets, ambition, and the kind of romance that moves slowly because apparently feelings need paperwork.",
  ];

  const years = [
    2023,
    2025,
    2022,
    2023,
    2021,
    2021,
    2022,
    2020,
    2021,
    2023,
    2024,
  ];

  const ageRatings = [
    "13+",
    "16+",
    "13+",
    "13+",
    "13+",
    "13+",
    "16+",
    "13+",
    "16+",
    "13+",
    "13+",
  ];

  const statuses = [
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
  ];

  return {
    ...item,
    title: titles[index] || item.title,
    description: descriptions[index] || item.description,
    year: years[index] || item.year,
    ageRating: ageRatings[index] || item.ageRating,
    status: statuses[index] || item.status,
    imageUrl: `/images/cdrama-${index + 1}.jpg`,
  };
});

export const topJDrama = buildItems("jdrama", 11, {
  collection: "top-j-drama",
  region: "japan",
}).map((item, index) => {
  const titles = [
    "Alice in Borderland",
    "First Love",
    "Silent",
    "My Love Mix-Up",
    "Good Morning Call",
    "Brush Up Life",
    "From Me to You",
    "Nodame Cantabile",
    "One Litre of Tears",
    "Trillion Game",
    "Pending Train",
  ];

  const descriptions = [
    "A Japanese survival thriller where people are trapped in an empty Tokyo and forced to play deadly games. Every choice becomes a test of intelligence, trust, and the will to stay alive.",
    "A romantic drama about first love, memory, and timing. Two people who once loved each other deeply are pulled apart by life, then quietly drawn back together years later.",
    "A quiet emotional drama about love, loss, and communication after a young woman reunites with her former boyfriend, who has lost his hearing. It is soft, painful, and very good at emotional sabotage.",
    "A school romance comedy about mistaken crushes, awkward confessions, and two boys slowly figuring out what their feelings mean. Sweet, chaotic, and far more charming than teenagers deserve to be.",
    "A romantic comedy about a girl who moves into the house of a popular boy from her school after an accident leaves her temporarily homeless. Naturally, feelings start causing structural damage.",
    "A clever life-retry drama about a woman who dies and gets a chance to live her life again from childhood. Instead of becoming some grand chosen one, she tries to fix the small mistakes that shaped everything.",
    "A gentle romance about a shy girl whose feelings slowly reach the boy she admires. A soft story about misunderstanding, kindness, and learning to connect with people.",
    "A music-centered romantic comedy about a chaotic piano prodigy and a perfectionist music student. Weird, funny, and full of classical music drama, because apparently even orchestras need emotional warfare.",
    "A heartbreaking drama based on the life of a girl diagnosed with a degenerative illness. It follows her struggle to keep living with dignity, love, and hope while her body slowly betrays her.",
    "A business drama about two ambitious men trying to build a trillion-dollar company from nothing. Strategy, ego, risk, and the usual capitalist circus in a shinier outfit.",
    "A sci-fi suspense drama about passengers trapped on a train that suddenly arrives in a devastated future. They must survive uncertainty, fear, and each other while trying to understand what happened.",
  ];

  const years = [
    2020,
    2022,
    2022,
    2021,
    2016,
    2023,
    2023,
    2006,
    2005,
    2023,
    2023,
  ];

  const ageRatings = [
    "18+",
    "16+",
    "13+",
    "13+",
    "13+",
    "13+",
    "13+",
    "13+",
    "13+",
    "13+",
    "16+",
  ];

  const statuses = [
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
    "completed",
  ];

  return {
    ...item,
    title: titles[index] || item.title,
    description: descriptions[index] || item.description,
    year: years[index] || item.year,
    ageRating: ageRatings[index] || item.ageRating,
    status: statuses[index] || item.status,
    imageUrl: `/images/jdrama-${index + 1}.jpg`,
  };
});

export const upcoming = buildItems("upcoming", 7, {
  collection: "upcoming",
  status: "upcoming",
  region: "south-korea",
}).map((item, index) => {
  const titles = [
    "The Remarried Empress",
    "Mercy for None",
    "The Murky Stream",
    "Moonlight Mystique",
    "Northward",
    "Glass Heart",
    "Romantics Anonymous",
  ];

  const descriptions = [
    "An upcoming Korean royal drama based on the popular web novel about an empress who chooses dignity, power, and a second marriage after betrayal. Palace politics, emotional warfare, and very expensive suffering.",
    "A Korean crime action drama about a former gangster who cuts ties with the underworld, only to be dragged back into violence after a family tragedy. Revenge arrives, because apparently peace was too much to ask.",
    "A Korean historical action drama set during the Joseon era, following people fighting through corruption, survival, and dangerous waters. Expect politics, betrayal, and men making terrible decisions in nice costumes.",
    "A Chinese fantasy romance about gods, demons, identity, and forbidden feelings. It mixes magical conflict with devotion, destiny, and the usual fantasy-romance habit of ruining everyone emotionally.",
    "A Chinese family and coming-of-age drama following people whose lives change through migration, distance, and the search for belonging. Softer, more grounded, and probably ready to hurt quietly.",
    "A Japanese music drama about young musicians chasing their dreams, breaking apart, and trying to find their way back to themselves. Loud feelings, louder guitars, and emotional damage with stage lighting.",
    "A Japanese romantic drama about two people whose carefully controlled lives are disrupted by unexpected love. Sweet, awkward, and exactly the kind of thing humans call healing while crying into tea.",
  ];

  const years = [
    2026,
    2025,
    2025,
    2025,
    2025,
    2025,
    2025,
  ];

  const regions = [
    "south-korea",
    "south-korea",
    "south-korea",
    "china",
    "china",
    "japan",
    "japan",
  ];

  const subtitles = [
    "english",
    "english",
    "russian",
    "english",
    "russian",
    "english",
    "russian",
  ];

  const ageRatings = [
    "16+",
    "18+",
    "16+",
    "13+",
    "13+",
    "13+",
    "13+",
  ];

  const labels = [
    "Coming 2026",
    "Coming Soon",
    "Coming Soon",
    "Coming Soon",
    "Coming Soon",
    "Coming Soon",
    "Coming Soon",
  ];

  return {
    ...item,
    title: titles[index] || item.title,
    description: descriptions[index] || item.description,
    year: years[index] || item.year,
    region: regions[index] || item.region,
    subtitle: subtitles[index] || item.subtitle,
    ageRating: ageRatings[index] || item.ageRating,
    status: "upcoming",
    episodeLabel: labels[index] || "Coming Soon",
    imageUrl: `/images/upcoming-${index + 1}.jpg`,
  };
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