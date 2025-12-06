import {
  ActionIcon,
  Card,
  Group,
  Loader,
  Pagination,
  Stack,
  Table as MantineTable,
  Text,
  Title,
  Tooltip,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { CheckIcon, XMarkIcon } from "@heroicons/react/24/solid";
import { EmotehistoryItem } from "../../../hooks/useEmotehistory";
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
}: TableProps) {
  const handleApprove = async (item: EmotehistoryItem) => {
    try {
      await onApprove?.(item);
      notifications.show({
        title: "Emote Approved",
        message: "Emote has been added successfully",
        color: "green",
      });
    } catch (error) {
      notifications.show({
        title: "Approval Failed",
        message: "Could not approve emote",
        color: "red",
      });
    }
  };

  const handleDeny = async (item: EmotehistoryItem) => {
    try {
      await onDeny?.(item);
      notifications.show({
        title: "Emote Denied",
        message: "Emote request has been denied",
        color: "orange",
      });
    } catch (error) {
      notifications.show({
        title: "Denial Failed",
        message: "Could not deny emote",
        color: "red",
      });
    }
  };

  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder>
      <Stack gap="md">
        <Title order={3} size="h4">
          {title}
        </Title>

        {loading ? (
          <Group justify="center" p="xl">
            <Loader size="lg" />
          </Group>
        ) : !history || history.length === 0 ? (
          <Text c="dimmed" ta="center" py="xl">
            No emote history yet
          </Text>
        ) : (
          <>
            <MantineTable highlightOnHover>
              <MantineTable.Thead>
                <MantineTable.Tr>
                  <MantineTable.Th>Emote</MantineTable.Th>
                  <MantineTable.Th>Emote ID</MantineTable.Th>
                  <MantineTable.Th>Type</MantineTable.Th>
                  <MantineTable.Th>User</MantineTable.Th>
                  <MantineTable.Th>Status</MantineTable.Th>
                  {(onApprove || onDeny) && (
                    <MantineTable.Th w={100}>Actions</MantineTable.Th>
                  )}
                </MantineTable.Tr>
              </MantineTable.Thead>
              <MantineTable.Tbody>
                {history.map((item, index) => (
                  <MantineTable.Tr key={index}>
                    <MantineTable.Td>
                      <Emote emoteId={item.emoteID} type={item.type} />
                    </MantineTable.Td>
                    <MantineTable.Td>
                      <Text size="sm" ff="monospace">
                        {item.emoteID}
                      </Text>
                    </MantineTable.Td>
                    <MantineTable.Td>
                      <Text size="sm">{item.type}</Text>
                    </MantineTable.Td>
                    <MantineTable.Td>
                      <Text size="sm">{item.userLogin}</Text>
                    </MantineTable.Td>
                    <MantineTable.Td>
                      <Text
                        size="sm"
                        c={item.changeType === "ADD" ? "green" : "red"}
                      >
                        {item.changeType}
                      </Text>
                    </MantineTable.Td>
                    {(onApprove || onDeny) && (
                      <MantineTable.Td>
                        <Group gap="xs">
                          {onApprove && (
                            <Tooltip label="Approve">
                              <ActionIcon
                                variant="subtle"
                                color="green"
                                onClick={() => handleApprove(item)}
                              >
                                <CheckIcon style={{ width: 16, height: 16 }} />
                              </ActionIcon>
                            </Tooltip>
                          )}
                          {onDeny && (
                            <Tooltip label="Deny">
                              <ActionIcon
                                variant="subtle"
                                color="red"
                                onClick={() => handleDeny(item)}
                              >
                                <XMarkIcon style={{ width: 16, height: 16 }} />
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
                  color="cyan"
                />
              </Group>
            )}
          </>
        )}
      </Stack>
    </Card>
  );
}
