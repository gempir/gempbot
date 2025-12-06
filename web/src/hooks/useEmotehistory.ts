import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

export interface EmotehistoryItem {
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date | null;
  id: number;
  channelTwitchID: string;
  type: string;
  changeType: string;
  emoteID: string;
  userLogin: string;
}

interface RawEmotehistoryItem {
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  ID: number;
  ChannelTwitchID: string;
  Type: string;
  ChangeType: string;
  EmoteID: string;
  UserLogin: string;
}

const PAGE_SIZE = 20;

type Return = {
  history: Array<EmotehistoryItem>;
  loading: boolean;
  page: number;
  totalPages: number;
  setPage: (page: number) => void;
  fetch: () => void;
};

export function useEmotehistory(type?: string): Return {
  const [page, setPage] = useState(1);
  const pageRef = useRef(page);
  pageRef.current = page;

  const [emotehistory, setEmotehistory] = useState<Array<EmotehistoryItem>>([]);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(false);
  const managing = useStore((state) => state.managing);
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);
  const scToken = useStore((state) => state.scToken);

  const fetch = () => {
    setLoading(true);
    const currentPage = pageRef.current;

    const endPoint = "/api/emotehistory";
    const searchParams = new URLSearchParams();
    searchParams.append("page", page.toString());
    searchParams.append("limit", PAGE_SIZE.toString());
    if (type) {
      searchParams.append("type", type);
    }

    doFetch(
      { apiBaseUrl, managing, scToken },
      Method.GET,
      endPoint,
      searchParams,
    )
      .then((resp) => {
        if (currentPage !== pageRef.current) {
          throw new Error("Page changed");
        }
        return resp;
      })
      .then((items: Array<RawEmotehistoryItem>) => {
        const mapped = items.map((item: RawEmotehistoryItem) => ({
          createdAt: new Date(item.CreatedAt),
          updatedAt: new Date(item.UpdatedAt),
          deletedAt: item.DeletedAt ? new Date(item.DeletedAt) : null,
          id: item.ID,
          channelTwitchID: item.ChannelTwitchID,
          type: item.Type,
          changeType: item.ChangeType,
          emoteID: item.EmoteID,
          userLogin: item.UserLogin,
        }));
        setEmotehistory(mapped);
        setTotalPages(Math.max(1, Math.ceil(items.length / PAGE_SIZE)));
      })
      .then(() => setLoading(false))
      .catch((err) => {
        if (err.message !== "Page changed") {
          setLoading(false);
        }
      });
  };

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(fetch, [managing, page, type]);

  return {
    history: emotehistory,
    fetch,
    loading,
    page,
    totalPages,
    setPage,
  };
}
