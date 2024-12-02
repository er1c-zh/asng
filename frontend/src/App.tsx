import { act, useCallback, useEffect, useRef, useState } from "react";
import "./App.css";
import CommandPanel from "./components/CommandPanel";
import Terminal from "./components/Terminal";
import Portal from "./components/Portal";
import Viewer from "./components/Viewer";
import StatusBar from "./components/StatusBar";
import KeyMessage from "./components/KeyMessage";
import { ServerStatus } from "../wailsjs/go/api/App";
import { LogInfo } from "../wailsjs/runtime/runtime";

function App() {
  const [appState, setAppState] = useState(Number);
  // command panel
  const [code, setCode] = useState("");
  const [showTerminal, setShowTerminal] = useState(false);
  const [blur0, setBlur0] = useState(false);
  const [blur1, setBlur1] = useState(false);

  const connectDone = () => {
    setAppState(1);
  };

  const terminalHandler = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === " ") {
        setShowTerminal(!showTerminal);
        if (!showTerminal) {
          setBlur0(true);
        } else {
          setBlur0(false);
        }
        e.preventDefault();
      }
    },
    [showTerminal]
  );

  useEffect(() => {
    document.addEventListener("keydown", terminalHandler);
    return () => {
      document.removeEventListener("keydown", terminalHandler);
    };
  }, [showTerminal]);

  useEffect(() => {
    ServerStatus().then((info) => {
      if (info.Connected) {
        setAppState(1);
      }
    });
  }, []);

  const statusBarRef = useRef<HTMLDivElement>(null);
  const [statusBarDimensions, setStatusBarDimensions] = useState({
    width: 0,
    height: 0,
  });
  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      setStatusBarDimensions({
        width: entries[0].contentRect.width,
        height: entries[0].contentRect.height,
      });
      console.log(entries[0].contentRect);
    });
    resizeObserver.observe(statusBarRef.current!);
    return () => {
      resizeObserver.disconnect();
    };
  }, [statusBarRef.current]);

  return (
    <div id="App" className="bg-gray-900 h-dvh max-h-dvh w-full">
      <div
        id="content"
        className={`flex flex-col w-full h-full overflow-hidden ${
          blur0 || blur1 ? "blur-sm" : ""
        }`}
      >
        <div ref={statusBarRef} className="flex w-full">
          <StatusBar appState={appState} />
        </div>
        <div
          className="flex w-full"
          style={{
            height: window.innerHeight - statusBarDimensions.height,
          }}
        >
          {appState === 0 ? (
            <Portal connectDoneCallback={connectDone} />
          ) : (
            <Viewer Code={code} />
          )}
        </div>
      </div>
      <CommandPanel setCode={setCode} setIsActive={setBlur1} />
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
