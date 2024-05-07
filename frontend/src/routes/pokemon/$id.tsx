import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/pokemon/$id")({
  component: () => <Pokemon />,
});

const Pokemon = () => {
  const { id } = Route.useParams();

  return <div>Hello Pokemon {id}</div>;
};
