import React from "react";
import type { NewSkin } from "../model/types";
import { UI_TEXT } from "../../../shared/config/uiText";
import { openExternal } from "../../../shared/lib/browser/openExternal";

type SkinCardProps = {
  skin: NewSkin;
  saving: boolean;
  saved: boolean;
  onSave: (skin: NewSkin) => Promise<void> | void;
};

export const SkinCard: React.FC<SkinCardProps> = ({ skin, saving, saved, onSave }) => (
  <div className="card">
    <div className="image-wrapper">
      {skin.imageUrl ? <img src={skin.imageUrl} alt={skin.title} className="card-image" /> : <div className="image-fallback">{skin.title}</div>}
    </div>
    <div className="card-body">
      <button className="card-link-button" type="button" onClick={() => openExternal(skin.pageUrl)}>
        <h2 className="title">{skin.title}</h2>
      </button>
      <p className="text">{skin.name}</p>
      <p className="text">
        {UI_TEXT.priceLabel}: {skin.priceText || "-"}
      </p>
      <p className="text">
        {UI_TEXT.listingsLabel}: {skin.sellListings ?? "-"}
      </p>
      <button
        className="toolbar-button toolbar-button-primary card-save-button"
        type="button"
        disabled={saving || saved}
        onClick={() => void onSave(skin)}
      >
        {saving ? UI_TEXT.saveSkinPending : saved ? UI_TEXT.saveSkinDone : UI_TEXT.saveSkin}
      </button>
    </div>
  </div>
);
