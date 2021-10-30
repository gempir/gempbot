import React from "react";
import { useEmotehistory } from "../../hooks/useEmotehistory";
import { Emote } from "../Emote/Emote";
import { RefreshIcon, ChevronLeftIcon, ChevronRightIcon } from "@heroicons/react/solid";


export function Emotehistory({channel}: {channel?: string}) {
    const [history, fetch, loading, page, increasePage, decreasePage] = useEmotehistory(channel);

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
                </tr>
            </thead>
            <tbody>
                {history.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                    <td className="text-center px-5"><Emote id={item.EmoteID} type={item.Type} /></td>
                    <td className="text-center px-10">{item.Type}</td>
                    <td className="text-center px-10">{item.ChangeType}</td>
                    <td className="p-3 text-center whitespace-nowrap">{item.UpdatedAt.toLocaleDateString()} {item.UpdatedAt.toLocaleTimeString()}</td>
                </tr>)}
            </tbody>
        </table>
    </div>
}
