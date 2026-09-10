import type { ProviderState, SettingsPatch, SettingsState } from "./types";

export function stateToPatch(state: SettingsState): SettingsPatch {
  return {
    autostart_on: state.autostart_on,
    require_work_item: state.require_work_item,
    hooks: providersToRecord(state.hooks),
    skills: providersToRecord(state.skills),
  };
}

export function mergePatch(state: SettingsState, patch: SettingsPatch): SettingsState {
  return {
    ...state,
    autostart_on: patch.autostart_on ?? state.autostart_on,
    require_work_item: patch.require_work_item ?? state.require_work_item,
    hooks: patch.hooks ? mergeProviders(state.hooks, patch.hooks) : state.hooks,
    skills: patch.skills ? mergeProviders(state.skills, patch.skills) : state.skills,
  };
}

function providersToRecord(providers: ProviderState[]): Record<string, boolean> {
  return Object.fromEntries(providers.map((provider) => [provider.binary, provider.enabled]));
}

function mergeProviders(providers: ProviderState[], wants: Record<string, boolean>): ProviderState[] {
  return providers.map((provider) =>
    provider.binary in wants ? { ...provider, enabled: wants[provider.binary] } : provider,
  );
}
