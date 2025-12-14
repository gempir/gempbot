import { Anchor, Box, Image, Text } from "@mantine/core";
import { QuestionMarkCircleIcon } from "@heroicons/react/24/outline";
import { useState } from "react";

interface EmoteProps {
  emoteId: string;
  type: string;
  size?: number;
}

export function Emote({ emoteId, type, size = 32 }: EmoteProps) {
  const [imageError, setImageError] = useState(false);

  const getEmoteUrl = () => {
    if (type === "BTTV") {
      return `https://cdn.betterttv.net/emote/${emoteId}/3x`;
    } else if (type === "7TV") {
      return `https://cdn.7tv.app/emote/${emoteId}/2x.webp`;
    }
    return "";
  };

  const getEmotePageUrl = () => {
    if (type === "BTTV") {
      return `https://betterttv.com/emotes/${emoteId}`;
    } else if (type === "7TV") {
      return `https://7tv.app/emotes/${emoteId}`;
    }
    return "";
  };

  const fallbackContent = (
    <Box
      w={size}
      h={size}
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        backgroundColor: "#e9ecef",
        borderRadius: "4px",
      }}
    >
      <QuestionMarkCircleIcon style={{ width: size * 0.5, height: size * 0.5, color: "#868e96" }} />
    </Box>
  );

  return (
    <Anchor href={getEmotePageUrl()} target="_blank" rel="noopener noreferrer">
      {imageError ? (
        fallbackContent
      ) : (
        <Image
          src={getEmoteUrl()}
          alt={`${type} emote`}
          w={size}
          h={size}
          fit="contain"
          onError={() => setImageError(true)}
        />
      )}
    </Anchor>
  );
}
