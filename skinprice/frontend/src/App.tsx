import React, { useEffect, useState } from "react";
import { GetSavedSkins, UpdateAllSavedSkinsPrices, UpdateSavedSkinPrice, type SavedSkinsResponse } from "./wailsjs/go/main/App";
import "./styles.css";

type SavedSkinsState = {
  items: SavedSkinsResponse["items"];
  loading: boolean;
  error: string | null;
};

const App: React.FC = () => {
  const [state, setState] = useState<SavedSkinsState>({
    items: [],
    loading: true,
    error: null,
  });
  const [currency, setCurrency] = useState("1");

  const loadSkins = () => GetSavedSkins({ limit: 50, offset: 0 })
    .then((response) => setState({ items: response.items, loading: false, error: null }));

  useEffect(() => {
    void loadSkins()
      .catch((err: unknown) => {
        setState({
          items: [],
          loading: false,
          error: err instanceof Error ? err.message : "Не удалось загрузить сохранённые скины",
        });
      });
  }, []);

  const refreshOne = (marketHashName: string) => UpdateSavedSkinPrice({ market_hash_name: marketHashName, currency }).then(loadSkins);
  const refreshAll = () => UpdateAllSavedSkinsPrices({ currency }).then(loadSkins);

  if (state.loading) {
    return <div className="app"><div className="container">Загрузка...</div></div>;
  }

  if (state.error) {
    return <div className="app"><div className="container">Ошибка: {state.error}</div></div>;
  }

  return (
    <div className="app">
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
    </div>
  );
};

export default App;
