import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "../../../app/router/routes";
import { saveSkin } from "../../../entities/skin/api/skinApi";
import type { NewSkin } from "../../../entities/skin/model/types";
import { useNewSkinsSearch } from "../../../features/search-skins/model/useNewSkinsSearch";
import { MIN_SEARCH_LENGTH } from "../../../shared/config/search";
import { UI_TEXT } from "../../../shared/config/uiText";
import { formatErrorMessage } from "../../../shared/lib/error/formatErrorMessage";
import { PageHeader } from "../../../shared/ui/page-header/PageHeader";
import { EmptyState, ErrorState, LoadingState, ToastAlert } from "../../../shared/ui/states/States";
import { NewSkinsResults } from "../../../widgets/new-skins-results/NewSkinsResults";
import { NewSkinsSearchPanel } from "../../../widgets/new-skins-search-panel/NewSkinsSearchPanel";

export const NewSkinsPage: React.FC = () => {
  const navigate = useNavigate();
  const { items, loading, loadingMore, error, hasMore, hasSearched, currentQuery, loadNewSkins, loadNextPage } = useNewSkinsSearch();
  const [savingIds, setSavingIds] = useState<Record<string, boolean>>({});
  const [savedIds, setSavedIds] = useState<Record<string, boolean>>({});
  const [notice, setNotice] = useState<{ type: "success" | "error"; text: string } | null>(null);
  const [searchLabel, setSearchLabel] = useState("");
  const [activeSource, setActiveSource] = useState<"steam" | "lisskins">("steam");

  const minSearchMessage = `Введите минимум ${MIN_SEARCH_LENGTH} символа для поиска.`;

  const onSearch = async () => {
    const rawValue = window.prompt(UI_TEXT.searchPrompt, currentQuery || searchLabel);
    if (rawValue === null) return;

    const value = rawValue.trim();
    if (value.length < MIN_SEARCH_LENGTH) {
      setNotice({ type: "error", text: minSearchMessage });
      return;
    }

    setNotice(null);
    setSearchLabel(value);
    await loadNewSkins(value);
  };

  const onChangeSource = (source: "steam" | "lisskins") => {
    setActiveSource(source);
    setNotice(null);
  };

  const onSave = async (skin: NewSkin) => {
    setSavingIds((prev) => ({ ...prev, [skin.id]: true }));
    setNotice(null);
    try {
      await saveSkin(skin);
      setSavedIds((prev) => ({ ...prev, [skin.id]: true }));
      setNotice({ type: "success", text: UI_TEXT.saveDoneShort });
      await loadNewSkins(currentQuery);
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errSave, err) });
    } finally {
      setSavingIds((prev) => ({ ...prev, [skin.id]: false }));
    }
  };

  return (
    <div className="new-skins-page">
      <PageHeader
        eyebrow={UI_TEXT.searchEyebrow}
        title={UI_TEXT.newSkinsTitle}
        actions={
          <div className="new-skins-toolbar">
            <button className="toolbar-button" type="button" onClick={() => navigate(ROUTES.home)}>
              {UI_TEXT.ctaToHome}
            </button>
          </div>
        }
      />
      <NewSkinsSearchPanel
        searchLabel={searchLabel}
        fallbackLabel={`Минимум ${MIN_SEARCH_LENGTH} символа`}
        activeSource={activeSource}
        onChangeSource={onChangeSource}
        onSearch={onSearch}
      />
      {notice && <ToastAlert type={notice.type} text={notice.text} />}
      {loading && <LoadingState text={UI_TEXT.loadingNew} />}
      {error && !loading && <ErrorState text={error} />}
      {!loading && !error && !hasSearched && <EmptyState text={minSearchMessage} />}
      {!loading && !error && hasSearched && items.length === 0 && <EmptyState text={UI_TEXT.notFoundSearch} />}
      {!loading && !error && items.length > 0 && (
        <NewSkinsResults
          items={items}
          loadingMore={loadingMore}
          hasMore={hasMore}
          savingIds={savingIds}
          savedIds={savedIds}
          onLoadMore={loadNextPage}
          onSave={onSave}
        />
      )}
    </div>
  );
};
