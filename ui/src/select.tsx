import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@fontsource/inter/400.css";
import "@fontsource/inter/500.css";
import "@fontsource/inter/600.css";
import "@fontsource/jetbrains-mono/400.css";
import "@fontsource/jetbrains-mono/500.css";
import "./styles/globals.css";
import { initTheme } from "./components/ThemeToggle";
import { PickerPage } from "./features/picker";

initTheme();

const root = document.getElementById("root");
if (!root) {
  throw new Error("missing #root element");
}

const sessionID = new URLSearchParams(window.location.search).get("session");

createRoot(root).render(
  <StrictMode>
    <PickerPage sessionID={sessionID} />
  </StrictMode>,
);
