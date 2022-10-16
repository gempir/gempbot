import React from "react";
import { MediaPage } from "../components/Media/MediaPage";
import { initializeStore } from "../service/initializeStore";

export default function Media() {
    return <MediaPage />
}

Media.getInitialProps = initializeStore