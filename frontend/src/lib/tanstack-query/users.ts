import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query";

export const getData = async (url: string) => {
  const res = await fetch(url);
  return res.json();
};

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
const CREDENTIALS = import.meta.env.VITE_CREDENTIALS;

export type User = {
  id: number;
  name: string;
  password: string;
};

// type UserData = {
//   count: number;
//   result: User;
// };

type UsersData = {
  count: number;
  // FIXME change data to result
  data: User[];
};

export const useUserList = () => {
  return useQuery({
    queryKey: ["users"],
    queryFn: async () => {
      try {
        const users: UsersData = await getData(`${BACKEND_URL}/api/auth/users`);
        return users;
      } catch (error) {
        console.log("Error while fetching Users");
        console.log(error);
      }
    },
  });
};

export const useUser = (id: number) => {
  const url = `${BACKEND_URL}/api/users/${id}`;
  return useQuery({
    queryKey: ["user", id],
    queryFn: async () => {
      try {
        const response = await fetch(url, { credentials: CREDENTIALS });
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.message);
        }
        return data as User;
      } catch (error) {
        throw error;
      }
    },
  });
};

export const useUserWhoami = () => {
  const url = `${BACKEND_URL}/api/auth/whoami/`;

  const { isError, error, data, isLoading, isFetched } = useQuery({
    queryKey: ["user"],
    queryFn: async () => {
      try {
        const response = await fetch(url, { credentials: CREDENTIALS });
        const data = await response.json();
        if (response.status === 400 || response.status === 404) {
          return null;
        }
        if (!response.ok) throw new Error(data.message);

        return data.user as User;
      } catch (error) {
        throw error;
      }
    },
    // FIXME problem - it retries 3 more times if it throws error
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

export const useInsertUser = () => {
  const url = `${BACKEND_URL}/api/auth/signup`;
  const queryClient = useQueryClient();
  return useMutation({
    async mutationFn(data: { name: string; password: string }) {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        throw new Error("Failed to creating a new User");
      }
      return await response.json();
    },
    async onSuccess(/*data*/) {
      await queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    // async onError() {
    //   console.log("Error creating a new User");
    // },
  });
};
