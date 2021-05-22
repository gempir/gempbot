import { useTitle } from "react-use";
import { Teaser } from "../Home/Teaser";

export function Home() {
    useTitle("bitraft");

    return <div>
        <Teaser />
    </div>
}