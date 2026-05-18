import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getPlaceholderBackground } from "../data/catalog";

export default function HeroCarousel({ slides = [] }) {
  const [activeIndex, setActiveIndex] = useState(0);

  useEffect(() => {
    if (!slides.length) return;

    const timer = setInterval(() => {
      setActiveIndex((prev) => (prev + 1) % slides.length);
    }, 5000);

    return () => clearInterval(timer);
  }, [slides]);

  const prevSlide = () => {
    setActiveIndex((prev) => (prev === 0 ? slides.length - 1 : prev - 1));
  };

  const nextSlide = () => {
    setActiveIndex((prev) => (prev + 1) % slides.length);
  };

  const current = slides[activeIndex];

  if (!current) return null;

  return (
    <section className="hero">
      <button className="hero__arrow hero__arrow--left" onClick={prevSlide}>
        ‹
      </button>

      <div
  className="hero__slide"
  style={
    current.imageUrl
      ? {
          backgroundImage: `linear-gradient(to right, rgba(0,0,0,0.78), rgba(0,0,0,0.18)), url(${current.imageUrl})`,
        }
      : {
          background: getPlaceholderBackground(current.imageLabel),
        }
  }
>
        <div className="hero__overlay">
          <span className="hero__tag">Featured</span>
          <h1 className="hero__title">{current.title}</h1>
          <p className="hero__subtitle">{current.subtitle}</p>

          <Link to={`/title/${current.linkedItemId}`} className="hero__cta">
            View Details
          </Link>
        </div>
      </div>

      <button className="hero__arrow hero__arrow--right" onClick={nextSlide}>
        ›
      </button>

      <div className="hero__dots">
        {slides.map((slide, index) => (
          <button
            key={slide.id}
            className={`hero__dot ${index === activeIndex ? "active" : ""}`}
            onClick={() => setActiveIndex(index)}
          />
        ))}
      </div>
    </section>
  );
}