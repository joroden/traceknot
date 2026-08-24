import { LogoMark } from "../../../../components/Logo";

export function FirstRunHero() {
  return (
    <div className="mx-auto my-auto flex max-w-[560px] flex-col items-center gap-4 text-center">
      <span className="grid size-12 place-items-center rounded-lg bg-violet-600 text-white">
        <LogoMark size={24} />
      </span>
      <h2 className="text-lg font-bold">Waiting for your first session</h2>
      <p className="text-sm text-zinc-400 light:text-zinc-500">
        Open your coding agent and start a conversation — the session will
        show up here automatically, along with cost, tokens, coverage, and
        what the money bought.
      </p>
    </div>
  );
}
