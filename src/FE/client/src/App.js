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
    setInput("");
    console.log("Input sent");

    //TODO
    //MAKE TOGGLE KMP AND BM
    //if toggle KMP
    getAnswerKMP(input);
    
    //if toggle BM
    // getAnswerBM(input.toLowerCase());
  }
  
  function getAnswerKMP(question){
    var encodedInput = encodeURIComponent(question);
    var url = `http://localhost:5000/answer/KMP/${encodedInput}`;

    axios.get(url, {
      responseType: 'json'
    }).then(response => {
      if(response.status === 200) {
        if (response.data == null) {
          setChatLog(prevLog => [...prevLog, { user: "gpt", message: "Pertanyaan tidak ditemukan, silakan tambahkan pertanyaan"}])
          console.log("Data not found");
        }

        else {
          // response.data is already json format
          if (response.data[0].answer != null){

            // Suggestion to user
            if (response.data[0].answer === "Pertanyaan tidak ditemukan, mungkin maksudnya:") {
              var suggestions = response.data[0].answer + `<br/>`;

              for (let i = 1; i < response.data.length; i++) {
                var capitalizeSuggestion = response.data[i].question.charAt(0).toUpperCase() + response.data[i].question.slice(1);
                var rowSuggestion = i.toString() + ". " + capitalizeSuggestion;
                suggestions = suggestions + rowSuggestion;

                if (i < response.data.length - 1) {
                  suggestions = suggestions + `<br/>`
                }
              }
              setChatLog(prevLog => [...prevLog, { user: "gpt", message: suggestions }]);
              console.log("Get response from listing suggestions");
            }

            else { // Other usage

              // Json data more than 1 (not only flag, used for listing questions)
              if (response.data.length > 1) {
                var questionList = response.data[0].answer + '<br/>'

                for (let i = 1; i < response.data.length; i++) {
                  var capitalizeList = response.data[i].question.charAt(0).toUpperCase() + response.data[i].question.slice(1);
                  var rowList = i.toString() + ". " + capitalizeList;
                  questionList = questionList + rowList;

                  if (i < response.data.length - 1){
                    questionList = questionList + `<br/>`;
                  }
                }
                setChatLog(prevLog => [...prevLog, { user: "gpt", message: questionList }]);
                console.log("Get response from listing questions");
              }
              
              // Casual case
              else {
                var gptResponse = response.data[0].answer;
                gptResponse = gptResponse.charAt(0).toUpperCase() + gptResponse.slice(1);
                setChatLog(prevLog => [...prevLog, { user: "gpt", message: gptResponse }]);
                console.log("Get response from the most normal case");
              }
            }
          }

          else { // response.data is in array of map
            console.log("Check response from array of map");

            // If array invalid as a map (usually as a text)
            if (response.data[0]["answer"] == null){
              var dataStr = "";

              // If array have leading number (ex: 1[{answer: blabla}])
              // Remove leading number and parse into json
              if (!response.data.startsWith('[')){
                dataStr = response.data.substring(response.data[0]);
                const data = JSON.parse(dataStr);
                setChatLog(prevLog => [...prevLog, { user: "gpt", message: data[0]["answer"]}]);
                console.log("Get response from parsing leading number");
              }
              
              // Array map more than 1, need to combined first (update question)
              else {
                const combinedArray = response.data.replace(/\]\[/g, ',');
                const parsedArray = JSON.parse(combinedArray);
                setChatLog(prevLog => [...prevLog, { user: "gpt", message: parsedArray[0]["answer"]}]);
                console.log("Get response from parsing more than 1 array map");
              }
            }

            else { // Valid array, directly retrieve the data
              setChatLog(prevLog => [...prevLog, { user: "gpt", message: response.data[0]["answer"]}]);
              console.log("Get Response from direct valid array map");
          }
        }
      }
    }
    else{
      console.log("No question found in database");
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
