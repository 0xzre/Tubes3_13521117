import './App.css';
import './normal.css';
import { useRef, useState, useEffect } from 'react'
import axios from "axios"

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
    setChatLog(prevLog => [...prevLog, { user: "me", message: `${input}` }]);
    console.log("input berhasil");
    setInput("");
    getAnswer(input.toLowerCase());
  }
  
  function getAnswer(question){
    var url = `http://localhost:5000/answer/${question}`;
    axios.get(url, {
      responseType: 'json'
    }).then(response => {
      if(response.status === 200){
        if (response.data == null){
          setChatLog(prevLog => [...prevLog, { user: "gpt", message: "Pertanyaan tidak ditemukan, silakan tambahkan pertanyaan"}])
        } else{
          if (response.data[0].answer === "Pertanyaan tidak ditemukan, mungkin maksud anda: \n"){
            var respon = response.data[0].answer
            for (let i = 1; i < response.data.length; i++) {
              let elmt = i.toString() + ". " + response.data[i].question;
              respon = respon + elmt
              if (i < response.data.length -1){
                respon = respon + '\n'
              }
            }
            setChatLog(prevLog => [...prevLog, { user: "gpt", message: respon}])
          } else{
            const firstAnswer = response.data[0].answer;
            setChatLog(prevLog => [...prevLog, { user: "gpt", message: firstAnswer}])
            console.log(response.data)
            console.log("dapat jawaban")
          }
        }
      }
      else{
        console.log("No question found in database")
      }
    })
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
