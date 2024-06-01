import { createFileRoute } from "@tanstack/react-router";
import { useAuth } from "../../contexts/auth-context";
import { useEffect } from "react";

export const Route = createFileRoute("/auth/signout")({
  component: SignOut,
});

function SignOut() {
  const { user, logout } = useAuth();

  useEffect(() => {
    logout();
    // TODO redirect
  }, []);

  return (
    <div>
      Signing Out
      <div>User useAuth: {JSON.stringify(user)}</div>
    </div>
  );
}
