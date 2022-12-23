import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

interface RawNomination {
    EmoteID: string
    ChannelTwitchID: string
    Votes: number
    EmoteCode: string
    NominatedBy: string
    CreatedAt: string
    UpdatedAt: string
}

export type Nomination = RawNomination & {
    CreatedAt: Date,
    UpdatedAt: Date,
}

const PAGE_SIZE = 20;

interface Return {
    nominations: Array<Nomination>,
    fetch: () => void,
    loading: boolean,
    page: number,
    increasePage: () => void,
    decreasePage: () => void,
}

export function useNominations(channel: string): Return {
    const [page, setPage] = useState(1);
    const pageRef = useRef(page);
    pageRef.current = page;

    const [nominations, setBlocks] = useState<Array<Nomination>>([]);
    const [loading, setLoading] = useState(false);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);

    const fetchNominations = () => {
        setLoading(true);

        const currentPage = pageRef.current;

        const endPoint = "/api/nominations";
        const searchParams = new URLSearchParams();
        searchParams.append("page", page.toString());
        searchParams.append("channel", channel);
        doFetch({ apiBaseUrl }, Method.GET, endPoint, searchParams).then((resp) => {
            if (currentPage !== pageRef.current) {
                throw new Error("Page changed");
            }

            return resp
        }).then(rawNoms => setBlocks(rawNoms.map((rawNom: RawNomination) => ({ ...rawNom, CreatedAt: new Date(rawNom.CreatedAt), UpdatedAt: new Date(rawNom.UpdatedAt) }))))
            .then(() => setLoading(false)).catch(err => {
                if (err.message !== "Page changed") {
                    throw err;
                }
            });
    };

    useEffect(fetchNominations, [page]);

    return {
        nominations: nominations,
        fetch: fetchNominations,
        loading: loading,
        page: page,
        increasePage: () => nominations.length === PAGE_SIZE ? setPage(page + 1) : undefined,
        decreasePage: () => page > 1 ? setPage(page - 1) : undefined,
    };
}