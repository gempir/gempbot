import {
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  TrashIcon,
  XMarkIcon,
} from "@heroicons/react/24/solid";
import {
  ActionIcon,
  Box,
  Button,
  Checkbox,
  FileButton,
  Group,
  Loader,
  Modal,
  Pagination,
  Select,
  Stack,
  Table,
  Text,
  Textarea,
  TextInput,
  Tooltip,
} from "@mantine/core";
import { useDebouncedValue } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { useMemo, useState } from "react";
import { type Block, useBlocks } from "../../hooks/useBlocks";
import {
  downloadCsv,
  exportBlocksAsCsv,
  getExportFilename,
  importBlocksFromCsv,
  parseEmoteIds,
} from "../../utils/emoteBlockHelpers";
import { Emote } from "../Emote/Emote";
import { EmotePreview } from "./EmotePreview";

export function Blocks() {
  const {
    blocks,
    page,
    totalPages,
    loading,
    addBlock,
    removeBlock,
    removeMultiple,
    setPage,
  } = useBlocks();

  // Form state
  const [newEmoteIds, setNewEmoteIds] = useState("");
  const [rewardType, setRewardType] = useState<string>("seventv");
  const [adding, setAdding] = useState(false);

  // Preview state
  const [previewEmoteId, setPreviewEmoteId] = useState("");
  const [debouncedPreview] = useDebouncedValue(previewEmoteId, 300);

  // Search/Filter state
  const [searchQuery, setSearchQuery] = useState("");
  const [typeFilter, setTypeFilter] = useState<string>("");
  const [debouncedSearch] = useDebouncedValue(searchQuery, 300);

  // Selection state
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [bulkDeleting, setBulkDeleting] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);

  // Import state
  const [importPreview, setImportPreview] = useState<Block[]>([]);
  const [importModalOpen, setImportModalOpen] = useState(false);

  // Parse emote IDs from input
  const parsedIds = useMemo(() => parseEmoteIds(newEmoteIds), [newEmoteIds]);

  // Filter blocks based on search and type
  const filteredBlocks = useMemo(() => {
    let filtered = blocks;

    if (debouncedSearch) {
      const query = debouncedSearch.toLowerCase();
      filtered = filtered.filter((block) =>
        block.EmoteID.toLowerCase().includes(query),
      );
    }

    if (typeFilter) {
      filtered = filtered.filter((block) => block.Type === typeFilter);
    }

    return filtered;
  }, [blocks, debouncedSearch, typeFilter]);

  // Selection helpers
  const toggleSelection = (blockKey: string) => {
    const newSet = new Set(selectedIds);
    if (newSet.has(blockKey)) {
      newSet.delete(blockKey);
    } else {
      newSet.add(blockKey);
    }
    setSelectedIds(newSet);
  };

  const toggleSelectAll = () => {
    if (
      selectedIds.size === filteredBlocks.length &&
      filteredBlocks.length > 0
    ) {
      setSelectedIds(new Set());
    } else {
      const allKeys = filteredBlocks.map(
        (block) => `${block.EmoteID}-${block.Type}`,
      );
      setSelectedIds(new Set(allKeys));
    }
  };

  const handleAdd = async () => {
    if (parsedIds.length === 0) {
      notifications.show({
        title: "error",
        message: "enter at least one emote id",
        color: "red",
      });
      return;
    }

    setAdding(true);
    try {
      const emoteIdsString = parsedIds.join(",");
      await addBlock(emoteIdsString, rewardType as "seventv");
      setNewEmoteIds("");
      setPreviewEmoteId("");
      notifications.show({
        title: "blocked",
        message: `${parsedIds.length} emote(s) added to block list`,
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "error",
        message: "could not add emote(s) to block list",
        color: "red",
      });
    } finally {
      setAdding(false);
    }
  };

  const handleRemove = async (block: Block) => {
    try {
      await removeBlock(block);
      notifications.show({
        title: "unblocked",
        message: "emote removed from block list",
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "error",
        message: "could not remove emote from block list",
        color: "red",
      });
    }
  };

  const handleBulkDelete = async () => {
    const blocksToDelete = filteredBlocks.filter((block) =>
      selectedIds.has(`${block.EmoteID}-${block.Type}`),
    );

    if (blocksToDelete.length === 0) return;

    setBulkDeleting(true);
    setDeleteConfirmOpen(false);

    try {
      await removeMultiple(blocksToDelete);
      setSelectedIds(new Set());
      notifications.show({
        title: "unblocked",
        message: `${blocksToDelete.length} emote(s) removed`,
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "error",
        message: "could not remove all emotes",
        color: "red",
      });
    } finally {
      setBulkDeleting(false);
    }
  };

  const handleExport = () => {
    const csv = exportBlocksAsCsv(blocks);
    downloadCsv(csv, getExportFilename());
    notifications.show({
      title: "exported",
      message: `${blocks.length} emote block(s) exported`,
      color: "green",
    });
  };

  const handleImportFile = (file: File | null) => {
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      const { blocks: importedBlocks, errors } = importBlocksFromCsv(content);

      if (errors.length > 0) {
        notifications.show({
          title: "import error",
          message: errors[0],
          color: "red",
        });
        return;
      }

      setImportPreview(importedBlocks);
      setImportModalOpen(true);
    };
    reader.readAsText(file);
  };

  const handleConfirmImport = async () => {
    if (importPreview.length === 0) return;

    setImportModalOpen(false);

    try {
      const sevenTvIds = importPreview
        .filter((b) => b.Type === "seventv")
        .map((b) => b.EmoteID);

      if (sevenTvIds.length > 0) {
        await addBlock(sevenTvIds.join(","), "seventv");
      }

      notifications.show({
        title: "imported",
        message: `${importPreview.length} emote block(s) imported`,
        color: "green",
      });
      setImportPreview([]);
    } catch (_error) {
      notifications.show({
        title: "import failed",
        message: "could not import all emotes",
        color: "red",
      });
    }
  };

  const handleEmoteIdsChange = (value: string) => {
    setNewEmoteIds(value);
    const ids = parseEmoteIds(value);
    setPreviewEmoteId(ids[0] || "");
  };

  const selectedCount = selectedIds.size;

  return (
    <Box maw={1000} mx="auto">
      <Stack gap="lg">
        {/* Header */}
        <Box>
          <Text size="lg" fw={600} ff="monospace" c="white">
            blocked_emotes
          </Text>
          <Text size="xs" c="dimmed" ff="monospace" mt={4}>
            prevent specific emotes from being added to your channel
          </Text>
        </Box>

        {/* Search and Action Bar */}
        <Box
          p="sm"
          style={{
            border: "1px solid var(--border-subtle)",
            backgroundColor: "var(--bg-elevated)",
          }}
        >
          <Group justify="space-between" wrap="nowrap">
            <Group style={{ flex: 1 }} gap="sm">
              <TextInput
                placeholder="search by emote id..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.currentTarget.value)}
                leftSection={
                  <MagnifyingGlassIcon style={{ width: 14, height: 14 }} />
                }
                rightSection={
                  searchQuery && (
                    <ActionIcon
                      variant="subtle"
                      onClick={() => setSearchQuery("")}
                      size="xs"
                    >
                      <XMarkIcon style={{ width: 12, height: 12 }} />
                    </ActionIcon>
                  )
                }
                size="xs"
                style={{ flex: 1, minWidth: 180 }}
              />
              <Select
                placeholder="type"
                value={typeFilter}
                onChange={(value) => setTypeFilter(value || "")}
                data={[
                  { value: "", label: "all" },
                  { value: "seventv", label: "7tv" },
                ]}
                size="xs"
                w={100}
                clearable
              />
            </Group>
            <Group gap="xs">
              <Tooltip label="export csv">
                <Button
                  variant="subtle"
                  size="xs"
                  leftSection={
                    <ArrowDownTrayIcon style={{ width: 14, height: 14 }} />
                  }
                  onClick={handleExport}
                  disabled={blocks.length === 0}
                  c="dimmed"
                >
                  export
                </Button>
              </Tooltip>
              <FileButton onChange={handleImportFile} accept="text/csv,.csv">
                {(props) => (
                  <Tooltip label="import csv">
                    <Button
                      {...props}
                      variant="subtle"
                      size="xs"
                      leftSection={
                        <ArrowUpTrayIcon style={{ width: 14, height: 14 }} />
                      }
                      c="dimmed"
                    >
                      import
                    </Button>
                  </Tooltip>
                )}
              </FileButton>
            </Group>
          </Group>
        </Box>

        {/* Add Block Form */}
        <Box
          p="md"
          style={{
            border: "1px solid var(--border-subtle)",
            backgroundColor: "var(--bg-elevated)",
          }}
        >
          <Stack gap="md">
            <Group justify="space-between" align="center">
              <Text size="sm" fw={600} ff="monospace" c="white">
                + add_block
              </Text>
              {parsedIds.length > 0 && (
                <Text size="xs" ff="monospace" c="terminal">
                  {parsedIds.length} detected
                </Text>
              )}
            </Group>

            <Group align="flex-start" wrap="nowrap">
              <Stack style={{ flex: 1 }} gap="xs">
                <Textarea
                  label="emote_ids"
                  placeholder="enter emote ids (comma, space, or newline separated)"
                  value={newEmoteIds}
                  onChange={(e) => handleEmoteIdsChange(e.currentTarget.value)}
                  minRows={3}
                  autosize
                  maxRows={6}
                  size="xs"
                />
              </Stack>

              <EmotePreview emoteId={debouncedPreview} type={rewardType} />
            </Group>

            <Group align="flex-end" justify="space-between">
              <Select
                label="type"
                data={[{ value: "seventv", label: "7tv" }]}
                value={rewardType}
                onChange={(value) => setRewardType(value || "seventv")}
                size="xs"
                w={100}
              />

              <Button
                leftSection={<PlusIcon style={{ width: 14, height: 14 }} />}
                onClick={handleAdd}
                loading={adding}
                size="xs"
                color="terminal"
                disabled={parsedIds.length === 0}
              >
                block{" "}
                {parsedIds.length > 1 ? `${parsedIds.length} emotes` : "emote"}
              </Button>
            </Group>
          </Stack>
        </Box>

        {/* Bulk Actions Bar */}
        {selectedCount > 0 && (
          <Box
            p="sm"
            style={{
              border: "1px solid var(--terminal-green)",
              backgroundColor: "var(--bg-surface)",
            }}
          >
            <Group justify="space-between">
              <Group gap="md">
                <Checkbox
                  checked={
                    selectedCount === filteredBlocks.length &&
                    filteredBlocks.length > 0
                  }
                  indeterminate={
                    selectedCount > 0 && selectedCount < filteredBlocks.length
                  }
                  onChange={toggleSelectAll}
                  label={
                    <Text size="xs" ff="monospace">
                      {selectedCount} selected
                    </Text>
                  }
                  size="xs"
                />
              </Group>
              <Group gap="xs">
                <Button
                  variant="subtle"
                  size="xs"
                  onClick={() => setSelectedIds(new Set())}
                  c="dimmed"
                >
                  deselect
                </Button>
                <Button
                  color="red"
                  size="xs"
                  leftSection={<TrashIcon style={{ width: 14, height: 14 }} />}
                  onClick={() => setDeleteConfirmOpen(true)}
                  loading={bulkDeleting}
                >
                  delete
                </Button>
              </Group>
            </Group>
          </Box>
        )}

        {/* Blocks Table */}
        <Box
          p="md"
          style={{
            border: "1px solid var(--border-subtle)",
            backgroundColor: "var(--bg-elevated)",
          }}
        >
          {loading ? (
            <Group justify="center" p="xl">
              <Loader size="sm" />
            </Group>
          ) : filteredBlocks.length === 0 ? (
            <Text c="dimmed" ta="center" py="xl" size="xs" ff="monospace">
              {blocks.length === 0
                ? "no blocked emotes"
                : "no emotes match filter"}
            </Text>
          ) : (
            <Stack gap="md">
              {(debouncedSearch || typeFilter) && (
                <Text size="xs" c="dimmed" ff="monospace">
                  {filteredBlocks.length} result
                  {filteredBlocks.length !== 1 ? "s" : ""}
                </Text>
              )}

              <Table highlightOnHover>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th w={40}>
                      <Checkbox
                        checked={
                          selectedCount === filteredBlocks.length &&
                          filteredBlocks.length > 0
                        }
                        indeterminate={
                          selectedCount > 0 &&
                          selectedCount < filteredBlocks.length
                        }
                        onChange={toggleSelectAll}
                        size="xs"
                      />
                    </Table.Th>
                    <Table.Th w={50}>img</Table.Th>
                    <Table.Th>emote_id</Table.Th>
                    <Table.Th w={80}>type</Table.Th>
                    <Table.Th w={60}>del</Table.Th>
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  {filteredBlocks.map((block) => {
                    const blockKey = `${block.EmoteID}-${block.Type}`;
                    return (
                      <Table.Tr
                        key={blockKey}
                        onClick={() => toggleSelection(blockKey)}
                        style={{ cursor: "pointer" }}
                      >
                        <Table.Td>
                          <Checkbox
                            checked={selectedIds.has(blockKey)}
                            onChange={() => toggleSelection(blockKey)}
                            size="xs"
                          />
                        </Table.Td>
                        <Table.Td>
                          <Emote
                            emoteId={block.EmoteID}
                            type={block.Type}
                            size={20}
                          />
                        </Table.Td>
                        <Table.Td>
                          <Text size="xs" ff="monospace" c="dimmed">
                            {block.EmoteID}
                          </Text>
                        </Table.Td>
                        <Table.Td>
                          <Text size="xs" ff="monospace" c="terminal">
                            {block.Type === "seventv" ? "7tv" : block.Type}
                          </Text>
                        </Table.Td>
                        <Table.Td onClick={(e) => e.stopPropagation()}>
                          <Tooltip label="remove">
                            <ActionIcon
                              variant="subtle"
                              color="red"
                              size="xs"
                              onClick={() => handleRemove(block)}
                            >
                              <TrashIcon style={{ width: 12, height: 12 }} />
                            </ActionIcon>
                          </Tooltip>
                        </Table.Td>
                      </Table.Tr>
                    );
                  })}
                </Table.Tbody>
              </Table>

              {totalPages > 1 && (
                <Group justify="center" mt="md">
                  <Pagination
                    total={totalPages}
                    value={page}
                    onChange={setPage}
                    color="terminal"
                    size="sm"
                  />
                </Group>
              )}
            </Stack>
          )}
        </Box>
      </Stack>

      {/* Delete Confirmation Modal */}
      <Modal
        opened={deleteConfirmOpen}
        onClose={() => setDeleteConfirmOpen(false)}
        title="confirm_delete"
        size="sm"
      >
        <Stack gap="md">
          <Text size="sm" ff="monospace">
            delete {selectedCount} emote block{selectedCount > 1 ? "s" : ""}?
          </Text>
          <Group justify="flex-end" gap="xs">
            <Button
              variant="subtle"
              size="xs"
              onClick={() => setDeleteConfirmOpen(false)}
            >
              cancel
            </Button>
            <Button color="red" size="xs" onClick={handleBulkDelete}>
              delete
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Import Preview Modal */}
      <Modal
        opened={importModalOpen}
        onClose={() => {
          setImportModalOpen(false);
          setImportPreview([]);
        }}
        title="import_preview"
        size="md"
      >
        <Stack gap="md">
          <Text size="sm" ff="monospace">
            import {importPreview.length} emote block
            {importPreview.length > 1 ? "s" : ""}?
          </Text>
          <Box
            p="sm"
            mah={200}
            style={{
              overflow: "auto",
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-surface)",
            }}
          >
            <Stack gap="xs">
              {importPreview.slice(0, 10).map((block) => (
                <Group key={`${block.Type}-${block.EmoteID}`} gap="sm">
                  <Text size="xs" ff="monospace" c="terminal">
                    {block.Type === "seventv" ? "7tv" : block.Type}
                  </Text>
                  <Text size="xs" ff="monospace" c="dimmed">
                    {block.EmoteID}
                  </Text>
                </Group>
              ))}
              {importPreview.length > 10 && (
                <Text size="xs" c="dimmed" ff="monospace">
                  ... +{importPreview.length - 10} more
                </Text>
              )}
            </Stack>
          </Box>
          <Group justify="flex-end" gap="xs">
            <Button
              variant="subtle"
              size="xs"
              onClick={() => {
                setImportModalOpen(false);
                setImportPreview([]);
              }}
            >
              cancel
            </Button>
            <Button color="terminal" size="xs" onClick={handleConfirmImport}>
              import
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Box>
  );
}
