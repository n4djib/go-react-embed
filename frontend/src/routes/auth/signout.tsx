import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useAuth } from "../../contexts/auth-context";
import { useEffect } from "react";

export const Route = createFileRoute("/auth/signout")({
  component: SignOut,
});

function SignOut() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    logout();
    // TODO redirect
    // // FIXME this loggs a warning
    navigate({ to: "/", replace: true });
  }, []);

  // return null;

  return (
    <div>
      Signing Out
      {/* <div>User useAuth: {JSON.stringify(user)}</div> */}
    </div>
  );
}
