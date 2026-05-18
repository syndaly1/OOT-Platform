import { Link, NavLink, useNavigate } from "react-router-dom";
import { createBrowseLink } from "../utils/browse";
import { getSelectedProfile } from "../api/profileApi";
import { isAuthenticated, logoutUser } from "../api/authApi";

export default function Navbar() {
  const navigate = useNavigate();
  const authed = isAuthenticated();
  const profile = getSelectedProfile();

  const handleLogout = () => {
    logoutUser();
    navigate("/");
  };

  return (
    <header className="navbar">
      <div className="navbar__left">
        <Link to="/" className="brand">
          DORA<span>BOX</span>
        </Link>
      </div>

      <nav className="navbar__center">
        <NavLink to="/" className={({ isActive }) => `nav-pill ${isActive ? "active" : ""}`}>
          Home
        </NavLink>

        <NavLink
          to={createBrowseLink({ collection: "all", sort: "popular" })}
          className={({ isActive }) => `nav-pill ${isActive ? "active" : ""}`}
        >
          Explore
        </NavLink>

        <NavLink to="/profiles" className={({ isActive }) => `nav-pill ${isActive ? "active" : ""}`}>
          Profiles
        </NavLink>
      </nav>

      <div className="navbar__right">
        {authed ? (
          <>
            <span className="navbar-profile">
              {profile.avatar} {profile.name}
            </span>
            <button className="nav-outline-btn nav-button-reset" onClick={handleLogout}>
              Logout
            </button>
          </>
        ) : (
          <>
            <Link to="/signin" className="nav-outline-btn">
              Sign In
            </Link>
            <Link to="/register" className="nav-solid-btn">
              Register
            </Link>
          </>
        )}
      </div>
    </header>
  );
}