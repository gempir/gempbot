import dayjs, { Dayjs } from "dayjs";
import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";


interface EmoteLogItemBase {
    EmoteID: string;
    EmoteCode: string;
    AddedBy: string;
    ChannelTwitchID: string;
    Type: string;
}

export type EmoteLogItem = EmoteLogItemBase & {
    CreatedAt: Dayjs;
}

export type RawEmoteLogItem = EmoteLogItemBase & {
    CreatedAt: string;
}

const PAGE_SIZE = 20;

type Return = {
    emoteLog: Array<EmoteLogItem>,
    increasePage: () => void,
    loading: boolean,
    page: number,
    decreasePage: () => void,
    fetch: () => void
};

export function useEmoteLog(channel?: string): Return {
    const [page, setPage] = useState(1);
    const pageRef = useRef(page);
    pageRef.current = page;

    const [emotelog, setEmotelog] = useState<Array<EmoteLogItem>>([]);
    const [loading, setLoading] = useState(false);
    const managing = useStore(state => state.managing);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const fetch = () => {
        setLoading(true);
        const currentPage = pageRef.current;

        const endPoint = "/api/emotelog";
        const searchParams = new URLSearchParams();
        if (channel) {
            searchParams.append("channel", channel);
        }
        searchParams.append("page", page.toString());
        searchParams.append("limit", PAGE_SIZE.toString());
        doFetch({ apiBaseUrl, managing, scToken }, Method.GET, endPoint, searchParams).then((resp) => {
            if (currentPage !== pageRef.current) {
                throw new Error("Page changed");
            }

            return resp
        }).then((items: Array<EmoteLogItem>) =>
            setEmotelog(
                items.map(
                    // @ts-ignore
                    (item: RawEmoteLogItem) => (
                        {
                            ...item,
                            CreatedAt: dayjs(item.CreatedAt),
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
    useEffect(fetch, [managing, page]);

    return { emoteLog: emotelog, fetch, loading, page, increasePage: () => emotelog.length === PAGE_SIZE ? setPage(page + 1) : undefined, decreasePage: () => page > 1 ? setPage(page - 1) : undefined };
}