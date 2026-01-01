import { CheckIcon, XMarkIcon } from "@heroicons/react/24/solid";
import {
  ActionIcon,
  Box,
  Group,
  Loader,
  Table as MantineTable,
  Pagination,
  Stack,
  Text,
  Tooltip,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import type { EmotehistoryItem } from "../../../hooks/useEmotehistory";
import { Emote } from "../../Emote/Emote";

interface TableProps {
  title: string;
  history: EmotehistoryItem[];
  page: number;
  totalPages: number;
  loading: boolean;
  onPageChange: (page: number) => void;
  onApprove?: (item: EmotehistoryItem) => Promise<void>;
  onDeny?: (item: EmotehistoryItem) => Promise<void>;
  approveLabel?: string;
  denyLabel?: string;
}

export function Table({
  title,
  history,
  page,
  totalPages,
  loading,
  onPageChange,
  onApprove,
  onDeny,
  approveLabel = "approve",
  denyLabel = "deny",
}: TableProps) {
  const handleApprove = async (item: EmotehistoryItem) => {
    try {
      await onApprove?.(item);
      notifications.show({
        title: approveLabel,
        message: `emote ${approveLabel}d successfully`,
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "error",
        message: `could not ${approveLabel} emote`,
        color: "red",
      });
    }
  };

  const handleDeny = async (item: EmotehistoryItem) => {
    try {
      await onDeny?.(item);
      notifications.show({
        title: denyLabel,
        message: `emote ${denyLabel === "deny" ? "rejected" : `${denyLabel}ed`}`,
        color: "orange",
      });
    } catch (_error) {
      notifications.show({
        title: "error",
        message: `could not ${denyLabel} emote`,
        color: "red",
      });
    }
  };

  return (
    <Box
      p="md"
      style={{
        border: "1px solid var(--border-subtle)",
        backgroundColor: "var(--bg-elevated)",
      }}
    >
      <Stack gap="md">
        <Text size="sm" fw={600} ff="monospace" c="white">
          {title}
        </Text>

        {loading ? (
          <Group justify="center" p="xl">
            <Loader size="sm" />
          </Group>
        ) : !history || history.length === 0 ? (
          <Text c="dimmed" ta="center" py="xl" size="xs" ff="monospace">
            no emote history
          </Text>
        ) : (
          <>
            <MantineTable highlightOnHover>
              <MantineTable.Thead>
                <MantineTable.Tr>
                  <MantineTable.Th w={50}>img</MantineTable.Th>
                  <MantineTable.Th>emote_id</MantineTable.Th>
                  <MantineTable.Th w={60}>type</MantineTable.Th>
                  <MantineTable.Th w={100}>user</MantineTable.Th>
                  {(onApprove || onDeny) && (
                    <MantineTable.Th w={80}>actions</MantineTable.Th>
                  )}
                </MantineTable.Tr>
              </MantineTable.Thead>
              <MantineTable.Tbody>
                {history.map((item, index) => (
                  <MantineTable.Tr key={index}>
                    <MantineTable.Td>
                      <Emote
                        emoteId={item.emoteID}
                        type={item.type}
                        size={20}
                      />
                    </MantineTable.Td>
                    <MantineTable.Td>
                      <Text size="xs" ff="monospace" c="dimmed">
                        {item.emoteID}
                      </Text>
                    </MantineTable.Td>
                    <MantineTable.Td>
                      <Text size="xs" ff="monospace" c="terminal">
                        {item.type === "seventv" ? "7tv" : item.type}
                      </Text>
                    </MantineTable.Td>
                    <MantineTable.Td>
                      <Text size="xs" ff="monospace" c="dimmed">
                        {item.userLogin}
                      </Text>
                    </MantineTable.Td>
                    {(onApprove || onDeny) && (
                      <MantineTable.Td>
                        <Group gap="xs">
                          {onApprove && (
                            <Tooltip label={approveLabel}>
                              <ActionIcon
                                variant="subtle"
                                color="green"
                                size="xs"
                                onClick={() => handleApprove(item)}
                              >
                                <CheckIcon style={{ width: 12, height: 12 }} />
                              </ActionIcon>
                            </Tooltip>
                          )}
                          {onDeny && (
                            <Tooltip label={denyLabel}>
                              <ActionIcon
                                variant="subtle"
                                color="red"
                                size="xs"
                                onClick={() => handleDeny(item)}
                              >
                                <XMarkIcon style={{ width: 12, height: 12 }} />
                              </ActionIcon>
                            </Tooltip>
                          )}
                        </Group>
                      </MantineTable.Td>
                    )}
                  </MantineTable.Tr>
                ))}
              </MantineTable.Tbody>
            </MantineTable>

            {totalPages > 1 && (
              <Group justify="center" mt="md">
                <Pagination
                  total={totalPages}
                  value={page}
                  onChange={onPageChange}
                  color="terminal"
                  size="sm"
                />
              </Group>
            )}
          </>
        )}
      </Stack>
    </Box>
  );
}
