import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";
import { EmoteType } from "./useEmotehistory";

export interface Block {
    ChannelTwitchID: string
    Type: EmoteType
    EmoteID: string
    CreatedAt: string
}

const PAGE_SIZE = 20;

export function useBlocks(): [Array<Block>, () => void, boolean, number, () => void, () => void] {
    const [page, setPage] = useState(1);
    const pageRef = useRef(page);
    pageRef.current = page;

    const [blocks, setBlocks] = useState<Array<Block>>([]);
    const [loading, setLoading] = useState(false);
    const managing = useStore(state => state.managing);

    const fetchBlocks = () => {
        setLoading(true);

        const currentPage = pageRef.current;
        
        const endPoint = "/api/blocks";
        const searchParams = new URLSearchParams();
        searchParams.append("page", page.toString());
        doFetch(Method.GET, endPoint, searchParams).then((resp) => {
            if (currentPage !== pageRef.current) {
                throw new Error("Page changed");
            }

            return resp
        }).then(setBlocks).then(() => setLoading(false)).catch(err => {
            if (err.message !== "Page changed") {
                throw err;
            }
        });
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchBlocks, [managing, page]);

    return [blocks, fetchBlocks, loading, page, () => blocks.length === PAGE_SIZE ? setPage(page + 1) : undefined, () => page > 1 ? setPage(page - 1) : undefined];
}