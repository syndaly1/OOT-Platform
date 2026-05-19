import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { registerUser } from "../api/authApi";

const accountPlans = [
  {
    id: "basic",
    title: "Basic",
    description: "Single user account with one main profile.",
  },
  {
    id: "family",
    title: "Family",
    description: "Best for multiple profiles inside one account.",
  },
  {
    id: "premium",
    title: "Premium",
    description: "Extended access with one main profile for this demo.",
  },
];

export default function RegisterPage() {
  const navigate = useNavigate();

  const [form, setForm] = useState({
    fullName: "",
    email: "",
    password: "",
    confirmPassword: "",
    accountType: "basic",
  });

  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (event) => {
    setError("");

    setForm((prev) => ({
      ...prev,
      [event.target.name]: event.target.value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (loading) return;

    setError("");

    if (form.password !== form.confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    setLoading(true);

    try {
      await registerUser(form);
      navigate("/profiles");
    } catch (err) {
      console.error("Registration failed:", err);

      const message = err?.message || "Registration failed.";

      if (message.includes("email already in use")) {
        setError("This email is already registered. Please sign in instead.");
      } else {
        setError(message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="auth-page">
      <div className="auth-card auth-card--wide">
        <div className="auth-card__intro">
          <h1>Create account</h1>
          <p>Select an account type and complete your registration.</p>
        </div>

        <form className="auth-form" onSubmit={handleSubmit}>
          <div className="form-grid">
            <label className="form-field">
              <span>Full Name</span>
              <input
                type="text"
                name="fullName"
                value={form.fullName}
                onChange={handleChange}
                placeholder="Enter your full name"
                required
              />
            </label>

            <label className="form-field">
              <span>Email</span>
              <input
                type="email"
                name="email"
                value={form.email}
                onChange={handleChange}
                placeholder="Enter your email"
                required
              />
            </label>

            <label className="form-field">
              <span>Password</span>
              <input
                type="password"
                name="password"
                value={form.password}
                onChange={handleChange}
                placeholder="Enter your password"
                required
              />
            </label>

            <label className="form-field">
              <span>Confirm Password</span>
              <input
                type="password"
                name="confirmPassword"
                value={form.confirmPassword}
                onChange={handleChange}
                placeholder="Confirm your password"
                required
              />
            </label>
          </div>

          <div className="plan-grid">
            {accountPlans.map((plan) => (
              <button
                key={plan.id}
                type="button"
                className={`plan-card ${
                  form.accountType === plan.id ? "selected" : ""
                }`}
                onClick={() =>
                  setForm((prev) => ({
                    ...prev,
                    accountType: plan.id,
                  }))
                }
              >
                <h3>{plan.title}</h3>
                <p>{plan.description}</p>
              </button>
            ))}
          </div>

          {error ? <div className="form-error">{error}</div> : null}

          <button type="submit" className="primary-btn" disabled={loading}>
            {loading ? "Registering..." : "Register"}
          </button>
        </form>
      </div>
    </main>
  );
}