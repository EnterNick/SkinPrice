import React, { useEffect, useRef, useState } from "react";
import { getNewSkins, getSavedSkins, saveSkin, updateAllSkinPrices, updateSkinPrice } from "./api/client";
import type { NewSkin, SavedSkin } from "./api/models";
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
  const [updatingSkinIds, setUpdatingSkinIds] = useState<Record<string, boolean>>({});
  const [updatedAtMap, setUpdatedAtMap] = useState<Record<string, string>>({});
  const [isUpdatingAll, setIsUpdatingAll] = useState(false);
  const [notice, setNotice] = useState<{ type: "success" | "warning" | "error"; text: string } | null>(null);

  const loadSkins = async () => {
    const response = await getSavedSkins();
    setState({ items: response.items, loading: false, error: null });
    return response;
  };

  useEffect(() => {
    void loadSkins().catch((err: unknown) => {
      const apiError = toApiError(err);
      setState({ items: [], loading: false, error: apiError.message || "Не удалось загрузить сохранённые скины" });
    });
  }, []);

  const refreshOne = async (skinId: string) => {
    setUpdatingSkinIds((prev) => ({ ...prev, [skinId]: true }));

    try {
      await updateSkinPrice(skinId, currency);
      await loadSkins();
      setUpdatedAtMap((prev) => ({ ...prev, [skinId]: new Date().toLocaleString("ru-RU") }));
    } catch (err: unknown) {
      const apiError = toApiError(err);
      window.alert(`Не удалось обновить цену: ${apiError.message}`);
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
        setUpdatedAtMap(
          Object.fromEntries(refreshed.items.map((skin) => [skin.id, refreshedAt])),
        );
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
      {state.items.map((skin) => {
        const isUpdating = Boolean(updatingSkinIds[skin.id]);

        return (
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
              <p className="text">Обновлено: {updatedAtMap[skin.id] || "-"}</p>
              <button type="button" disabled={isUpdating || isUpdatingAll} onClick={() => void refreshOne(skin.id)}>
                {isUpdating ? "Обновляем..." : "Обновить"}
              </button>
            </div>
          </div>
        );
      })}
    </div>
  );
};

const NewSkinsPage: React.FC = () => {
    const [query, setQuery] = useState("");
    const [debouncedQuery, setDebouncedQuery] = useState("");
    const [items, setItems] = useState<NewSkin[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [savingIds, setSavingIds] = useState<Record<string, boolean>>({});
    const [savedIds, setSavedIds] = useState<Record<string, boolean>>({});
    const [notice, setNotice] = useState<{ type: "success" | "error"; text: string } | null>(null);
    const requestIdRef = useRef(0);

    useEffect(() => {
      const timeout = window.setTimeout(() => {
        setDebouncedQuery(query.trim());
      }, 400);

      return () => window.clearTimeout(timeout);
    }, [query]);

    const loadNewSkins = async (searchValue: string) => {
      const requestId = requestIdRef.current + 1;
      requestIdRef.current = requestId;
      setLoading(true);
      setError(null);

      try {
        const response = await getNewSkins(searchValue);
        if (requestId !== requestIdRef.current) {
          return;
        }

        const normalizedQuery = searchValue.toLowerCase();
        const filteredItems = normalizedQuery
          ? response.items.filter((item) => item.title.toLowerCase().includes(normalizedQuery) || item.name.toLowerCase().includes(normalizedQuery))
          : response.items;

        setItems(filteredItems);
        setLoading(false);
      } catch (err: unknown) {
        if (requestId !== requestIdRef.current) {
          return;
        }

        const apiError = toApiError(err);
        setError(apiError.message || "Не удалось загрузить новые скины");
        setItems([]);
        setLoading(false);
      }
    };

    useEffect(() => {
      void loadNewSkins(debouncedQuery);
    }, [debouncedQuery]);

    const onSave = async (id: string) => {
      setSavingIds((prev) => ({ ...prev, [id]: true }));
      setNotice(null);

      try {
        await saveSkin(id);
        setSavedIds((prev) => ({ ...prev, [id]: true }));
        setNotice({ type: "success", text: "Скин сохранён." });
        await loadNewSkins(debouncedQuery);
      } catch (err: unknown) {
        const apiError = toApiError(err);
        setNotice({ type: "error", text: `Ошибка сохранения: ${apiError.message}` });
        window.alert(`Не удалось сохранить скин: ${apiError.message}`);
      } finally {
        setSavingIds((prev) => ({ ...prev, [id]: false }));
      }
    };

    return (
      <div className="new-skins-page">
        <div className="new-skins-toolbar">
          <input
            className="search-input"
            type="search"
            placeholder="Поиск по названию"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        {notice && <div className="status-box">{notice.text}</div>}
        {loading && <div className="status-box">Загрузка новых скинов...</div>}
        {error && !loading && <div className="status-box">Ошибка: {error}</div>}
        {!loading && !error && items.length === 0 && <div className="status-box">По запросу ничего не найдено.</div>}
        {!loading && !error && items.length > 0 && (
          <div className="container">
            {items.map((skin) => (
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
                  <p className="text">Лотов: {skin.sellListings ?? "-"}</p>
                  <button
                    type="button"
                    disabled={Boolean(savingIds[skin.id]) || Boolean(savedIds[skin.id])}
                    onClick={() => void onSave(skin.id)}
                  >
                    {savingIds[skin.id] ? "Сохраняем..." : savedIds[skin.id] ? "Сохранено" : "Сохранить"}
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    );
};

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
