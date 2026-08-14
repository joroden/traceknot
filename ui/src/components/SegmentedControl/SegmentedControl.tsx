export interface SegmentedOption<T extends string> {
  value: T;
  label: string;
}

export interface SegmentedControlProps<T extends string> {
  options: SegmentedOption<T>[];
  value: T;
  onChange: (value: T) => void;
}

export function SegmentedControl<T extends string>({
  options,
  value,
  onChange,
}: SegmentedControlProps<T>) {
  return (
    <div className="inline-flex rounded-lg border border-zinc-800 bg-zinc-900 p-0.5 light:border-zinc-200 light:bg-white">
      {options.map((option) => (
        <button
          key={option.value}
          type="button"
          className={
            option.value === value
              ? "cursor-pointer rounded-md bg-zinc-800 px-3 py-1.5 text-sm font-semibold text-zinc-100 transition-colors light:bg-zinc-200 light:text-zinc-900"
              : "cursor-pointer rounded-md px-3 py-1.5 text-sm text-zinc-400 transition-colors hover:text-zinc-100 light:text-zinc-500 light:hover:text-zinc-900"
          }
          onClick={() => onChange(option.value)}
        >
          {option.label}
        </button>
      ))}
    </div>
  );
}
