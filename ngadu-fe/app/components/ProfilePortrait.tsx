const palettes = [
  { bg: "#f4dfc6", skin: "#d99870", hair: "#3b2d29", shirt: "#638c68" },
  { bg: "#dcebd7", skin: "#f0bd91", hair: "#49332a", shirt: "#d57c55" },
  { bg: "#f5e2cf", skin: "#b97856", hair: "#292c32", shirt: "#6b82a5" },
  { bg: "#dce8f3", skin: "#e6ad82", hair: "#382b2a", shirt: "#e6b94e" },
  { bg: "#f2dfe1", skin: "#e8b08c", hair: "#3b2927", shirt: "#9c7396" },
  { bg: "#e6e6d2", skin: "#c58462", hair: "#332b2b", shirt: "#498b86" },
];

type Props = { variant?: number; className?: string; name?: string };

export default function ProfilePortrait({ variant = 0, className = "", name = "Illustrated player" }: Props) {
  const colors = palettes[((variant % palettes.length) + palettes.length) % palettes.length];
  return (
    <span className={`profile-portrait ${className}`} role="img" aria-label={`${name} illustrated avatar`}>
      <svg viewBox="0 0 80 80" aria-hidden="true">
        <circle cx="40" cy="40" r="40" fill={colors.bg} />
        <path d="M12 80c2-16 12-24 28-24s26 8 28 24" fill={colors.shirt} />
        <path d="M33 49h14v14H33z" fill={colors.skin} />
        <ellipse cx="40" cy="35" rx="20" ry="23" fill={colors.skin} />
        <path d="M20 34c-2-18 7-27 21-27 14 0 22 9 19 28l-5-9c-7 1-19-2-27-8-1 8-4 12-8 16z" fill={colors.hair} />
        <path d="M31 37h2m14 0h2" stroke="#372a27" strokeWidth="3" strokeLinecap="round" />
        <path d="M35 47c3 3 7 3 10 0" fill="none" stroke="#8a4e43" strokeWidth="2" strokeLinecap="round" />
        <circle cx="18" cy="37" r="3" fill={colors.skin} /><circle cx="62" cy="37" r="3" fill={colors.skin} />
      </svg>
    </span>
  );
}
