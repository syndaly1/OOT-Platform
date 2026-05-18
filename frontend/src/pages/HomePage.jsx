import HeroCarousel from "../components/HeroCarousel";
import SectionCarousel from "../components/SectionCarousel";
import ContinueWatchingSection from "../components/ContinueWatchingSection";
import {
  heroSlides,
  latestUpdates,
  topKDrama,
  topCDrama,
  topJDrama,
  upcoming,
} from "../data/catalog";

export default function HomePage() {
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