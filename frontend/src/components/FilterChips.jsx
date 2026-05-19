export default function FilterChips({
  label,
  options,
  value,
  selectedValue,
  onChange,
  onSelect,
}) {
  const activeValue = value ?? selectedValue;
  const handleChange = onChange ?? onSelect;

  return (
    <div className="filter-row">
      {label ? <div className="filter-row__label">{label}</div> : null}

      <div className="filter-row__chips">
        {options.map((option) => (
          <button
            key={option.value}
            type="button"
            className={`filter-chip ${
              activeValue === option.value ? "active" : ""
            }`}
            onClick={() => handleChange?.(option.value)}
          >
            {option.label}
          </button>
        ))}
      </div>
    </div>
  );
}