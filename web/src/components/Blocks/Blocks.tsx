import {
  ActionIcon,
  Button,
  Card,
  Container,
  Group,
  Loader,
  Pagination,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
  Tooltip,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { PlusIcon, TrashIcon } from "@heroicons/react/24/solid";
import { useState } from "react";
import { Emote } from "../Emote/Emote";
import { useBlocks } from "../../hooks/useBlocks";

export function Blocks() {
  const { blocks, page, totalPages, loading, addBlock, removeBlock, setPage } =
    useBlocks();
  const [newEmoteId, setNewEmoteId] = useState("");
  const [rewardType, setRewardType] = useState<string>("BTTV");
  const [adding, setAdding] = useState(false);

  const handleAdd = async () => {
    if (!newEmoteId.trim()) {
      notifications.show({
        title: "Validation Error",
        message: "Please enter an emote ID",
        color: "red",
      });
      return;
    }

    setAdding(true);
    try {
      await addBlock(newEmoteId, rewardType as "BTTV" | "7TV");
      setNewEmoteId("");
      notifications.show({
        title: "Emote Blocked",
        message: "Emote has been added to the block list",
        color: "green",
      });
    } catch (error) {
      notifications.show({
        title: "Failed to Block",
        message: "Could not add emote to block list",
        color: "red",
      });
    } finally {
      setAdding(false);
    }
  };

  const handleRemove = async (emoteId: string, type: string) => {
    try {
      await removeBlock(emoteId, type as "BTTV" | "7TV");
      notifications.show({
        title: "Emote Unblocked",
        message: "Emote has been removed from the block list",
        color: "green",
      });
    } catch (error) {
      notifications.show({
        title: "Failed to Unblock",
        message: "Could not remove emote from block list",
        color: "red",
      });
    }
  };

  return (
    <Container size="xl">
      <Stack gap="lg">
        <div>
          <Title order={1} mb="xs">
            Blocked Emotes
          </Title>
          <Text c="dimmed">
            Prevent specific emotes from being added to your channel
          </Text>
        </div>

        {/* Add Block Form */}
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Stack gap="md">
            <Title order={3} size="h4">
              Block New Emote
            </Title>

            <Group align="flex-end">
              <TextInput
                label="Emote ID"
                placeholder="Enter emote ID"
                value={newEmoteId}
                onChange={(e) => setNewEmoteId(e.currentTarget.value)}
                style={{ flex: 1 }}
                onKeyDown={(e) => {
                  if (e.key === "Enter") handleAdd();
                }}
              />

              <Select
                label="Type"
                data={[
                  { value: "BTTV", label: "BTTV" },
                  { value: "7TV", label: "7TV" },
                ]}
                value={rewardType}
                onChange={(value) => setRewardType(value || "BTTV")}
                w={120}
              />

              <Button
                leftSection={<PlusIcon style={{ width: 16, height: 16 }} />}
                onClick={handleAdd}
                loading={adding}
                color="purple"
              >
                Block
              </Button>
            </Group>
          </Stack>
        </Card>

        {/* Blocks Table */}
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          {loading ? (
            <Group justify="center" p="xl">
              <Loader size="lg" />
            </Group>
          ) : blocks.length === 0 ? (
            <Text c="dimmed" ta="center" py="xl">
              No blocked emotes yet
            </Text>
          ) : (
            <Stack gap="md">
              <Table highlightOnHover>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Emote</Table.Th>
                    <Table.Th>Emote ID</Table.Th>
                    <Table.Th>Type</Table.Th>
                    <Table.Th w={100}>Actions</Table.Th>
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  {blocks.map((block) => (
                    <Table.Tr key={`${block.emoteID}-${block.type}`}>
                      <Table.Td>
                        <Emote emoteId={block.emoteID} type={block.type} />
                      </Table.Td>
                      <Table.Td>
                        <Text size="sm" ff="monospace">
                          {block.emoteID}
                        </Text>
                      </Table.Td>
                      <Table.Td>
                        <Text size="sm">{block.type}</Text>
                      </Table.Td>
                      <Table.Td>
                        <Tooltip label="Remove block">
                          <ActionIcon
                            variant="subtle"
                            color="red"
                            onClick={() =>
                              handleRemove(block.emoteID, block.type)
                            }
                          >
                            <TrashIcon style={{ width: 16, height: 16 }} />
                          </ActionIcon>
                        </Tooltip>
                      </Table.Td>
                    </Table.Tr>
                  ))}
                </Table.Tbody>
              </Table>

              {totalPages > 1 && (
                <Group justify="center" mt="md">
                  <Pagination
                    total={totalPages}
                    value={page}
                    onChange={setPage}
                    color="purple"
                  />
                </Group>
              )}
            </Stack>
          )}
        </Card>
      </Stack>
    </Container>
  );
}
