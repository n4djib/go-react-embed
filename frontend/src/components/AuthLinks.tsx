import { Link } from "@tanstack/react-router";
import { useAuth } from "../contexts/auth-context";
import { useUserWhoami } from "../lib/tanstack-query/users";
import { useEffect } from "react";

const activeProps = {
  className: "font-bold underline",
};

const AuthLinks = () => {
  const { user, setUser } = useAuth();
  const { data, isLoading } = useUserWhoami();

  if (isLoading)
    <div className="flex gap-2 items-center ml-auto">Loading...</div>;

  useEffect(() => {
    setUser(data || null);
  }, [data]);

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
      <Link to="/auth/signout" activeProps={activeProps}>
        Sign Out
      </Link>
      <br />
    </div>
  );
};

export default AuthLinks;
