import { Box, Group, Switch, Text } from "@mantine/core";

interface ToggleProps {
  label: string;
  description?: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
}

export function Toggle({
  label,
  description,
  checked,
  onChange,
  disabled,
}: ToggleProps) {
  return (
    <Group align="flex-start" wrap="nowrap" justify="space-between" gap="lg">
      <Box style={{ flex: 1 }}>
        <Group gap="xs" mb={description ? 4 : 0}>
          <Box
            className={
              checked ? "status-dot status-online" : "status-dot status-offline"
            }
          />
          <Text size="xs" fw={600} ff="monospace" c="white">
            {label}
          </Text>
        </Group>
        {description && (
          <Text size="xs" c="dimmed" ff="monospace" lh={1.5}>
            {description}
          </Text>
        )}
      </Box>

      <Switch
        checked={checked}
        onChange={(event) => onChange(event.currentTarget.checked)}
        disabled={disabled}
        size="sm"
        color="terminal"
        onLabel="on"
        offLabel="off"
        styles={{
          track: {
            borderRadius: 0,
          },
          thumb: {
            borderRadius: 0,
          },
        }}
      />
    </Group>
  );
}
