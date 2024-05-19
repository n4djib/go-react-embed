import {
  // useMutation,
  // useQueryClient,
  useQuery,
} from "@tanstack/react-query";

export const getData = async (url: string) => {
  const res = await fetch(url);
  return res.json();
};

const baseUrl =
  import.meta.env.VITE_APP_URL + ":" + import.meta.env.VITE_APP_PORT;

type Pokemon = {
  id: number;
  name: string;
  image: string;
};

type PokemonsData = {
  count: number;
  data: Pokemon[];
};

export const usePokemonList = () => {
  return useQuery({
    queryKey: ["pokemons"],
    queryFn: async () => {
      try {
        const pokemons: PokemonsData = await getData(`${baseUrl}/api/pokemons`);
        console.log("pokemons:", pokemons);
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
    queryKey: ["pokeon", id],
    queryFn: async () => {
      try {
        // const pokemon: Pokemon = await getData(`${baseUrl}/api/pokemons/${id}`);
        const pokemon = await getData(`${baseUrl}/api/pokemons/${id}`);
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
