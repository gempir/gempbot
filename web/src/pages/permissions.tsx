import React from "react";
import { Permissions as PermissionsPage } from "../components/Permissions/Permissions";
import { initializeStore } from "../service/initializeStore";

export default function Permissions() {
    return <PermissionsPage />
}

export const getServerSideProps = initializeStore;