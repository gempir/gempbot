import { Box, Loader, Stack, Text } from "@mantine/core";
import { useUserConfig } from "../../hooks/useUserConfig";
import { UserPermissions } from "./UserPermissions";

export function Permissions() {
  const [userConfig, setUserConfig, , loading, errorMessage] = useUserConfig();

  if (loading || !userConfig) {
    return (
      <Box maw={800} mx="auto">
        <Stack align="center" justify="center" h={400}>
          <Loader size="sm" />
          <Text c="dimmed" size="xs" ff="monospace">
            loading permissions...
          </Text>
        </Stack>
      </Box>
    );
  }

  return (
    <Box maw={800} mx="auto">
      <Stack gap="lg">
        {/* Header */}
        <Box>
          <Text size="lg" fw={600} ff="monospace" c="white">
            user_permissions
          </Text>
          <Text size="xs" c="dimmed" ff="monospace" mt={4}>
            control who can access and manage your bot settings
          </Text>
        </Box>

        <UserPermissions
          userConfig={userConfig}
          setUserConfig={setUserConfig}
          errorMessage={errorMessage}
          loading={loading}
        />
      </Stack>
    </Box>
  );
}
