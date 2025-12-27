import {
  ArrowRightIcon,
  ChatBubbleLeftIcon,
  ShieldCheckIcon,
  TrophyIcon,
  UserGroupIcon,
} from "@heroicons/react/24/solid";
import { Box, Grid, Group, Stack, Text, UnstyledButton } from "@mantine/core";
import { Link } from "@tanstack/react-router";
import { createLoginUrl } from "../factory/createLoginUrl";
import { useStore } from "../store";

const features = [
  {
    icon: TrophyIcon,
    title: "channel_point_rewards",
    description:
      "manage 7tv emotes via channel point redemptions. viewers add their favorite emotes automatically.",
    href: "/rewards",
    shortcut: "r",
  },
  {
    icon: ShieldCheckIcon,
    title: "emote_management",
    description:
      "block unwanted emotes and manage editor permissions. keep your emote list clean.",
    href: "/blocks",
    shortcut: "b",
  },
  {
    icon: ChatBubbleLeftIcon,
    title: "prediction_bot",
    description:
      "announce predictions in chat automatically with customizable messages and formatting.",
    href: "/bot",
    shortcut: "p",
  },
  {
    icon: UserGroupIcon,
    title: "user_permissions",
    description:
      "grant editor access to trusted users. control who manages your bot settings.",
    href: "/permissions",
    shortcut: "u",
  },
];

export function Teaser() {
  const isLoggedIn = useStore((state) => Boolean(state.scToken));
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);
  const twitchClientId = useStore((state) => state.twitchClientId);
  const loginUrl = createLoginUrl(apiBaseUrl, twitchClientId);

  return (
    <Box maw={900} mx="auto" py="xl">
      <Stack gap="xl">
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
              system online
            </Text>
          </Group>
          <Text
            size="xl"
            fw={700}
            ff="monospace"
            style={{ color: "var(--terminal-green)" }}
          >
            gempbot
          </Text>
          <Text size="sm" c="dimmed" ff="monospace" mt="xs">
            twitch bot & overlay management system
          </Text>
        </Box>

        {/* Features Grid */}
        <Grid gutter="md">
          {features.map((feature) => {
            const Icon = feature.icon;
            const isClickable = isLoggedIn;

            const content = (
              <Box
                p="md"
                className="card-interactive"
                style={{
                  border: "1px solid var(--border-subtle)",
                  backgroundColor: "var(--bg-elevated)",
                  cursor: isClickable ? "pointer" : "default",
                  height: "100%",
                  transition: "border-color 0.15s ease, background-color 0.15s ease",
                }}
                onMouseEnter={(e) => {
                  if (isClickable) {
                    e.currentTarget.style.borderColor = "var(--terminal-green)";
                    e.currentTarget.style.backgroundColor = "var(--bg-surface)";
                  }
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = "var(--border-subtle)";
                  e.currentTarget.style.backgroundColor = "var(--bg-elevated)";
                }}
              >
                <Stack gap="sm">
                  <Group justify="space-between" align="flex-start">
                    <Group gap="xs">
                      <Icon
                        style={{
                          width: 16,
                          height: 16,
                          color: "var(--terminal-green)",
                        }}
                      />
                      <Text size="sm" fw={600} ff="monospace" c="white">
                        {feature.title}
                      </Text>
                    </Group>
                    {isClickable && (
                      <Text size="xs" c="dimmed" ff="monospace">
                        [{feature.shortcut}]
                      </Text>
                    )}
                  </Group>

                  <Text size="xs" c="dimmed" ff="monospace" lh={1.5}>
                    {feature.description}
                  </Text>

                  {isClickable && (
                    <Group gap={4} mt="xs">
                      <ArrowRightIcon
                        style={{
                          width: 10,
                          height: 10,
                          color: "var(--terminal-green)",
                        }}
                      />
                      <Text size="xs" ff="monospace" c="terminal">
                        open
                      </Text>
                    </Group>
                  )}
                </Stack>
              </Box>
            );

            return (
              <Grid.Col key={feature.title} span={{ base: 12, sm: 6 }}>
                {isClickable ? (
                  <UnstyledButton
                    component={Link}
                    to={feature.href}
                    style={{ display: "block", height: "100%" }}
                  >
                    {content}
                  </UnstyledButton>
                ) : (
                  content
                )}
              </Grid.Col>
            );
          })}
        </Grid>

        {/* Call to Action for logged out users */}
        {!isLoggedIn && (
          <Box
            p="lg"
            style={{
              border: "1px solid var(--terminal-green)",
              backgroundColor: "var(--bg-elevated)",
            }}
          >
            <Stack gap="md" align="center">
              <Box ta="center">
                <Text size="sm" fw={600} ff="monospace" c="white" mb="xs">
                  {">"} ready to get started?
                </Text>
                <Text size="xs" c="dimmed" ff="monospace">
                  authenticate with twitch to access all features
                </Text>
              </Box>

              <UnstyledButton
                component="a"
                href={loginUrl.toString()}
                px="lg"
                py="sm"
                style={{
                  backgroundColor: "var(--terminal-green)",
                  color: "var(--bg-base)",
                  fontFamily: "'JetBrains Mono', monospace",
                  fontSize: "0.75rem",
                  fontWeight: 600,
                  textTransform: "uppercase",
                  letterSpacing: "0.1em",
                  border: "none",
                  cursor: "pointer",
                  transition: "opacity 0.15s ease",
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.opacity = "0.9";
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.opacity = "1";
                }}
              >
                login with twitch
              </UnstyledButton>
            </Stack>
          </Box>
        )}

        {/* System Status Footer */}
        <Box
          pt="md"
          style={{ borderTop: "1px solid var(--border-subtle)" }}
        >
          <Group justify="space-between">
            <Group gap="xl">
              <Group gap="xs">
                <Box className="status-dot status-online" />
                <Text size="xs" c="dimmed" ff="monospace">
                  api
                </Text>
              </Group>
              <Group gap="xs">
                <Box className="status-dot status-online" />
                <Text size="xs" c="dimmed" ff="monospace">
                  7tv
                </Text>
              </Group>
              <Group gap="xs">
                <Box className="status-dot status-online" />
                <Text size="xs" c="dimmed" ff="monospace">
                  twitch
                </Text>
              </Group>
            </Group>
            <Text size="xs" c="dimmed" ff="monospace" style={{ opacity: 0.5 }}>
              v2.0.0
            </Text>
          </Group>
        </Box>
      </Stack>
    </Box>
  );
}
