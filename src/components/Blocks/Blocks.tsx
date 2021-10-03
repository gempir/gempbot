import { useBlocks } from "../../hooks/useBlocks";
import { Emote } from "../Emote/Emote";

export function Blocks() {
    const [blocks] = useBlocks();

    return <div className="p-4">
        <div className="p-4 bg-gray-800 rounded shadow relative">
        <h2 className="mb-4 text-xl">Blocks</h2>
        <table className="w-full">
            <thead>
                <tr className="border-b-8 border-transparent">
                    <th />
                    <th className="px-5">Emote</th>
                    <th className="text-left pl-5">EmoteID</th>
                    <th className="px-5">Platform</th>
                </tr>
            </thead>
            <tbody>
                {blocks.map(block => <tr key={block.ChannelTwitchID + block.EmoteID + block.EmoteID}>
                    <th></th>
                    <th><Emote id={block.EmoteID} type={block.Type}  /></th>
                    <th>{block.EmoteID}</th>
                    <th>{block.Type}</th>
                </tr>)}
            </tbody>
        </table>
    </div>
    </div>;
}