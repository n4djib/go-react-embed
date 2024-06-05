import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Button } from "@material-tailwind/react";
import { useAuth } from "../contexts/auth-context";

export const Route = createFileRoute("/")({
  component: () => <Index />,
});

const baseUrl = import.meta.env.VITE_BACKEND_URL;

function Index() {
  const [data, setData] = useState<string>("");

  const { user } = useAuth();

  const fetchData = async () => {
    const response = await fetch(baseUrl + "/ping");
    const data = await response.text();
    setData(data);
  };

  useEffect(() => {
    fetchData().catch((err) => console.log(err));
  }, []);

  return (
    <>
      <h1 className="text-2xl font-bold">Home Page</h1>
      <div>Data fetched from (/ping) : {data}</div>
      <br />
      <div className="flex w-max gap-2">
        <Button color="indigo" onClick={fetchData}>
          Fetch
        </Button>
        <Button color="red" onClick={() => setData("")}>
          Unset
        </Button>
      </div>
      <br />
      <div>{JSON.stringify(user)}</div>
    </>
  );
}
