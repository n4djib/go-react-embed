import { Link, Outlet, createRootRoute } from "@tanstack/react-router";

const activeProps = {
  className: "font-bold underline",
};

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="flex bg-gray-200 p-2 shadow-sm text-nowrap">
        <Link to="/">
          <h1 className="text-xl font-bold mr-3">GO-REACT</h1>
        </Link>
        <div className="gap-2 items-center hidden sm:flex">
          <Link to="/" activeProps={activeProps}>
            Home
          </Link>
          {/* <Link to="/profile" activeProps={activeProps}>
            Profile
          </Link> */}
          <Link
            to="/pokemons"
            activeProps={activeProps}
            activeOptions={{
              includeSearch: false,
            }}
            search={{
              limit: 10,
              offset: 0,
            }}
          >
            Pokemon
          </Link>
          {/* <Link to="/posts/1/edit" activeProps={activeProps}>
            Edit Post
          </Link> */}
          <Link to="/echart" activeProps={activeProps}>
            Echart
          </Link>
        </div>
        <div className="flex gap-2 items-center ml-auto">
          <Link to="/auth/signin" activeProps={activeProps}>
            Sign In
          </Link>
          <Link to="/auth/signup" activeProps={activeProps}>
            Sign Up
          </Link>
        </div>
      </div>

      <div className="p-2">
        <Outlet />
      </div>
    </>
  ),
});
