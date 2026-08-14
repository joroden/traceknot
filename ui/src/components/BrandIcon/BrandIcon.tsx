import { GitHubMark } from "./GitHubMark";
import { JiraMark } from "./JiraMark";

export interface BrandIconProps {
  provider: string;
  size?: number;
}

export function BrandIcon({ provider, size = 14 }: BrandIconProps) {
  if (provider === "github") {
    return <GitHubMark size={size} />;
  }
  return <JiraMark size={size} />;
}
