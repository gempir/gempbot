import { UserIcon } from "@heroicons/react/24/solid";
import { Button } from "@mantine/core";
import { createLoginUrl } from "../../factory/createLoginUrl";
import { useStore } from "../../store";

export function Login() {
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);
  const twitchClientId = useStore((state) => state.twitchClientId);
  const isLoggedIn = useStore((state) => Boolean(state.scToken));
  const url = createLoginUrl(apiBaseUrl, twitchClientId);

  return (
    <Button
      component="a"
      href={url.toString()}
      variant={isLoggedIn ? "subtle" : "gradient"}
      gradient={{ from: "cyan", to: "blue", deg: 90 }}
      size="md"
      fullWidth
      leftSection={<UserIcon style={{ width: 20, height: 20 }} />}
    >
      {isLoggedIn ? "Logged in" : "Login with Twitch"}
    </Button>
  );
}
