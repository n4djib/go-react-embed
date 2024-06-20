import React, { createContext, useContext, useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useUserWhoami } from "../lib/tanstack-query/users";

export type ContextUserType = {
  id: number;
  name: string;
  roles: string[];
};

type loginData = {
  name: string;
  password: string;
};

type AuthContextType = {
  user: ContextUserType | null;
  login: (userData: loginData) => void;
  logout: () => void;
  isLoading: boolean;
};

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);

// FIXME how to use ENV of Golang
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
const CREDENTIALS = import.meta.env.VITE_CREDENTIALS;

export default function AuthContextProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, setUser] = useState<ContextUserType | null>(null);
  const { data: whoami, isLoading } = useUserWhoami();

  useEffect(() => {
    if (whoami) {
      const user: ContextUserType = {
        id: whoami.id,
        name: whoami.name,
        roles: whoami.roles,
      };
      setUser(user);
    }
  }, [whoami]);

  const login = async (data: loginData) => {
    try {
      const response = await fetch(`${BACKEND_URL}/api/auth/signin`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: CREDENTIALS,
        body: JSON.stringify(data),
      });
      const result = await response.json();

      if (response.ok) {
        const contextUser: ContextUserType = {
          id: result.user.id,
          name: result.user.name,
          roles: result.roles,
        };
        setUser(contextUser);
        toast.success(result.message);
      } else {
        toast.error(result?.message || "err");
      }
    } catch (error) {
      console.log("error:", error);
    }
  };

  const logout = async () => {
    try {
      await fetch(`${BACKEND_URL}/api/auth/signout`, {
        credentials: CREDENTIALS,
      });
      setUser(null);
      toast.success("Logged out successfully");
    } catch (error) {
      toast.error("Failed to logout");
      console.log("error:", error);
    }
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// Custom hook to use the Auth Context
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
