import React from "react";
import { Blocks as BlocksPage } from "../components/Blocks/Blocks";
import { initializeStore } from "../service/initializeStore";

export default function Blocks() {
    return <BlocksPage />
}

export const getServerSideProps = initializeStore;