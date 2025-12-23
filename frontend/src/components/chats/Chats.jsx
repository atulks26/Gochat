import { useState } from "react"
import ChatList from "./chatlist/ChatList";
import ChatFocus from "./chatFocus/ChatFocus";

const Chats = () => {
  const [chatFocus, setChatFocus] = useState(null);

  return (
    <div>
        <ChatList setChatFocus={setChatFocus}/>
        <ChatFocus focusedChat={chatFocus}/>
    </div>
  )
}

export default Chats