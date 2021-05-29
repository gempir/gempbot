import {
    BrowserRouter as Router,
    Route, Switch
} from "react-router-dom";
import { Navbar } from "./components/Navbar";
import { Dashboard } from "./components/Dashboard/Dashboard";
import { Home } from "./components/Home/Home";


export function App() {
    return <Router>
        <Navbar />
        <Switch>
            <Route path="/dashboard">
                <Dashboard />
            </Route>
            <Route path="/">
                <Home />
            </Route>
        </Switch>
    </Router>
}