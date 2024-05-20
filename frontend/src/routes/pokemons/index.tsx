import { Link, createFileRoute } from "@tanstack/react-router";
import { usePokemonList } from "../../lib/tanstack-query/pokemons";
import { Button, Spinner } from "@material-tailwind/react";
import { useEffect, useState } from "react";
import {
  List,
  ListItem,
  ListItemPrefix,
  Avatar,
  Card,
  Typography,
  IconButton,
} from "@material-tailwind/react";
import { capitalizeFirstLetter } from "../../lib/utils";
import { ArrowLeft, ArrowRight } from "lucide-react";

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
  const [offsets, setOffsets] = useState<number[]>([]);
  // const navigate = useNavigate();

  const {
    data: pokemons,
    isLoading,
    refetch,
    error,
  } = usePokemonList({ limit, offset });

  useEffect(() => {
    refetch();
  }, [limit, offset]);

  const prev = offset <= 0 ? true : false;
  const next = pokemons && offset + limit >= pokemons.count ? true : false;

  const offsetsList = () => {
    if (!pokemons) return;
    const offsets: number[] = [];
    let currentOffset = 0;
    while (pokemons.count > currentOffset) {
      offsets.push(currentOffset);
      currentOffset += limit;
    }
    console.log("offsets:::", offsets);
    return offsets;
  };

  useEffect(() => {
    const offsets = offsetsList();
    if (offsets) {
      setOffsets(offsets);
    }
  }, [pokemons && pokemons.count]);

  // onClick={() => {
  //   navigate({
  //     to: "/pokemons",
  //     replace: false,
  //     search: { limit: limit, offset: newOffset - limit },
  //   });
  // }}

  if (isLoading) return <Spinner className="h-12 w-12" />;

  if (error) return <div>Failed to fetch.</div>;

  return (
    <>
      <h1 className="text-2xl font-bold">Pokemons List</h1>
      <Card className="w-80">
        <List>
          {pokemons &&
            pokemons.data &&
            pokemons.data.map((pokemon) => (
              <Link
                to="/pokemons/$id"
                params={{ id: `${pokemon.id}` }}
                key={pokemon.id}
              >
                <ListItem>
                  <ListItemPrefix>
                    <Avatar
                      variant="circular"
                      alt={pokemon.name}
                      src={pokemon.image}
                    />
                  </ListItemPrefix>
                  <div>
                    <Typography variant="h6" color="blue-gray">
                      {capitalizeFirstLetter(pokemon.name)}
                    </Typography>
                    {/* <Typography variant="small" color="gray" 
                    className="font-normal">Description</Typography> */}
                  </div>
                </ListItem>
              </Link>
            ))}
        </List>
      </Card>

      <div className="flex items-center gap-2 mt-4">
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
        <div className="flex items-center gap-2">
          {offsets.map((offset, index) => (
            <Link
              href="/pokemons"
              search={{
                limit,
                offset: offset,
              }}
            >
              <IconButton>{index}</IconButton>
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
    </>
  );
}
