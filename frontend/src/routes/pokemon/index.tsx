import { Link, createFileRoute } from "@tanstack/react-router";
import { usePokemonList } from "../../lib/tanstack-query/pokemons";

export const Route = createFileRoute("/pokemon/")({
  component: () => <Pokemons />,
});

function Pokemons() {
  const { data: pokemons, isLoading } = usePokemonList();

  if (isLoading) return <div>Loading...</div>;

  return (
    <>
      <h1 className="text-2xl font-bold">All pokemons</h1>
      {pokemons
        ? pokemons.data &&
          pokemons.data.map((pokemon: any) => (
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
