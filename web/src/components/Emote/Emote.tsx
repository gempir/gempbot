import { EmoteType } from "../../hooks/useEmotehistory";

export function Emote({ id, type }: { id: string, type: EmoteType }): JSX.Element {
    let url = `https://cdn.betterttv.net/emote/${id}/1x`;
    let hrefUrl = `https://betterttv.com/emotes/${id}`;

    if (type === EmoteType.SEVENTV) {
        url = `https://cdn.7tv.app/emote/${id}/1x`
        hrefUrl = `https://7tv.app/emotes/${id}`;
    }

    return <a href={hrefUrl} target="_blank">
        <img className="inline-block" style={{ minWidth: 28 }} src={url} alt={id} />
    </a>;
}