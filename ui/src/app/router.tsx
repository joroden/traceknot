import { lazy, Suspense } from "react";
import type { ComponentType, LazyExoticComponent } from "react";
import { createBrowserRouter, Navigate } from "react-router";
import { AppShell } from "./AppShell";

const Dashboard = lazy(() =>
  import("../features/dashboard").then((mod) => ({ default: mod.DashboardPage })),
);
const WorkItems = lazy(() =>
  import("../features/work-items").then((mod) => ({ default: mod.WorkItemsPage })),
);
const SessionDetail = lazy(() =>
  import("../features/session-detail").then((mod) => ({
    default: mod.SessionDetailPage,
  })),
);
const Unclaimed = lazy(() =>
  import("../features/unclaimed").then((mod) => ({ default: mod.UnclaimedPage })),
);

function lazyWrap(Component: LazyExoticComponent<ComponentType>) {
  return function LazyRoute() {
    return (
      <Suspense fallback={<div />}>
        <Component />
      </Suspense>
    );
  };
}

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppShell />,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: "dashboard", Component: lazyWrap(Dashboard) },
      { path: "work-items", Component: lazyWrap(WorkItems) },
      { path: "sessions/:id", Component: lazyWrap(SessionDetail) },
      { path: "unclaimed", Component: lazyWrap(Unclaimed) },
    ],
  },
]);

export interface PageTitle {
  path: string;
  title: string;
}

export const pageTitles: PageTitle[] = [
  { path: "/dashboard", title: "Dashboard" },
  { path: "/sessions/:id", title: "Session" },
  { path: "/work-items", title: "Work items" },
  { path: "/unclaimed", title: "Unclaimed sessions" },
];

export function titleForPathname(pathname: string): string {
  const exact = pageTitles.find((entry) => entry.path === pathname);
  if (exact) {
    return exact.title;
  }
  const param = pageTitles.find((entry) => {
    if (!entry.path.includes(":")) {
      return false;
    }
    const pattern = new RegExp(
      `^${entry.path.replace(/:[^/]+/g, "[^/]+")}$`,
    );
    return pattern.test(pathname);
  });
  return param?.title ?? "traceknot";
}
