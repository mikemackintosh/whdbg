import GettingStarted from "./pages/GettingStarted.js";
import Listener from "./pages/Listener.js";

const routes = [
  { path: "/", component: GettingStarted },
  { path: "/_/:listener", component: Listener },
];

export default routes;
