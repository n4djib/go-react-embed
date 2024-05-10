import { Link, Outlet, createRootRoute } from "@tanstack/react-router";

const activeProps = {
  className: "font-bold underline",
};

export const Route = createRootRoute({
  component: () => (
    <>
      <ul className="flex gap-2 items-center bg-gray-200 p-2 shadow-sm">
        <h1 className="text-xl font-bold mr-3">GO-REACT</h1>
        <li>
          <Link to="/" activeProps={activeProps}>
            Home
          </Link>
        </li>
        <li>
          <Link to="/profile" activeProps={activeProps}>
            Profile
          </Link>
        </li>
        <li>
          <Link to="/pokemon" activeProps={activeProps}>
            Pokemon
          </Link>
        </li>
        <li>
          <Link to="/posts/1/edit" activeProps={activeProps}>
            Edit Post
          </Link>
        </li>
        <li>
          <Link to="/echart" activeProps={activeProps}>
            Echart
          </Link>
        </li>
      </ul>
      <div className="p-2">
        <Outlet />
      </div>
    </>
  ),
});
