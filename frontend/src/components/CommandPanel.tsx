import { models } from "../../wailsjs/go/models";
import "../App.css";
import {
  MutableRefObject,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { CommandMatch } from "../../wailsjs/go/api/App";

interface CommandPanelProps {
  setCode: React.Dispatch<React.SetStateAction<string>>;
  setIsActive: React.Dispatch<React.SetStateAction<boolean>>;
}

function CommandPanel(props: CommandPanelProps) {
  const [cmd, setCmd] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const [candidators, setCandidators] = useState<models.StockMetaItem[]>([]);
  const [focusIndex, setFocusIndex] = useState(0);
  const [isActive, setIsActive] = useState(false);

  const inputHandler = useCallback(
    (e: KeyboardEvent) => {
      if (isActive && e.key === "Escape") {
        e.preventDefault();
        setIsActive(false);
        props.setIsActive(false);
        setCmd("");
        setFocusIndex(0);
      } else if (!isActive && /^[0-9a-zA-Z]+$/.test(e.key)) {
        setIsActive(true);
        props.setIsActive(true);
        setFocusIndex(0);
      } else if (isActive && e.key === "Enter") {
        e.preventDefault();
        setIsActive(false);
        props.setIsActive(false);
        setCmd("");
        if (candidators.length < focusIndex) {
          setFocusIndex(0);
          return;
        }
        props.setCode(candidators[focusIndex].Code);
        setFocusIndex(0);
      } else if (isActive && e.key === "ArrowUp") {
        e.preventDefault();
        if (candidators.length === 0) {
          return;
        }
        setFocusIndex((focusIndex - 1) % candidators.length);
      } else if (isActive && (e.key === "ArrowDown" || e.key === "Tab")) {
        e.preventDefault();
        if (candidators.length === 0) {
          return;
        }
        setFocusIndex((focusIndex + 1) % candidators.length);
      }

      if (e.key === "Escape") {
        // FIXME hijack esc from system to avoid exit from full screen
        e.preventDefault();
      }
    },
    [isActive, candidators, focusIndex]
  );
  useEffect(() => {
    if (cmd.length === 0) {
      setCandidators([]);
    } else {
      CommandMatch(cmd).then((c) => {
        setCandidators(c);
      });
    }
  }, [cmd]);
  useEffect(() => {
    document.addEventListener("keydown", inputHandler);
    return () => {
      document.removeEventListener("keydown", inputHandler);
    };
  }, [inputHandler]);
  useEffect(() => {
    inputRef.current?.focus();
  });

  return (
    <div
      id="command-panel-root"
      className={`fixed top-0 left-0 flex w-full h-full ${
        isActive ? "" : "hidden"
      }`}
    >
      <div className="flex flex-col mx-auto w-1/3 min-w-64 bg-gray-600 mt-36 h-fit rounded border-gray-600 border-4">
        <input
          value={cmd}
          ref={inputRef}
          className="w-full rounded-t text-4xl bg-gray-800 
           px-8 py-4 text-left overflow-x-auto
           focus:outline-none"
          autoComplete="off"
          autoCorrect="off"
          autoCapitalize="off"
          onChange={(e) => {
            setCmd(e.target.value);
          }}
        ></input>
        <div className="text-left text-2xl px-4 py-2 bg-gray-700 rounded-b">
          {candidators.map((c, i) => {
            return (
              <div
                key={i}
                className={`text-2xl rounded ${
                  focusIndex == i ? "bg-yellow-600" : ""
                }`}
              >
                <div className={`px-4 py-2`}>
                  {i} {c.Code} {c.Desc}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default CommandPanel;
