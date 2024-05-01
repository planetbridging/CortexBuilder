import logo from "./logo.svg";
import "./App.css";

import { useToast } from "@chakra-ui/react";

import OHome from "./OHome";

function App() {
  const toast = useToast();
  return <OHome toast={toast} />;
}

export default App;
