import { OverlayEditPage } from "../../../components/Overlay/OverlayEditPage";
import { initializeStore } from "../../../service/initializeStore";

export default function OverlaysEditPage() {
    return <OverlayEditPage />
}

export const getServerSideProps = initializeStore;