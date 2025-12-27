import { QuestionMarkCircleIcon } from "@heroicons/react/24/outline";
import { Box, Image, Stack, Text } from "@mantine/core";
import { useState } from "react";

interface EmotePreviewProps {
  emoteId: string;
  type: string;
}

export function EmotePreview({ emoteId, type }: EmotePreviewProps) {
  const [imageError, setImageError] = useState(false);

  const getEmoteUrl = () => {
    if (type?.toLowerCase() === "seventv") {
      return `https://cdn.7tv.app/emote/${emoteId}/2x.avif`;
    }
    return "";
  };

  const emoteUrl = getEmoteUrl();

  if (!emoteId || !emoteUrl) {
    return (
      <Box
        w={80}
        h={80}
        style={{
          border: "1px solid var(--border-subtle)",
          backgroundColor: "var(--bg-surface)",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Stack gap={4} align="center">
          <QuestionMarkCircleIcon
            style={{ width: 24, height: 24, color: "var(--text-tertiary)" }}
          />
          <Text size="xs" c="dimmed" ff="monospace">
            preview
          </Text>
        </Stack>
      </Box>
    );
  }

  return (
    <Box
      w={80}
      h={80}
      style={{
        border: "1px solid var(--border-subtle)",
        backgroundColor: "var(--bg-surface)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      {imageError ? (
        <Stack gap={4} align="center">
          <QuestionMarkCircleIcon
            style={{ width: 24, height: 24, color: "var(--text-tertiary)" }}
          />
          <Text size="xs" c="dimmed" ff="monospace">
            error
          </Text>
        </Stack>
      ) : (
        <Image
          src={emoteUrl}
          alt="Emote preview"
          w={48}
          h={48}
          fit="contain"
          onError={() => setImageError(true)}
        />
      )}
    </Box>
  );
}
