import { createFileRoute } from "@tanstack/react-router";
import { usePokemon } from "../../lib/tanstack-query/pokemons";

export const Route = createFileRoute("/pokemons/$id")({
  component: () => <PokemonList />,
});

const PokemonList = () => {
  const { id } = Route.useParams();
  const { data: pokemon, isLoading, error } = usePokemon(parseInt(id));

  if (isLoading) return <div>Loading...</div>;

  if (error) return <div>Failed to fetch.</div>;

  return (
    <>
      {pokemon && (
        <>
          <p>ID: {pokemon.data.id}</p>
          <p>Name: {pokemon.data.name}</p>
          <div>
            <img
              src={pokemon.data.image}
              alt={`pokemon image ${pokemon.data.image}`}
              className="w-60 h-60"
            />
          </div>
        </>
      )}
    </>
  );
};
