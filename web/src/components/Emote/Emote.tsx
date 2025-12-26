import { QuestionMarkCircleIcon } from "@heroicons/react/24/outline";
import { Anchor, Box, Image } from "@mantine/core";
import { useState } from "react";

interface EmoteProps {
  emoteId: string;
  type: string;
  size?: number;
}

export function Emote({ emoteId, type, size = 32 }: EmoteProps) {
  const [imageError, setImageError] = useState(false);

  const normalizedType = type?.toUpperCase();

  const getEmoteUrl = () => {
    console.log('Emote component - ID:', emoteId, 'Type:', type, 'Normalized:', normalizedType);
    if (normalizedType === "BTTV") {
      return `https://cdn.betterttv.net/emote/${emoteId}/1x.webp`;
    } else if (normalizedType === "7TV") {
      return `https://cdn.7tv.app/emote/${emoteId}/1x.avif`;
    }
    console.log('No matching type found for:', type);
    return "";
  };

  const getEmotePageUrl = () => {
    if (normalizedType === "BTTV") {
      return `https://betterttv.com/emotes/${emoteId}`;
    } else if (normalizedType === "7TV") {
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
      <QuestionMarkCircleIcon
        style={{ width: size * 0.5, height: size * 0.5, color: "#868e96" }}
      />
    </Box>
  );

  const emoteUrl = getEmoteUrl();
  const emotePageUrl = getEmotePageUrl();

  console.log('Constructed URLs - Emote:', emoteUrl, 'Page:', emotePageUrl);

  // If no valid URLs, just show the content without a link
  if (!emoteUrl || !emotePageUrl || !emoteId || !emoteId.trim()) {
    console.log('Showing fallback - emoteUrl:', !!emoteUrl, 'emotePageUrl:', !!emotePageUrl, 'emoteId:', emoteId);
    return fallbackContent;
  }

  return (
    <Anchor href={emotePageUrl} target="_blank" rel="noopener noreferrer">
      {imageError ? (
        fallbackContent
      ) : (
        <Image
          src={emoteUrl}
          alt={`${type} emote`}
          w={size}
          h={size}
          fit="contain"
          onError={(e) => {
            console.error('Image failed to load:', emoteUrl, e);
            setImageError(true);
          }}
        />
      )}
    </Anchor>
  );
}
