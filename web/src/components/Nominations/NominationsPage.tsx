import { ArrowPathIcon, ArrowUpCircleIcon } from "@heroicons/react/24/solid";
import { createLoginUrl } from "../../factory/createLoginUrl";
import { EmoteType } from "../../hooks/useEmotehistory";
import { useNominations } from "../../hooks/useNominations";
import { useStore } from "../../store";
import { Emote } from "../Emote/Emote";

export function NominationsPage({ channel }: { channel: string }): JSX.Element {
    const { nominations, fetch, loading, makeVote } = useNominations(channel);
    const scToken = useStore(state => state.scToken);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const twitchClientId = useStore(state => state.twitchClientId);
    const url = createLoginUrl(apiBaseUrl, twitchClientId);

    const handleVote = (emoteID: string) => {
        if (!scToken) {
            window.localStorage.setItem("redirect", window.location.pathname);
            window.location.href = url.toString();
        }

        makeVote(emoteID);
    };

    return <div className="p-4 w-full flex gap-4">
        <div className="p-4 bg-gray-800 rounded shadow relative select-none">
            <div className="text-2xl flex gap-5 w-full" onClick={fetch}>
                <ArrowPathIcon className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
            </div>
            {nominations.length === 0 && !loading && <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 font-bold text-5xl text-slate-600">nothing yet</div>}
            <table className="w-full table-auto">
                <thead>
                    <tr>
                        <th className="min-w-[6em] max-w-[8em]">Emote</th>
                        <th className="min-w-[6em]">Code</th>
                        <th className="min-w-[6em]">Votes</th>
                        <th className="min-w-[6em]">Nominated By</th>
                        <th className="min-w-[12em]">Created At</th>
                        <th className="min-w-[6em]"></th>
                    </tr>
                </thead>
                <tbody>
                    {nominations.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                        <td className="text-center px-5"><Emote id={item.EmoteID} type={EmoteType.SEVENTV} /></td>
                        <td className="text-center px-10">{item.EmoteCode}</td>
                        <td className="text-center px-10">{item.Votes.length}</td>
                        <td className="text-center px-10">{item.NominatedBy}</td>
                        <td className="p-3 text-center whitespace-nowrap">{item.CreatedAt.toLocaleDateString()} {item.CreatedAt.toLocaleTimeString()}</td>
                        <td className="text-center px-10"><ArrowUpCircleIcon onClick={() => handleVote(item.EmoteID)} className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} /></td>
                    </tr>)}
                </tbody>
            </table>
        </div>
    </div>;
}