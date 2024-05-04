import { useEffect, useState } from "react";

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
      <div>Data fetched from (/api): {data}</div>
    </>
  );
}

export default App;
