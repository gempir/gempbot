import { ChevronLeftIcon, ChevronRightIcon, MinusCircleIcon, ArrowPathIcon, StopIcon } from "@heroicons/react/24/solid";
import { useEmotehistory } from "../../hooks/useEmotehistory";
import { Emote } from "../Emote/Emote";

export function Table({ channel, added, removeable, blockable, title }: { channel?: string, added: boolean, removeable: boolean, blockable: boolean, title: string }) {
    const [history, fetch, loading, page, increasePage, decreasePage, remove, block] = useEmotehistory(added, channel);

    return <div className="p-4 bg-gray-800 rounded shadow relative select-none min-w-[40rem] min-h-[20rem]">
        <div className="flex gap-5 w-full items-center justify-between" onClick={fetch}>
            <div className="text-2xl flex items-center gap-5">
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
            <h3 className="text-lg text-gray-400">{title}</h3>
        </div>
        {history.length === 0 && !loading && <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 font-bold text-5xl text-slate-600">nothing yet</div>}
        <table className="w-full table-auto">
            <thead>
                <tr>
                    <th className="min-w-[6em] max-w-[8em]">Emote</th>
                    <th className="min-w-[6em]">Type</th>
                    <th className="min-w-[12em]">Updated At</th>
                    {removeable && <th className="min-w-[5em]">Remove</th>}
                    {blockable && <th className="min-w-[5em]">Block</th>}
                </tr>
            </thead>
            <tbody>
                {history.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                    <td className="text-center px-5"><Emote id={item.EmoteID} type={item.Type} /></td>
                    <td className="text-center px-10">{item.Type}</td>
                    <td className="p-3 text-center whitespace-nowrap">{item.UpdatedAt.toLocaleDateString()} {item.UpdatedAt.toLocaleTimeString()}</td>
                    {removeable && !item.Blocked &&
                        <td className="text-center px-5 cursor-pointer hover:text-blue-500 group" onClick={() => remove(item.EmoteID)}>
                            <MinusCircleIcon className="h-6 mx-auto" /><span className="absolute z-50 hidden p-2 mx-10 -my-12 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Remove emote from history, preventing future removal of it. Make sure you reduce the slots or free up a slot</span>
                        </td>
                    }
                    {blockable && !item.Blocked &&
                        <td className="text-center px-5 cursor-pointer hover:text-blue-500 group" onClick={() => block(item.EmoteID)}>
                            <StopIcon className="h-6 mx-auto" /><span className="absolute z-50 hidden p-2 mx-10 -my-12 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Block emote and remove it from the channel. Don't fill the slot after it has been removed.</span>
                        </td>
                    }
                </tr>)}
            </tbody>
        </table>
    </div>
}