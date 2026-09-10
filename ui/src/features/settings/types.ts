export interface ProviderState {
  binary: string;
  enabled: boolean;
}

export interface SettingsState {
  autostart_on: boolean;
  require_work_item: boolean;
  hooks: ProviderState[];
  skills: ProviderState[];
}

export interface SettingsPatch {
  autostart_on?: boolean;
  require_work_item?: boolean;
  hooks?: Record<string, boolean>;
  skills?: Record<string, boolean>;
}
