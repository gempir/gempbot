import { useTitle } from "react-use";
import { store } from "../../store";
import { Dashboard } from "../Dashboard/Dashboard";
import { Teaser } from "../Home/Teaser";

export function Home() {
    useTitle("bitraft");

    const scToken = store.useState(store => store.scToken);
    if (scToken) {
        return <Dashboard />
    }

    return <div>
        <Teaser />
    </div>
}