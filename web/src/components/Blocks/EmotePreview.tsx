import { Box, Stack, Text } from "@mantine/core";
import { Emote } from "../Emote/Emote";

interface EmotePreviewProps {
  emoteId: string;
  type: string;
}

export function EmotePreview({ emoteId, type }: EmotePreviewProps) {
  if (!emoteId || !emoteId.trim()) {
    return (
      <Box
        w={64}
        h={64}
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          backgroundColor: "#f1f3f5",
          borderRadius: "4px",
          border: "1px dashed #ced4da",
        }}
      >
        <Text size="xs" c="dimmed">
          Preview
        </Text>
      </Box>
    );
  }

  return (
    <Stack gap="xs" align="center">
      <Emote emoteId={emoteId.trim()} type={type} size={64} />
      <Text size="xs" c="dimmed">
        Preview
      </Text>
    </Stack>
  );
}
