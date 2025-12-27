import { Stack } from "@mantine/core";
import { useEmotehistory } from "../../../hooks/useEmotehistory";
import { Table } from "./Table";

export function Emotehistory() {
  const seventvHistory = useEmotehistory("seventv");

  return (
    <Stack gap="lg">
      <Table
        title="7tv_emote_history"
        history={seventvHistory.history}
        page={seventvHistory.page}
        totalPages={seventvHistory.totalPages}
        loading={seventvHistory.loading}
        onPageChange={seventvHistory.setPage}
      />
    </Stack>
  );
}
