import { useEffect, useState } from "react";
import EventService from "../../service/EventService";
import { EventMessage, Record } from "../../types/Events";
import { Meta } from "../Meta";
import { Records } from "../Records";

export function Home() {
    const [joinedChannels, setJoinedChannels] = useState(0);
    const [activeChannels, setActiveChannels] = useState(0);
    const [records, setRecords] = useState<Array<Record>>([])

    useEffect(() => {
        new EventService((message: EventMessage) => {
            setJoinedChannels(message.joinedChannels);
            setActiveChannels(message.activeChannels);
            setRecords(message.records);
        });
    }, []);

    return <div>
        <Meta activeChannels={activeChannels} joinedChannels={joinedChannels} />
        <Records records={records} />
    </div>
}