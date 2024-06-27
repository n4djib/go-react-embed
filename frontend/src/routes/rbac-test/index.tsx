import { createFileRoute } from "@tanstack/react-router";
import { UserWithRoles } from "../../lib/rbac";
import { useAuth } from "../../contexts/auth-context";
import { useEffect } from "react";

export const Route = createFileRoute("/rbac-test/")({
  component: () => <Index />,
});

function Index() {
  const { user: authUser, rbac } = useAuth();

  useEffect(() => {
    if (!authUser || !rbac) return;

    const user: UserWithRoles = {
      id: authUser.id,
      roles: authUser.roles,
    };
    const resource = { id: 3, owner: 3 };

    // const date1 = new Date().getTime();
    const allowed = rbac.IsAllowed(user, resource, "edit_user");
    // const date2 = new Date().getTime();

    // const diff = date2 - date1;
    // console.log("duration [ms]:", diff);

    console.log("allowed:", allowed);
  }, [rbac]);

  return (
    <div>
      Hello /rbac-test/!
      <br />
      <br />
    </div>
  );
}
