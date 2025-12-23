import "./ChatList.css"
import { useEffect, useState } from "react"

const ChatList = ({setChatFocus}) => {
    const [chatList, setChatList] = useState([]);

    return (
        <div className="chat-list">
            {!chatList.length ? ( <div><p>No chats to display</p></div> ) : (
                <ul>
                    {chatList.map((chat, index) => {
                        const preview = chat.lastMessage ?? "";

                        return (
                            <li className="chat-item" key={index} onClick={() => {setChatFocus(chat.id)}}>
                                <div className="chat-item-body">
                                    <div>{chat.username}</div>
                                    {preview && <div>{preview}</div>}
                                </div>
                            </li>
                        )
                    })}
                </ul>
            )}
        </div>
    )
}

export default ChatList