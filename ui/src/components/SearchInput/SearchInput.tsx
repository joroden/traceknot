import { Search } from "lucide-react";

export interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  autoFocus?: boolean;
}

export function SearchInput({ value, onChange, placeholder, autoFocus }: SearchInputProps) {
  return (
    <div className="flex items-center gap-1.5 rounded-md border border-zinc-700 bg-zinc-900 px-2 py-1.5 text-xs text-zinc-100 focus-within:border-violet-500 light:border-zinc-300 light:bg-white light:text-zinc-900">
      <Search className="h-3.5 w-3.5 shrink-0 text-zinc-500" />
      <input
        type="text"
        value={value}
        autoFocus={autoFocus}
        placeholder={placeholder}
        onChange={(event) => onChange(event.target.value)}
        className="w-full bg-transparent outline-none placeholder:text-zinc-500"
      />
    </div>
  );
}
