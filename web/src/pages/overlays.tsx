import { OverlaysPage } from "../components/Overlays/OverlaysPage";
import { initializeStore } from "../service/initializeStore";

export default function Overlays() {
    return <OverlaysPage />
}

export const getServerSideProps = initializeStore;