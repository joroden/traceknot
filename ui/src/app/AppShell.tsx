import { useLocation } from "react-router";
import { Outlet } from "react-router";
import { Sidebar } from "./Sidebar";
import { titleForPathname } from "./router";

export function AppShell() {
  const location = useLocation();
  return (
    <div className="flex h-full">
      <Sidebar />
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-[52px] items-center border-b border-zinc-800 bg-zinc-950 px-6 light:border-zinc-200 light:bg-zinc-50">
          <h1 className="text-base font-semibold">
            {titleForPathname(location.pathname)}
          </h1>
        </header>
        <main className="flex min-h-0 flex-1 flex-col overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
