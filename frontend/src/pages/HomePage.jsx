import { useEffect, useState } from "react";
import HeroCarousel from "../components/HeroCarousel";
import SectionCarousel from "../components/SectionCarousel";
import ContinueWatchingSection from "../components/ContinueWatchingSection";
import {
  getAllContent,
  getFeaturedContent,
  getTrendingContent,
} from "../api/contentApi";

function makeHeroSlides(items) {
  return items.slice(0, 4).map((item) => ({
    id: item.id,
    title: item.title,
    description: item.description,
    imageUrl: item.imageUrl || item.posterUrl,
    imageLabel: item.imageLabel || item.title,
    tag: item.collection || item.genre || "Featured",
    to: `/title/${item.id}`,
  }));
}

function getByRegion(items, region) {
  return items.filter((item) => item.region === region);
}

function getUpcoming(items) {
  return items.filter((item) => item.status === "upcoming");
}

export default function HomePage() {
  const [allItems, setAllItems] = useState([]);
  const [featured, setFeatured] = useState([]);
  const [trending, setTrending] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadHomeData() {
      try {
        const [all, featuredItems, trendingItems] = await Promise.all([
          getAllContent(),
          getFeaturedContent(10),
          getTrendingContent(10),
        ]);

        setAllItems(all);
        setFeatured(featuredItems.length ? featuredItems : all.slice(0, 4));
        setTrending(trendingItems.length ? trendingItems : all.slice(0, 10));
      } catch (error) {
        console.error("Failed to load home content:", error);
      } finally {
        setLoading(false);
      }
    }

    loadHomeData();
  }, []);

  const latestUpdates = [...allItems]
    .sort((a, b) => String(b.updatedAt || "").localeCompare(String(a.updatedAt || "")))
    .slice(0, 11);

  const topKDrama = getByRegion(allItems, "south-korea").slice(0, 11);
  const topCDrama = getByRegion(allItems, "china").slice(0, 11);
  const topJDrama = getByRegion(allItems, "japan").slice(0, 11);
  const upcoming = getUpcoming(allItems).slice(0, 7);

  const heroSlides = makeHeroSlides(featured.length ? featured : trending);

  if (loading) {
    return (
      <main className="page page--home">
        <div className="page__content">
          <section className="continue-empty">
            <h2>Loading Dorabox...</h2>
            <p>Backend is waking up. Naturally, it needs drama too.</p>
          </section>
        </div>
      </main>
    );
  }

  return (
    <main className="page page--home">
      <HeroCarousel slides={heroSlides} />

      <div className="page__content">
        <ContinueWatchingSection />

        <SectionCarousel
          title="Latest Update"
          items={latestUpdates}
          browseParams={{
            collection: "latest-update",
            sort: "last-update",
          }}
        />

        <SectionCarousel
          title="Top K-Drama"
          items={topKDrama}
          browseParams={{
            collection: "top-k-drama",
            region: "south-korea",
            sort: "popular",
          }}
        />

        <SectionCarousel
          title="Top C-Drama"
          items={topCDrama}
          browseParams={{
            collection: "top-c-drama",
            region: "china",
            sort: "popular",
          }}
        />

        <SectionCarousel
          title="Top J-Drama"
          items={topJDrama}
          browseParams={{
            collection: "top-j-drama",
            region: "japan",
            sort: "popular",
          }}
        />

        <SectionCarousel
          title="Upcoming"
          items={upcoming}
          browseParams={{
            collection: "upcoming",
            status: "upcoming",
            sort: "last-update",
          }}
        />
      </div>
    </main>
  );
}