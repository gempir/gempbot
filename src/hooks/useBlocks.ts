import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";
import { EmoteType } from "./useEmotehistory";

interface RawBlock {
    ChannelTwitchID: string
    Type: EmoteType
    EmoteID: string
    CreatedAt: string
}

export type Block = RawBlock & {
    CreatedAt: Date,
}

const PAGE_SIZE = 20;

interface Return {
    blocks: Array<Block>,
    fetch: () => void,
    loading: boolean,
    page: number,
    increasePage: () => void,
    decreasePage: () => void,
    block: (emoteIds: string, type: EmoteType) => void,
}

export function useBlocks(): Return {
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
        }).then(rawBlocks => setBlocks(rawBlocks.map((rawBlock: RawBlock) => ({ ...rawBlock, CreatedAt: new Date(rawBlock.CreatedAt) }))))
            .then(() => setLoading(false)).catch(err => {
                if (err.message !== "Page changed") {
                    throw err;
                }
            });
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchBlocks, [managing, page]);

    const block = (emoteIds: string, type: EmoteType) => {
        setLoading(true);

        const endPoint = "/api/blocks";
        const searchParams = new URLSearchParams();
        doFetch(Method.PATCH, endPoint, searchParams, { emoteIds: emoteIds, type: type }).then(fetchBlocks).catch(err => {
            setLoading(false);
            throw err;
        });
    };

    return {
        blocks: blocks,
        fetch: fetchBlocks,
        loading: loading,
        page: page,
        increasePage: () => blocks.length === PAGE_SIZE ? setPage(page + 1) : undefined,
        decreasePage: () => page > 1 ? setPage(page - 1) : undefined,
        block: block,
    };
}