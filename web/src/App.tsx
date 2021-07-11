import {
    BrowserRouter as Router,
    Route, Switch
} from "react-router-dom";
import { useTitle } from "react-use";
import { Home } from "./components/Home/Home";
import { Permissions } from "./components/Permissions/Permissions";
import { Privacy } from "./components/Privacy/Privacy";
import { Rewards } from "./components/Rewards/Rewards";
import { Sidebar } from "./components/Sidebar";
import { Teaser } from "./components/Teaser";
import { store } from "./store";


export function App() {
    useTitle("bitraft");

    const scToken = store.useState(store => store.scToken);
    if (!scToken) {
        return <Teaser />
    }

    return <Router>
        <div className="flex">
            <Sidebar />
            <Switch>
                <Route path="/privacy">
                    <Privacy />
                </Route>
                <Route path="/rewards">
                    <Rewards />
                </Route>
                <Route path="/permissions">
                    <Permissions />
                </Route>
                <Route path="/">
                    <Home />
                </Route>
            </Switch>
        </div>
    </Router>
}