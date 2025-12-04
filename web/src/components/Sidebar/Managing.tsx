import { Select } from "@mantine/core";
import { UserGroupIcon } from "@heroicons/react/24/solid";
import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { setCookie } from "../../service/cookie";

export function Managing() {
  const { userConfig } = useUserConfig();
  const managing = useStore((state) => state.managing);
  const setManaging = useStore((state) => state.setManaging);

  const handleChange = (value: string | null) => {
    const newValue = value && value.trim() !== "" ? value : null;
    setManaging(newValue);
    setCookie("managing", value || "");
  };

  const channels = userConfig?.Protected.EditorFor.sort() || [];
  const options = [
    { value: "", label: "You (Own Channel)" },
    ...channels.map((channel) => ({
      value: channel,
      label: channel,
    })),
  ];

  return (
    <Select
      label="Managing Channel"
      placeholder="Select a channel"
      data={options}
      value={managing || ""}
      onChange={handleChange}
      leftSection={<UserGroupIcon style={{ width: 16, height: 16 }} />}
      searchable
      clearable
    />
  );
}
