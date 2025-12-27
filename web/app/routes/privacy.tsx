import { Anchor, Box, Divider, List, Stack, Text } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/privacy")({
  component: Privacy,
});

function Privacy() {
  return (
    <Box maw={700} mx="auto" py="lg">
      <Stack gap="lg">
        {/* Header */}
        <Box>
          <Text size="lg" fw={600} ff="monospace" c="white">
            privacy_policy
          </Text>
          <Text size="xs" c="dimmed" ff="monospace" mt={4}>
            last updated: 2024
          </Text>
        </Box>

        <Divider color="var(--border-subtle)" />

        <Stack gap="lg">
          {/* Data Collection */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              data_collection
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" mb="sm" lh={1.6}>
              gempbot collects and stores the following information when you use
              our service:
            </Text>
            <List
              size="xs"
              spacing="xs"
              styles={{
                itemLabel: { fontFamily: "'JetBrains Mono', monospace" },
              }}
            >
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  your twitch user id and username
                </Text>
              </List.Item>
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  channel point reward configurations you create
                </Text>
              </List.Item>
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  emote history and blocked emotes for your channel
                </Text>
              </List.Item>
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  user permissions you grant to other users
                </Text>
              </List.Item>
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  overlay data created through our editor
                </Text>
              </List.Item>
            </List>
          </Box>

          {/* Data Usage */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              data_usage
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" mb="sm" lh={1.6}>
              your data is used exclusively to provide the bot and overlay
              management services:
            </Text>
            <List size="xs" spacing="xs">
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  manage channel point rewards and emote redemptions
                </Text>
              </List.Item>
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  store and serve your custom overlay configurations
                </Text>
              </List.Item>
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  authenticate your access to the dashboard
                </Text>
              </List.Item>
              <List.Item>
                <Text size="xs" c="dimmed" ff="monospace">
                  enforce permissions and access controls
                </Text>
              </List.Item>
            </List>
          </Box>

          {/* Third-Party Services */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              third_party_services
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" mb="sm" lh={1.6}>
              gempbot integrates with the following third-party services:
            </Text>
            <Stack gap="xs">
              <Box>
                <Text size="xs" ff="monospace" c="terminal">
                  twitch:
                </Text>
                <Text size="xs" c="dimmed" ff="monospace" ml="md">
                  authentication and channel point reward management
                </Text>
              </Box>
              <Box>
                <Text size="xs" ff="monospace" c="terminal">
                  7tv:
                </Text>
                <Text size="xs" c="dimmed" ff="monospace" ml="md">
                  emote management functionality
                </Text>
              </Box>
            </Stack>
            <Text size="xs" c="dimmed" ff="monospace" mt="sm" lh={1.6}>
              these services may collect their own data according to their
              respective privacy policies.
            </Text>
          </Box>

          {/* Data Storage */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              data_storage
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" lh={1.6}>
              all data is stored securely in our database. we implement
              appropriate technical and organizational measures to protect your
              personal information against unauthorized access, alteration,
              disclosure, or destruction.
            </Text>
          </Box>

          {/* Data Deletion */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              data_deletion
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" lh={1.6}>
              you can request deletion of your data at any time by contacting
              us. when you delete a channel point reward or overlay, the
              associated data is immediately removed from our systems.
            </Text>
          </Box>

          {/* Cookies */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              cookies
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" lh={1.6}>
              gempbot uses cookies to maintain your authentication session.
              these cookies are essential for the service to function and are
              not used for tracking or analytics purposes.
            </Text>
          </Box>

          {/* Changes */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              policy_changes
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" lh={1.6}>
              we may update this privacy policy from time to time. we will
              notify you of any changes by posting the new privacy policy on
              this page and updating the "last updated" date.
            </Text>
          </Box>

          {/* Contact */}
          <Box
            p="md"
            style={{
              border: "1px solid var(--terminal-green)",
              backgroundColor: "var(--bg-surface)",
            }}
          >
            <Text size="sm" fw={600} ff="monospace" c="white" mb="sm">
              contact
            </Text>
            <Text size="xs" c="dimmed" ff="monospace" lh={1.6}>
              if you have questions about this privacy policy or our data
              practices, contact us through{" "}
              <Anchor
                href="https://twitch.tv/gempir"
                target="_blank"
                rel="noopener noreferrer"
                c="terminal"
                size="xs"
              >
                twitch.tv/gempir
              </Anchor>
            </Text>
          </Box>
        </Stack>
      </Stack>
    </Box>
  );
}
