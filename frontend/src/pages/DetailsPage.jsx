import { Link, useParams } from "react-router-dom";
import { getItemById, getPlaceholderBackground } from "../data/catalog";

function formatDuration(seconds) {
  const minutes = Math.round(seconds / 60);
  return `${minutes} min`;
}

function formatRegion(region) {
  const regions = {
    "south-korea": "South Korea",
    china: "China",
    japan: "Japan",
  };

  return regions[region] || region;
}

function formatStatus(status) {
  const statuses = {
    ongoing: "Ongoing",
    completed: "Completed",
    upcoming: "Upcoming",
  };

  return statuses[status] || status;
}

export default function DetailsPage() {
  const { id } = useParams();
  const item = getItemById(id);

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

  const heroStyle = item.imageUrl
    ? {
        backgroundImage: `linear-gradient(to right, rgba(0,0,0,0.92), rgba(0,0,0,0.62), rgba(0,0,0,0.25)), url(${item.imageUrl})`,
      }
    : {
        background: getPlaceholderBackground(item.imageLabel),
      };

  const posterStyle = item.imageUrl
    ? {
        backgroundImage: `url(${item.imageUrl})`,
      }
    : {
        background: "rgba(0, 0, 0, 0.28)",
      };

  return (
    <main className="details-page">
      <section className="details-hero" style={heroStyle}>
        <div className="details-hero__overlay">
          <div className="details-poster" style={posterStyle}>
            {!item.imageUrl && <span>{item.imageLabel}</span>}
          </div>

          <div className="details-info">
            <span className="hero__tag">{item.collection}</span>

            <h1>{item.title}</h1>

            <div className="details-meta">
              <span>{item.year}</span>
              <span>{formatDuration(item.durationSec)}</span>
              <span>{item.ageRating}</span>
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