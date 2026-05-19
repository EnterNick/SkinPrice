import React from "react";
import { CardGrid } from "../../shared/ui/card-grid/CardGrid";
import { LoadingState, StatusBox } from "../../shared/ui/states/States";
import { UI_TEXT } from "../../shared/config/uiText";
import { SkinCard } from "../../entities/skin/ui/SkinCard";
import type { NewSkin } from "../../entities/skin/model/types";

type NewSkinsResultsProps = {
  items: NewSkin[];
  loading: boolean;
  loadingMore: boolean;
  hasMore: boolean;
  savingIds: Record<string, boolean>;
  savedIds: Record<string, boolean>;
  onLoadMore: () => Promise<void> | void;
  onSave: (skin: NewSkin) => Promise<void> | void;
};

export const NewSkinsResults: React.FC<NewSkinsResultsProps> = ({
  items,
  loading,
  loadingMore,
  hasMore,
  savingIds,
  savedIds,
  onLoadMore,
  onSave,
}) => (
  <>
    {loading && (
      <CardGrid
        className="search-results"
        items={Array.from({ length: 8 }, (_, idx) => idx)}
        renderItem={(idx) => (
          <div key={idx} className="card card-skeleton" aria-hidden="true">
            <div className="skeleton-image" />
            <div className="card-body">
              <div className="skeleton-line skeleton-line-title" />
              <div className="skeleton-line" />
              <div className="skeleton-line" />
              <div className="skeleton-line skeleton-line-short" />
            </div>
          </div>
        )}
      />
    )}
    <CardGrid
      className="search-results"
      items={items}
      renderItem={(skin) => (
        <SkinCard
          key={skin.id}
          skin={skin}
          saving={Boolean(savingIds[skin.id])}
          saved={Boolean(savedIds[skin.id])}
          onSave={onSave}
        />
      )}
    />
    {!loading && hasMore && !loadingMore && (
      <div className="load-more-wrap">
        <button className="toolbar-button toolbar-button-primary" type="button" onClick={() => void onLoadMore()}>
          {UI_TEXT.loadMore}
        </button>
      </div>
    )}
    {!loading && loadingMore && <LoadingState text={UI_TEXT.loadingMore} />}
    {!loading && !hasMore && items.length > 0 && <StatusBox text={UI_TEXT.listFullyLoaded} />}
  </>
);
