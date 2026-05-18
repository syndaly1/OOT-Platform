import { Link, useParams } from "react-router-dom";
import { getItemById, getPlaceholderBackground } from "../data/catalog";

function formatDuration(seconds) {
  const minutes = Math.round(seconds / 60);
  return `${minutes} min`;
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

  return (
    <main className="details-page">
      <section
        className="details-hero"
        style={{ background: getPlaceholderBackground(item.imageLabel) }}
      >
        <div className="details-hero__overlay">
          <div className="details-poster">
            <span>{item.imageLabel}</span>
          </div>

          <div className="details-info">
            <span className="hero__tag">{item.collection}</span>

            <h1>{item.title}</h1>

            <div className="details-meta">
              <span>{item.year}</span>
              <span>{formatDuration(item.durationSec)}</span>
              <span>{item.ageRating}</span>
              <span>{item.status}</span>
              <span>{item.region}</span>
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