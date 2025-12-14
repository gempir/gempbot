import type { Block } from "../hooks/useBlocks";

/**
 * Parse emote IDs from a string input
 * Supports comma, space, and newline separation
 */
export function parseEmoteIds(input: string): string[] {
  if (!input || !input.trim()) {
    return [];
  }

  // Split by comma, space, newline, or combination
  const ids = input
    .split(/[,\s\n]+/)
    .map((id) => id.trim())
    .filter((id) => id.length > 0);

  return Array.from(new Set(ids)); // Remove duplicates
}

/**
 * Basic validation for emote ID
 */
export function validateEmoteId(id: string): boolean {
  if (!id || !id.trim()) {
    return false;
  }

  // Basic validation: alphanumeric and some special chars, 3-64 chars
  const cleaned = id.trim();
  return (
    cleaned.length >= 3 &&
    cleaned.length <= 64 &&
    /^[a-zA-Z0-9_-]+$/.test(cleaned)
  );
}

/**
 * Export blocks as CSV string
 */
export function exportBlocksAsCsv(blocks: Block[]): string {
  // CSV format: type,emoteId
  const rows = blocks.map((block) => `${block.Type},${block.EmoteID}`);
  return rows.join("\n");
}

/**
 * Import blocks from CSV string
 * Returns parsed blocks and any validation errors
 */
export function importBlocksFromCsv(csvString: string): {
  blocks: Block[];
  errors: string[];
} {
  const errors: string[] = [];
  const blocks: Block[] = [];

  try {
    const lines = csvString.trim().split("\n");

    if (lines.length === 0) {
      errors.push("CSV file is empty");
      return { blocks: [], errors };
    }

    lines.forEach((line, index) => {
      const trimmedLine = line.trim();
      if (!trimmedLine) return; // Skip empty lines

      const parts = trimmedLine.split(",");

      if (parts.length < 2) {
        errors.push(
          `Line ${index + 1}: Invalid format (expected: type,emoteId)`,
        );
        return;
      }

      const type = parts[0].trim();
      const emoteId = parts[1].trim();

      // Validate type
      if (type !== "BTTV" && type !== "7TV") {
        errors.push(
          `Line ${index + 1}: Invalid type "${type}" (must be BTTV or 7TV)`,
        );
        return;
      }

      // Validate emote ID
      if (!emoteId) {
        errors.push(`Line ${index + 1}: Missing emote ID`);
        return;
      }

      // Create block
      blocks.push({
        EmoteID: emoteId,
        Type: type,
        ChannelTwitchID: "",
        CreatedAt: new Date(),
      } as Block);
    });

    if (blocks.length === 0 && errors.length === 0) {
      errors.push("No valid blocks found in CSV");
    }
  } catch (error) {
    errors.push(
      `Invalid CSV: ${error instanceof Error ? error.message : "Unknown error"}`,
    );
  }

  return { blocks, errors };
}

/**
 * Download a CSV string as a file
 */
export function downloadCsv(csvString: string, filename: string): void {
  const blob = new Blob([csvString], { type: "text/csv" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

/**
 * Format date for filename
 */
export function getExportFilename(): string {
  const date = new Date();
  const dateStr = date.toISOString().split("T")[0]; // YYYY-MM-DD
  return `emote-blocks-${dateStr}.csv`;
}
