import { Stack } from '@mantine/core';
import { useEmotehistory } from "../../../hooks/useEmotehistory";
import { Table } from "./Table";

export function Emotehistory() {
    const bttvHistory = useEmotehistory('BTTV');
    const seventvHistory = useEmotehistory('7TV');

    return (
        <Stack gap="lg">
            <Table
                title="BTTV Emote History"
                history={bttvHistory.history}
                page={bttvHistory.page}
                totalPages={bttvHistory.totalPages}
                loading={bttvHistory.loading}
                onPageChange={bttvHistory.setPage}
            />

            <Table
                title="7TV Emote History"
                history={seventvHistory.history}
                page={seventvHistory.page}
                totalPages={seventvHistory.totalPages}
                loading={seventvHistory.loading}
                onPageChange={seventvHistory.setPage}
            />
        </Stack>
    );
}
