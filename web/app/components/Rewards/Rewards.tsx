import { Box, Group, Loader, Stack, Text } from "@mantine/core";
import { useUserConfig } from "../../hooks/useUserConfig";
import { Emotehistory } from "./RewardForms/Emotehistory";
import { SevenTvForm } from "./RewardForms/SevenTvForm";

export function Rewards() {
  const [userConfig, , , loading] = useUserConfig();

  if (loading || !userConfig) {
    return (
      <Box maw={900} mx="auto">
        <Stack align="center" justify="center" h={400}>
          <Loader size="sm" />
          <Text c="dimmed" size="xs" ff="monospace">
            loading rewards...
          </Text>
        </Stack>
      </Box>
    );
  }

  return (
    <Box maw={900} mx="auto">
      <Stack gap="lg">
        {/* Header */}
        <Box>
          <Group gap="xs" mb="xs">
            <Box
              style={{
                width: 8,
                height: 8,
                backgroundColor: "var(--terminal-green)",
                boxShadow: "0 0 8px var(--terminal-green)",
              }}
            />
            <Text size="xs" c="dimmed" ff="monospace" tt="uppercase" style={{ letterSpacing: "0.1em" }}>
              rewards active
            </Text>
          </Group>
          <Text size="lg" fw={600} ff="monospace" c="white">
            channel_point_rewards
          </Text>
          <Text size="xs" c="dimmed" ff="monospace" mt={4}>
            configure emote rewards for 7tv integration
          </Text>
        </Box>

        <SevenTvForm userConfig={userConfig} />

        <Emotehistory />
      </Stack>
    </Box>
  );
}
