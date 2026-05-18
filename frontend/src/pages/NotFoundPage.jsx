import { Link } from "react-router-dom";

export default function NotFoundPage() {
  return (
    <main className="not-found-page">
      <h1>404</h1>
      <p>Page not found.</p>
      <Link to="/" className="primary-btn primary-btn--inline">
        Go Home
      </Link>
    </main>
  );
}