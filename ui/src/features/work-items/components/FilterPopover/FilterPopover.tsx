import type { ReactNode, KeyboardEvent } from "react";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import { Filter } from "lucide-react";

export interface FilterPopoverProps {
  active: boolean;
  children: ReactNode;
}

export function FilterPopover({ active, children }: FilterPopoverProps) {
  const stopMenuKeyHandling = (event: KeyboardEvent) => {
    event.stopPropagation();
  };

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        onClick={(event) => event.stopPropagation()}
        className="relative rounded p-0.5 text-zinc-500 outline-none hover:text-zinc-200 light:hover:text-zinc-800"
      >
        <Filter className="h-3 w-3" />
        {active ? (
          <span className="absolute -right-0.5 -top-0.5 h-1.5 w-1.5 rounded-full bg-violet-500" />
        ) : null}
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          align="start"
          onCloseAutoFocus={(event) => event.preventDefault()}
          onKeyDown={stopMenuKeyHandling}
          className="z-20 min-w-[200px] rounded-md border border-zinc-700 bg-zinc-900 p-2.5 shadow-[0_12px_40px_rgba(0,0,0,0.5)] light:border-zinc-300 light:bg-white"
        >
          {children}
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  );
}
