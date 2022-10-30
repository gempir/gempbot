import React from "react";
import { Rewards as RewardsPage } from "../components/Rewards/Rewards";
import { initializeStore } from "../service/initializeStore";

export default function Rewards() {
    return <RewardsPage />
}

export const getServerSideProps = initializeStore;