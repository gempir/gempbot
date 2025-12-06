import { Container, Loader, Stack, Text, Title } from "@mantine/core";
import { useUserConfig } from "../../hooks/useUserConfig";
import { UserPermissions } from "./UserPermissions";

export function Permissions() {
  const [userConfig, setUserConfig, , loading, errorMessage] = useUserConfig();

  if (loading || !userConfig) {
    return (
      <Container size="xl">
        <Stack align="center" justify="center" h={400}>
          <Loader size="lg" />
          <Text c="dimmed">Loading permissions...</Text>
        </Stack>
      </Container>
    );
  }

  return (
    <Container size="xl">
      <Stack gap="lg">
        <div>
          <Title order={1} mb="xs">
            User Permissions
          </Title>
          <Text c="dimmed">
            Control who can access and manage your bot settings
          </Text>
        </div>

        <UserPermissions
          userConfig={userConfig}
          setUserConfig={setUserConfig}
          errorMessage={errorMessage}
          loading={loading}
        />
      </Stack>
    </Container>
  );
}
