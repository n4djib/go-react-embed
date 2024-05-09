import { Link, createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/pokemon/")({
  component: () => <Pokemons />,
});

type Pokemon = {
  id: number;
  name: string;
  url: string;
};

type PokemonsList = {
  count: number;
  results: Pokemon[];
};

function Pokemons() {
  const [pokemons, setPokemons] = useState<PokemonsList | null>(null);

  const fetchData = async () => {
    const response = await fetch("http://localhost:8080/api/pokemons");
    const data = await response.json();
    setPokemons(data);
    console.log("fetching");
  };

  useEffect(() => {
    fetchData().catch((err) => console.log(err));
  }, []);

  return (
    <>
      <h1 className="text-2xl font-bold">All pokemons</h1>
      {pokemons
        ? pokemons.results &&
          pokemons.results.map((pokemon: any) => (
            <div key={pokemon.id}>
              <Link to={`/pokemon/${pokemon.id}`}>
                Pokemon #{pokemon.id} {pokemon.name}
              </Link>
            </div>
          ))
        : "Fetching..."}
    </>
  );
}
