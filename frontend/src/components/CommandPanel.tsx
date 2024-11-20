import { models } from "../../wailsjs/go/models";
import "../App.css";
import { useEffect, useRef, useState } from "react";
import { CommandMatch } from "../../wailsjs/go/api/App";

interface CommandPanelProps {
  setIsShow: React.Dispatch<React.SetStateAction<boolean>>;
  IsShow: boolean;
  setCode: React.Dispatch<React.SetStateAction<string>>;
}

function CommandPanel(props: CommandPanelProps) {
  const [cmd, setCmd] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const [candidators, setCandidators] = useState<models.StockMetaItem[]>([]);
  const [focusIndex, setFocusIndex] = useState(0);
  const inputHandler = (e: KeyboardEvent) => {
    if (props.IsShow && e.key === "Escape") {
      e.preventDefault();
      props.setIsShow(false);
      setCmd("");
      setFocusIndex(0);
    } else if (!props.IsShow && /^[0-9a-zA-Z]+$/.test(e.key)) {
      props.setIsShow(true);
      setFocusIndex(0);
    } else if (props.IsShow && e.key === "Enter") {
      e.preventDefault();
      props.setIsShow(false);
      setCmd("");
      if (candidators.length < focusIndex) {
        setFocusIndex(0);
        return;
      }
      props.setCode(candidators[focusIndex].Code);
      setFocusIndex(0);
    } else if (props.IsShow && e.key === "ArrowUp") {
      e.preventDefault();
      if (candidators.length === 0) {
        return;
      }
      setFocusIndex((focusIndex - 1) % candidators.length);
    } else if (props.IsShow && (e.key === "ArrowDown" || e.key === "Tab")) {
      e.preventDefault();
      if (candidators.length === 0) {
        return;
      }
      setFocusIndex((focusIndex + 1) % candidators.length);
    }
  };
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
    inputRef.current?.focus();
    document.addEventListener("keydown", inputHandler);
    return () => {
      document.removeEventListener("keydown", inputHandler);
    };
  });
  return (
    <div id="command-panel-root" className="flex">
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
