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
  const [reloading, setReloading] = useState(false);

  const loadSkins = useCallback(async () => {
    setState((prev) => ({ ...prev, loading: true, error: null }));
    try {
      const response = await getSavedSkins();
      setState({ items: response.items, loading: false, error: null });
    } catch (err: unknown) {
      const apiError = toApiError(err);
      setState({ items: [], loading: false, error: apiError.message });
    }
  }, []);

  useEffect(() => {
    void loadSkins();
  }, [loadSkins]);

  const refreshOne = async (skinId: string) => {
    setReloading(true);
    try {
      await updateSkinPrice(skinId, currency);
      await loadSkins();
    } finally {
      setReloading(false);
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
        <button onClick={() => void refreshAll()} disabled={reloading}>
          {reloading ? "Обновление..." : "Обновить цены всех"}
        </button>
      </div>

      <table className="skins-table">
        <thead>
          <tr>
            <th>Название</th>
            <th>Цена</th>
            <th>Обновлено</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {state.items.map((skin) => (
            <tr key={skin.id}>
              <td>
                <a href={skin.pageUrl} target="_blank" rel="noreferrer">
                  {skin.title}
                </a>
              </td>
              <td>{skin.priceText || "-"}</td>
              <td>-</td>
              <td>
                <button onClick={() => void refreshOne(skin.id)} disabled={reloading}>
                  Обновить
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
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
