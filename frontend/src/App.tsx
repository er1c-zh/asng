import { useEffect, useState } from "react";
import "./App.css";
import CommandPanel from "./components/CommandPanel";
import Terminal from "./components/Terminal";
import Portal from "./components/Portal";
import Viewer from "./components/Viewer";
import StatusBar from "./components/StatusBar";
import KeyMessage from "./components/KeyMessage";

function App() {
  const [appState, setAppState] = useState(Number);
  // command panel
  const [cmdPanelShow, setCmdPanelShow] = useState(false);
  const [code, setCode] = useState("");
  const [showTerminal, setShowTerminal] = useState(false);
  const connectDone = () => {
    setAppState(1);
  };
  const terminalHandler = (e: KeyboardEvent) => {
    if (e.key === " ") {
      setShowTerminal(!showTerminal);
      e.preventDefault();
    }
  };
  useEffect(() => {
    document.addEventListener("keydown", terminalHandler);
    return () => {
      document.removeEventListener("keydown", terminalHandler);
    };
  });

  return (
    <div id="App" className="bg-gray-900 h-dvh w-full">
      <div
        id="content"
        className={`h-full w-full flex flex-col ${
          cmdPanelShow ? "blur-sm" : ""
        }`}
      >
        <StatusBar Components={appState === 0 ? [] : [KeyMessage]} />
        <div className="flex h-full w-full">
          {appState === 0 ? (
            <Portal connectDoneCallback={connectDone} />
          ) : (
            <Viewer Code={code} />
          )}
        </div>
      </div>
      <div
        id="command-panel"
        className={`${
          cmdPanelShow ? "" : "hidden"
        } w-full h-full fixed top-0 left-0 border-2 border-gray-500`}
      >
        <CommandPanel
          IsShow={cmdPanelShow}
          setIsShow={setCmdPanelShow}
          setCode={setCode}
        />
      </div>
      <div
        className={`fixed top-0 left-0 w-full h-full border-2 border-gray-500 opacity-75 ${
          showTerminal ? "" : "hidden"
        }`}
      >
        <Terminal />
      </div>
    </div>
  );
}

export default App;
