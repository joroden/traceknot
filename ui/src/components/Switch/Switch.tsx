export interface SwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
  label: string;
}

export function Switch({ checked, onChange, disabled = false, label }: SwitchProps) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      aria-label={label}
      disabled={disabled}
      onClick={() => onChange(!checked)}
      className={
        checked
          ? "relative h-6 w-10 flex-shrink-0 cursor-pointer rounded-full bg-violet-600 transition-colors disabled:cursor-default disabled:opacity-60"
          : "relative h-6 w-10 flex-shrink-0 cursor-pointer rounded-full bg-zinc-700 transition-colors disabled:cursor-default disabled:opacity-60 light:bg-zinc-300"
      }
    >
      <span
        className={
          checked
            ? "absolute top-0.5 left-4.5 size-5 rounded-full bg-white transition-[left]"
            : "absolute top-0.5 left-0.5 size-5 rounded-full bg-white transition-[left]"
        }
      />
    </button>
  );
}
