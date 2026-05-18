import React from "react";
import { BrowserRouter } from "react-router-dom";

type RouterProviderProps = {
  children: React.ReactNode;
};

export const RouterProvider: React.FC<RouterProviderProps> = ({ children }) => (
  <BrowserRouter>{children}</BrowserRouter>
);
