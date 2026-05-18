import React from "react";
import { RouterProvider } from "./app/providers/router";
import { AppRouter } from "./app/router/AppRouter";

const App: React.FC = () => (
  <RouterProvider>
    <div className="app">
      <main className="layout-content">
        <AppRouter />
      </main>
    </div>
  </RouterProvider>
);

export default App;
