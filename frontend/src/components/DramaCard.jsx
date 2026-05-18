import { Link } from "react-router-dom";
import { getPlaceholderBackground } from "../data/catalog";

export default function DramaCard({ item, to }) {
  const cardLink = to || `/title/${item.id}`;

  const mediaStyle = item.imageUrl
    ? {
        backgroundImage: `linear-gradient(to bottom, rgba(0,0,0,0.05), rgba(0,0,0,0.55)), url(${item.imageUrl})`,
      }
    : {
        background: getPlaceholderBackground(item.imageLabel),
      };

  return (
    <Link to={cardLink} className="drama-card">
      <div className="drama-card__media" style={mediaStyle}>
        {!item.imageUrl && (
          <div className="drama-card__placeholder-label">
            {item.imageLabel}
          </div>
        )}

        <span className="drama-card__badge">{item.episodeLabel}</span>
      </div>

      <div className="drama-card__body">
        <h3 className="drama-card__title">{item.title}</h3>
      </div>
    </Link>
  );
}