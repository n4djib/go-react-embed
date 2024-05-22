import {
  // useMutation,
  // useQueryClient,
  useQuery,
} from "@tanstack/react-query";

export const getData = async (url: string) => {
  const res = await fetch(url);
  return res.json();
};

const baseUrl = import.meta.env.VITE_BACKEND_URL;

export type Pokemon = {
  id: number;
  name: string;
  image: string;
};

type PokemonData = {
  data: Pokemon;
};

export type PokemonsData = {
  count: number;
  limit: number;
  offset: number;
  data: Pokemon[];
};

type UsePokemonListProps = {
  limit: number;
  offset: number;
};

export const usePokemonList = ({ limit, offset }: UsePokemonListProps) => {
  const url = `${baseUrl}/api/pokemons?limit=${limit}&offset=${offset}`;
  return useQuery({
    queryKey: ["pokemons"],
    queryFn: async () => {
      try {
        const pokemons: PokemonsData = await getData(url);
        return pokemons;
      } catch (error) {
        console.log("Error while fetching Pokemons");
        console.log(error);
      }
    },
  });
};

export const usePokemon = (id: number) => {
  return useQuery({
    queryKey: ["pokemon", id],
    queryFn: async () => {
      try {
        const pokemon: PokemonData = await getData(
          `${baseUrl}/api/pokemons/${id}`
        );
        return pokemon;
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
