import React, { useEffect, useState } from "react";
import { getSavedSkins, updateAllSkinPrices, updateSkinPrice } from "./api/client";
import type { SavedSkin } from "./api/models";
import { toApiError } from "./api/errors";
import "./styles.css";

type SavedSkinsState = {
  items: SavedSkin[];
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

  const loadSkins = () => getSavedSkins()
    .then((response) => setState({ items: response.items, loading: false, error: null }));

  useEffect(() => {
    void loadSkins()
      .catch((err: unknown) => {
        setState({
          items: [],
          loading: false,
          error: toApiError(err).message,
        });
      });
  }, []);

  const refreshOne = (skinId: string) => updateSkinPrice(skinId, currency).then(loadSkins);
  const refreshAll = () => updateAllSkinPrices(currency).then(loadSkins);

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
              <button onClick={() => void refreshOne(skin.id)}>Обновить цену</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default App;
