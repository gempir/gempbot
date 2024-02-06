import { IframeOverlayPage } from "../../components/Overlay/IframeOverlayPage";
import { initializeStoreWithProps } from "../../service/initializeStore";

export default function Overlay() {
    return <IframeOverlayPage />
}

export const getServerSideProps = initializeStoreWithProps({renderFullLayout: false});