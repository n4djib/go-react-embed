import {
  // useMutation,
  // useQueryClient,
  useQuery,
} from "@tanstack/react-query";

export const getData = async (url: string) => {
  const res = await fetch(url);
  return res.json();
};

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export type Pokemon = {
  id: number;
  name: string;
  image: string;
};

type PokemonResult = {
  result: Pokemon;
};

export type PokemonsResult = {
  count: number;
  limit: number;
  offset: number;
  result: Pokemon[];
};

type UsePokemonListProps = {
  limit: number;
  offset: number;
};

// TODO tanstack query can handle pagination
export const usePokemonList = ({ limit, offset }: UsePokemonListProps) => {
  const url = `${BACKEND_URL}/api/pokemons?limit=${limit}&offset=${offset}`;
  return useQuery({
    queryKey: ["pokemons"],
    queryFn: async () => {
      try {
        const pokemons: PokemonsResult = await getData(url);
        return pokemons;
      } catch (error) {
        console.log("Error while fetching Pokemons");
        console.log(error);
      }
    },
  });
};

export const usePokemon = (id: number) => {
  const url = `${BACKEND_URL}/api/pokemons/${id}`;
  return useQuery({
    queryKey: ["pokemon", id],
    queryFn: async () => {
      try {
        const response = await fetch(url);
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.message);
        }
        return data as PokemonResult;
      } catch (error) {
        console.log("Error while fetching Pokemon " + id);
        console.log(error);
      }
    },
  });
};

//   return useMutation({
//     async mutationFn(data: any) {},
//     async onSuccess() {},
//     async onError() {},
//   });
