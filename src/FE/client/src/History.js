import "./History.css"

function History(props) {
    return(
        <>
            <div className="history-container" key={props.key}>
                <div className="question-container">
                <div className="logo q">Q</div>
                <div className="question">{props.question}</div>
                </div>
                <div className="answer-container">
                <div className="logo a">A</div>
                <div className="answer">{props.answer}</div>
                </div>
            </div>
        </>
            
    );
};

export default History;