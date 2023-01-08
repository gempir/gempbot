import { ArrowDownCircleIcon, ArrowPathIcon, ArrowUpCircleIcon, LinkIcon, SquaresPlusIcon, StopIcon } from "@heroicons/react/24/solid";
import Link from "next/link";
import React, { useEffect, useState } from "react";
import seedrandom from "seedrandom";
import { createLoginUrl } from "../../factory/createLoginUrl";
import { EmoteType } from "../../hooks/useEmotehistory";
import { Nomination, useNominations } from "../../hooks/useNominations";
import { useStore } from "../../store";
import { Election } from "../../types/Election";
import { Emote } from "../Emote/Emote";
import { ElectionStatus } from "./ElectionStatus";

export function NominationsView({ channel, election, tableMode = false }: { channel: string, election?: Election, tableMode?: boolean }): JSX.Element {
    const { nominations, fetch, loading, vote, unvote, block, downvote, undownvote } = useNominations(channel);
    const scToken = useStore(state => state.scToken);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scTokenContent = useStore(state => state.scTokenContent);
    const twitchClientId = useStore(state => state.twitchClientId);
    const managing = useStore(state => state.managing);
    const url = createLoginUrl(apiBaseUrl, twitchClientId);
    const [seed, setSeed] = useState(scToken ?? "");

    useEffect(() => {
        if (!scToken) {
            let newSeed = window.localStorage.getItem("seed");
            if (!newSeed) {
                newSeed = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
                window.localStorage.setItem("seed", newSeed);
            }
            setSeed(newSeed);
        }

    }, [scToken]);

    useEffect(fetch, [channel]);

    const handleVote = (e: React.MouseEvent<HTMLDivElement | SVGSVGElement>, nomination: Nomination) => {
        e.preventDefault();
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

    const handleDownvote = (e: React.MouseEvent<HTMLDivElement | SVGSVGElement>, nomination: Nomination) => {
        e.preventDefault();
        if (!scToken) {
            window.localStorage.setItem("redirect", window.location.pathname);
            window.location.href = url.toString();
        }

        if (nomination.Downvotes.some(value => value.VoteBy === scTokenContent?.UserID)) {
            undownvote(nomination.EmoteID);
        } else {
            downvote(nomination.EmoteID);
        }
    };

    const rng = seedrandom(seed);
    const shuffledNominations = [...nominations].sort((a, b) => {
        if (a.EmoteID === b.EmoteID) {
            return 0;
        }
        return a.EmoteID < b.EmoteID ? -1 : 1;
    }).sort((a, b) => 0.5 - rng());

    const blockable = scTokenContent?.Login === channel || managing === channel;

    const hideVotes = nominations.map(nom => nom.Votes.map(vote => vote.VoteBy == scTokenContent?.UserID)).flat().filter(Boolean).length >= (election?.VoteAmount ?? 3);
    const hideDownvotes = nominations.map(nom => nom.Downvotes.map(vote => vote.VoteBy == scTokenContent?.UserID)).flat().filter(Boolean).length >= (election?.VoteAmount ?? 3);

    return <div className="flex flex-col gap-4">
        <ElectionStatus election={election} />
        <div className="flex gap-4 min-h-[20em]">
            <div className="p-4 bg-gray-800 rounded shadow relative select-none w-full">
                <div className="flex gap-5 items-center mb-5 justify-start">
                    <h2 className="text-xl">Nominations</h2>
                    <div className="text-2xl flex gap-5 select-none" onClick={fetch}>
                        <ArrowPathIcon className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
                    </div>
                    <div>
                        <a className={"relative group"} href={`https://twitch.tv/popout/${channel}/chat`} target="_blank" rel="noreferrer">
                            <SquaresPlusIcon className="h-6 hover:text-blue-500 cursor-pointer" />
                            <span className="absolute z-50 hidden p-2 mx-10 -my-10 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Nominate with Channel Points</span>
                        </a>
                    </div>
                    {tableMode &&
                        <div>
                            <a href={`/nominations/${channel}`} className={"relative group"} target="_blank" rel="noreferrer">
                                <>
                                    <LinkIcon className="h-6 hover:text-blue-500 cursor-pointer" />
                                    <span className="absolute z-50 hidden p-2 mx-10 -my-10 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Public page</span>
                                </>
                            </a>
                        </div>
                    }
                </div>
                {tableMode && <>
                    {nominations.length === 0 && !loading && <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 font-bold text-5xl text-slate-600">nothing yet</div>}
                    <table className="w-full table-auto">
                        <thead>
                            <tr>
                                <th className="text-left">Total</th>
                                <th className="text-center min-w-[6em]">Votes</th>
                                <th className="text-center min-w-[6em]">Downvotes</th>
                                <th className="min-w-[6em] max-w-[8em]">Emote</th>
                                <th className="min-w-[6em] max-w-[250px] truncate">Code</th>
                                <th className="min-w-[6em] max-w-[250px] truncate">Nominated By</th>
                                <th className="min-w-[6em]"></th>
                                <th className="min-w-[6em]"></th>
                                {blockable && <th className="min-w-[6em]"></th>}
                            </tr>
                        </thead>
                        <tbody>
                            {nominations.map((item, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                                {tableMode && <td className="text-center">{item.Votes.length - item.Downvotes.length}</td>}
                                {tableMode && <td className="text-center">{item.Votes.length}</td>}
                                {tableMode && <td className="text-center">{item.Downvotes.length}</td>}
                                <td className="text-center px-5"><Emote id={item.EmoteID} type={EmoteType.SEVENTV} /></td>
                                <td className="text-center px-10 max-w-[250px] truncate">{item.EmoteCode}</td>
                                <td className="text-center px-10 max-w-[250px] truncate">{item.NominatedBy}</td>
                                <td className="text-center px-10"><ArrowUpCircleIcon onClick={(e) => handleVote(e, item)} className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "" + (item.Votes.some(value => value.VoteBy === scTokenContent?.UserID) ? "text-blue-600" : ""))} /></td>
                                <td className="text-center px-10"><ArrowDownCircleIcon onClick={(e) => handleDownvote(e, item)} className={"h-6 hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "" + (item.Downvotes.some(value => value.VoteBy === scTokenContent?.UserID) ? "text-red-600" : ""))} /></td>
                                {blockable && <td className="text-center px-5 cursor-pointer hover:text-blue-500 group" onClick={() => block(item.EmoteID)}>
                                    <StopIcon className="h-6 mx-auto" /><span className="absolute z-50 hidden p-2 mx-10 -my-12 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">Block emote and remove it from election</span>
                                </td>}
                            </tr>)}
                        </tbody>
                    </table>
                </>}
                {!tableMode && <div className="w-full">
                    {shuffledNominations.length === 0 && !loading && <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 font-bold text-5xl text-slate-600">nothing yet</div>}
                    <div className="flex flex-wrap gap-3 w-full">
                        {shuffledNominations.map((item, index) => <div className={`text-center p-3 cursor-pointer flex flex-col gap-3 border border-transparent hover:border-gray-500 group relative ${(item.Votes.some(value => value.VoteBy === scTokenContent?.UserID) ? "hover:border-blue-600" : "")}`} key={index}>
                            <Emote size={2} id={item.EmoteID} type={EmoteType.SEVENTV} />
                            <ArrowUpCircleIcon onClick={(e) => handleVote(e, item)} className={"h-6 absolute top-0 right-1 hover:text-blue-500 cursor-pointer group" + (hideVotes ? "" : " group-hover:block") + (loading ? " animate-spin" : "") + (item.Votes.some(value => value.VoteBy === scTokenContent?.UserID) ? " text-blue-600 block" : " hidden")} />
                            <ArrowDownCircleIcon onClick={(e) => handleDownvote(e, item)} className={"h-6 absolute bottom-10 right-1 hover:text-red-500 cursor-pointer group" + (hideDownvotes ? "" : " group-hover:block") + (loading ? " animate-spin" : "") + (item.Downvotes.some(value => value.VoteBy === scTokenContent?.UserID) ? " text-red-600 block" : " hidden")} />
                            <span className="absolute z-50 hidden p-2 -mx-4 -my-14 text-center text-sm bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">by {item.NominatedBy}</span>
                            <span className="truncate max-w-xs pb-1">{item.EmoteCode}</span>
                        </div>
                        )}
                    </div>
                </div>}
            </div>
        </div>
    </div>;
}
