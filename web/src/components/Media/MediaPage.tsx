import { MediaPlayer } from "./MediaPlayer";


export function MediaPage({ channel = "" }: { channel?: string }): JSX.Element {
    return <div className="p-4 w-full">
        <MediaPlayer channel={channel} />
    </div>;
}