import React from "react";
import { Election as ElectionPage } from "../components/Election/Election";
import { initializeStore } from "../service/initializeStore";

export default function Election() {
    return <ElectionPage />
}

export const getServerSideProps = initializeStore;