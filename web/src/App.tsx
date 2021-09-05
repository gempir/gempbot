import {
    BrowserRouter as Router,
    Route, Switch
} from "react-router-dom";
import { EmotehistoryPage } from "./components/Emotehistory/EmotehistoryPage";
import { Home } from "./components/Home/Home";
import { Permissions } from "./components/Permissions/Permissions";
import { Privacy } from "./components/Privacy/Privacy";
import { Rewards } from "./components/Rewards/Rewards";
import { Sidebar } from "./components/Sidebar/Sidebar";
import { Teaser } from "./components/Teaser";
import { store } from "./store";


export function App() {
    const scToken = store.useState(store => store.scToken);

    return <Router>
        <div className="flex">
            <Sidebar />
            <Switch>
                <Route path="/privacy">
                    <Privacy />
                </Route>
                <Route path="/emotehistory/:channel">
                    <EmotehistoryPage />
                </Route>
                <Route path="/rewards">
                    <Rewards />
                </Route>
                <Route path="/permissions">
                    <Permissions />
                </Route>
                <Route path="/">
                    {scToken && <Home />}
                    {!scToken && <Teaser />}
                </Route>
            </Switch>
        </div>
    </Router>
}