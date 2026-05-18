import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";

export const openExternal = (url: string) => {
  if (!url) return;
  BrowserOpenURL(url);
};
