import {
    BrowserRouter as Router,
    Link, Route, Switch
} from "react-router-dom";
import { Home } from "./components/Home/Home";
import { Navbar } from "./components/Navbar";
import { Privacy } from "./components/Privacy/Privacy";


export function App() {
    return <Router>
        <Navbar />
        <Switch>
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