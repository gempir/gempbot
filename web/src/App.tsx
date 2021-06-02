import {
    BrowserRouter as Router,
    Route, Switch, Link
} from "react-router-dom";
import { Navbar } from "./components/Navbar";
import { Dashboard } from "./components/Dashboard/Dashboard";
import { Home } from "./components/Home/Home";
import { Privacy } from "./components/Privacy/Privacy";


export function App() {
    return <Router>
        <Navbar />
        <Switch>
            <Route path="/dashboard">
                <Dashboard />
            </Route>
            <Route path="/privacy">
                <Privacy />
            </Route>
            <Route path="/">
                <Home />
            </Route>
        </Switch>
        <Link to="/privacy" className="fixed bottom-3 right-3 hover:text-gray-400">
            Privacy
        </Link>
    </Router>
}