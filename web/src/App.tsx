import {
    BrowserRouter as Router,
    Route, Switch
} from "react-router-dom";
import { Home } from "./components/Home/Home";
import { Permissions } from "./components/Permissions/Permissions";
import { Privacy } from "./components/Privacy/Privacy";
import { Rewards } from "./components/Rewards/Rewards";
import { Sidebar } from "./components/Sidebar";


export function App() {
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