"use client";

import { useState } from "react";

const placeholder = "/default-profile.svg";

type Props = { photoUrl: string | null; name: string; className: string };

function Photo({ photoUrl, name, className }: Props) {
  const [failed, setFailed] = useState(false);
  return (
    // Signed private URLs are fetched directly; no image optimizer or URL allowlist is needed.
    // eslint-disable-next-line @next/next/no-img-element
    <img src={failed || !photoUrl?.trim() ? placeholder : photoUrl} alt={`${name}'s profile`}
      className={className} width={86} height={86} referrerPolicy="no-referrer"
      onError={() => { if (!failed && photoUrl) setFailed(true); }} />
  );
}

export default function ProfilePhoto(props: Props) {
  return <Photo key={props.photoUrl} {...props} />;
}
