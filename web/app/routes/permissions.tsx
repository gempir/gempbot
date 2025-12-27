import { createFileRoute } from "@tanstack/react-router";
import { Permissions as PermissionsPage } from "../components/Permissions/Permissions";

export const Route = createFileRoute("/permissions")({
  component: Permissions,
});

function Permissions() {
  return <PermissionsPage />;
}
