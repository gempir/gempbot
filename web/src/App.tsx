import { useContext, useEffect, useState } from "react";
import styled from "styled-components";
import EventService from "./service/EventService";
import { store } from "./store";
import { EventMessage } from "./types/Events";

const AppContainer = styled.main`
    
`

export function App() {
    const { state } = useContext(store);
    const [joinedChannels, setJoinedChannels] = useState(0);
    const [activeChannels, setActiveChannels] = useState(0);

    useEffect(() => {
        new EventService(state.apiBaseUrl, (message: EventMessage) => {
            setJoinedChannels(message.joinedChannels);
            setActiveChannels(message.activeChannels);
        });
    }, [state.apiBaseUrl])

    return <AppContainer>
        Joined Channels: {joinedChannels} | Active Channels: {activeChannels}
    </AppContainer>
}