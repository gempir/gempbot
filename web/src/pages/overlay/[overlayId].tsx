import React from "react";
import { initializeStore } from "../../service/initializeStore";
import { IframeOverlayPage } from "../../components/Overlay/IframeOverlayPage";

export default function Overlay() {
    return <IframeOverlayPage />
}

export const getServerSideProps = initializeStore;