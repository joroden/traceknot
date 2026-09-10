import { SettingsCard, SettingRow } from "../../../../components/SettingsCard";
import { SegmentedControl } from "../../../../components/SegmentedControl";

export interface WorkItemCardProps {
  requireWorkItem: boolean;
  onChange: (required: boolean) => void;
}

export function WorkItemCard({ requireWorkItem, onChange }: WorkItemCardProps) {
  return (
    <SettingsCard title="Work items">
      <SettingRow
        title="Work item assignment"
        description="When required, your agent won't be able to start working on a session until a work item is selected for it."
        control={
          <SegmentedControl
            options={[
              { value: "optional", label: "Optional" },
              { value: "required", label: "Required" },
            ]}
            value={requireWorkItem ? "required" : "optional"}
            onChange={(value) => onChange(value === "required")}
          />
        }
      />
    </SettingsCard>
  );
}
