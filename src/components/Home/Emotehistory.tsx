import React from "react";
import { Table } from "./Table";


export function Emotehistory({ channel }: { channel?: string }) {

    return <div className="flex gap-4">
        <Table channel={channel} added={true} removeable={true} blockable={true} />
        <Table channel={channel} added={false} removeable={false} blockable={false} />
    </div>
}
