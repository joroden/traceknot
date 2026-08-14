export function LoadingList() {
  return (
    <div className="flex flex-col gap-1.5">
      {[0, 1, 2, 3].map((index) => (
        <div
          key={index}
          className="flex items-center gap-2.5 rounded-lg border border-zinc-800 px-3.5 py-2.5 light:border-zinc-200"
        >
          <span className="block size-3.5 flex-shrink-0 animate-pulse rounded-md bg-zinc-800 light:bg-zinc-200" />
          <span className="block h-[11px] w-[110px] flex-shrink-0 animate-pulse rounded-md bg-zinc-800 light:bg-zinc-200" />
          <span className="block h-[11px] flex-1 animate-pulse rounded-md bg-zinc-800 light:bg-zinc-200" />
        </div>
      ))}
    </div>
  );
}
