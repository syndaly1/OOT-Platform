export default function FilterChips({ label, options, selectedValue, onSelect }) {
  return (
    <div className="filter-row">
      {label ? <div className="filter-row__label">{label}</div> : null}

      <div className="filter-row__chips">
        {options.map((option) => (
          <button
            key={option.value}
            className={`filter-chip ${selectedValue === option.value ? "active" : ""}`}
            onClick={() => onSelect(option.value)}
          >
            {option.label}
          </button>
        ))}
      </div>
    </div>
  );
}