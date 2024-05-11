import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/pokemon/$id")({
  component: () => <PokemonList />,
});

type Pokemon = {
  id: number;
  name: string;
  image: string;
};

type Data = {
  data: Pokemon;
};

const PokemonList = () => {
  const { id } = Route.useParams();
  const [pokemon, setPokemon] = useState<Data | null>(null);

  const fetchData = async () => {
    const response = await fetch(`http://localhost:8080/api/pokemons/${id}`);
    const data = await response.json();
    setPokemon(data);
  };

  useEffect(() => {
    fetchData().catch((err) => console.log(err));
  }, []);

  return (
    <>
      {pokemon ? (
        <div>
          <p>ID: {pokemon.data.id}</p>
          <p>Name: {pokemon.data.name}</p>
          <div>
            <img
              src={pokemon.data.image}
              alt={`pokemon image ${pokemon.data.name}`}
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
