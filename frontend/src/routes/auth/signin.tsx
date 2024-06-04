import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { z } from "zod";
import ErrorMessage from "../../components/ErrorMessage";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, User } from "lucide-react";
import { useEffect, useState } from "react";
import { ContextUserType, useAuth } from "../../contexts/auth-context";

// FIXME import size is huge
import { Button, Input } from "@material-tailwind/react";
import toast from "react-hot-toast";
// const Button = lazy(() =>
//   import("@material-tailwind/react").then((module) => ({
//     default: module.Button,
//   }))
// );

export const Route = createFileRoute("/auth/signin")({
  component: SignIn,
});

const Schema = z.object({
  name: z
    .string()
    .min(2, "Username must be at least 2 characters")
    .max(45, "Username must be less than 45 characters"),
  password: z
    .string()
    .min(6, "Password must be at least 6 characters")
    .max(50, "Password must be less than 50 characters"),
});

type InputType = z.infer<typeof Schema>;

function SignIn() {
  const [isVisiblePass, setIsVisiblePass] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<InputType>({
    resolver: zodResolver(Schema),
  });
  const navigate = useNavigate();
  const { user, login } = useAuth();

  useEffect(() => {
    if (user) {
      navigate({ to: "/", replace: true });
    }
  }, [user]);

  const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
  const CREDENTIALS = import.meta.env.VITE_CREDENTIALS;

  const signIn: SubmitHandler<InputType> = async (data) => {
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
        };

        login(contextUser);

        toast.success(result.message);
      } else {
        toast.error(result?.message || "err");
      }
    } catch (error) {
      console.log("error:", error);
    }
  };

  return (
    <>
      <h2 className="text-xl font-bold mb-3">Sign In Form</h2>
      <form
        onSubmit={handleSubmit(signIn)}
        className="flex flex-col w-96 gap-2"
      >
        <Input
          crossOrigin={undefined}
          label="Username"
          {...register("name")}
          error={!!errors.name}
          icon={<User />}
        />
        <ErrorMessage err={errors.name} />

        <Input
          crossOrigin={undefined}
          label="Password"
          {...register("password")}
          error={!!errors.password}
          type={isVisiblePass ? "text" : "password"}
          icon={
            isVisiblePass ? (
              <EyeOff
                className="w-4 cursor-pointer"
                onClick={() => setIsVisiblePass(false)}
              />
            ) : (
              <Eye
                className="w-4 cursor-pointer"
                onClick={() => setIsVisiblePass(true)}
              />
            )
          }
        />
        <ErrorMessage err={errors.password} />

        <Button color="black" type="submit">
          {isSubmitting ? "Please wait" : "Sign In"}
        </Button>
      </form>
    </>
  );
}
