import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/posts/$id/edit")({
  component: () => <EditPost />,
});

function EditPost() {
  const { id } = Route.useParams();

  return <div>Hello /posts/$id/edit! #{id}</div>;
}
