import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getContentById } from "../api/contentApi";
import { getPlaceholderBackground } from "../data/catalog";

function formatDuration(seconds) {
  const minutes = Math.round((seconds || 3600) / 60);
  return `${minutes} min`;
}

function formatRegion(region) {
  const regions = {
    "south-korea": "South Korea",
    china: "China",
    japan: "Japan",
  };

  return regions[region] || region || "Unknown";
}

function formatStatus(status) {
  const statuses = {
    ongoing: "Ongoing",
    completed: "Completed",
    upcoming: "Upcoming",
  };

  return statuses[status] || status || "Ongoing";
}

export default function DetailsPage() {
  const { id } = useParams();

  const [item, setItem] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadDetails() {
      try {
        const data = await getContentById(id);
        setItem(data);
      } catch (error) {
        console.error("Failed to load details:", error);
        setItem(null);
      } finally {
        setLoading(false);
      }
    }

    loadDetails();
  }, [id]);

  if (loading) {
    return (
      <main className="not-found-page">
        <h1>Loading title...</h1>
      </main>
    );
  }

  if (!item) {
    return (
      <main className="not-found-page">
        <h1>Title not found</h1>
        <Link to="/" className="primary-btn primary-btn--inline">
          Go Home
        </Link>
      </main>
    );
  }

  const imageUrl = item.imageUrl || item.posterUrl;

  const heroStyle = imageUrl
    ? {
        backgroundImage: `linear-gradient(to right, rgba(0,0,0,0.92), rgba(0,0,0,0.62), rgba(0,0,0,0.25)), url(${imageUrl})`,
      }
    : {
        background: getPlaceholderBackground(item.imageLabel),
      };

  const posterStyle = imageUrl
    ? {
        backgroundImage: `url(${imageUrl})`,
      }
    : {
        background: "rgba(0, 0, 0, 0.28)",
      };

  return (
    <main className="details-page">
      <section className="details-hero" style={heroStyle}>
        <div className="details-hero__overlay">
          <div className="details-poster" style={posterStyle}>
            {!imageUrl && <span>{item.imageLabel}</span>}
          </div>

          <div className="details-info">
            <span className="hero__tag">{item.collection || item.genre}</span>

            <h1>{item.title}</h1>

            <div className="details-meta">
              <span>{item.year || item.releaseYear}</span>
              <span>{formatDuration(item.durationSec)}</span>
              <span>{item.ageRating || "13+"}</span>
              <span>{formatStatus(item.status)}</span>
              <span>{formatRegion(item.region)}</span>
            </div>

            <p>{item.description}</p>

            <div className="details-actions">
              <Link to={`/watch/${item.id}`} className="primary-btn primary-btn--inline">
                Watch Now
              </Link>

              <Link to="/profiles" className="secondary-btn">
                Change Profile
              </Link>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}