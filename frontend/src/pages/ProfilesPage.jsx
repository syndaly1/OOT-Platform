import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getCurrentUser } from "../api/authApi";
import { createProfile, getProfiles, selectProfile } from "../api/profileApi";

const profileTypes = [
  { label: "Adult", value: "adult" },
  { label: "Family", value: "family" },
  { label: "Kid", value: "kid" },
];

function formatAccountType(type) {
  const labels = {
    basic: "Basic",
    family: "Family",
    premium: "Premium",
  };

  return labels[type] || "Basic";
}

export default function ProfilesPage() {
  const navigate = useNavigate();
  const user = getCurrentUser();

  const accountType = user?.accountType || "basic";
  const isFamilyAccount = accountType === "family";

  const [profiles, setProfiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [form, setForm] = useState({
    name: "",
    type: "adult",
    avatar: "👤",
  });

  const pageText = useMemo(() => {
    if (!user) {
      return "Sign in or register to select a profile.";
    }

    if (isFamilyAccount) {
      return "Select or create a profile. Each profile keeps its own watch history and continue watching list.";
    }

    return `${formatAccountType(accountType)} account uses one main profile.`;
  }, [user, isFamilyAccount, accountType]);

  useEffect(() => {
    async function load() {
      try {
        const data = await getProfiles();
        setProfiles(data);
      } catch (err) {
        console.error("Failed to load profiles:", err);
        setProfiles([]);
      } finally {
        setLoading(false);
      }
    }

    load();
  }, []);

  const handleSelect = (profileId) => {
    selectProfile(profileId);
    navigate("/");
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    setError("");

    if (!isFamilyAccount) {
      setError("Only Family accounts can create multiple profiles.");
      return;
    }

    if (!form.name.trim()) {
      setError("Profile name is required.");
      return;
    }

    try {
      const newProfile = await createProfile(form);

      setProfiles((prev) => [...prev, newProfile]);
      setForm({
        name: "",
        type: "adult",
        avatar: "👤",
      });
    } catch (err) {
      console.error("Failed to create profile:", err);
      setError(err?.message || "Failed to create profile.");
    }
  };

  return (
    <main className="profiles-page">
      <section className="profiles-hero">
        <h1>Who is watching?</h1>
        <p>{pageText}</p>
      </section>

      {loading ? (
        <section className="continue-empty">
          <h2>Loading profiles...</h2>
          <p>Profiles are waking up. Slowly, because apparently they unionized.</p>
        </section>
      ) : (
        <section className="profiles-grid">
          {profiles.map((profile) => (
            <button
              key={profile.id || profile.profileId}
              type="button"
              className="profile-card"
              onClick={() => handleSelect(profile.id || profile.profileId)}
            >
              <div className="profile-card__avatar">{profile.avatar}</div>
              <h3>{profile.name}</h3>
              <p>{profile.type}</p>
            </button>
          ))}
        </section>
      )}

      {isFamilyAccount ? (
        <section className="create-profile-panel">
          <h2>Create profile</h2>

          <form className="create-profile-form" onSubmit={handleSubmit}>
            <label className="form-field">
              <span>Profile name</span>
              <input
                value={form.name}
                onChange={(event) =>
                  setForm((prev) => ({
                    ...prev,
                    name: event.target.value,
                  }))
                }
                placeholder="Example: Ingkar"
              />
            </label>

            <label className="form-field">
              <span>Avatar</span>
              <input
                value={form.avatar}
                onChange={(event) =>
                  setForm((prev) => ({
                    ...prev,
                    avatar: event.target.value,
                  }))
                }
                placeholder="🌙"
              />
            </label>

            <div className="profile-type-row">
              {profileTypes.map((type) => (
                <button
                  key={type.value}
                  type="button"
                  className={`filter-chip ${
                    form.type === type.value ? "active" : ""
                  }`}
                  onClick={() =>
                    setForm((prev) => ({
                      ...prev,
                      type: type.value,
                    }))
                  }
                >
                  {type.label}
                </button>
              ))}
            </div>

            {error ? <div className="form-error">{error}</div> : null}

            <button className="primary-btn" type="submit">
              Add Profile
            </button>
          </form>
        </section>
      ) : (
        <section className="create-profile-panel">
          <h2>Single profile account</h2>
          <p>
            This account type does not create extra profiles. Choose Family
            during registration to use multiple profiles.
          </p>
          {error ? <div className="form-error">{error}</div> : null}
        </section>
      )}
    </main>
  );
}