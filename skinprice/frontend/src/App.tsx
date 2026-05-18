import React, { useState } from "react";
import { saveSkin, updateAllSkinPrices, updateSkinPrice } from "./api/client";
import { DataTable } from "./components/DataTable";
import { EmptyState, ErrorState, LoadingState, ToastAlert } from "./components/States";
import { ROUTE_TITLE, UI_TEXT } from "./constants/uiText";
import { useNewSkinsSearch } from "./hooks/useNewSkinsSearch";
import { useSavedSkins } from "./hooks/useSavedSkins";
import { formatErrorMessage } from "./utils/error";
import "./styles.css";

type RoutePath = "/" | "/new" | "404";

const navigate = (to: Exclude<RoutePath, "404">) => {
  if (window.location.pathname !== to) {
    window.history.pushState({}, "", to);
    window.dispatchEvent(new PopStateEvent("popstate"));
  }
};

const useRoute = (): RoutePath => {
  const getPath = (): RoutePath => {
    const path = window.location.pathname;
    if (path === "/" || path === "/new") return path;
    return "404";
  };

  const [path, setPath] = useState<RoutePath>(getPath);
  React.useEffect(() => {
    const onChange = () => setPath(getPath());
    window.addEventListener("popstate", onChange);
    return () => window.removeEventListener("popstate", onChange);
  }, []);

  return path;
};

const SavedSkinsPage: React.FC = () => {
  const { items, loading, error, loadSkins } = useSavedSkins();
  const [currency, setCurrency] = useState("1");
  const [updatingSkinIds, setUpdatingSkinIds] = useState<Record<string, boolean>>({});
  const [updatedAtMap, setUpdatedAtMap] = useState<Record<string, string>>({});
  const [isUpdatingAll, setIsUpdatingAll] = useState(false);
  const [notice, setNotice] = useState<{ type: "success" | "warning" | "error"; text: string } | null>(null);

  const refreshOne = async (skinId: string) => {
    setUpdatingSkinIds((prev) => ({ ...prev, [skinId]: true }));
    try {
      await updateSkinPrice(skinId, currency);
      await loadSkins();
      setUpdatedAtMap((prev) => ({ ...prev, [skinId]: new Date().toLocaleString("ru-RU") }));
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errUpdateOne, err) });
    } finally {
      setUpdatingSkinIds((prev) => ({ ...prev, [skinId]: false }));
    }
  };

  const refreshAll = async () => {
    setIsUpdatingAll(true);
    setNotice(null);
    try {
      await updateAllSkinPrices(currency);
      try {
        const refreshed = await loadSkins();
        const refreshedAt = new Date().toLocaleString("ru-RU");
        setUpdatedAtMap(Object.fromEntries(refreshed.items.map((skin) => [skin.id, refreshedAt])));
        setNotice({ type: "success", text: UI_TEXT.successUpdatedAll.replace("{count}", String(refreshed.items.length)) });
      } catch (reloadError) {
        const message = formatErrorMessage("", reloadError).trim();
        setNotice({ type: "warning", text: UI_TEXT.warningPartialUpdated.replace("{message}", message) });
      }
    } catch (updateError) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errUpdateAll, updateError) });
    } finally {
      setIsUpdatingAll(false);
    }
  };

  if (loading) return <LoadingState text={UI_TEXT.loadingSaved} />;
  if (error) return <ErrorState text={error} />;
  if (items.length === 0) return <EmptyState text={UI_TEXT.notFoundSaved} />;

  return (
    <div className="saved-skins-page">
      <div className="saved-skins-toolbar">
        <select value={currency} onChange={(e) => setCurrency(e.target.value)}><option value="1">USD</option><option value="5">RUB</option><option value="3">EUR</option></select>
        <button type="button" disabled={isUpdatingAll} onClick={() => void refreshAll()}>{isUpdatingAll ? UI_TEXT.updateAllPending : UI_TEXT.updateAll}</button>
        <button type="button" onClick={() => navigate("/new")}>{UI_TEXT.ctaAddSkins}</button>
      </div>
      {notice && <ToastAlert type={notice.type} text={notice.text} />}
      <DataTable
        className="saved-skins-page"
        items={items}
        renderItem={(skin) => {
          const isUpdating = Boolean(updatingSkinIds[skin.id]);
          return <div key={skin.id} className="card"><div className="image-wrapper"><img src={skin.imageUrl} alt={skin.title} className="card-image" /></div><div className="card-body"><a href={skin.pageUrl} target="_blank" rel="noreferrer"><h2 className="title">{skin.title}</h2></a><p className="text">{skin.name}</p><p className="text">{UI_TEXT.priceLabel}: {skin.priceText || "-"}</p><p className="text">{UI_TEXT.updatedLabel}: {updatedAtMap[skin.id] || "-"}</p><button type="button" disabled={isUpdating || isUpdatingAll} onClick={() => void refreshOne(skin.id)}>{isUpdating ? UI_TEXT.updateOnePending : UI_TEXT.updateOne}</button></div></div>;
        }}
      />
    </div>
  );
};

const NewSkinsPage: React.FC = () => {
  const { query, setQuery, debouncedQuery, items, loading, error, loadNewSkins } = useNewSkinsSearch();
  const [savingIds, setSavingIds] = useState<Record<string, boolean>>({});
  const [savedIds, setSavedIds] = useState<Record<string, boolean>>({});
  const [notice, setNotice] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const onSave = async (id: string) => {
    setSavingIds((prev) => ({ ...prev, [id]: true }));
    setNotice(null);
    try {
      await saveSkin(id);
      setSavedIds((prev) => ({ ...prev, [id]: true }));
      setNotice({ type: "success", text: UI_TEXT.saveDoneShort });
      await loadNewSkins(debouncedQuery);
    } catch (err: unknown) {
      setNotice({ type: "error", text: formatErrorMessage(UI_TEXT.errSave, err) });
    } finally {
      setSavingIds((prev) => ({ ...prev, [id]: false }));
    }
  };

  return <div className="new-skins-page"><div className="new-skins-toolbar"><button type="button" onClick={() => navigate("/")}>{UI_TEXT.ctaToHome}</button></div><div className="new-skins-toolbar"><input className="search-input" type="search" placeholder={UI_TEXT.searchPlaceholder} value={query} onChange={(e) => setQuery(e.target.value)} /></div>{notice && <ToastAlert type={notice.type} text={notice.text} />}{loading && <LoadingState text={UI_TEXT.loadingNew} />}{error && !loading && <ErrorState text={error} />}{!loading && !error && items.length === 0 && <EmptyState text={UI_TEXT.notFoundSearch} />}{!loading && !error && items.length > 0 && <DataTable items={items} renderItem={(skin) => <div key={skin.id} className="card"><div className="image-wrapper"><img src={skin.imageUrl} alt={skin.title} className="card-image" /></div><div className="card-body"><a href={skin.pageUrl} target="_blank" rel="noreferrer"><h2 className="title">{skin.title}</h2></a><p className="text">{skin.name}</p><p className="text">{UI_TEXT.priceLabel}: {skin.priceText || "-"}</p><p className="text">{UI_TEXT.listingsLabel}: {skin.sellListings ?? "-"}</p><button type="button" disabled={Boolean(savingIds[skin.id]) || Boolean(savedIds[skin.id])} onClick={() => void onSave(skin.id)}>{savingIds[skin.id] ? UI_TEXT.saveSkinPending : savedIds[skin.id] ? UI_TEXT.saveSkinDone : UI_TEXT.saveSkin}</button></div></div>} />}</div>;
};

const NotFoundPage: React.FC = () => <div className="container"><div className="card"><div className="card-body"><h2 className="title">404</h2><p className="text">{UI_TEXT.pageNotFound}</p></div></div></div>;

const App: React.FC = () => {
  const path = useRoute();
  return <div className="app"><header className="layout-header"><nav className="layout-nav"><button className={`nav-link ${path === "/" ? "nav-link-active" : ""}`} type="button" onClick={() => navigate("/")}>{ROUTE_TITLE["/"]}</button><button className={`nav-link ${path === "/new" ? "nav-link-active" : ""}`} type="button" onClick={() => navigate("/new")}>{ROUTE_TITLE["/new"]}</button></nav></header><main className="layout-content">{path === "/" && <SavedSkinsPage />}{path === "/new" && <NewSkinsPage />}{path === "404" && <NotFoundPage />}</main></div>;
};

export default App;
