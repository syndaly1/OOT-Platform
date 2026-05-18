import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createProfile, getProfiles, selectProfile } from "../api/profileApi";

const profileTypes = [
  { label: "Adult", value: "adult" },
  { label: "Family", value: "family" },
  { label: "Kid", value: "kid" },
];

export default function ProfilesPage() {
  const navigate = useNavigate();

  const [profiles, setProfiles] = useState([]);
  const [form, setForm] = useState({
    name: "",
    type: "adult",
    avatar: "👤",
  });

  useEffect(() => {
    async function load() {
      const data = await getProfiles();
      setProfiles(data);
    }

    load();
  }, []);

  const handleSelect = (profileId) => {
    selectProfile(profileId);
    navigate("/");
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!form.name.trim()) {
      alert("Profile name is required.");
      return;
    }

    const newProfile = await createProfile(form);

    setProfiles((prev) => [...prev, newProfile]);
    setForm({
      name: "",
      type: "adult",
      avatar: "👤",
    });
  };

  return (
    <main className="profiles-page">
      <section className="profiles-hero">
        <h1>Who is watching?</h1>
        <p>Select a profile. Each profile keeps its own watch history and continue watching list.</p>
      </section>

      <section className="profiles-grid">
        {profiles.map((profile) => (
          <button
            key={profile.id}
            className="profile-card"
            onClick={() => handleSelect(profile.id)}
          >
            <div className="profile-card__avatar">{profile.avatar}</div>
            <h3>{profile.name}</h3>
            <p>{profile.type}</p>
          </button>
        ))}
      </section>

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
                className={`filter-chip ${form.type === type.value ? "active" : ""}`}
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

          <button className="primary-btn" type="submit">
            Add Profile
          </button>
        </form>
      </section>
    </main>
  );
}