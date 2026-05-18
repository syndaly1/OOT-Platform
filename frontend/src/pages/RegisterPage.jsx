import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { registerUser } from "../api/authApi";

const accountPlans = [
  {
    id: "basic",
    title: "Basic",
    description: "Single user account with a simple setup.",
  },
  {
    id: "family",
    title: "Family",
    description: "Best for multiple profiles inside one account.",
  },
  {
    id: "premium",
    title: "Premium",
    description: "Extended access and advanced account features.",
  },
];

export default function RegisterPage() {
  const navigate = useNavigate();

  const [form, setForm] = useState({
    fullName: "",
    email: "",
    password: "",
    confirmPassword: "",
    accountType: "family",
  });

  const handleChange = (event) => {
    setForm((prev) => ({
      ...prev,
      [event.target.name]: event.target.value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (form.password !== form.confirmPassword) {
      alert("Passwords do not match.");
      return;
    }

    await registerUser(form);
    navigate("/profiles");
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
                className={`plan-card ${form.accountType === plan.id ? "selected" : ""}`}
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

          <button type="submit" className="primary-btn">
            Register
          </button>
        </form>
      </div>
    </main>
  );
}