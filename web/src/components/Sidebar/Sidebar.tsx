import {
  Anchor,
  Box,
  Divider,
  Group,
  NavLink,
  Stack,
  Text,
  ThemeIcon,
} from "@mantine/core";
import {
  ChatBubbleLeftIcon,
  HomeIcon,
  ShieldCheckIcon,
  TrophyIcon,
  UserGroupIcon,
} from "@heroicons/react/24/solid";
import { usePathname } from "next/navigation";
import Link from "next/link";
import { Login } from "./Login";
import { Managing } from "./Managing";
import { useStore } from "../../store";

export function Sidebar() {
  const pathname = usePathname();
  const isLoggedIn = useStore((state) => Boolean(state.scToken));

  const navLinks = [
    { href: "/", label: "Home", icon: HomeIcon },
    {
      href: "/rewards",
      label: "Rewards",
      icon: TrophyIcon,
      requiresAuth: true,
    },
    {
      href: "/permissions",
      label: "Permissions",
      icon: UserGroupIcon,
      requiresAuth: true,
    },
    {
      href: "/bot",
      label: "Bot",
      icon: ChatBubbleLeftIcon,
      requiresAuth: true,
    },
    {
      href: "/blocks",
      label: "Blocks",
      icon: ShieldCheckIcon,
      requiresAuth: true,
    },
  ];

  return (
    <Stack h="100%" justify="space-between">
      <Stack gap="md">
        {/* Brand */}
        <Box mb="md">
          <Text
            size="xl"
            fw={700}
            variant="gradient"
            gradient={{ from: "purple", to: "indigo", deg: 90 }}
          >
            gempbot
          </Text>
          <Text size="xs" c="dimmed">
            Twitch Bot & Overlay Manager
          </Text>
        </Box>

        <Divider />

        {/* Login */}
        <Login />

        {/* Managing Channel Selector */}
        {isLoggedIn && (
          <>
            <Managing />
            <Divider />
          </>
        )}

        {/* Navigation Links */}
        <Stack gap="xs">
          {navLinks.map((link) => {
            if (link.requiresAuth && !isLoggedIn) return null;

            const Icon = link.icon;
            const isActive = pathname === link.href;

            return (
              <NavLink
                key={link.href}
                component={Link}
                href={link.href}
                label={link.label}
                active={isActive}
                leftSection={
                  <ThemeIcon
                    variant="subtle"
                    size="md"
                    color={isActive ? "purple" : "gray"}
                  >
                    <Icon style={{ width: "70%", height: "70%" }} />
                  </ThemeIcon>
                }
              />
            );
          })}
        </Stack>
      </Stack>

      {/* Footer */}
      <Box pt="md">
        <Divider mb="sm" />
        <Anchor
          component={Link}
          href="/privacy"
          size="sm"
          c="dimmed"
          ta="center"
          w="100%"
        >
          Privacy Policy
        </Anchor>
      </Box>
    </Stack>
  );
}
