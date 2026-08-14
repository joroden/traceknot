export function formatRelativeTime(unixMs: number | null): string {
  if (unixMs === null || unixMs <= 0) {
    return "";
  }
  const elapsedMs = Date.now() - unixMs;
  if (elapsedMs < 0) {
    return "now";
  }
  const seconds = Math.floor(elapsedMs / 1000);
  if (seconds < 60) {
    return "now";
  }
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) {
    return `${minutes}m ago`;
  }
  const hours = Math.floor(minutes / 60);
  if (hours < 24) {
    return `${hours}h ago`;
  }
  const days = Math.floor(hours / 24);
  if (days < 30) {
    return `${days}d ago`;
  }
  return new Date(unixMs).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

export function formatUSD(value: number): string {
  if (value >= 1000) {
    return `$${value.toFixed(0)}`;
  }
  if (value >= 100) {
    return `$${value.toFixed(1)}`;
  }
  if (value > 0 && value < 0.01) {
    return "<$0.01";
  }
  return `$${value.toFixed(2)}`;
}

export function formatCompactCount(value: number): string {
  if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(1)}M`;
  }
  if (value >= 1000) {
    return `${(value / 1000).toFixed(1)}K`;
  }
  return String(value);
}

export function formatPct(value: number | null): string {
  if (value === null) {
    return "—";
  }
  return `${value.toFixed(0)}%`;
}

export function formatDeltaPct(value: number | null): string {
  if (value === null) {
    return "";
  }
  const sign = value > 0 ? "+" : "";
  return `${sign}${value.toFixed(0)}%`;
}

export function formatBucketLabel(unixMs: number, granularityMs: number): string {
  const hourMs = 60 * 60 * 1000;
  const dayMs = 24 * hourMs;
  if (granularityMs <= hourMs) {
    return new Date(unixMs).toLocaleTimeString(undefined, {
      hour: "2-digit",
      minute: "2-digit",
    });
  }
  if (granularityMs <= dayMs) {
    return new Date(unixMs).toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
    });
  }
  return new Date(unixMs).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

export function formatDuration(ms: number | null | undefined): string {
  if (ms === null || ms === undefined || ms < 0 || !Number.isFinite(ms)) {
    return "—";
  }
  if (ms < 1000) {
    return `${ms.toFixed(0)}ms`;
  }
  const totalSeconds = ms / 1000;
  if (totalSeconds < 60) {
    return `${totalSeconds.toFixed(1)}s`;
  }
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = Math.round(totalSeconds % 60);
  if (minutes < 60) {
    return `${minutes}m ${seconds}s`;
  }
  const hours = Math.floor(minutes / 60);
  const restMinutes = minutes % 60;
  if (hours < 24) {
    return `${hours}h ${restMinutes}m`;
  }
  const days = Math.floor(hours / 24);
  return `${days}d ${hours % 24}h`;
}

export function formatTokens(value: number): string {
  if (value <= 0) {
    return "—";
  }
  return formatCompactCount(value);
}

export function formatProviderLabel(provider: string): string {
  return provider
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

export function formatTimestamp(unixMs: number | null | undefined): string {
  if (!unixMs || unixMs <= 0) {
    return "—";
  }
  return new Date(unixMs).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}
