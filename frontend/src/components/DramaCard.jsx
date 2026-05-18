import { Link } from "react-router-dom";
import { getPlaceholderBackground } from "../data/catalog";

export default function DramaCard({ item, to }) {
  const cardLink = to || `/title/${item.id}`;

  return (
    <Link to={cardLink} className="drama-card">
      <div
        className="drama-card__media"
        style={{ background: getPlaceholderBackground(item.imageLabel) }}
      >
        <div className="drama-card__placeholder-label">{item.imageLabel}</div>
        <span className="drama-card__badge">{item.episodeLabel}</span>
      </div>

      <div className="drama-card__body">
        <h3 className="drama-card__title">{item.title}</h3>
      </div>
    </Link>
  );
}