import {
  ChatBubbleLeftIcon,
  HomeIcon,
  ShieldCheckIcon,
  TrophyIcon,
  UserGroupIcon,
} from "@heroicons/react/24/solid";
import {
  Box,
  Divider,
  NavLink,
  Stack,
  Text,
  UnstyledButton,
} from "@mantine/core";
import { Link, useRouterState } from "@tanstack/react-router";
import { useStore } from "../../store";
import { Login } from "./Login";
import { Managing } from "./Managing";

export function Sidebar() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });
  const isLoggedIn = useStore((state) => Boolean(state.scToken));

  const navLinks = [
    { href: "/", label: "home", icon: HomeIcon },
    {
      href: "/rewards",
      label: "rewards",
      icon: TrophyIcon,
      requiresAuth: true,
    },
    {
      href: "/permissions",
      label: "permissions",
      icon: UserGroupIcon,
      requiresAuth: true,
    },
    {
      href: "/bot",
      label: "bot",
      icon: ChatBubbleLeftIcon,
      requiresAuth: true,
    },
    {
      href: "/blocks",
      label: "blocks",
      icon: ShieldCheckIcon,
      requiresAuth: true,
    },
  ];

  return (
    <Stack h="100%" justify="space-between" gap={0}>
      <Stack gap="sm">
        {/* Brand */}
        <Box py="xs">
          <Text
            size="sm"
            fw={700}
            ff="monospace"
            style={{
              color: "var(--terminal-green)",
              letterSpacing: "0.1em",
            }}
          >
            gempbot
          </Text>
          <Text size="xs" c="dimmed" ff="monospace" mt={2}>
            v2.0 {"// twitch bot"}
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
        <Box>
          <Text
            size="xs"
            c="dimmed"
            fw={600}
            tt="uppercase"
            mb="xs"
            style={{ letterSpacing: "0.1em" }}
          >
            navigation
          </Text>
          <Stack gap={2}>
            {navLinks.map((link) => {
              if (link.requiresAuth && !isLoggedIn) return null;

              const Icon = link.icon;
              const isActive = pathname === link.href;

              return (
                <NavLink
                  key={link.href}
                  component={Link}
                  to={link.href}
                  label={
                    <Text size="xs" ff="monospace">
                      {link.label}
                    </Text>
                  }
                  active={isActive}
                  leftSection={
                    <Icon
                      style={{
                        width: 14,
                        height: 14,
                        opacity: isActive ? 1 : 0.5,
                        color: isActive
                          ? "var(--terminal-green)"
                          : "var(--text-secondary)",
                      }}
                    />
                  }
                  styles={{
                    root: {
                      borderLeft: isActive
                        ? "2px solid var(--terminal-green)"
                        : "2px solid transparent",
                      backgroundColor: isActive
                        ? "var(--bg-surface)"
                        : "transparent",
                      padding: "0.375rem 0.5rem",
                    },
                    label: {
                      color: isActive
                        ? "var(--terminal-green)"
                        : "var(--text-secondary)",
                    },
                  }}
                />
              );
            })}
          </Stack>
        </Box>
      </Stack>

      {/* Footer */}
      <Box>
        <Divider mb="sm" />
        <Stack gap={4}>
          <UnstyledButton
            component={Link}
            to="/privacy"
            style={{
              display: "block",
              padding: "0.25rem 0",
            }}
          >
            <Text size="xs" c="dimmed" ff="monospace">
              [privacy]
            </Text>
          </UnstyledButton>
          <Text size="xs" c="dimmed" ff="monospace" style={{ opacity: 0.5 }}>
            © 2024 gempir
          </Text>
        </Stack>
      </Box>
    </Stack>
  );
}
