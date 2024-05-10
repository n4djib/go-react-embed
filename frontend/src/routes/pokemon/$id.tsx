import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/pokemon/$id")({
  component: () => <PokemonList />,
});

type Pokemon = {
  id: number;
  name: string;
  url: string;
};

const PokemonList = () => {
  const { id } = Route.useParams();
  const [pokemon, setPokemon] = useState<Pokemon | null>(null);

  const fetchData = async () => {
    const response = await fetch(`http://localhost:8080/api/pokemons/${id}`);
    const data = await response.json();
    setPokemon(data);
    // console.log("fetching");
  };

  useEffect(() => {
    fetchData().catch((err) => console.log(err));
  }, []);

  // console.log(pokemon);

  return (
    <>
      {pokemon ? (
        <div>
          <p>ID: {pokemon.id}</p>
          <p>Name: {pokemon.name}</p>
          <div>
            <img
              src={pokemon.url}
              alt={`pokemon image ${pokemon.name}`}
              className="w-60 h-60"
            />
          </div>
        </div>
      ) : (
        "Loading..."
      )}
    </>
  );
};
