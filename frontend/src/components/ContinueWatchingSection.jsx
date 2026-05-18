import { useEffect, useState } from "react";
import DramaCard from "./DramaCard";
import { getSelectedProfileId } from "../api/profileApi";
import { getContinueWatching } from "../api/playbackApi";

export default function ContinueWatchingSection() {
  const [items, setItems] = useState([]);

  useEffect(() => {
    async function load() {
      const profileId = getSelectedProfileId();
      const progress = await getContinueWatching(profileId);
      setItems(progress);
    }

    load();

    const handler = () => load();
    window.addEventListener("storage", handler);

    return () => window.removeEventListener("storage", handler);
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
        {items.map((entry) => (
          <div key={`${entry.profileId}-${entry.contentId}`} className="continue-card-wrap">
            <DramaCard item={entry.item} to={`/watch/${entry.contentId}`} />

            <div className="progress-line">
              <span
                style={{
                  width: `${Math.min(
                    100,
                    Math.round((entry.currentPositionSec / entry.durationSec) * 100)
                  )}%`,
                }}
              />
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}