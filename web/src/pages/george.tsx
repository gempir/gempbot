import React from "react";
import { George as GeorgePage } from "../components/George/George";
import { initializeStore } from "../service/initializeStore";

export default function George() {
    return <GeorgePage />
}

export const getServerSideProps = initializeStore;