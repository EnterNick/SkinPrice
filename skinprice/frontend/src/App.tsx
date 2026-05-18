import React, { useCallback, useEffect, useState } from "react";
import { getSavedSkins, updateAllSkinPrices, updateSkinPrice } from "./api/client";
import type { SavedSkin } from "./api/models";
import { toApiError } from "./api/errors";
import "./styles.css";

type SavedSkinsState = {
  items: SavedSkin[];
  loading: boolean;
  error: string | null;
};

type RoutePath = "/" | "/new" | "404";

const routeTitle: Record<Exclude<RoutePath, "404">, string> = {
  "/": "Сохранённые скины",
  "/new": "Новые скины",
};

const navigate = (to: Exclude<RoutePath, "404">) => {
  if (window.location.pathname !== to) {
    window.history.pushState({}, "", to);
    window.dispatchEvent(new PopStateEvent("popstate"));
  }
};

const useRoute = (): RoutePath => {
  const getPath = (): RoutePath => {
    const path = window.location.pathname;
    if (path === "/" || path === "/new") {
      return path;
    }
    return "404";
  };

  const [path, setPath] = useState<RoutePath>(getPath);

  useEffect(() => {
    const onChange = () => setPath(getPath());
    window.addEventListener("popstate", onChange);
    return () => window.removeEventListener("popstate", onChange);
  }, []);

  return path;
};

const SavedSkinsPage: React.FC = () => {
  const [state, setState] = useState<SavedSkinsState>({ items: [], loading: true, error: null });
  const [currency, setCurrency] = useState("1");
  const [isUpdatingAll, setIsUpdatingAll] = useState(false);
  const [notice, setNotice] = useState<{ type: "success" | "warning" | "error"; text: string } | null>(null);

  const loadSkins = async () => {
    const response = await getSavedSkins();
    setState({ items: response.items, loading: false, error: null });
    return response;
  };

  useEffect(() => {
    void loadSkins().catch((err: unknown) => {
      setState({
        items: [],
        loading: false,
        error: err instanceof Error ? err.message : "Не удалось загрузить сохранённые скины",
      });
    });
  }, []);

  const refreshOne = (skinId: string) => updateSkinPrice(skinId, currency).then(loadSkins);

  const refreshAll = async () => {
    setIsUpdatingAll(true);
    setNotice(null);

    try {
      await updateAllSkinPrices(currency);

      try {
        const refreshed = await loadSkins();
        setNotice({
          type: "success",
          text: `Готово: цены обновлены для ${refreshed.items.length} скинов.`,
        });
      } catch (reloadError) {
        const apiError = toApiError(reloadError);
        setNotice({
          type: "warning",
          text: `Частичный успех: цены обновлены, но не удалось синхронизировать список (${apiError.message}).`,
        });
      }
    } catch (updateError) {
      const apiError = toApiError(updateError);
      setNotice({
        type: "error",
        text: `Ошибка обновления цен: ${apiError.message}`,
      });
    } finally {
      setIsUpdatingAll(false);
    }
  };

  const refreshAll = async () => {
    setReloading(true);
    try {
      await updateAllSkinPrices(currency);
      await loadSkins();
    } finally {
      setReloading(false);
    }
  };

  if (state.loading) return <div className="status-box">Загрузка сохранённых скинов...</div>;
  if (state.error) return <div className="status-box">Ошибка: {state.error}</div>;
  if (state.items.length === 0) return <div className="status-box">Нет сохранённых скинов.</div>;

  return (
    <div className="saved-skins-page">
      <div className="saved-skins-toolbar">
        <select value={currency} onChange={(e) => setCurrency(e.target.value)}>
          <option value="1">USD</option>
          <option value="5">RUB</option>
          <option value="3">EUR</option>
        </select>
        <button type="button" disabled={isUpdatingAll} onClick={() => void refreshAll()}>
          {isUpdatingAll ? "Обновление..." : "Обновить все цены"}
        </button>
      </div>
      {notice && <p className={`text notice-${notice.type}`}>{notice.text}</p>}
      {state.items.map((skin) => (
        <div key={skin.id} className="card">
          <div className="image-wrapper">
            <img src={skin.imageUrl} alt={skin.title} className="card-image" />
          </div>
          <div className="card-body">
            <a href={skin.pageUrl} target="_blank" rel="noreferrer">
              <h2 className="title">{skin.title}</h2>
            </a>
            <p className="text">{skin.name}</p>
            <p className="text">Цена: {skin.priceText || "-"}</p>
            <button type="button" onClick={() => void refreshOne(skin.id)}>
              Обновить цену
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};

const NewSkinsPage: React.FC = () => (
  <div className="container">
    <div className="card">
      <div className="card-body">
        <h2 className="title">NewSkinsPage</h2>
        <p className="text">Заглушка страницы новых скинов.</p>
      </div>
    </div>
  </div>
);

const NotFoundPage: React.FC = () => (
  <div className="container">
    <div className="card">
      <div className="card-body">
        <h2 className="title">404</h2>
        <p className="text">Страница не найдена.</p>
      </div>
    </div>
  </div>
);

const App: React.FC = () => {
  const path = useRoute();

  return (
    <div className="app">
      <header className="layout-header">
        <nav className="layout-nav">
          <button
            className={`nav-link ${path === "/" ? "nav-link-active" : ""}`}
            type="button"
            onClick={() => navigate("/")}
          >
            {routeTitle["/"]}
          </button>
          <button
            className={`nav-link ${path === "/new" ? "nav-link-active" : ""}`}
            type="button"
            onClick={() => navigate("/new")}
          >
            {routeTitle["/new"]}
          </button>
        </nav>
      </header>
      <main className="layout-content">
        {path === "/" && <SavedSkinsPage />}
        {path === "/new" && <NewSkinsPage />}
        {path === "404" && <NotFoundPage />}
      </main>
    </div>
  );
};

export default App;
