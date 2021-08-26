import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { store } from "../store";

export enum ChangeType {
    ADD = "add",
    REMOVE = "remove",
    REMOVED_RANDOM = "removed_random",
}

export enum EmoteType {
    BTTV = "bttv",
    SEVENTV = "seventv",
}

export interface RawEmoteHistoryItem {
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    ID: number;
    ChannelTwitchID: string;
    Type: EmoteType;
    ChangeType: ChangeType;
    EmoteID: string;
}

export interface EmoteHistoryItem {
    CreatedAt: Date;
    UpdatedAt: Date;
    DeletedAt: Date | null;
    ID: number;
    ChannelTwitchID: string;
    Type: EmoteType;
    ChangeType: ChangeType;
    EmoteID: string;
}


const PAGE_SIZE = 20;

export function useEmoteHistory(): [Array<EmoteHistoryItem>, () => void, boolean, number, () => void, () => void] {
    const [page, setPage] = useState(1);
    const pageRef = useRef(page);
    pageRef.current = page;

    const [emoteHistory, setEmoteHistory] = useState<Array<EmoteHistoryItem>>([]);
    const [loading, setLoading] = useState(false);
    const managing = store.useState(s => s.managing);

    const fetchPredictions = () => {
        setLoading(true);

        const currentPage = pageRef.current;

        let endPoint = "/api/emoteHistory";
        const searchParams = new URLSearchParams();
        searchParams.append("page", page.toString());
        doFetch(Method.GET, endPoint, searchParams).then((resp) => {
            if (currentPage !== pageRef.current) {
                throw new Error("Page changed");
            }

            return resp
        }).then((items: Array<RawEmoteHistoryItem>) =>
            setEmoteHistory(
                items.map(
                    (item: RawEmoteHistoryItem) => (
                        {
                            ...item,
                            CreatedAt: new Date(item.CreatedAt),
                            UpdatedAt: new Date(item.CreatedAt),
                            DeletedAt: item.CreatedAt ? new Date(item.CreatedAt) : null,
                        }
                    )
                )
            )
        ).then(() => setLoading(false)).catch(err => {
            if (err.message !== "Page changed") {
                throw err;
            }
        });
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchPredictions, [managing, page]);

    return [emoteHistory, fetchPredictions, loading, page, () => emoteHistory.length === PAGE_SIZE ? setPage(page + 1) : undefined, () => page > 1 ? setPage(page - 1) : undefined];
}