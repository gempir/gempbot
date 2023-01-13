import { EmoteType } from "../../hooks/useEmotehistory";

export type EmoteSize = 1 | 2 | 3 | 4;

export function Emote({ id, type = EmoteType.SEVENTV, size = 1 }: { id: string, type?: EmoteType, size?: EmoteSize }): JSX.Element {
    let url = `https://cdn.betterttv.net/emote/${id}/${size}x`;
    let hrefUrl = `https://betterttv.com/emotes/${id}`;

    if (type === EmoteType.SEVENTV) {
        url = `https://cdn.7tv.app/emote/${id}/${size}x`
        hrefUrl = `https://7tv.app/emotes/${id}`;
    }

    return <a href={hrefUrl} target="_blank">
        <img className="inline-block" style={{ minWidth: 28 }} src={url} alt={id} />
    </a>;
}