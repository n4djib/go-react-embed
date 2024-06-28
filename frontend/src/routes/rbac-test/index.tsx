import { createFileRoute } from "@tanstack/react-router";
import { UserWithRoles } from "../../lib/rbac/types";
import { useAuth } from "../../contexts/auth-context";
import { useEffect } from "react";
import { Button } from "@material-tailwind/react";

export const Route = createFileRoute("/rbac-test/")({
  component: () => <Index />,
});

function Index() {
  const { user: authUser, rbac } = useAuth();

  useEffect(() => {
    handleClick();
  }, [rbac]);

  const handleClick = () => {
    // console.clear();
    if (!authUser || !rbac) return;

    const user: UserWithRoles = {
      id: authUser.id,
      roles: authUser.roles,
      // roles: [
      //   // "ADMIN",
      //   "USER",
      // ],
    };
    const resource = { id: 3, owner: 3 };

    const allowed = rbac.IsAllowed(user, resource, "edit_user");

    console.log("allowed:", allowed);
    console.log(" ");
  };

  return (
    <div>
      Hello /rbac-test/!
      <br />
      <br />
      <Button color="indigo" onClick={handleClick}>
        Check RBAC
      </Button>
    </div>
  );
}
