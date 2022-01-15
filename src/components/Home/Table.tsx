import { ChevronLeftIcon, ChevronRightIcon, MinusCircleIcon, RefreshIcon, StopIcon } from "@heroicons/react/solid";
import { useEmotehistory } from "../../hooks/useEmotehistory";
import { Emote } from "../Emote/Emote";

export function Table({ channel, added, removeable, blockable }: { channel?: string, added: boolean, removeable: boolean, blockable: boolean }) {
    const [history, fetch, loading, page, increasePage, decreasePage, remove, block] = useEmotehistory(added, channel);

    return <div className="p-4 bg-gray-800 rounded shadow relative select-none">
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
            <RefreshIcon className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
        </div>
        <table className="w-full table-auto">
            <thead>
                <tr>
                    <th>Emote</th>
                    <th>Type</th>
                    <th>Change Type</th>
                    <th>Updated At</th>
                    {removeable && <th>Remove</th>}
                    {blockable && <th>Block</th>}
                </tr>
            </thead>
            <tbody>
                {history.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                    <td className="text-center px-5"><Emote id={item.EmoteID} type={item.Type} /></td>
                    <td className="text-center px-10">{item.Type}</td>
                    <td className="text-center px-10">{item.ChangeType}</td>
                    <td className="p-3 text-center whitespace-nowrap">{item.UpdatedAt.toLocaleDateString()} {item.UpdatedAt.toLocaleTimeString()}</td>
                    {removeable && !item.Blocked &&
                        <td className="text-center px-5 cursor-pointer hover:text-blue-500 group" onClick={() => remove(item.EmoteID)}>
                            <MinusCircleIcon className="h-6" /><span className="absolute z-50 hidden p-2 mx-10 -my-12 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Remove emote from history, preventing future removal of it. Make sure you reduce the slots or free up a slot</span>
                        </td>
                    }
                    {blockable && !item.Blocked &&
                        <td className="text-center px-5 cursor-pointer hover:text-blue-500 group" onClick={() => block(item.EmoteID)}>
                            <StopIcon className="h-6" /><span className="absolute z-50 hidden p-2 mx-10 -my-12 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Block emote and remove it from the channel</span>
                        </td>
                    }
                </tr>)}
            </tbody>
        </table>
    </div>
}