import { createFileRoute } from "@tanstack/react-router";
import { Bot as BotPage } from "../components/Bot/Bot";

export const Route = createFileRoute("/bot")({
  component: Bot,
});

function Bot() {
  return <BotPage />;
}
