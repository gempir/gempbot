import React from "react";
import { Bot as BotPage } from "../components/Bot/Bot";
import { initializeStore } from "../service/initializeStore";

export default function Bot() {
    return <BotPage />
}

Bot.getInitialProps = initializeStore