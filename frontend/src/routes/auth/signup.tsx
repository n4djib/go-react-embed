import { zodResolver } from "@hookform/resolvers/zod";
// FIXME import size is huge
import { Button, Input } from "@material-tailwind/react";
import { createFileRoute } from "@tanstack/react-router";
import { Eye, EyeOff, User } from "lucide-react";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";
import { useInsertUser } from "../../lib/tanstack-query/users";
import ErrorMessage from "../../components/ErrorMessage";
import toast from "react-hot-toast";

export const Route = createFileRoute("/auth/signup")({
  component: () => <SignUp />,
});

// TODO Check if user is available by calling api
const Schema = z
  .object({
    name: z
      .string()
      .min(2, "Username must be at least 2 characters")
      .max(45, "Username must be less than 45 characters"),
    password: z
      .string()
      .min(6, "Password must be at least 6 characters")
      .max(50, "Password must be less than 50 characters"),
    confirmPassword: z
      .string()
      .min(6, "Password must be at least 6 characters")
      .max(50, "Password must be less than 50 characters"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Password and Confirm Password doesn't match!",
    path: ["confirmPassword"], // path is accepting only one field!
    // and maybe it triggers on confirmPassword changes
  });

type InputType = z.infer<typeof Schema>;

function SignUp() {
  const [isVisiblePass, setIsVisiblePass] = useState(false);
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<InputType>({
    resolver: zodResolver(Schema),
  });

  const {
    mutate: insertUser,
    data: createdUser,
    isSuccess,
    error,
  } = useInsertUser();

  const signUp: SubmitHandler<InputType> = async (data) => {
    // const resp =
    await insertUser({
      name: data.name,
      password: data.password,
    });
    if (isSuccess) {
      toast.success(createdUser.message);
    } else {
      toast.error(error?.message || "err");
    }
  };

  return (
    <>
      <h2 className="text-xl font-bold mb-3">Sign Up Form</h2>
      <form
        onSubmit={handleSubmit(signUp)}
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

        <Input
          crossOrigin={undefined}
          label="Confirm Password"
          {...register("confirmPassword")}
          error={!!errors.confirmPassword}
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
        <ErrorMessage err={errors.confirmPassword} />

        <Button color="black" type="submit">
          {isSubmitting ? "Please wait" : "Submit"}
        </Button>
      </form>
    </>
  );
}
