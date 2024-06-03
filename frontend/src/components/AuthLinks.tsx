import { Link } from "@tanstack/react-router";
import { ContextUserType, useAuth } from "../contexts/auth-context";
import { useUserWhoami } from "../lib/tanstack-query/users";
import { useEffect } from "react";
import toast from "react-hot-toast";

const activeProps = {
  className: "font-bold underline",
};

const AuthLinks = () => {
  const { user, login } = useAuth();
  const { data, isLoading } = useUserWhoami();
  // const { logout } = useAuth();
  // const navigate = useNavigate();

  if (isLoading)
    <div className="flex gap-2 items-center ml-auto">Loading...</div>;

  useEffect(() => {
    if (data) {
      const user: ContextUserType = {
        id: data.id,
        name: data.name,
      };
      login(user);
    }
  }, [data]);

  const handleSignOut = () => {
    // logout();
    // navigate({ to: "/", replace: true });
    toast.success("Logged out successfully");
  };

  if (!user) {
    return (
      <div className="flex gap-2 items-center ml-auto">
        <Link to="/auth/signin" activeProps={activeProps}>
          Sign In
        </Link>
        <Link to="/auth/signup" activeProps={activeProps}>
          Sign Up
        </Link>
      </div>
    );
  }

  return (
    <div className="flex gap-2 items-center ml-auto">
      <Link
        to="/auth/signout"
        activeProps={activeProps}
        onClick={handleSignOut}
      >
        Sign Out
      </Link>
      {/* <div className="cursor-pointe" onClick={handleSignOut}>
        Sign Out
      </div> */}
    </div>
  );
};

export default AuthLinks;
