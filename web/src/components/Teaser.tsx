import {
  ChatBubbleLeftIcon,
  ShieldCheckIcon,
  TrophyIcon,
  UserGroupIcon,
} from "@heroicons/react/24/solid";
import {
  Button,
  Card,
  Container,
  Grid,
  Group,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useStore } from "../store";

const features = [
  {
    icon: TrophyIcon,
    title: "Channel Point Rewards",
    description:
      "Automatically manage BTTV and 7TV emotes with channel point redemptions. Let your viewers add their favorite emotes!",
    color: "cyan",
    href: "/rewards",
  },
  {
    icon: ShieldCheckIcon,
    title: "Emote Management",
    description:
      "Block unwanted emotes and manage permissions for editors. Keep your emote list clean and organized.",
    color: "violet",
    href: "/blocks",
  },
  {
    icon: ChatBubbleLeftIcon,
    title: "Prediction Bot",
    description:
      "Automatically announce predictions in chat with customizable messages and formatting.",
    color: "grape",
    href: "/bot",
  },
  {
    icon: UserGroupIcon,
    title: "User Permissions",
    description:
      "Grant editor access to trusted users and control who can manage your bot settings and predictions.",
    color: "blue",
    href: "/permissions",
  },
];

export function Teaser() {
  const isLoggedIn = useStore((state) => Boolean(state.scToken));

  return (
    <Container size="xl">
      <Stack gap="xl" py="xl">
        {/* Features Grid */}
        <Grid gutter="lg">
          {features.map((feature) => {
            const Icon = feature.icon;

            return (
              <Grid.Col key={feature.title} span={{ base: 12, sm: 6, md: 6 }}>
                <Card
                  shadow="sm"
                  padding="lg"
                  radius="md"
                  withBorder
                  h="100%"
                  component={(isLoggedIn ? "a" : "div") as any}
                  {...(isLoggedIn ? { href: feature.href } : {})}
                  style={{
                    cursor: isLoggedIn ? "pointer" : "default",
                    transition: "transform 0.2s",
                    textDecoration: "none",
                  }}
                  onMouseEnter={(e: React.MouseEvent<HTMLDivElement>) => {
                    if (isLoggedIn) {
                      e.currentTarget.style.transform = "translateY(-4px)";
                    }
                  }}
                  onMouseLeave={(e: React.MouseEvent<HTMLDivElement>) => {
                    e.currentTarget.style.transform = "translateY(0)";
                  }}
                >
                  <Stack gap="md">
                    <Group>
                      <div
                        style={{
                          background: `var(--mantine-color-${feature.color}-9)`,
                          borderRadius: "12px",
                          padding: "12px",
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                        }}
                      >
                        <Icon
                          style={{ width: 24, height: 24, color: "white" }}
                        />
                      </div>
                      <Title order={3} size="h4">
                        {feature.title}
                      </Title>
                    </Group>

                    <Text size="sm" c="dimmed">
                      {feature.description}
                    </Text>
                  </Stack>
                </Card>
              </Grid.Col>
            );
          })}
        </Grid>

        {/* Call to Action */}
        {!isLoggedIn && (
          <Card shadow="sm" padding="xl" radius="md" withBorder mt="xl">
            <Stack gap="md" ta="center">
              <Title order={2} size="h3">
                Ready to get started?
              </Title>
              <Text c="dimmed">
                Login with your Twitch account to start managing your channel
              </Text>
              <Group justify="center" mt="md">
                <Button
                  size="md"
                  variant="gradient"
                  gradient={{ from: "purple", to: "indigo", deg: 90 }}
                >
                  Login with Twitch
                </Button>
              </Group>
            </Stack>
          </Card>
        )}
      </Stack>
    </Container>
  );
}
