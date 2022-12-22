import { useRoutes } from "@solidjs/router";
import { lazy } from "solid-js";

export default useRoutes([
  {
    path: "/apps",
    component: lazy(() => import("/@/pages/apps")),
  },
]);
