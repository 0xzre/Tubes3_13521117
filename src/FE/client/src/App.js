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
    //TODO
    //MAKE TOGGLE KMP AND BM
    //if toggle KMP
    getAnswerKMP(input);
    
    //if toggle BM
    // getAnswerBM(input.toLowerCase());
  }
  
  function getAnswerKMP(question){
    var encodedInput = encodeURIComponent(question)
    var url = `http://localhost:5000/answer/KMP/${encodedInput}`;
    axios.get(url, {
      responseType: 'json'
    }).then(response => {
      if(response.status === 200)
      {
        if (response.data == null)
        {
          setChatLog(prevLog => [...prevLog, { user: "gpt", message: "Pertanyaan tidak ditemukan, silakan tambahkan pertanyaan"}])
        }
        else
        {
          if (response.data[0].answer != null){
            if (response.data[0].answer === "Pertanyaan tidak ditemukan, mungkin maksudnya:"){
              var answer = response.data[0].answer + `<br/>`
              for (let i = 1; i < response.data.length; i++) {
                var capitalized = response.data[i].question.charAt(0).toUpperCase() + response.data[i].question.slice(1);
                var elmt = i.toString() + ". " + capitalized
                answer = answer + elmt
                if (i < response.data.length - 1){
                  answer = answer + `<br/>`
                }
              }
              setChatLog(prevLog => [...prevLog, { user: "gpt", message: answer }])
            }
            else{
              var respon = response.data[0].answer
              respon = respon.charAt(0).toUpperCase() + respon.slice(1);
              setChatLog(prevLog => [...prevLog, { user: "gpt", message: respon }])
            }
          }
          else {
            if (response.data[0]["answer"] == null){
              var dataStr = "";
              if (!response.data.startsWith('[')){
                dataStr = response.data.substring(response.data[0]);
                const data = JSON.parse(dataStr);
                setChatLog(prevLog => [...prevLog, { user: "gpt", message: data[0]["answer"]}])
                console.log("dapat jawaban");
              } else {
                const dataArray = response.data.replace(/\]\[/g, ',')
                const data = JSON.parse(dataArray);
                setChatLog(prevLog => [...prevLog, { user: "gpt", message: data[0]["answer"]}])
                console.log("dapat jawaban");
              }
            }
            else {
              setChatLog(prevLog => [...prevLog, { user: "gpt", message: response.data[0]["answer"]}])
              console.log("dapat jawaban");
          }
        }
      }
    }
    else{
      console.log("No question found in database")
    }
    })
    .catch(error => {
      console.log(error);
    });

      
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
                placeholder='Type your input here'
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
              <div dangerouslySetInnerHTML={{__html: message.message}} />
              </div>
            </div>
          </div>
  )
}

export default App;
