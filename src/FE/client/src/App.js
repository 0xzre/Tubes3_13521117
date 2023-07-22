import './App.css';
import History from './History';
import { useRef, useState, useEffect } from 'react'
import axios from "axios"
import DOMPurify from 'dompurify';
const calculator = require('./calculator')
const date = require('./date')

function App() {

  const [input, setInput] = useState("");
  const [chatLog, setChatLog] = useState([]);
  const [history, setHistory] = useState([]);

  async function handleSubmit(e){
    e.preventDefault();

    setChatLog(prevLog => [...prevLog, { user: "me", message: `${input}` }]);
    
    setInput("");
    console.log("Input sent:");
    if (input === "") {
      setChatLog(prevLog => [...prevLog, { user: "gpt", message: `Masukkan input!` }]);
      // setHistory(prevHistory => [...prevHistory, {question: input, answer: `Masukkan input!` }]);
      return
    }

    // Check if there is more than 1 sentences with delimiter ". "
    if (input.includes(". ")){
      const sentences = input.split(". ").map(sentence => sentence.trim());
      for (let i = 0; i < sentences.length; i++) {
        const sentence = sentences[i];
        console.log(sentence);
        await getAnswer(sentence);
    }
  }

  // Only 1 sentence
  else {
    await getAnswer(input);
  }
  
async function getAnswer(question){

  const regexCalc = /^(\s*[\d()+\/*^+-]+\s*)*\??$/;
  const regexDate = /(0?[1-9]|[1-2][0-9]|3[0-1])\/(0?[1-9]|1[0-2])\/\d{4}/g;

  if (regexDate.test(question)) {
    if (date.isValidDate(question)) { // if Valid input date dd/mm/yyyy or d/m/yyyy
      const answer = date.getDay(question)
      setChatLog(prevLog => [...prevLog, { user: "gpt", message: `Tanggal tersebut adalah hari ${answer}` }]);
    }
    
    else if (regexCalc.test(question)) { // not valid as a date, but valid as calculator (cannot retrieve minus sign)
      const answer = calculator.calculate(question)
      setChatLog(prevLog => [...prevLog, { user: "gpt", message: `Perintah dianggap operasi matematika dengan hasil ${answer}` }]);

    } else { // not valid as both
      setChatLog(prevLog => [...prevLog, { user: "gpt", message: `Input yang mengandung tanda '/' haruslah operasi matematika atau tanggal yang valid, bukan perintah` }]);
    }
  
  } else if (regexCalc.test(question)) { // valid as calculator (cannot retrieve minus sign)
    const answer = calculator.calculate(question)
    setChatLog(prevLog => [...prevLog, { user: "gpt", message: `Perintah dianggap operasi matematika dengan hasil ${answer}` }]);
  
  } else{ // check if a valid question

    if (question.includes('/')){ // cannot retrieve backslash (/) character into endpoints
      setChatLog(prevLog => [...prevLog, { user: "gpt", message: `Input yang mengandung tanda '/' haruslah operasi matematika atau tanggal yang valid, bukan perintah` }]);
    
    } else {
      var encodedInput = encodeURIComponent(question);
      var url;
      
      //if toggle KMP
      if(selectedAlgorithm === 'kmp') {
        url = `/response/KMP/${encodedInput}`;
        console.log("KMP");
      }
      
      //if toggle BM
      else if(selectedAlgorithm === 'bm'){
        url = `/response/BM/${encodedInput}`;
        console.log("BM");
      }

      axios.get("https://cuakgpt-api.fly.dev" + url, {
        responseType: 'json'
      }).then(response => {
        if(response.status === 200) {
          if (response.data == null) {
            setChatLog(prevLog => [...prevLog, { user: "gpt", message: "Pertanyaan tidak ditemukan, silakan tambahkan pertanyaan"}]);
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
      await new Promise(resolve => setTimeout(resolve, 1500));
    }
      }
  };
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

  useEffect(() => {
    const url = '/response/history';
    axios.get("https://cuakgpt-api.fly.dev" + url, {
      responseType: 'json'
    }).then(response => {
      if(response.status === 200) {
        setHistory(response.data);
      }
      else{
        console.log("No history found in database");
      }
    })
    .catch(error => {
      console.log(error);
    });
  }, [chatLog]);

  const [selectedAlgorithm, setSelectedAlgorithm] = useState('bm');

  function handleToggleAlgorithm() {
    setSelectedAlgorithm(selectedAlgorithm === 'kmp' ? 'bm' : 'kmp');
  }

  return (
    <div className="App">
      <aside className="sidemenu">
        <div className="side-menu-button" onClick={clearLog}>
          <span>
            +
          </span>
          New chat
        </div>
        <div class="algorithm-toggle-container">
          <span class="toggle-label">BM</span>
          <label class="kmp-bm-toggle">
            <input type="checkbox" id="kmp-bm-toggle-checkbox" class="kmp-bm-toggle-checkbox" onChange={handleToggleAlgorithm}/>
            <span class="kmp-bm-toggle-slider"></span>
          </label>
          <span class="toggle-label">KMP</span>
        </div>

        <div className='question-answer-history-container'>
          <div className='question-answer-history-title'>
            History
          </div>
          <div className='question-answer-history'>
            {history.toReversed().map((questionAnswer, index) => (
              <History key={index} question={questionAnswer.question} answer={questionAnswer.answer} />
            ))}
          </div>
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
  const sanitizedMessage = DOMPurify.sanitize(message.message);
  const sanitizedUser = DOMPurify.sanitize(message.user);

  return (
    <div className={`chat-message ${sanitizedUser === "gpt" && "chatgpt"}`}>
      <div className="chat-message-center">
        <div className={`avatar ${sanitizedUser === "gpt" && "chatgpt"}`}>
          {/* Avatar content */}
        </div>
        <div className="message">
          <div dangerouslySetInnerHTML={{ __html: sanitizedMessage }} />
        </div>
      </div>
    </div>
  );
};

export default App;
