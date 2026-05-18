import { useRef } from "react";
import { Link } from "react-router-dom";
import DramaCard from "./DramaCard";
import { createBrowseLink } from "../utils/browse";

export default function SectionCarousel({ title, items, browseParams }) {
  const trackRef = useRef(null);
  const sectionLink = createBrowseLink(browseParams);

  const scrollByAmount = (amount) => {
    if (!trackRef.current) return;

    trackRef.current.scrollBy({
      left: amount,
      behavior: "smooth",
    });
  };

  return (
    <section className="section-carousel">
      <div className="section-carousel__head">
        <Link to={sectionLink} className="section-carousel__title-link">
          <h2 className="section-carousel__title">{title}</h2>
          <span className="section-carousel__title-arrow">›</span>
        </Link>

        <div className="section-carousel__controls">
          <button onClick={() => scrollByAmount(-640)} className="circle-btn">
            ‹
          </button>
          <button onClick={() => scrollByAmount(640)} className="circle-btn">
            ›
          </button>
        </div>
      </div>

      <div className="section-carousel__track" ref={trackRef}>
        {items.map((item) => (
          <DramaCard key={item.id} item={item} to={`/title/${item.id}`} />
        ))}
      </div>
    </section>
  );
}