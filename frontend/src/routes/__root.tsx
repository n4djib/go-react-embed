import { Link, Outlet, createRootRoute } from "@tanstack/react-router";
import AuthLinks from "../components/AuthLinks";

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
              limit: 7,
              offset: 0,
            }}
          >
            Pokemons
          </Link>
          <Link to="/buttons" activeProps={activeProps}>
            Buttons
          </Link>
          <Link to="/rbac-test" activeProps={activeProps}>
            RBAC
          </Link>
          {/* <Link to="/posts/1/edit" activeProps={activeProps}>
            Edit Post
          </Link> */}
          {/* <Link to="/echart" activeProps={activeProps}>
            Echart
          </Link> */}
        </div>
        <AuthLinks />
      </div>

      <div className="p-2">
        <Outlet />
      </div>
    </>
  ),
});
