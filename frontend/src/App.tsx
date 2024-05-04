import { useEffect, useState } from "react";

import { Button } from "@material-tailwind/react";

function App() {
  const [data, setData] = useState<string>("");

  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("http://localhost:8080/api");
      const data = await response.text();
      setData(data);
    };
    fetchData().catch((err) => console.log(err));
  }, []);

  return (
    <>
      <h1 className="text-3xl font-bold underline">Home Page</h1>
      <div>Data fetched from (/api) : {data}</div>
      <br />
      <br />
      <div className="flex w-max gap-4">
        <Button color="blue">color blue</Button>
        <Button color="red">color red</Button>
        <Button color="green">color green</Button>
        <Button color="amber">color amber</Button>
      </div>
    </>
  );
}

export default App;
