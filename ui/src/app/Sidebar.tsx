import { useEffect, useState } from "react";
import { Inbox, Layers, LayoutDashboard, PanelLeftClose, PanelLeftOpen } from "lucide-react";
import { NavLink } from "react-router";
import { LogoMark } from "../components/Logo";
import { ThemeToggle } from "../components/ThemeToggle";

const STORAGE_KEY = "traceknot.sidebarCollapsed";

const NAV_ITEMS = [
  { to: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { to: "/work-items", label: "Work Items", icon: Layers },
  { to: "/unclaimed", label: "Unclaimed", icon: Inbox },
];

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem(STORAGE_KEY) === "1");

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, collapsed ? "1" : "0");
  }, [collapsed]);

  return (
    <aside
      className={`flex shrink-0 flex-col border-r border-zinc-800 bg-zinc-900 transition-[width] duration-200 light:border-zinc-200 light:bg-white ${
        collapsed ? "w-[52px]" : "w-[216px]"
      }`}
    >
      <div
        className={`flex h-[52px] items-center border-b border-zinc-800 light:border-zinc-200 ${
          collapsed ? "justify-center" : "gap-2 px-4"
        }`}
      >
        {!collapsed && (
          <>
            <span className="grid size-[26px] place-items-center rounded-md bg-violet-600 text-white">
              <LogoMark size={16} />
            </span>
            <span className="font-mono text-sm font-medium tracking-tight">traceknot</span>
          </>
        )}
        <button
          type="button"
          onClick={() => setCollapsed((prev) => !prev)}
          className={`cursor-pointer rounded p-1.5 text-zinc-500 transition-colors hover:bg-zinc-800 hover:text-zinc-200 light:hover:bg-zinc-100 light:hover:text-zinc-800 ${
            collapsed ? "" : "ml-auto"
          }`}
          title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          {collapsed ? <PanelLeftOpen size={15} /> : <PanelLeftClose size={15} />}
        </button>
      </div>
      <nav className="flex flex-col gap-0.5 px-2 py-3">
        {NAV_ITEMS.map(({ to, label, icon: Icon }) => (
          <NavLink
            key={to}
            to={to}
            title={collapsed ? label : undefined}
            className={({ isActive }) =>
              isActive
                ? `flex items-center rounded-md bg-zinc-800 text-sm font-medium text-zinc-100 no-underline transition-colors light:bg-zinc-100 light:text-zinc-900 ${
                    collapsed ? "justify-center py-2" : "gap-2.5 px-2.5 py-2"
                  }`
                : `flex items-center rounded-md text-sm text-zinc-400 no-underline transition-colors hover:bg-zinc-800 hover:text-zinc-100 light:text-zinc-500 light:hover:bg-zinc-100 light:hover:text-zinc-900 ${
                    collapsed ? "justify-center py-2" : "gap-2.5 px-2.5 py-2"
                  }`
            }
          >
            <Icon size={15} />
            {!collapsed && label}
          </NavLink>
        ))}
      </nav>
      <div
        className={`mt-auto flex items-center border-t border-zinc-800 py-3 light:border-zinc-200 ${
          collapsed ? "justify-center" : "px-4"
        }`}
      >
        <ThemeToggle />
      </div>
    </aside>
  );
}
