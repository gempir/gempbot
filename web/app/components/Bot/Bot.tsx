import { Box, Group, Loader, Stack, Switch, Text } from "@mantine/core";
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
        title: "updated",
        message: `prediction announcements ${checked ? "enabled" : "disabled"}`,
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "error",
        message: "failed to update settings",
        color: "red",
      });
    } finally {
      setUpdating(false);
    }
  };

  if (loading) {
    return (
      <Box maw={700} mx="auto">
        <Box
          p="lg"
          style={{
            border: "1px solid var(--border-subtle)",
            backgroundColor: "var(--bg-elevated)",
          }}
        >
          <Group justify="center" p="xl">
            <Loader size="sm" />
          </Group>
        </Box>
      </Box>
    );
  }

  return (
    <Box maw={700} mx="auto">
      <Stack gap="lg">
        {/* Header */}
        <Box>
          <Text size="lg" fw={600} ff="monospace" c="white">
            bot_settings
          </Text>
          <Text size="xs" c="dimmed" ff="monospace" mt={4}>
            configure automated bot features for your channel
          </Text>
        </Box>

        {/* Prediction Announcements */}
        <Box
          p="md"
          style={{
            border: "1px solid var(--border-subtle)",
            backgroundColor: "var(--bg-elevated)",
          }}
        >
          <Group align="flex-start" wrap="nowrap" justify="space-between">
            <Stack gap="xs" style={{ flex: 1 }}>
              <Group gap="xs">
                <Box
                  className={
                    config?.predictionAnnouncements
                      ? "status-dot status-online"
                      : "status-dot status-offline"
                  }
                />
                <Text size="sm" fw={600} ff="monospace" c="white">
                  prediction_announcements
                </Text>
              </Group>
              <Text size="xs" c="dimmed" ff="monospace" lh={1.5}>
                automatically announce when predictions are created in your
                channel. the bot will post a message in chat with prediction
                details.
              </Text>
            </Stack>

            <Switch
              checked={config?.predictionAnnouncements || false}
              onChange={(event) => handleToggle(event.currentTarget.checked)}
              disabled={updating}
              size="md"
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
        </Box>

        {/* Info Section */}
        <Box
          p="md"
          style={{
            border: "1px solid var(--border-subtle)",
            backgroundColor: "var(--bg-surface)",
          }}
        >
          <Stack gap="sm">
            <Text
              size="xs"
              fw={600}
              ff="monospace"
              c="dimmed"
              tt="uppercase"
              style={{ letterSpacing: "0.1em" }}
            >
              info
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" lh={1.6}>
              when enabled, gempbot monitors your channel for new predictions
              and automatically posts an announcement in chat. this increases
              viewer engagement by notifying everyone when a prediction starts.
            </Text>
          </Stack>
        </Box>

        {/* Status Footer */}
        <Box pt="md" style={{ borderTop: "1px solid var(--border-subtle)" }}>
          <Group gap="xl">
            <Group gap="xs">
              <Box
                className={
                  config?.predictionAnnouncements
                    ? "status-dot status-online"
                    : "status-dot status-offline"
                }
              />
              <Text size="xs" c="dimmed" ff="monospace">
                predictions:{" "}
                {config?.predictionAnnouncements ? "active" : "inactive"}
              </Text>
            </Group>
          </Group>
        </Box>
      </Stack>
    </Box>
  );
}
