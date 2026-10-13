const socials = [
  {
    name: "TikTok",
    className: "social-tiktok",
    icon: (
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M19.6 7.2a6.7 6.7 0 0 1-4.1-1.4v8.1a6.1 6.1 0 1 1-5.2-6.1v3.3a2.8 2.8 0 1 0 1.9 2.7V2.5h3.3a6.7 6.7 0 0 0 4.1 4.2v.5Z" />
      </svg>
    ),
  },
  {
    name: "Instagram",
    className: "social-instagram",
    icon: (
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <rect x="3" y="3" width="18" height="18" rx="5" />
        <circle cx="12" cy="12" r="4" />
        <circle className="social-icon-dot" cx="17.6" cy="6.7" r="1" />
      </svg>
    ),
  },
  {
    name: "GitHub",
    className: "social-github",
    icon: (
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 .9a11.1 11.1 0 0 0-3.5 21.6c.6.1.8-.2.8-.5v-2c-3.1.7-3.8-1.3-3.8-1.3-.5-1.3-1.2-1.6-1.2-1.6-1-.7.1-.7.1-.7 1.1.1 1.7 1.2 1.7 1.2 1 .1 1.9-.8 2.3-1.3a3 3 0 0 1 .9-1.9c-2.5-.3-5.1-1.3-5.1-5.5 0-1.2.4-2.1 1.1-2.9-.1-.3-.5-1.4.1-2.9 0 0 .9-.3 3 1.1a10.4 10.4 0 0 1 5.5 0c2.1-1.4 3-1.1 3-1.1.6 1.5.2 2.6.1 2.9.7.8 1.1 1.7 1.1 2.9 0 4.2-2.6 5.2-5.1 5.5.4.4.8 1 .8 2.1V22c0 .3.2.6.8.5A11.1 11.1 0 0 0 12 .9Z" />
      </svg>
    ),
  },
] as const;

export default function AppSocialRail() {
  return (
    <nav className="social-rail" aria-label="Social media">
      {socials.map((social) => (
        <span
          className={`social-icon ${social.className}`}
          key={social.name}
          role="img"
          aria-label={social.name}
          title={social.name}
        >
          {social.icon}
        </span>
      ))}
    </nav>
  );
}
