import { useContext, useEffect, useState } from "react";
import styled from "styled-components";
import { Meta } from "./components/Meta";
import { Records } from "./components/Records";
import EventService from "./service/EventService";
import { store } from "./store";
import { EventMessage, Record } from "./types/Events";
import { Login } from "./components/Login";

const AppContainer = styled.main`
    
`

export function App() {
    const { state } = useContext(store);
    const [joinedChannels, setJoinedChannels] = useState(0);
    const [activeChannels, setActiveChannels] = useState(0);
    const [records, setRecords] = useState<Array<Record>>([]);

    useEffect(() => {
        new EventService(state.apiBaseUrl, (message: EventMessage) => {
            setJoinedChannels(message.joinedChannels);
            setActiveChannels(message.activeChannels);
            setRecords(message.records);
        });
    }, [state.apiBaseUrl]);

    return <AppContainer>
        <Meta activeChannels={activeChannels} joinedChannels={joinedChannels} />
        <Records records={records} />
        <Login />
    </AppContainer>
}