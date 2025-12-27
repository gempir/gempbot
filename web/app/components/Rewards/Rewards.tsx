import { Container, Loader, Stack, Text, Title } from "@mantine/core";
import { useUserConfig } from "../../hooks/useUserConfig";
import { Emotehistory } from "./RewardForms/Emotehistory";
import { SevenTvForm } from "./RewardForms/SevenTvForm";

export function Rewards() {
  const [userConfig, , , loading] = useUserConfig();

  if (loading || !userConfig) {
    return (
      <Container size="xl">
        <Stack align="center" justify="center" h={400}>
          <Loader size="lg" />
          <Text c="dimmed">Loading rewards...</Text>
        </Stack>
      </Container>
    );
  }

  return (
    <Container size="xl">
      <Stack gap="lg">
        <div>
          <Title order={1} mb="xs">
            Channel Point Rewards
          </Title>
          <Text c="dimmed">Configure emote rewards for 7TV</Text>
        </div>

        <SevenTvForm userConfig={userConfig} />

        <Emotehistory />
      </Stack>
    </Container>
  );
}
