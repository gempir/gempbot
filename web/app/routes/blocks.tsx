import { createFileRoute } from "@tanstack/react-router";
import { Blocks as BlocksPage } from "../components/Blocks/Blocks";

export const Route = createFileRoute("/blocks")({
  component: Blocks,
});

function Blocks() {
  return <BlocksPage />;
}
