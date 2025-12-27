import {
  Card,
  Container,
  Group,
  Loader,
  Stack,
  Switch,
  Text,
  Title,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { useBotConfig } from "../../hooks/useBotConfig";

export function Bot() {
  const { config, loading, updateConfig } = useBotConfig();
  const [updating, setUpdating] = useState(false);

  const handleToggle = async (checked: boolean) => {
    setUpdating(true);
    try {
      await updateConfig({ predictionAnnouncements: checked });
      notifications.show({
        title: "Settings Updated",
        message: `Prediction announcements ${checked ? "enabled" : "disabled"}`,
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "Update Failed",
        message: "Failed to update bot settings",
        color: "red",
      });
    } finally {
      setUpdating(false);
    }
  };

  if (loading) {
    return (
      <Container size="lg">
        <Card shadow="sm" padding="xl" radius="md" withBorder>
          <Group justify="center" p="xl">
            <Loader size="lg" />
          </Group>
        </Card>
      </Container>
    );
  }

  return (
    <Container size="lg">
      <Stack gap="lg">
        <div>
          <Title order={1} mb="xs">
            Bot Settings
          </Title>
          <Text c="dimmed">
            Configure automated bot features for your channel
          </Text>
        </div>

        <Card shadow="sm" padding="xl" radius="md" withBorder>
          <Group align="flex-start" wrap="nowrap">
            <Stack gap="sm" style={{ flex: 1 }}>
              <Group justify="space-between" align="flex-start">
                <div>
                  <Title order={3} size="h4">
                    Prediction Announcements
                  </Title>
                  <Text size="sm" c="dimmed" mt={4}>
                    Automatically announce when predictions are created in your
                    channel. The bot will post a message in chat with prediction
                    details.
                  </Text>
                </div>

                <Switch
                  checked={config?.predictionAnnouncements || false}
                  onChange={(event) =>
                    handleToggle(event.currentTarget.checked)
                  }
                  disabled={updating}
                  size="lg"
                  color="cyan"
                  onLabel="ON"
                  offLabel="OFF"
                />
              </Group>
            </Stack>
          </Group>
        </Card>

        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Stack gap="xs">
            <Title order={4} size="h5">
              About Prediction Announcements
            </Title>
            <Text size="sm" c="dimmed">
              When enabled, gempbot will monitor your channel for new
              predictions and automatically post an announcement in chat. This
              helps increase viewer engagement by notifying everyone when a
              prediction starts.
            </Text>
          </Stack>
        </Card>
      </Stack>
    </Container>
  );
}
