import { ArrowPathIcon, ArrowUpCircleIcon, SquaresPlusIcon, StopIcon } from "@heroicons/react/24/solid";
import { useEffect } from "react";
import { createLoginUrl } from "../../factory/createLoginUrl";
import { EmoteType } from "../../hooks/useEmotehistory";
import { Nomination, useNominations } from "../../hooks/useNominations";
import { useStore } from "../../store";
import { Election } from "../../types/Election";
import { Emote } from "../Emote/Emote";
import { ElectionStatus } from "./ElectionStatus";

export function NominationsView({ channel, election }: { channel: string, election?: Election }): JSX.Element {
    const { nominations, fetch, loading, vote, unvote, block } = useNominations(channel);
    const scToken = useStore(state => state.scToken);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scTokenContent = useStore(state => state.scTokenContent);
    const twitchClientId = useStore(state => state.twitchClientId);
    const managing = useStore(state => state.managing);
    const url = createLoginUrl(apiBaseUrl, twitchClientId);
    useEffect(fetch, [channel]);

    const handleVote = (nomination: Nomination) => {
        if (!scToken) {
            window.localStorage.setItem("redirect", window.location.pathname);
            window.location.href = url.toString();
        }

        if (nomination.Votes.some(value => value.VoteBy === scTokenContent?.UserID)) {
            unvote(nomination.EmoteID);
        } else {
            vote(nomination.EmoteID);
        }
    };

    const blockable = scTokenContent?.Login === channel || managing === channel;

    return <div className="flex flex-col gap-3">
        <div className="p-4 bg-gray-800 rounded shadow relative select-none">
            <ElectionStatus election={election} />
        </div>
        <div className="flex gap-3 min-h-[20em]">
            <div className="p-4 bg-gray-800 rounded shadow relative select-none">
                <div className="flex gap-5 items-center mb-5">
                    <h2 className="text-xl">Nominations</h2>
                    <div className="text-2xl flex gap-5 w-full select-none" onClick={fetch}>
                        <ArrowPathIcon className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
                    </div>
                    <div>
                        <a className={"relative group"} href={`https://twitch.tv/popout/${channel}/chat`} target="_blank" rel="noreferrer">
                            <SquaresPlusIcon className="h-6 hover:text-blue-500 cursor-pointer" />
                            <span className="absolute z-50 hidden p-2 mx-10 -my-10 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Nominate with Channel Points</span>
                        </a>
                    </div>
                </div>
                {nominations.length === 0 && !loading && <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 font-bold text-5xl text-slate-600">nothing yet</div>}
                <table className="w-full table-auto">
                    <thead>
                        <tr>
                            <th className="text-left">Votes</th>
                            <th className="min-w-[6em] max-w-[8em]">Emote</th>
                            <th className="min-w-[6em] max-w-[250px] truncate">Code</th>
                            <th className="min-w-[6em] max-w-[250px] truncate">Nominated By</th>
                            <th className="min-w-[6em]"></th>
                            {blockable && <th className="min-w-[6em]"></th>}
                        </tr>
                    </thead>
                    <tbody>
                        {nominations.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                            <td className="text-center">{item.Votes.length}</td>
                            <td className="text-center px-5"><Emote id={item.EmoteID} type={EmoteType.SEVENTV} /></td>
                            <td className="text-center px-10 max-w-[250px] truncate">{item.EmoteCode}</td>
                            <td className="text-center px-10 max-w-[250px] truncate">{item.NominatedBy}</td>
                            <td className="text-center px-10"><ArrowUpCircleIcon onClick={() => handleVote(item)} className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "" + (item.Votes.some(value => value.VoteBy === scTokenContent?.UserID) ? "text-blue-600" : ""))} /></td>
                            {blockable && <td className="text-center px-5 cursor-pointer hover:text-blue-500 group" onClick={() => block(item.EmoteID)}>
                                <StopIcon className="h-6 mx-auto" /><span className="absolute z-50 hidden p-2 mx-10 -my-12 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Block emote and remove it from election</span>
                            </td>}
                        </tr>)}
                    </tbody>
                </table>
            </div>
        </div>
    </div>;
}