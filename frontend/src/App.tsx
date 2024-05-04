import { useEffect, useState } from "react";
import viteLogo from "/vite.svg";
import "./App.css";

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
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
      </div>
      Data fetched from /api : {data}
    </>
  );
}

export default App;
