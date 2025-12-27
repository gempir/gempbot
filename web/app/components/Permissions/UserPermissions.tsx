import { PlusIcon, TrashIcon } from "@heroicons/react/24/solid";
import {
  ActionIcon,
  Box,
  Button,
  Checkbox,
  Group,
  Loader,
  Stack,
  Table,
  Text,
  TextInput,
  Tooltip,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useEffect, useState } from "react";
import type {
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
        user: "",
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
    value: string | boolean,
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
      title: "saved",
      message: "permissions updated",
      color: "green",
    });
  };

  if (loading) {
    return (
      <Box
        p="lg"
        style={{
          border: "1px solid var(--border-subtle)",
          backgroundColor: "var(--bg-elevated)",
        }}
      >
        <Group justify="center" p="xl">
          <Loader size="sm" />
        </Group>
      </Box>
    );
  }

  return (
    <Stack gap="lg">
      {/* Permissions Table */}
      <Box
        p="md"
        style={{
          border: "1px solid var(--border-subtle)",
          backgroundColor: "var(--bg-elevated)",
        }}
      >
        <Stack gap="md">
          <Group justify="space-between">
            <Text size="sm" fw={600} ff="monospace" c="white">
              permission_table
            </Text>
            <Button
              leftSection={<PlusIcon style={{ width: 12, height: 12 }} />}
              onClick={handleAddRow}
              color="terminal"
              size="xs"
            >
              add user
            </Button>
          </Group>

          {rows.length === 0 ? (
            <Text c="dimmed" ta="center" py="xl" size="xs" ff="monospace">
              no permissions configured. add users to grant access.
            </Text>
          ) : (
            <Table highlightOnHover>
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>username</Table.Th>
                  <Table.Th w={80}>editor</Table.Th>
                  <Table.Th w={100}>predictions</Table.Th>
                  <Table.Th w={60}>del</Table.Th>
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
                        size="xs"
                        styles={{
                          input: {
                            fontFamily: "'JetBrains Mono', monospace",
                            fontSize: "0.8125rem",
                            padding: 0,
                          },
                        }}
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
                        color="terminal"
                        size="xs"
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
                        color="terminal"
                        size="xs"
                      />
                    </Table.Td>
                    <Table.Td>
                      <Tooltip label="remove">
                        <ActionIcon
                          variant="subtle"
                          color="red"
                          size="xs"
                          onClick={() => handleRemoveRow(index)}
                        >
                          <TrashIcon style={{ width: 12, height: 12 }} />
                        </ActionIcon>
                      </Tooltip>
                    </Table.Td>
                  </Table.Tr>
                ))}
              </Table.Tbody>
            </Table>
          )}

          {errorMessage && (
            <Text c="red" size="xs" ff="monospace">
              error: {errorMessage}
            </Text>
          )}
        </Stack>
      </Box>

      {/* Permission Types Info */}
      <Box
        p="md"
        style={{
          border: "1px solid var(--border-subtle)",
          backgroundColor: "var(--bg-surface)",
        }}
      >
        <Stack gap="sm">
          <Text
            size="xs"
            fw={600}
            ff="monospace"
            c="dimmed"
            tt="uppercase"
            style={{ letterSpacing: "0.1em" }}
          >
            permission_types
          </Text>
          <Stack gap="xs">
            <Group gap="xs">
              <Text size="xs" ff="monospace" c="terminal" w={80}>
                editor:
              </Text>
              <Text size="xs" c="dimmed" ff="monospace">
                modify bot settings, rewards, and manage blocked emotes
              </Text>
            </Group>
            <Group gap="xs">
              <Text size="xs" ff="monospace" c="terminal" w={80}>
                predictions:
              </Text>
              <Text size="xs" c="dimmed" ff="monospace">
                create and manage predictions in your channel
              </Text>
            </Group>
          </Stack>
        </Stack>
      </Box>
    </Stack>
  );
}
