import { useContext, useEffect, useState } from "react";
import styled from "styled-components";
import { Meta } from "./components/Meta";
import { Records } from "./components/Records";
import {
    BrowserRouter as Router,
    Switch,
    Route,
} from "react-router-dom";
import EventService from "./service/EventService";
import { store } from "./store";
import { EventMessage, Record } from "./types/Events";
import { Login } from "./components/Login";
import { Dashboard } from "./components/Dashboard";

const AppContainer = styled.main`
    strong {
        display: block;
        width: 100%;
        text-align: center;
        color: var(--danger);
    }
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
        <strong>Don't use this yet for your channel points, it will break soon since it's still in development. I'm working on a stable version.</strong>
        <Router>
            <Login />
            <Switch>
                <Route path="/dashboard">
                    <Dashboard />
                </Route>
                <Route path="/">
                    <Meta activeChannels={activeChannels} joinedChannels={joinedChannels} />
                    <Records records={records} />
                </Route>
            </Switch>
        </Router>

    </AppContainer>
}