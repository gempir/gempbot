import { useEffect, useState } from "react";
import {
    BrowserRouter as Router,
    Route, Switch
} from "react-router-dom";
import { Dashboard } from "./components/Dashboard";
import { Navbar } from "./components/Navbar";
import { Meta } from "./components/Meta";
import { Records } from "./components/Records";
import EventService from "./service/EventService";
import { EventMessage, Record } from "./types/Events";


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

    return <Router>
        <Navbar />
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
}