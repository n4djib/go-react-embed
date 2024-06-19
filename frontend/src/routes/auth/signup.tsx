import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Eye, EyeOff, User } from "lucide-react";
import { useEffect, useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";
import { useInsertUser } from "../../lib/tanstack-query/users";
import ErrorMessage from "../../components/ErrorMessage";
import { Button, Input } from "@material-tailwind/react";

export const Route = createFileRoute("/auth/signup")({
  component: () => <SignUp />,
});

// TODO Check if user is available by calling api checkName
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
  });
// .refine(
//   async (data) => {
//     if (data.name) {
//       try {
//         console.log("---running refine:", data.name);
//         const url = `${BACKEND_URL}/api/auth/check-name/${data.name}`;
//         const response = await fetch( url
//           // , { credentials: CREDENTIALS }
//         );
//         const res = await response.json();
//         if (!response.ok) { throw new Error(res.message);  }
//         return res;
//       } catch (error) { throw error; }
//     }  return true;
//   },
//   {
//     message: "user name allready taken", path: ["name"],
//   }
// );

type InputType = z.infer<typeof Schema>;

function SignUp() {
  const [isVisiblePass, setIsVisiblePass] = useState(true);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<InputType>({ resolver: zodResolver(Schema) });

  const {
    // isSuccess,
    // isPending,
    // error,
    mutate: insertUser,
    data: user,
  } = useInsertUser({ onSuccess: null, onError: null });

  const navigate = useNavigate();

  useEffect(() => {
    if (user) {
      console.log("createdUser", user);
      navigate({ to: "/auth/signin", replace: true });
    }
  }, [user]);

  const signUp: SubmitHandler<InputType> = async (data) => {
    await insertUser({
      name: data.name,
      password: data.password,
    });
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
          autoComplete="false"
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
