import { createFileRoute } from "@tanstack/react-router";
import { Rewards as RewardsPage } from "../components/Rewards/Rewards";

export const Route = createFileRoute("/rewards")({
  component: Rewards,
});

function Rewards() {
  return <RewardsPage />;
}
