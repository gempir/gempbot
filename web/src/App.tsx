import { useEffect, useState } from "react";
import {
    BrowserRouter as Router,

    Route, Switch
} from "react-router-dom";
import styled from "styled-components";
import { Dashboard } from "./components/Dashboard";
import { Login } from "./components/Login";
import { Meta } from "./components/Meta";
import { Records } from "./components/Records";
import EventService from "./service/EventService";
import { EventMessage, Record } from "./types/Events";

const AppContainer = styled.main``;

export function App() {
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

    return <AppContainer>
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