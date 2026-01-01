import { Stack } from "@mantine/core";
import type { EmotehistoryItem } from "../../../hooks/useEmotehistory";
import { useEmotehistory } from "../../../hooks/useEmotehistory";
import { doFetch, Method } from "../../../service/doFetch";
import { useStore } from "../../../store";
import { Table } from "./Table";

export function Emotehistory() {
  const seventvHistory = useEmotehistory("seventv");
  const managing = useStore((state) => state.managing);
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);
  const scToken = useStore((state) => state.scToken);

  const handleRemove = async (item: EmotehistoryItem) => {
    const searchParams = new URLSearchParams();
    searchParams.append("emoteId", item.emoteID);

    await doFetch(
      { apiBaseUrl, managing, scToken },
      Method.DELETE,
      "/api/emotehistory",
      searchParams,
    );

    // Refresh the history
    seventvHistory.fetch();
  };

  const handleBlock = async (item: EmotehistoryItem) => {
    const searchParams = new URLSearchParams();
    searchParams.append("emoteId", item.emoteID);

    await doFetch(
      { apiBaseUrl, managing, scToken },
      Method.PATCH,
      "/api/emotehistory",
      searchParams,
    );

    // Refresh the history
    seventvHistory.fetch();
  };

  return (
    <Stack gap="lg">
      <Table
        title="7tv_emote_history"
        history={seventvHistory.history}
        page={seventvHistory.page}
        totalPages={seventvHistory.totalPages}
        loading={seventvHistory.loading}
        onPageChange={seventvHistory.setPage}
        onApprove={handleRemove}
        onDeny={handleBlock}
        approveLabel="remove"
        denyLabel="block"
      />
    </Stack>
  );
}
