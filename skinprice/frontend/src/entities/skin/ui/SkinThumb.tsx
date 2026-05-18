import React from "react";

type SkinThumbProps = {
  imageUrl: string;
  title: string;
  fallbackText?: string;
};

export const SkinThumb: React.FC<SkinThumbProps> = ({ imageUrl, title, fallbackText = "?" }) => (
  <div className="skin-thumb">
    {imageUrl ? <img src={imageUrl} alt={title} className="skin-thumb-image" /> : <div className="skin-thumb-fallback">{fallbackText}</div>}
  </div>
);
