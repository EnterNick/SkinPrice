import React, { useEffect, useState } from "react";
import { GetSavedSkins, UpdateAllSavedSkinsPrices, UpdateSavedSkinPrice } from "./wailsjs/go/main/App";
import { skins } from "./wailsjs/go/models";
import "./styles.css";

type SavedSkinsState = {
  items: skins.SavedSkinsResponse["items"];
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
  const [state, setState] = useState<SavedSkinsState>({
    items: [],
    loading: true,
    error: null,
  });
  const [currency, setCurrency] = useState("1");

  const loadSkins = () =>
    GetSavedSkins({ limit: 50, offset: 0 }).then((response) =>
      setState({ items: response.items, loading: false, error: null }),
    );

  useEffect(() => {
    void loadSkins().catch((err: unknown) => {
      setState({
        items: [],
        loading: false,
        error: err instanceof Error ? err.message : "Не удалось загрузить сохранённые скины",
      });
    });
  }, []);

  const refreshOne = (marketHashName: string) =>
    UpdateSavedSkinPrice({ market_hash_name: marketHashName, currency }).then(loadSkins);
  const refreshAll = () => UpdateAllSavedSkinsPrices({ currency }).then(loadSkins);

  if (state.loading) {
    return <div className="container">Загрузка...</div>;
  }

  if (state.error) {
    return <div className="container">Ошибка: {state.error}</div>;
  }

  return (
    <div className="container">
      <div style={{ marginBottom: 16, display: "flex", gap: 8 }}>
        <select value={currency} onChange={(e) => setCurrency(e.target.value)}>
          <option value="1">USD</option>
          <option value="5">RUB</option>
          <option value="3">EUR</option>
        </select>
        <button onClick={() => void refreshAll()}>Обновить цены всех</button>
      </div>
      {state.items.map((skin) => (
        <div key={skin.market_hash_name} className="card">
          <div className="image-wrapper">
            <img src={skin.icon_url} alt={skin.display_name} className="card-image" />
          </div>
          <div className="card-body">
            <a href={skin.page_url} target="_blank" rel="noreferrer">
              <h2 className="title">{skin.display_name}</h2>
            </a>
            <p className="text">{skin.market_hash_name}</p>
            <p className="text">Цена: {skin.price_text || "-"}</p>
            <button onClick={() => void refreshOne(skin.market_hash_name)}>Обновить цену</button>
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
