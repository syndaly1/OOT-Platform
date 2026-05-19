import { useEffect, useState } from "react";
import { Link, NavLink, useNavigate } from "react-router-dom";
import { createBrowseLink } from "../utils/browse";
import { getSelectedProfile } from "../api/profileApi";
import { getCurrentUser, isAuthenticated, logoutUser } from "../api/authApi";

function getNavbarState() {
  return {
    authed: isAuthenticated(),
    user: getCurrentUser(),
    profile: getSelectedProfile(),
  };
}

export default function Navbar() {
  const navigate = useNavigate();
  const [state, setState] = useState(getNavbarState);

  useEffect(() => {
    const refresh = () => {
      setState(getNavbarState());
    };

    window.addEventListener("storage", refresh);
    window.addEventListener("auth-changed", refresh);
    window.addEventListener("profile-changed", refresh);

    return () => {
      window.removeEventListener("storage", refresh);
      window.removeEventListener("auth-changed", refresh);
      window.removeEventListener("profile-changed", refresh);
    };
  }, []);

  const handleLogout = async () => {
    await logoutUser();
    setState(getNavbarState());
    navigate("/");
  };

  const displayName =
    state.profile?.name ||
    state.user?.username ||
    state.user?.fullName ||
    state.user?.email ||
    "Account";

  const avatar = state.profile?.avatar || "👤";

  return (
    <header className="navbar">
      <div className="navbar__left">
        <Link to="/" className="brand">
          DORA<span>BOX</span>
        </Link>
      </div>

      <nav className="navbar__center">
        <NavLink
          to="/"
          className={({ isActive }) => `nav-pill ${isActive ? "active" : ""}`}
        >
          Home
        </NavLink>

        <NavLink
          to={createBrowseLink({ collection: "all", sort: "popular" })}
          className={({ isActive }) => `nav-pill ${isActive ? "active" : ""}`}
        >
          Explore
        </NavLink>

        <NavLink
          to="/profiles"
          className={({ isActive }) => `nav-pill ${isActive ? "active" : ""}`}
        >
          Profiles
        </NavLink>
      </nav>

      <div className="navbar__right">
        {state.authed ? (
          <>
            <Link to="/profiles" className="navbar-profile">
              <span>{avatar}</span>
              <span>{displayName}</span>
            </Link>

            <button
              type="button"
              className="nav-outline-btn nav-button-reset"
              onClick={handleLogout}
            >
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