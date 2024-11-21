import { useEffect, useState } from "react";
import { EventsOn, LogInfo } from "../../wailsjs/runtime/runtime";
import { api, models } from "../../wailsjs/go/models";
import { ServerStatus } from "../../wailsjs/go/api/App";
import KeyMessage from "./KeyMessage";

type StatusBarProps = {
  appState: number;
};
function StatusBar(props: StatusBarProps) {
  const [time, setTime] = useState("");
  const [serverStatus, setServerStatus] = useState<models.ServerStatus>(
    models.ServerStatus.createFrom({})
  );
  useEffect(() => {
    const ticker = setInterval(() => {
      setTime(new Date().toLocaleTimeString());
    }, 500);
    const cancel = EventsOn(
      api.MsgKey.serverStatus,
      (info: models.ServerStatus) => {
        setServerStatus(info);
      }
    );
    ServerStatus().then((info) => {
      setServerStatus(info);
    });
    return () => {
      clearInterval(ticker);
      cancel();
    };
  }, []);
  return (
    <div className="flex flex-row w-full bg-gray-800">
      <div className="flex flex-col h-full w-auto max-w-48">
        <div className="w-full bg-yellow-900 text-center">{time}</div>
        <div
          className={`w-full overflow-x-hidden truncate ... ${
            serverStatus.Connected ? "bg-green-700" : "bg-red-900"
          } text-left px-2`}
        >
          {serverStatus.ServerInfo
            ? serverStatus.ServerInfo
            : serverStatus.Connected
            ? "Connected"
            : "Disconnected"}
        </div>
      </div>
      <div>{props.appState > 0 ? <KeyMessage /> : null}</div>
    </div>
  );
}

export default StatusBar;
