import { MediaPlayer } from "./MediaPlayer";
import { MediaQueue } from "./MediaQueue";


export function MediaPage({ channel = "" }: { channel?: string }): JSX.Element {
    return <div className="p-4 w-full flex gap-4">
        <MediaPlayer channel={channel} />
        <MediaQueue />
    </div>;
}