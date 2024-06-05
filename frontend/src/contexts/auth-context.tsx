import React, { createContext, useContext, useState } from "react";

export type ContextUserType = {
  id: number;
  name: string;
  roles: string[] | null;
};

type AuthContextType = {
  user: ContextUserType | null;
  login: (userData: ContextUserType) => void;
  logout: () => void;
};

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
const CREDENTIALS = import.meta.env.VITE_CREDENTIALS;

export default function AuthContextProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, setUser] = useState<ContextUserType | null>(null);

  const login = (userData: ContextUserType) => {
    setUser(userData);
  };

  const logout = async () => {
    await fetch(`${BACKEND_URL}/api/auth/signout`, {
      credentials: CREDENTIALS,
    });
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
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
