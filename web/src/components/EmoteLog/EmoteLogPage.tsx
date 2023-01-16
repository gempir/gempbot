import { ArrowPathIcon, ChevronLeftIcon, ChevronRightIcon } from "@heroicons/react/24/solid";
import { useEmoteLog } from "../../hooks/useEmoteLog";
import { Emote } from "../Emote/Emote";

export function EmoteLogPage({ channel }: { channel: string }): JSX.Element {
    const { fetch, increasePage, decreasePage, emoteLog, loading, page } = useEmoteLog(channel);

    return <div className="p-4">
        <div className="p-4 bg-gray-800 rounded shadow relative select-none min-h-[20rem]">
            <div className="text-2xl flex gap-5 w-full" onClick={fetch}>
                <div className="flex gap-2 items-center">
                    <div onClick={decreasePage} className="cursor-pointer hover:text-blue-500">
                        <ChevronLeftIcon className="h-6" />
                    </div>
                    <div className="text-base w-4 text-center">
                        {page}
                    </div>
                    <div onClick={increasePage} className="cursor-pointer hover:text-blue-500">
                        <ChevronRightIcon className="h-6" />
                    </div>
                </div>
                <ArrowPathIcon className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
            </div>
            {emoteLog.length === 0 && !loading && <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 font-bold text-5xl text-slate-600">nothing yet</div>}
            <table className="w-full table-auto">
                <thead>
                    <tr>
                        <th className="min-w-[6em] max-w-[8em]">Emote</th>
                        <th className="min-w-[6em]">Code</th>
                        <th className="min-w-[12em]">By</th>
                        <th className="min-w-[6em]">Type</th>
                    </tr>
                </thead>
                <tbody>
                    {emoteLog.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                        <td className="text-center px-5"><Emote id={item.EmoteID} /></td>
                        <td className="text-center px-10">{item.EmoteCode}</td>
                        <td className="text-center px-10">{item.AddedBy}</td>
                        <td className="text-center px-10">{item.Type === "seventv" ? "Redemption" : "Election"}</td>
                        <td className="p-3 text-center whitespace-nowrap">{item.CreatedAt.format('L LT')}</td>
                    </tr>)}
                </tbody>
            </table>
        </div>
    </div>
}