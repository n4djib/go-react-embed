import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { ContextUserType } from "../../contexts/auth-context";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
const CREDENTIALS = import.meta.env.VITE_CREDENTIALS;

export type User = {
  id: number;
  name: string;
  password: string;
};

type UserResult = {
  result: User;
};

type UsersResult = {
  count: number;
  result: User[];
};

export const useUserList = () => {
  return useQuery({
    queryKey: ["users"],
    queryFn: async () => {
      try {
        const url = `${BACKEND_URL}/api/auth/users`;
        const response = await fetch(url, { credentials: CREDENTIALS });
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.message);
        }
        return data as UsersResult;
      } catch (error) {
        throw error;
      }
    },
  });
};

export const useUser = (id: number) => {
  return useQuery({
    queryKey: ["user", id],
    queryFn: async () => {
      try {
        const url = `${BACKEND_URL}/api/users/${id}`;
        const response = await fetch(url, { credentials: CREDENTIALS });
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.message);
        }
        return data as UserResult;
      } catch (error) {
        throw error;
      }
    },
  });
};

export const useCheckName = (name: string) => {
  return useQuery({
    queryKey: ["user", name],
    queryFn: async () => {
      try {
        const url = `${BACKEND_URL}/api/auth/check-name/${name}`;
        const response = await fetch(url, { credentials: CREDENTIALS });
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.message);
        }
        return data;
      } catch (error) {
        throw error;
      }
    },
  });
};

export const useUserWhoami = () => {
  const { isError, error, data, isLoading, isFetched } = useQuery({
    queryKey: ["user", "whoami"],
    queryFn: async () => {
      try {
        const url = `${BACKEND_URL}/api/auth/whoami/`;
        const response = await fetch(url, { credentials: CREDENTIALS });
        const data = await response.json();

        if (response.status === 400 || response.status === 404) {
          return null;
        }
        if (!response.ok) throw new Error(data.message);

        const user: ContextUserType = data.user;
        return user;
      } catch (error) {
        throw error;
      }
    },
    // retry: false,
  });

  return {
    isError,
    error,
    data,
    isLoading,
    isFetched,
  };
};

export const useInsertUser = ({ onSuccess, onError }: any) => {
  const defaultOnSuccess = onSuccess
    ? onSuccess
    : async (data: any) => {
        await queryClient.invalidateQueries({ queryKey: ["users"] });
        toast.success(data.message);
      };
  const defaultOnError = onError
    ? onError
    : async (error: any) => {
        toast.error(error.message);
      };

  const queryClient = useQueryClient();

  return useMutation({
    async mutationFn(data: { name: string; password: string }) {
      try {
        const url = `${BACKEND_URL}/api/auth/signup`;
        const response = await fetch(url, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: CREDENTIALS,
          body: JSON.stringify(data),
        });
        const res = await response.json();

        if (!response.ok) {
          throw new Error("Failed to create a new User");
        }

        return res;
      } catch (error) {
        throw error;
      }
    },
    onSuccess: defaultOnSuccess,
    onError: defaultOnError,
  });
};

export const useRbacData = () => {
  return useQuery({
    queryKey: ["rbac-data"],
    queryFn: async () => {
      try {
        const url = `${BACKEND_URL}/api/auth/get-rbac`;
        const response = await fetch(url, { credentials: CREDENTIALS });
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.message);
        }
        return data;
      } catch (error) {
        throw error;
      }
    },
  });
};
