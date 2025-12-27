import { ChevronDownIcon } from "@heroicons/react/24/solid";
import { Select, Text } from "@mantine/core";
import { useUserConfig } from "../../hooks/useUserConfig";
import { setCookie } from "../../service/cookie";
import { useStore } from "../../store";

export function Managing() {
  const [userConfig] = useUserConfig();
  const managing = useStore((state) => state.managing);
  const setManaging = useStore((state) => state.setManaging);

  const handleChange = (value: string | null) => {
    const newValue = value && value.trim() !== "" ? value : null;
    setManaging(newValue);
    setCookie("managing", value || "");
  };

  const channels = userConfig?.Protected.EditorFor.sort() || [];
  const options = [
    { value: "", label: "own channel" },
    ...channels.map((channel) => ({
      value: channel,
      label: channel.toLowerCase(),
    })),
  ];

  return (
    <Select
      label={
        <Text size="xs" c="dimmed" fw={600} tt="uppercase" style={{ letterSpacing: "0.1em" }}>
          managing
        </Text>
      }
      placeholder="select channel"
      data={options}
      value={managing || ""}
      onChange={handleChange}
      rightSection={<ChevronDownIcon style={{ width: 12, height: 12 }} />}
      searchable
      clearable
      size="xs"
      styles={{
        input: {
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.75rem",
        },
      }}
    />
  );
}
