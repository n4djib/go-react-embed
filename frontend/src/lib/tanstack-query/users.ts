import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query";

export const getData = async (url: string) => {
  const res = await fetch(url);
  return res.json();
};

const baseUrl =
  import.meta.env.VITE_APP_URL + ":" + import.meta.env.VITE_APP_PORT;

type User = {
  id: number;
  name: string;
  password: string;
};

type UsersData = {
  count: number;
  data: User[];
};

export const useUserList = () => {
  return useQuery({
    queryKey: ["users"],
    queryFn: async () => {
      try {
        const users: UsersData = await getData(`${baseUrl}/api/auth/users`);
        return users;
      } catch (error) {
        console.log("Error while fetching Users");
        console.log(error);
      }
    },
  });
};

export const useUser = (id: number) => {
  return useQuery({
    queryKey: ["user", id],
    queryFn: async () => {
      try {
        const user: User = await getData(`${baseUrl}/api/auth/users/${id}`);
        return user;
      } catch (error) {
        console.log("Error while fetching User " + id);
        console.log(error);
      }
    },
  });
};

export const useInsertUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    async mutationFn(data: { name: string; password: string }) {
      const response = await fetch(`${baseUrl}/api/auth/signup`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        throw new Error("Failed to insert record ---");
      }
      return await response.json();
    },
    async onSuccess() {
      await queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    async onError() {
      console.log("Error creating a new User ----");
    },
  });
};
