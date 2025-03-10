import { Link, createFileRoute } from "@tanstack/react-router";
import {
  PokemonsResult,
  usePokemonList,
  Pokemon,
} from "../../lib/tanstack-query/pokemons";
import { useEffect, useState } from "react";
import { capitalizeFirstLetter } from "../../lib/utils";
import { ArrowLeft, ArrowRight } from "lucide-react";
import {
  List,
  ListItem,
  ListItemPrefix,
  // Avatar,
  Card,
  Typography,
  IconButton,
  Button,
  Spinner,
} from "@material-tailwind/react";

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

  const {
    data: pokemons,
    isLoading,
    refetch,
    error,
  } = usePokemonList({ limit, offset });

  useEffect(() => {
    refetch();
  }, [limit, offset]);

  if (isLoading) return <Spinner className="h-12 w-12" />;
  if (error) return <div>Failed to fetch.</div>;

  return (
    <>
      <h1 className="text-2xl font-bold">Pokemons List</h1>
      <Card className="w-80">
        <List>
          {pokemons &&
            pokemons.result &&
            pokemons.result.map((pokemon) => (
              <PokemonListItem pokemon={pokemon} key={pokemon.id} />
            ))}
        </List>
      </Card>
      <Pagination pokemons={pokemons!} limit={limit} offset={offset} />
    </>
  );
}

function PokemonListItem({ pokemon }: { pokemon: Pokemon }) {
  return (
    <Link to="/pokemons/$id" params={{ id: `${pokemon.id}` }}>
      <ListItem>
        <ListItemPrefix>
          {/* FIXME overflow is not working in firefox */}
          {/* <Avatar
            variant="square"
            alt={pokemon.name}
            src={pokemon.image}
            className="overflow-visible "
          /> */}
          <img alt={pokemon.name} src={pokemon.image} height={50} width={50} />
        </ListItemPrefix>
        <div className="ml-4">
          <Typography variant="h6" color="blue-gray">
            {capitalizeFirstLetter(pokemon.name)}
          </Typography>
          {/* <Typography variant="small" color="gray" 
                    className="font-normal">Description</Typography> */}
        </div>
      </ListItem>
    </Link>
  );
}

function Pagination({
  pokemons,
  limit,
  offset,
}: {
  pokemons: PokemonsResult;
  limit: number;
  offset: number;
}) {
  const [offsets, setOffsets] = useState<number[]>([]);
  const count = pokemons.count;

  useEffect(() => {
    const offsets = offsetsList();
    if (offsets) {
      setOffsets(offsets);
    }
  }, [pokemons && pokemons.count]);

  const offsetsList = () => {
    if (!pokemons) {
      return;
    }
    const offsets: number[] = [];
    let currentOffset = 0;
    while (count > currentOffset) {
      offsets.push(currentOffset);
      currentOffset += limit;
    }
    return offsets;
  };

  if (!pokemons) return <Spinner />;

  const prev = offset <= 0 ? true : false;
  const next = offset + limit >= count ? true : false;

  return (
    <div className="flex items-center gap-1 mt-4">
      <Link
        disabled={prev}
        href="/pokemons"
        search={{
          limit: limit,
          offset: offset < 0 ? 0 : offset - limit,
        }}
      >
        <Button
          variant="text"
          className="flex items-center gap-2"
          disabled={prev}
        >
          <ArrowLeft className="h-4 w-4" />
          Previous
        </Button>
      </Link>
      <div className="flex items-center gap-1">
        {offsets.map((currentOffset, index) => (
          <Link
            href="/pokemons"
            search={{
              limit: limit,
              offset: currentOffset,
            }}
            key={index}
          >
            <IconButton variant={offset == currentOffset ? "filled" : "text"}>
              {index + 1}
            </IconButton>
          </Link>
        ))}
      </div>
      <Link
        disabled={next}
        href="/pokemons"
        search={{
          limit: limit,
          offset: offset + limit,
        }}
      >
        <Button
          variant="text"
          className="flex items-center gap-2"
          disabled={next}
        >
          Next
          <ArrowRight className="h-4 w-4" />
        </Button>
      </Link>
    </div>
  );
}
