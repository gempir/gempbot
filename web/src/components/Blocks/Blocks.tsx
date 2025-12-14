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
  Badge,
  Button,
  Card,
  Checkbox,
  Container,
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
  Title,
  Tooltip,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useDebouncedValue } from "@mantine/hooks";
import { useState, useMemo } from "react";
import { type Block, useBlocks } from "../../hooks/useBlocks";
import {
  parseEmoteIds,
  downloadJson,
  exportBlocksAsJson,
  importBlocksFromJson,
  getExportFilename,
} from "../../utils/emoteBlockHelpers";
import { Emote } from "../Emote/Emote";
import { EmotePreview } from "./EmotePreview";

export function Blocks() {
  const { blocks, page, totalPages, loading, addBlock, removeBlock, removeMultiple, setPage } =
    useBlocks();

  // Form state
  const [newEmoteIds, setNewEmoteIds] = useState("");
  const [rewardType, setRewardType] = useState<string>("7TV");
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
  const [importFile, setImportFile] = useState<File | null>(null);
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
        block.EmoteID.toLowerCase().includes(query)
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
    if (selectedIds.size === filteredBlocks.length && filteredBlocks.length > 0) {
      setSelectedIds(new Set());
    } else {
      const allKeys = filteredBlocks.map((block) => `${block.EmoteID}-${block.Type}`);
      setSelectedIds(new Set(allKeys));
    }
  };

  const handleAdd = async () => {
    if (parsedIds.length === 0) {
      notifications.show({
        title: "Validation Error",
        message: "Please enter at least one emote ID",
        color: "red",
      });
      return;
    }

    setAdding(true);
    try {
      const emoteIdsString = parsedIds.join(",");
      await addBlock(emoteIdsString, rewardType as "BTTV" | "7TV");
      setNewEmoteIds("");
      setPreviewEmoteId("");
      notifications.show({
        title: "Emotes Blocked",
        message: `${parsedIds.length} emote(s) have been added to the block list`,
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "Failed to Block",
        message: "Could not add emote(s) to block list",
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
        title: "Emote Unblocked",
        message: "Emote has been removed from the block list",
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "Failed to Unblock",
        message: "Could not remove emote from block list",
        color: "red",
      });
    }
  };

  const handleBulkDelete = async () => {
    const blocksToDelete = filteredBlocks.filter((block) =>
      selectedIds.has(`${block.EmoteID}-${block.Type}`)
    );

    if (blocksToDelete.length === 0) return;

    setBulkDeleting(true);
    setDeleteConfirmOpen(false);

    try {
      await removeMultiple(blocksToDelete);
      setSelectedIds(new Set());
      notifications.show({
        title: "Emotes Unblocked",
        message: `${blocksToDelete.length} emote(s) removed from the block list`,
        color: "green",
      });
    } catch (_error) {
      notifications.show({
        title: "Failed to Unblock",
        message: "Could not remove all emotes. Some may have been deleted.",
        color: "red",
      });
    } finally {
      setBulkDeleting(false);
    }
  };

  const handleExport = () => {
    const json = exportBlocksAsJson(blocks);
    downloadJson(json, getExportFilename());
    notifications.show({
      title: "Export Successful",
      message: `Exported ${blocks.length} emote block(s)`,
      color: "green",
    });
  };

  const handleImportFile = (file: File | null) => {
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      const { blocks: importedBlocks, errors } = importBlocksFromJson(content);

      if (errors.length > 0) {
        notifications.show({
          title: "Import Error",
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
      // Group by type for bulk import
      const bttvIds = importPreview
        .filter((b) => b.Type === "BTTV")
        .map((b) => b.EmoteID);
      const sevenTvIds = importPreview
        .filter((b) => b.Type === "7TV")
        .map((b) => b.EmoteID);

      if (bttvIds.length > 0) {
        await addBlock(bttvIds.join(","), "BTTV");
      }
      if (sevenTvIds.length > 0) {
        await addBlock(sevenTvIds.join(","), "7TV");
      }

      notifications.show({
        title: "Import Successful",
        message: `Imported ${importPreview.length} emote block(s)`,
        color: "green",
      });
      setImportPreview([]);
      setImportFile(null);
    } catch (_error) {
      notifications.show({
        title: "Import Failed",
        message: "Could not import all emotes",
        color: "red",
      });
    }
  };

  // Update preview when typing
  const handleEmoteIdsChange = (value: string) => {
    setNewEmoteIds(value);
    const ids = parseEmoteIds(value);
    setPreviewEmoteId(ids[0] || "");
  };

  const selectedCount = selectedIds.size;

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

        {/* Search and Action Bar */}
        <Card shadow="sm" padding="md" radius="md" withBorder>
          <Group justify="space-between" wrap="nowrap">
            <Group style={{ flex: 1 }}>
              <TextInput
                placeholder="Search by emote ID..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.currentTarget.value)}
                leftSection={<MagnifyingGlassIcon style={{ width: 16, height: 16 }} />}
                rightSection={
                  searchQuery && (
                    <ActionIcon
                      variant="subtle"
                      onClick={() => setSearchQuery("")}
                      size="sm"
                    >
                      <XMarkIcon style={{ width: 14, height: 14 }} />
                    </ActionIcon>
                  )
                }
                style={{ flex: 1, minWidth: 200 }}
              />
              <Select
                placeholder="Type"
                value={typeFilter}
                onChange={(value) => setTypeFilter(value || "")}
                data={[
                  { value: "", label: "All Types" },
                  { value: "BTTV", label: "BTTV" },
                  { value: "7TV", label: "7TV" },
                ]}
                w={120}
                clearable
              />
            </Group>
            <Group gap="xs">
              <Tooltip label="Export blocks as JSON">
                <Button
                  variant="light"
                  leftSection={<ArrowDownTrayIcon style={{ width: 16, height: 16 }} />}
                  onClick={handleExport}
                  disabled={blocks.length === 0}
                >
                  Export
                </Button>
              </Tooltip>
              <FileButton onChange={handleImportFile} accept="application/json">
                {(props) => (
                  <Tooltip label="Import blocks from JSON">
                    <Button
                      {...props}
                      variant="light"
                      leftSection={<ArrowUpTrayIcon style={{ width: 16, height: 16 }} />}
                    >
                      Import
                    </Button>
                  </Tooltip>
                )}
              </FileButton>
            </Group>
          </Group>
        </Card>

        {/* Add Block Form */}
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Stack gap="md">
            <Group justify="space-between" align="center">
              <Title order={3} size="h4">
                Block New Emote(s)
              </Title>
              {parsedIds.length > 0 && (
                <Badge color="cyan" variant="light">
                  {parsedIds.length} emote{parsedIds.length > 1 ? "s" : ""} detected
                </Badge>
              )}
            </Group>

            <Group align="flex-start" wrap="nowrap">
              <Stack style={{ flex: 1 }} gap="xs">
                <Textarea
                  label="Emote ID(s)"
                  placeholder="Enter one or more emote IDs (comma, space, or newline separated)"
                  value={newEmoteIds}
                  onChange={(e) => handleEmoteIdsChange(e.currentTarget.value)}
                  minRows={3}
                  autosize
                  maxRows={6}
                />
                <Text size="xs" c="dimmed">
                  Supports comma, space, or newline separated IDs
                </Text>
              </Stack>

              <EmotePreview
                emoteId={debouncedPreview}
                type={rewardType}
              />
            </Group>

            <Group align="flex-end">
              <Select
                label="Type"
                data={[
                  { value: "7TV", label: "7TV" },
                  { value: "BTTV", label: "BTTV" },
                ]}
                value={rewardType}
                onChange={(value) => setRewardType(value || "7TV")}
                w={120}
              />

              <Button
                leftSection={<PlusIcon style={{ width: 16, height: 16 }} />}
                onClick={handleAdd}
                loading={adding}
                color="cyan"
                disabled={parsedIds.length === 0}
              >
                Block {parsedIds.length > 1 ? `${parsedIds.length} Emotes` : "Emote"}
              </Button>
            </Group>
          </Stack>
        </Card>

        {/* Bulk Actions Bar */}
        {selectedCount > 0 && (
          <Card
            shadow="sm"
            padding="md"
            radius="md"
            withBorder
            style={{ borderColor: "var(--mantine-color-cyan-6)" }}
          >
            <Group justify="space-between">
              <Group gap="md">
                <Checkbox
                  checked={selectedCount === filteredBlocks.length && filteredBlocks.length > 0}
                  indeterminate={selectedCount > 0 && selectedCount < filteredBlocks.length}
                  onChange={toggleSelectAll}
                  label={`${selectedCount} selected`}
                />
              </Group>
              <Group gap="xs">
                <Button
                  variant="subtle"
                  onClick={() => setSelectedIds(new Set())}
                >
                  Deselect All
                </Button>
                <Button
                  color="red"
                  leftSection={<TrashIcon style={{ width: 16, height: 16 }} />}
                  onClick={() => setDeleteConfirmOpen(true)}
                  loading={bulkDeleting}
                >
                  Delete Selected
                </Button>
              </Group>
            </Group>
          </Card>
        )}

        {/* Blocks Table */}
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          {loading ? (
            <Group justify="center" p="xl">
              <Loader size="lg" />
            </Group>
          ) : filteredBlocks.length === 0 ? (
            <Text c="dimmed" ta="center" py="xl">
              {blocks.length === 0
                ? "No blocked emotes yet"
                : "No emotes match your search"}
            </Text>
          ) : (
            <Stack gap="md">
              {debouncedSearch || typeFilter ? (
                <Text size="sm" c="dimmed">
                  {filteredBlocks.length} result{filteredBlocks.length !== 1 ? "s" : ""}
                </Text>
              ) : null}

              <Table highlightOnHover>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th w={40}>
                      <Checkbox
                        checked={selectedCount === filteredBlocks.length && filteredBlocks.length > 0}
                        indeterminate={selectedCount > 0 && selectedCount < filteredBlocks.length}
                        onChange={toggleSelectAll}
                      />
                    </Table.Th>
                    <Table.Th w={60}>Preview</Table.Th>
                    <Table.Th>Emote ID</Table.Th>
                    <Table.Th w={100}>Type</Table.Th>
                    <Table.Th w={100}>Actions</Table.Th>
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
                          />
                        </Table.Td>
                        <Table.Td>
                          <Emote emoteId={block.EmoteID} type={block.Type} size={24} />
                        </Table.Td>
                        <Table.Td>
                          <Text size="sm" ff="monospace">
                            {block.EmoteID}
                          </Text>
                        </Table.Td>
                        <Table.Td>
                          <Badge color={block.Type === "7TV" ? "violet" : "orange"} variant="light">
                            {block.Type}
                          </Badge>
                        </Table.Td>
                        <Table.Td onClick={(e) => e.stopPropagation()}>
                          <Tooltip label="Remove block">
                            <ActionIcon
                              variant="subtle"
                              color="red"
                              onClick={() => handleRemove(block)}
                            >
                              <TrashIcon style={{ width: 16, height: 16 }} />
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
                    color="cyan"
                  />
                </Group>
              )}
            </Stack>
          )}
        </Card>
      </Stack>

      {/* Delete Confirmation Modal */}
      <Modal
        opened={deleteConfirmOpen}
        onClose={() => setDeleteConfirmOpen(false)}
        title="Confirm Bulk Delete"
      >
        <Stack gap="md">
          <Text>
            Are you sure you want to delete {selectedCount} emote block{selectedCount > 1 ? "s" : ""}?
          </Text>
          <Group justify="flex-end">
            <Button variant="subtle" onClick={() => setDeleteConfirmOpen(false)}>
              Cancel
            </Button>
            <Button color="red" onClick={handleBulkDelete}>
              Delete
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
        title="Import Preview"
        size="lg"
      >
        <Stack gap="md">
          <Text>
            Import {importPreview.length} emote block{importPreview.length > 1 ? "s" : ""}?
          </Text>
          <Card withBorder mah={300} style={{ overflow: "auto" }}>
            <Stack gap="xs">
              {importPreview.slice(0, 10).map((block, idx) => (
                <Group key={idx} gap="sm">
                  <Badge color={block.Type === "7TV" ? "violet" : "orange"} variant="light">
                    {block.Type}
                  </Badge>
                  <Text size="sm" ff="monospace">
                    {block.EmoteID}
                  </Text>
                </Group>
              ))}
              {importPreview.length > 10 && (
                <Text size="sm" c="dimmed">
                  ... and {importPreview.length - 10} more
                </Text>
              )}
            </Stack>
          </Card>
          <Group justify="flex-end">
            <Button
              variant="subtle"
              onClick={() => {
                setImportModalOpen(false);
                setImportPreview([]);
              }}
            >
              Cancel
            </Button>
            <Button color="cyan" onClick={handleConfirmImport}>
              Import
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Container>
  );
}
