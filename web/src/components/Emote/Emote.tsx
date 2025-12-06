import { Anchor, Image } from "@mantine/core";

interface EmoteProps {
  emoteId: string;
  type: string;
}

export function Emote({ emoteId, type }: EmoteProps) {
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

  return (
    <Anchor href={getEmotePageUrl()} target="_blank" rel="noopener noreferrer">
      <Image
        src={getEmoteUrl()}
        alt={`${type} emote`}
        w={32}
        h={32}
        fit="contain"
      />
    </Anchor>
  );
}
