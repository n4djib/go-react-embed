import { Link } from "@tanstack/react-router";
import { useAuth } from "../contexts/auth-context";

const activeProps = {
  className: "font-bold underline",
};

const AuthLinks = () => {
  const { user, isLoading } = useAuth();

  if (isLoading)
    <div className="flex gap-2 items-center ml-auto">Loading...</div>;

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
    </div>
  );
};

export default AuthLinks;
