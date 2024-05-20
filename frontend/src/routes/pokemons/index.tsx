import { Link, createFileRoute, useNavigate } from "@tanstack/react-router";
import { usePokemonList } from "../../lib/tanstack-query/pokemons";
import { Button } from "@material-tailwind/react";
import { useEffect } from "react";

type PaginationParams = {
  limit: number;
  offset: number;
};

// TODO should we validate it using zod
export const Route = createFileRoute("/pokemons/")({
  validateSearch: (search: Record<string, unknown>): PaginationParams => {
    return {
      limit: search.limit as number,
      offset: search.offset as number,
    };
  },
  component: Pokemons,
});

function Pokemons() {
  const { limit, offset } = Route.useSearch();
  const navigate = useNavigate();

  const {
    data: pokemons,
    isLoading,
    refetch,
    error,
  } = usePokemonList({ limit, offset });

  useEffect(() => {
    refetch();
  }, [limit, offset]);

  if (isLoading) return <div>Loading...</div>;

  if (error) return <div>Failed to fetch.</div>;

  return (
    <>
      <h1 className="text-2xl font-bold">All pokemons</h1>
      {pokemons &&
        pokemons.data &&
        pokemons.data.map((pokemon) => (
          <div key={pokemon.id}>
            <Link to={`/pokemons/${pokemon.id}`}>
              Pokemon #{pokemon.id} {pokemon.name}
            </Link>
          </div>
        ))}

      <div className="flex gap-2 mt-2">
        <Button
          disabled={offset <= 0 ? true : false}
          onClick={() => {
            let newOffset = offset;
            if (offset < 0) newOffset = 0;
            navigate({
              to: "/pokemons",
              replace: false,
              search: { limit: limit, offset: newOffset - limit },
            });
          }}
        >
          Previous
        </Button>
        <Button
          disabled={pokemons && offset + limit >= pokemons.count ? true : false}
          onClick={() => {
            let newOffset = offset;
            if (pokemons && offset + limit >= pokemons.count) newOffset = 0;
            navigate({
              to: "/pokemons",
              replace: false,
              search: { limit: limit, offset: newOffset + limit },
            });
          }}
        >
          Next
        </Button>
      </div>
    </>
  );
}
