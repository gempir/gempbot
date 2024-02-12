import { OverlaysPage } from "../../components/Overlay/OverlaysPage";
import { initializeStore } from "../../service/initializeStore";

export default function OverlaysPageRoute() {
    return <OverlaysPage />
}

export const getServerSideProps = initializeStore;