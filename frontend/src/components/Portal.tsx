import { useEffect, useState } from "react";
import { Init } from "../../wailsjs/go/api/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { api } from "../../wailsjs/go/models";

type ProtalProps = {
  connectDoneCallback: () => void;
};
function Portal(props: ProtalProps) {
  const [text, setText] = useState("Connect");
  const [doing, setDoing] = useState(false);

  const connect = () => {
    if (doing) {
      return;
    }
    Init().then(() => {
      setText("Connecting...");
      setDoing(true);
    });
  };

  useEffect(() => {
    const cancel = EventsOn(api.MsgKey.init, (initDone: boolean) => {
      if (initDone) {
        props.connectDoneCallback();
      } else {
        setText("Retry connect");
        setDoing(false);
      }
    });
    return () => cancel();
  }, []);

  return (
    <div className="w-full h-full flex flex-col">
      <button
        className={`${
          doing ? "bg-gray-600" : "bg-lime-900 hover:bg-lime-800"
        } w-32 h-12 rounded-lg m-auto`}
        onClick={connect}
      >
        {text}
      </button>
    </div>
  );
}

export default Portal;
