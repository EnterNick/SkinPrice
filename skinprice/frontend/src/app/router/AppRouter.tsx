import React from "react";
import { Route, Routes } from "react-router-dom";
import { ROUTES } from "./routes";
import { SavedSkinsPage } from "../../pages/saved-skins-page/ui/SavedSkinsPage";
import { NewSkinsPage } from "../../pages/new-skins-page/ui/NewSkinsPage";
import { NotFoundPage } from "../../pages/not-found-page/ui/NotFoundPage";

export const AppRouter: React.FC = () => (
  <Routes>
    <Route path={ROUTES.home} element={<SavedSkinsPage />} />
    <Route path={ROUTES.newSkins} element={<NewSkinsPage />} />
    <Route path="*" element={<NotFoundPage />} />
  </Routes>
);
