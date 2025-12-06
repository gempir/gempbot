import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

interface RawBlock {
  ChannelTwitchID: string;
  Type: string;
  EmoteID: string;
  CreatedAt: string;
}

export type Block = RawBlock & {
  CreatedAt: Date;
};

const PAGE_SIZE = 20;

interface Return {
  blocks: Array<Block>;
  fetch: () => void;
  loading: boolean;
  page: number;
  totalPages: number;
  increasePage: () => void;
  decreasePage: () => void;
  addBlock: (emoteIds: string, type: string) => void;
  removeBlock: (block: Block) => void;
  setPage: (page: number) => void;
}

export function useBlocks(): Return {
  const [page, setPage] = useState(1);
  const pageRef = useRef(page);
  pageRef.current = page;

  const [blocks, setBlocks] = useState<Array<Block>>([]);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(false);
  const managing = useStore((state) => state.managing);
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);
  const scToken = useStore((state) => state.scToken);

  const fetchBlocks = () => {
    setLoading(true);

    const currentPage = pageRef.current;

    const endPoint = "/api/blocks";
    const searchParams = new URLSearchParams();
    searchParams.append("page", page.toString());
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
      .then((rawBlocks) => {
        const mapped = rawBlocks.map((rawBlock: RawBlock) => ({
          ...rawBlock,
          CreatedAt: new Date(rawBlock.CreatedAt),
        }));
        setBlocks(mapped);
        setTotalPages(Math.max(1, Math.ceil(rawBlocks.length / PAGE_SIZE)));
      })
      .then(() => setLoading(false))
      .catch((err) => {
        if (err.message !== "Page changed") {
          throw err;
        }
      });
  };

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(fetchBlocks, []);

  const addBlock = (emoteIds: string, type: string) => {
    setLoading(true);

    const endPoint = "/api/blocks";
    const searchParams = new URLSearchParams();
    doFetch(
      { apiBaseUrl, managing, scToken },
      Method.PATCH,
      endPoint,
      searchParams,
      { emoteIds: emoteIds, type: type },
    )
      .then(fetchBlocks)
      .catch((err) => {
        setLoading(false);
        throw err;
      });
  };

  const removeBlock = (block: Block) => {
    setLoading(true);

    const endPoint = "/api/blocks";
    const searchParams = new URLSearchParams();
    doFetch(
      { apiBaseUrl, managing, scToken },
      Method.DELETE,
      endPoint,
      searchParams,
      block,
    )
      .then(fetchBlocks)
      .catch((err) => {
        setLoading(false);
        throw err;
      });
  };

  return {
    blocks: blocks,
    fetch: fetchBlocks,
    loading: loading,
    page: page,
    totalPages: totalPages,
    setPage: setPage,
    increasePage: () =>
      blocks.length === PAGE_SIZE ? setPage(page + 1) : undefined,
    decreasePage: () => (page > 1 ? setPage(page - 1) : undefined),
    addBlock: addBlock,
    removeBlock: removeBlock,
  };
}
