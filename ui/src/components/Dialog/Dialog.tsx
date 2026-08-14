import type { ReactNode } from "react";
import * as DialogPrimitive from "@radix-ui/react-dialog";

export interface DialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  children: ReactNode;
  widthClassName?: string;
}

export function Dialog({
  open,
  onOpenChange,
  title,
  description,
  children,
  widthClassName = "w-[min(440px,calc(100vw-48px))]",
}: DialogProps) {
  return (
    <DialogPrimitive.Root open={open} onOpenChange={onOpenChange}>
      <DialogPrimitive.Portal>
        <DialogPrimitive.Overlay className="fixed inset-0 bg-black/60 backdrop-blur-sm" />
        <DialogPrimitive.Content
          className={`fixed left-1/2 top-1/2 ${widthClassName} -translate-x-1/2 -translate-y-1/2 rounded-lg border border-zinc-700 bg-zinc-900 p-5 shadow-[0_12px_40px_rgba(0,0,0,0.5)] light:border-zinc-300 light:bg-white`}
        >
          <DialogPrimitive.Title className="text-sm font-semibold outline-none">
            {title}
          </DialogPrimitive.Title>
          {description ? (
            <DialogPrimitive.Description className="mt-2 break-words text-sm text-zinc-400 outline-none light:text-zinc-500">
              {description}
            </DialogPrimitive.Description>
          ) : null}
          {children}
        </DialogPrimitive.Content>
      </DialogPrimitive.Portal>
    </DialogPrimitive.Root>
  );
}
