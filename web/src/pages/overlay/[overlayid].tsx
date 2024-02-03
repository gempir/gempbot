import React from "react";
import { initializeStore } from "../../service/initializeStore";
import { OverlayPage } from "../../components/Overlay/OverlayPage";

export default function Overlay() {
    return <OverlayPage />
}

export const getServerSideProps = initializeStore;