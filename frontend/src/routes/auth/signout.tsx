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
    navigate({ to: "/", replace: true });
  }, []);

  // toast.success("Logged out successfully");

  return null;

  // return (
  //   <div>
  //     Signing Out
  //   </div>
  // );
}
