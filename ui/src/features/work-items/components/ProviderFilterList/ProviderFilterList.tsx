import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import { Check } from "lucide-react";

export interface ProviderFilterOption {
  value: string;
  label: string;
}

export interface ProviderFilterListProps {
  value: string;
  onChange: (value: string) => void;
  options: ProviderFilterOption[];
  placeholder: string;
}

const itemClass =
  "flex cursor-pointer items-center justify-between gap-2 rounded px-2 py-1.5 text-xs text-zinc-300 outline-none data-[highlighted]:bg-zinc-800 light:text-zinc-700 light:data-[highlighted]:bg-zinc-100";

export function ProviderFilterList({ value, onChange, options, placeholder }: ProviderFilterListProps) {
  return (
    <DropdownMenu.RadioGroup value={value} onValueChange={onChange}>
      <DropdownMenu.RadioItem value="" className={itemClass}>
        {placeholder}
        <DropdownMenu.ItemIndicator>
          <Check className="h-3.5 w-3.5" />
        </DropdownMenu.ItemIndicator>
      </DropdownMenu.RadioItem>
      {options.map((option) => (
        <DropdownMenu.RadioItem key={option.value} value={option.value} className={itemClass}>
          {option.label}
          <DropdownMenu.ItemIndicator>
            <Check className="h-3.5 w-3.5" />
          </DropdownMenu.ItemIndicator>
        </DropdownMenu.RadioItem>
      ))}
    </DropdownMenu.RadioGroup>
  );
}
