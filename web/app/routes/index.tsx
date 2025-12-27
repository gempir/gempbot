import { createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";
import { Teaser } from "../components/Teaser";
import { useStore } from "../store";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  const _isLoggedIn = useStore((s) => !!s.scToken);

  useEffect(() => {
    const path = window.localStorage.getItem("redirect");

    if (path) {
      window.localStorage.removeItem("redirect");
      window.location.href = path;
    }
  }, []);

  return (
    <div className="p-4 w-full max-h-screen flex gap-4">
      <Teaser />
    </div>
  );
}
