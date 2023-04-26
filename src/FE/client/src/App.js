import './App.css';
import './normal.css';
import { useRef, useState, useEffect } from 'react'

function App() {

  const [input, setInput] = useState("");
  const [chatLog, setChatLog] = useState([{
    user: "gpt",
    message: "hewo aim bongt"
  }, {
    user: "me",
    message: "hewo bongt"
  }, {
    user: "gpt",
    message: "hewo aim bongt"
  }
  ]);

  function handleSubmit(e){
    e.preventDefault();
    setChatLog([...chatLog, { user: "me", message: `${input}` }]);
    setInput("");

  }

  function clearLog(){
    setChatLog([]);
  }

  const newMessageRef = useRef(null)

  const scrollToBottom = () => {
    newMessageRef.current?.scrollIntoView({ behaviour: "smooth"})
  }

  useEffect(() => {
    scrollToBottom()}
    ,[chatLog]
  );

  return (
    <div className="App">
      <aside className="sidemenu">
        <div className="side-menu-button" onClick={clearLog}>
          <span>
            +
          </span>
          New chat
        </div>
      </aside>
      <section className="chatbox">
        <div className="chat-log">
          {
            chatLog.map((message, index) => (
              <ChatMessage key={index} message={message} />
            ))
          }
          <div ref={newMessageRef} />
        
        </div>
        <div
          className="chat-input-holder">
            <form onSubmit={handleSubmit}>
              <input
                value={input}
                onChange={(e) => setInput(e.target.value) }
                className="chat-input-textarea"
                rows="1"
                >
              </input>
            </form>
        </div>

      </section>

    </div>
  );
}

const ChatMessage = ({ message }) => {
  return(
          <div className={`chat-message ${message.user === "gpt" && "chatgpt" }`}>
            <div className="chat-message-center">
              <div className= {`avatar ${message.user === "gpt" && "chatgpt" }`}>
                
              </div>
              <div className="message">
                {message.message}
              </div>
            </div>
          </div>
  )
}

export default App;
