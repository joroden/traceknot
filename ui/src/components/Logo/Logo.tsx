export interface LogoMarkProps {
  size?: number;
  className?: string;
}

export function LogoMark({ size = 16, className }: LogoMarkProps) {
  return (
    <svg width={size} height={size} viewBox="0 0 64 64" fill="none" className={className} aria-hidden="true">
      <path
        d="M20 44 C6 44 6 20 20 20 C30 20 30 32 40 32 C50 32 50 20 44 20"
        stroke="currentColor"
        strokeWidth={5}
        strokeLinecap="round"
      />
      <path
        d="M44 44 C58 44 58 20 44 20 C34 20 34 32 24 32"
        stroke="currentColor"
        strokeWidth={5}
        strokeLinecap="round"
      />
      <circle cx={20} cy={44} r={3} fill="currentColor" />
      <circle cx={44} cy={44} r={3} fill="currentColor" />
    </svg>
  );
}
