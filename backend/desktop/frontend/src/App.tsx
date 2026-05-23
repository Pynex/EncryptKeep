import { useCallback, useState } from "react";
import "./App.css";
import { UnlockView } from "./components/UnlockView";
import { VaultView } from "./components/VaultView";

type Screen = "unlock" | "vault";

function App() {
  const [screen, setScreen] = useState<Screen>("unlock");

  const goVault = useCallback(() => setScreen("vault"), []);
  const goUnlock = useCallback(() => setScreen("unlock"), []);

  return (
    <div id="App">
      {screen === "unlock" ? (
        <UnlockView onUnlocked={goVault} />
      ) : (
        <VaultView onLocked={goUnlock} />
      )}
    </div>
  );
}

export default App;
