import { useEffect, useState } from "react";
import DramaCard from "./DramaCard";
import { getSelectedProfileId } from "../api/profileApi";
import { getContinueWatching } from "../api/playbackApi";

export default function ContinueWatchingSection() {
  const [items, setItems] = useState([]);

  useEffect(() => {
    async function load() {
      try {
        const profileId = getSelectedProfileId();

        if (!profileId) {
          setItems([]);
          return;
        }

        const progress = await getContinueWatching(profileId);
        const safeProgress = progress.filter((entry) => Boolean(entry.item));

        setItems(safeProgress);
      } catch (error) {
        console.error("Failed to load continue watching:", error);
        setItems([]);
      }
    }

    load();

    const handler = () => load();

    window.addEventListener("storage", handler);
    window.addEventListener("auth-changed", handler);
    window.addEventListener("profile-changed", handler);

    return () => {
      window.removeEventListener("storage", handler);
      window.removeEventListener("auth-changed", handler);
      window.removeEventListener("profile-changed", handler);
    };
  }, []);

  if (!items.length) {
    return (
      <section className="continue-empty">
        <h2>Continue Watching</h2>
        <p>Start watching something and your progress will appear here.</p>
      </section>
    );
  }

  return (
    <section className="section-carousel">
      <div className="section-carousel__head">
        <h2 className="section-carousel__title">Continue Watching</h2>
      </div>

      <div className="section-carousel__track">
        {items.map((entry) => {
          const percent =
            entry.durationSec > 0
              ? Math.round((entry.currentPositionSec / entry.durationSec) * 100)
              : Math.round(entry.watchedPercent || 0);

          return (
            <div
              key={`${entry.profileId}-${entry.contentId}-${entry.episodeId || "episode"}`}
              className="continue-card-wrap"
            >
              <DramaCard item={entry.item} to={`/watch/${entry.contentId}`} />

              <div className="progress-line">
                <span
                  style={{
                    width: `${Math.min(100, percent)}%`,
                  }}
                />
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}