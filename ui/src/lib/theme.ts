const STORAGE_KEY = "traceknot-theme";

export type Theme = "dark" | "light";

export function getTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "dark" || stored === "light") {
    return stored;
  }
  return "dark";
}

export function applyTheme(theme: Theme): void {
  if (theme === "light") {
    document.documentElement.dataset.theme = "light";
  } else {
    delete document.documentElement.dataset.theme;
  }
  localStorage.setItem(STORAGE_KEY, theme);
}

export function toggleTheme(): Theme {
  const next: Theme = getTheme() === "dark" ? "light" : "dark";
  applyTheme(next);
  return next;
}
