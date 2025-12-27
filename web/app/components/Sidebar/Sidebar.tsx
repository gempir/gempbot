import {
  ChatBubbleLeftIcon,
  HomeIcon,
  ShieldCheckIcon,
  TrophyIcon,
  UserGroupIcon,
} from "@heroicons/react/24/solid";
import {
  Anchor,
  Box,
  Divider,
  NavLink,
  Stack,
  Text,
  ThemeIcon,
} from "@mantine/core";
import { Link, useRouterState } from "@tanstack/react-router";
import { useStore } from "../../store";
import { Login } from "./Login";
import { Managing } from "./Managing";

export function Sidebar() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });
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
      <Stack gap="sm">
        {/* Brand */}
        <Box>
          <Text
            size="xl"
            fw={700}
            style={{
              background: "linear-gradient(90deg, #00fa91 0%, #3b82f6 100%)",
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
              backgroundClip: "text",
            }}
          >
            gempbot
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
                to={link.href}
                label={link.label}
                active={isActive}
                leftSection={
                  <ThemeIcon
                    variant="subtle"
                    size="md"
                    color={isActive ? "cyan" : "gray"}
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
          to="/privacy"
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
