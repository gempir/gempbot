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
  return cleaned.length >= 3 && cleaned.length <= 64 && /^[a-zA-Z0-9_-]+$/.test(cleaned);
}

/**
 * Export blocks as JSON string
 */
export function exportBlocksAsJson(blocks: Block[]): string {
  return JSON.stringify(blocks, null, 2);
}

/**
 * Import blocks from JSON string
 * Returns parsed blocks and any validation errors
 */
export function importBlocksFromJson(jsonString: string): {
  blocks: Block[];
  errors: string[];
} {
  const errors: string[] = [];
  let blocks: Block[] = [];

  try {
    const parsed = JSON.parse(jsonString);

    if (!Array.isArray(parsed)) {
      errors.push("JSON must be an array of blocks");
      return { blocks: [], errors };
    }

    blocks = parsed
      .map((item, index) => {
        // Validate required fields
        if (!item.EmoteID || typeof item.EmoteID !== "string") {
          errors.push(`Block ${index + 1}: Missing or invalid EmoteID`);
          return null;
        }

        if (!item.Type || (item.Type !== "BTTV" && item.Type !== "7TV")) {
          errors.push(`Block ${index + 1}: Invalid Type (must be BTTV or 7TV)`);
          return null;
        }

        // Create block with required fields
        return {
          EmoteID: item.EmoteID,
          Type: item.Type,
          ChannelTwitchID: item.ChannelTwitchID || "",
          CreatedAt: item.CreatedAt ? new Date(item.CreatedAt) : new Date(),
        } as Block;
      })
      .filter((block): block is Block => block !== null);

    if (blocks.length === 0 && errors.length === 0) {
      errors.push("No valid blocks found in JSON");
    }
  } catch (error) {
    errors.push(`Invalid JSON: ${error instanceof Error ? error.message : "Unknown error"}`);
  }

  return { blocks, errors };
}

/**
 * Download a JSON string as a file
 */
export function downloadJson(jsonString: string, filename: string): void {
  const blob = new Blob([jsonString], { type: "application/json" });
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
  return `emote-blocks-${dateStr}.json`;
}
