import { Badge } from "../Badge";

export interface ModelBadgesProps {
  models: string[] | null;
}

export function ModelBadges({ models }: ModelBadgesProps) {
  if (!models || models.length === 0) {
    return <span className="text-xs text-zinc-600">—</span>;
  }
  return (
    <div className="flex flex-wrap items-center gap-1">
      {models.map((model) => (
        <Badge key={model} tone="neutral">
          {model}
        </Badge>
      ))}
    </div>
  );
}
