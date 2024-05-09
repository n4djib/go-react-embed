import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Button } from "@material-tailwind/react";

export const Route = createFileRoute("/")({
  component: () => <Index />,
});

function Index() {
  const [data, setData] = useState<string>("");

  const fetchData = async () => {
    const response = await fetch("http://localhost:8080/api");
    const data = await response.text();
    setData(data);
    console.log("fetching");
  };

  useEffect(() => {
    fetchData().catch((err) => console.log(err));
  }, []);

  return (
    <>
      <h1 className="text-2xl font-bold">Home Page</h1>
      <div>Data fetched from (/api) : {data}</div>
      <br />
      <div className="flex w-max gap-2">
        <Button color="indigo" onClick={fetchData}>
          Fetch
        </Button>
        <Button color="red" onClick={() => setData("")}>
          Unset
        </Button>
        {/* <Button color="green">color green</Button> */}
        {/* <Button color="amber">color amber</Button> */}
      </div>
    </>
  );
}
