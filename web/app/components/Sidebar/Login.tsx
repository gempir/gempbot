import { ArrowRightOnRectangleIcon, CheckCircleIcon } from "@heroicons/react/24/solid";
import { Button, Group, Text } from "@mantine/core";
import { createLoginUrl } from "../../factory/createLoginUrl";
import { useStore } from "../../store";

export function Login() {
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);
  const twitchClientId = useStore((state) => state.twitchClientId);
  const isLoggedIn = useStore((state) => Boolean(state.scToken));
  const scTokenContent = useStore((state) => state.scTokenContent);
  const url = createLoginUrl(apiBaseUrl, twitchClientId);

  if (isLoggedIn) {
    return (
      <Group gap="xs" py="xs">
        <CheckCircleIcon
          style={{
            width: 12,
            height: 12,
            color: "var(--terminal-green)",
          }}
        />
        <Text size="xs" ff="monospace" c="dimmed">
          logged in as{" "}
          <Text span c="white" inherit>
            {scTokenContent?.login || "user"}
          </Text>
        </Text>
      </Group>
    );
  }

  return (
    <Button
      component="a"
      href={url.toString()}
      variant="outline"
      color="terminal"
      size="xs"
      fullWidth
      leftSection={<ArrowRightOnRectangleIcon style={{ width: 14, height: 14 }} />}
      styles={{
        root: {
          borderColor: "var(--terminal-green)",
          color: "var(--terminal-green)",
          "&:hover": {
            backgroundColor: "var(--terminal-green)",
            color: "var(--bg-base)",
          },
        },
      }}
    >
      login with twitch
    </Button>
  );
}
