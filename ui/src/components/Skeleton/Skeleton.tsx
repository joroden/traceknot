export interface SkeletonProps {
  className?: string;
}

export function Skeleton({ className = "" }: SkeletonProps) {
  return <div className={`animate-pulse rounded-md bg-zinc-900 light:bg-zinc-100 ${className}`} />;
}
