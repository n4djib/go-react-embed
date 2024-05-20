import { Badge, Button } from "@material-tailwind/react";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/buttons")({
  component: buttonsExamples,
});

function buttonsExamples() {
  return (
    <>
      <h1 className="text-2xl font-bold">Buttons Examples</h1>
      <div className="flex gap-2 mt-8 flex-wrap">
        <Badge color="green">
          <Button
            variant="outlined"
            color="black"
            className="flex items-center gap-2"
          >
            <img
              src="https://docs.material-tailwind.com/icons/google.svg"
              alt="metamask"
              className="h-4 w-4"
            />
            black outlined
          </Button>
        </Badge>
        <Button color="white">white</Button>
        <Button color="black">black</Button>
        <Button color="gray">gray</Button>
        <Button variant="gradient" color="gray">
          gray gradient
        </Button>
        <Button color="blue-gray">blue-gray</Button>
        <Button color="indigo">indigo</Button>
        <Button loading color="blue">
          blue loading
        </Button>
        <Button color="blue" className="rounded-full">
          blue
        </Button>
        <Button color="light-blue">light-blue</Button>
        <Button color="cyan">cyan</Button>
        <Button color="brown">brown</Button>
        <Button color="teal">teal</Button>
        <Button color="green">green</Button>
        <Button color="light-green">light-green</Button>
        <Button color="lime">lime</Button>
        <Button color="yellow">yellow</Button>
        <Button color="amber">amber</Button>
        <Button color="orange">orange</Button>
        <Button color="deep-orange">deep-orange</Button>
        <Button color="red">red</Button>
        <Button color="pink">pink</Button>
        <Button color="purple">purple</Button>
        <Button size="lg" color="deep-purple">
          deep-purple
        </Button>
      </div>
    </>
  );
}
