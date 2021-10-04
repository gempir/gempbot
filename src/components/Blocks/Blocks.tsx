import React, { useState } from "react";
import { useBlocks } from "../../hooks/useBlocks";
import { EmoteType } from "../../hooks/useEmotehistory";
import { Emote } from "../Emote/Emote";
import { ChevronLeft } from "../../icons/ChevronLeft";
import { ChevronRight } from "../../icons/ChevronRight";
import { Refresh } from "../../icons/Refresh";

export function Blocks() {
    const { blocks, block, loading, increasePage, decreasePage, page, fetch } = useBlocks();

    const [newEmoteType, setNewEmoteType] = useState<EmoteType>(EmoteType.SEVENTV);
    const [newEmoteID, setNewEmoteID] = useState<string>("");

    const blockEmote = () => {
        if (newEmoteID === "") {
            return;
        }

        block(newEmoteID, newEmoteType);
        setNewEmoteID("");
    };

    return <div className="p-4">
        <div className="p-4 bg-gray-800 rounded shadow relative">
            <div className="flex gap-5 items-center mb-5">
                <h2 className="text-xl">Blocks</h2>
                <div className="text-2xl flex gap-5 w-full select-none" onClick={fetch}>
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
            </div>
            <table className={"w-full" + (loading ? " animate-pulse opacity-10" : "")}>
                <thead>
                    <tr className="border-b-8 border-transparent">
                        <th />
                        <th className="px-5">Emote</th>
                        <th className="text-left pl-5">EmoteID</th>
                        <th className="px-5">Type</th>
                        <th className="px-5">Created</th>
                    </tr>
                </thead>
                <tbody>
                    {blocks.map(block => <tr key={block.ChannelTwitchID + block.EmoteID + block.EmoteID}>
                        <th></th>
                        <th><Emote id={block.EmoteID} type={block.Type} /></th>
                        <th>{block.EmoteID}</th>
                        <th>{block.Type}</th>
                        <th>{block.CreatedAt.toLocaleDateString()} {block.CreatedAt.toLocaleTimeString()}</th>
                    </tr>)}
                </tbody>
            </table>
            <div className="mt-5 flex gap-5">
                <input type="text" placeholder="EmoteId,EmoteId2,EmoteId3" className="w-full p-1 bg-transparent leading-6 rounded" value={newEmoteID} onChange={e => setNewEmoteID(e.target.value)} />
                <select className="p-1 pr-10 bg-transparent leading-6 rounded appearance-none" onChange={e => setNewEmoteType(e.target.value as EmoteType)} value={newEmoteType}>
                    <option>{EmoteType.SEVENTV}</option>
                    <option>{EmoteType.BTTV}</option>
                </select>
                <button className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block cursor-pointer" onClick={blockEmote}>block</button>
            </div>
        </div>
    </div>;
}