import { YouTubeIcon } from "../../icons/YouTube";
import { buildYouTubeThumbnailUrl, getYouTubeID } from "../../service/youtube";
import { Queue } from "../../types/Media";

export function MediaQueue({ queue }: { queue: Queue }) {
    return <div className="p-2 bg-gray-800 rounded shadow relative min-w-[15rem] flex flex-col gap-2">
        {queue.map((item, index) => <>
            <div className="p-2 bg-gray-900 w-full relative">
                <div className="text-gray-400 absolute top-1 right-1">{item.Author}</div>
                <img src={buildYouTubeThumbnailUrl(getYouTubeID(item.Url) ?? "")} className="max-h-12 h-auto w-auto max-w-12" />
                <a href={item.Url} target="_blank" rel="noreferrer" className="text-gray-400">
                    <YouTubeIcon className="w-5 absolute bottom-1 right-1" />
                </a>
            </div>
        </>)}
    </div>
}