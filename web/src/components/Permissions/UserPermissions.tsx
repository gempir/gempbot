import {
  ActionIcon,
  Button,
  Card,
  Checkbox,
  Group,
  Loader,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
  Tooltip,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { PlusIcon, TrashIcon } from "@heroicons/react/24/solid";
import { useEffect, useState } from "react";
import {
  Permission,
  SetUserConfig,
  UserConfig,
} from "../../hooks/useUserConfig";

type PermissionRow = {
  user: string;
  editor: boolean;
  prediction: boolean;
  isNew?: boolean;
};

export function UserPermissions({
  userConfig,
  setUserConfig,
  errorMessage,
  loading,
}: {
  userConfig: UserConfig;
  setUserConfig: SetUserConfig;
  errorMessage?: string;
  loading?: boolean;
}) {
  const [rows, setRows] = useState<PermissionRow[]>([]);
  const [newUserCounter, setNewUserCounter] = useState(0);

  useEffect(() => {
    const permRows = Object.entries(userConfig.Permissions || {}).map(
      ([user, perm]) => ({
        user,
        editor: perm.Editor,
        prediction: perm.Prediction,
      }),
    );
    setRows(permRows);
  }, [userConfig.Permissions]);

  const handleAddRow = () => {
    setRows([
      ...rows,
      {
        user: ``,
        editor: false,
        prediction: true,
        isNew: true,
      },
    ]);
    setNewUserCounter(newUserCounter + 1);
  };

  const handleRemoveRow = (index: number) => {
    const newRows = rows.filter((_, i) => i !== index);
    setRows(newRows);
    handleSave(newRows);
  };

  const handleUpdateRow = (
    index: number,
    field: keyof PermissionRow,
    value: any,
  ) => {
    const newRows = [...rows];
    newRows[index] = { ...newRows[index], [field]: value };
    setRows(newRows);
  };

  const handleSave = (rowsToSave: PermissionRow[] = rows) => {
    const perms: Record<string, Permission> = {};

    for (const row of rowsToSave) {
      if (row.user.trim()) {
        perms[row.user.toLowerCase().trim()] = {
          Editor: row.editor,
          Prediction: row.prediction,
        };
      }
    }

    setUserConfig({ ...userConfig, Permissions: perms });

    notifications.show({
      title: "Permissions Updated",
      message: "User permissions have been saved",
      color: "green",
    });
  };

  if (loading) {
    return (
      <Card shadow="sm" padding="xl" radius="md" withBorder>
        <Group justify="center" p="xl">
          <Loader size="lg" />
        </Group>
      </Card>
    );
  }

  return (
    <Stack gap="lg">
      <Card shadow="sm" padding="lg" radius="md" withBorder>
        <Stack gap="md">
          <Group justify="space-between">
            <div>
              <Title order={3} size="h4">
                User Permissions
              </Title>
              <Text size="sm" c="dimmed">
                Manage who can edit your bot settings and make predictions
              </Text>
            </div>
            <Button
              leftSection={<PlusIcon style={{ width: 16, height: 16 }} />}
              onClick={handleAddRow}
              color="purple"
            >
              Add User
            </Button>
          </Group>

          {rows.length === 0 ? (
            <Text c="dimmed" ta="center" py="xl">
              No permissions configured. Add users to grant them access.
            </Text>
          ) : (
            <Table highlightOnHover>
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>Twitch Username</Table.Th>
                  <Table.Th>Editor</Table.Th>
                  <Table.Th>Predictions</Table.Th>
                  <Table.Th w={100}>Actions</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>
                {rows.map((row, index) => (
                  <Table.Tr key={index}>
                    <Table.Td>
                      <TextInput
                        placeholder="username"
                        value={row.user}
                        onChange={(e) =>
                          handleUpdateRow(index, "user", e.currentTarget.value)
                        }
                        onBlur={() => handleSave()}
                        variant="unstyled"
                      />
                    </Table.Td>
                    <Table.Td>
                      <Checkbox
                        checked={row.editor}
                        onChange={(e) => {
                          const newRows = [...rows];
                          newRows[index] = {
                            ...newRows[index],
                            editor: e.currentTarget.checked,
                          };
                          setRows(newRows);
                          handleSave(newRows);
                        }}
                        color="purple"
                      />
                    </Table.Td>
                    <Table.Td>
                      <Checkbox
                        checked={row.prediction}
                        onChange={(e) => {
                          const newRows = [...rows];
                          newRows[index] = {
                            ...newRows[index],
                            prediction: e.currentTarget.checked,
                          };
                          setRows(newRows);
                          handleSave(newRows);
                        }}
                        color="purple"
                      />
                    </Table.Td>
                    <Table.Td>
                      <Tooltip label="Remove user">
                        <ActionIcon
                          variant="subtle"
                          color="red"
                          onClick={() => handleRemoveRow(index)}
                        >
                          <TrashIcon style={{ width: 16, height: 16 }} />
                        </ActionIcon>
                      </Tooltip>
                    </Table.Td>
                  </Table.Tr>
                ))}
              </Table.Tbody>
            </Table>
          )}

          {errorMessage && (
            <Text c="red" size="sm">
              {errorMessage}
            </Text>
          )}
        </Stack>
      </Card>

      <Card shadow="sm" padding="lg" radius="md" withBorder>
        <Stack gap="xs">
          <Title order={4} size="h5">
            Permission Types
          </Title>
          <Text size="sm" c="dimmed">
            <strong>Editor:</strong> Users can modify bot settings, rewards, and
            manage blocked emotes
          </Text>
          <Text size="sm" c="dimmed">
            <strong>Predictions:</strong> Users can create and manage
            predictions in your channel
          </Text>
        </Stack>
      </Card>
    </Stack>
  );
}
