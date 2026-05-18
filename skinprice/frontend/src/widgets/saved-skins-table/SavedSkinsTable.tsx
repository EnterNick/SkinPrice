import React from "react";
import { UI_TEXT } from "../../shared/config/uiText";
import { openExternal } from "../../shared/lib/browser/openExternal";
import { formatUpdatedAt } from "../../shared/lib/date/formatUpdatedAt";
import { SkinThumb } from "../../entities/skin/ui/SkinThumb";
import type { SavedSkin } from "../../entities/skin/model/types";

type SavedSkinsTableProps = {
  items: SavedSkin[];
  isUpdatingAll: boolean;
  updatingSkinIds: Record<string, boolean>;
  onRefreshOne: (skinId: string) => Promise<void> | void;
};

export const SavedSkinsTable: React.FC<SavedSkinsTableProps> = ({ items, isUpdatingAll, updatingSkinIds, onRefreshOne }) => (
  <div className="table-shell">
    <table className="skins-table">
      <thead>
        <tr>
          <th>Скин</th>
          <th>Название</th>
          <th>{UI_TEXT.priceLabel}</th>
          <th>{UI_TEXT.updatedLabel}</th>
          <th>Ссылка</th>
          <th>Действие</th>
        </tr>
      </thead>
      <tbody>
        {items.map((skin) => {
          const isUpdating = Boolean(updatingSkinIds[skin.id]);
          return (
            <tr key={skin.id}>
              <td>
                <div className="skin-cell">
                  <SkinThumb imageUrl={skin.imageUrl} title={skin.title} />
                  <div>
                    <div className="skin-title">{skin.title}</div>
                    <div className="skin-subtitle">{skin.name}</div>
                  </div>
                </div>
              </td>
              <td>{skin.title}</td>
              <td>{skin.priceText || "-"}</td>
              <td>{formatUpdatedAt(skin.updatedAt)}</td>
              <td>
                <button className="table-link-button" type="button" onClick={() => openExternal(skin.pageUrl)}>
                  {UI_TEXT.openLink}
                </button>
              </td>
              <td>
                <button
                  className="table-action-button"
                  type="button"
                  disabled={isUpdating || isUpdatingAll}
                  onClick={() => void onRefreshOne(skin.id)}
                >
                  {isUpdating ? UI_TEXT.updateOnePending : UI_TEXT.updateOne}
                </button>
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  </div>
);
