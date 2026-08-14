import { useState } from "react";
import { Moon, Sun } from "lucide-react";
import { applyTheme, getTheme, toggleTheme } from "../../lib/theme";

export function ThemeToggle() {
  const [theme, setTheme] = useState(() => getTheme());

  const cycle = () => {
    const next = toggleTheme();
    setTheme(next);
  };

  return (
    <button
      type="button"
      className="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg border border-zinc-700 bg-zinc-800 text-zinc-400 transition-colors hover:border-violet-500 hover:text-zinc-100 light:border-zinc-300 light:bg-zinc-100 light:text-zinc-500 light:hover:text-zinc-900"
      onClick={cycle}
      aria-label={theme === "dark" ? "Switch to light theme" : "Switch to dark theme"}
      title={theme === "dark" ? "Switch to light theme" : "Switch to dark theme"}
    >
      {theme === "dark" ? <Sun size={13} /> : <Moon size={13} />}
    </button>
  );
}

export function initTheme(): void {
  applyTheme(getTheme());
}
