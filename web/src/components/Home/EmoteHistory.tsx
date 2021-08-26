import { EmoteType, useEmoteHistory } from "../../hooks/useEmoteHistory";
import { ChevronLeft } from "../../icons/ChevronLeft";
import { ChevronRight } from "../../icons/ChevronRight";
import { Refresh } from "../../icons/Refresh";


export function EmoteHistory() {
    const [history, fetch, loading, page, increasePage, decreasePage] = useEmoteHistory();

    return <div className="p-4 bg-gray-800 rounded shadow w-2/6 relative select-none">
        <div className="text-2xl flex gap-5 w-full" onClick={fetch}>
            <div className="flex gap-2 items-center">
                <div onClick={decreasePage} className="cursor-pointer hover:text-blue-500">
                    <ChevronLeft />
                </div>
                <div className="text-base w-4 text-center">
                    {page}
                </div>
                <div onClick={increasePage} className="cursor-pointer hover:text-blue-500">
                    <ChevronRight />
                </div>
            </div>
            <Refresh className={"hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
        </div>
        <table className="w-full">
            <thead>
                <tr>
                    <th>Emote</th>
                    <th>Type</th>
                    <th>Change Type</th>
                    <th>Updated At</th>
                </tr>
            </thead>
            <tbody>
                {history.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                    <td className="text-center"><Emote id={item.EmoteID} type={item.Type} /></td>
                    <td className="text-center">{item.Type}</td>
                    <td className="text-center">{item.ChangeType}</td>
                    <td className="p-3 text-center">{item.UpdatedAt.toLocaleDateString()} {item.UpdatedAt.toLocaleTimeString()}</td>
                </tr>)}
            </tbody>
        </table>
    </div>
}

function Emote({ id, type }: { id: string, type: EmoteType }): JSX.Element {
    const url = `https://cdn.betterttv.net/emote/${id}/1x`;

    if (type === EmoteType.SEVENTV) {
        //
    }

    return <img className="inline-block" src={url} alt={id} />
}

// CreatedAt: Date;
// UpdatedAt: Date;
// DeletedAt: Date | null;
// ID: number;
// ChannelTwitchID: string;
// Type: string;
// ChangeType: string;
// EmoteID: string;