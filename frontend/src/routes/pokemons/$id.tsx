import { createFileRoute } from "@tanstack/react-router";
// import { useEffect, useState } from "react";
import { usePokemon } from "../../lib/tanstack-query/pokemons";

export const Route = createFileRoute("/pokemons/$id")({
  component: () => <PokemonList />,
});

// type Pokemon = {
//   id: number;
//   name: string;
//   image: string;
// };

// type Data = {
//   data: Pokemon;
// };

const PokemonList = () => {
  const { id } = Route.useParams();
  const { data: pokemon, isLoading } = usePokemon(parseInt(id));

  if (isLoading) return <div>Loading...</div>;

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
