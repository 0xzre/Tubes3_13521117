function calculate(input) {
    const operands = [];
    const operators = [];
    input = input.replace(/\s/g, '');
    if (!validateCalculation(input)){
        return "Input tidak valid!";
    }
    const operation = input.match(/\d+|\+|\-|\*|\//g);
    operation.forEach(token => {
        if (!isNaN(token)) {
            operands.push(Number(token));
        } else {
            switch (token) {
                case '(':
                    operators.push(token);
                    break;
                    case ')':
                        while (operators[operators.length - 1] !== '(') {
                            const operator = operators.pop();
                            const rightOperand = operands.pop();
                            const leftOperand = operands.pop();
                            const result = evaluate(leftOperand, rightOperand, operator);
                            operands.push(result);
                        }
                        operators.pop();
                        break;
                        case '*':
                            case '/':
                                case '+':
                                    case '-':
                                        while (operators.length > 0 && getPrecedence(operators[operators.length - 1]) >= getPrecedence(token)) {
                                            const operator = operators.pop();
                        const rightOperand = operands.pop();
                        const leftOperand = operands.pop();
                        const result = evaluate(leftOperand, rightOperand, operator);
                        operands.push(result);
                    }
                    operators.push(token);
                    break;
                    default:
                        return "Terdapat operator tidak valid!";
                    }
                }
            });
            
            while (operators.length > 0) {
                const operator = operators.pop();
                const rightOperand = operands.pop();
                const leftOperand = operands.pop();
                const result = evaluate(leftOperand, rightOperand, operator);
                operands.push(result);
            }
            
            return operands.pop().toString();
        }
        
        function getPrecedence(operator) {
            if (operator === '+' || operator === '-') {
                return 1;
            } else if (operator === '*' || operator === '/') {
                return 2;
            } else {
                return 0;
            }
        }
        
        function evaluate(leftOperand, rightOperand, operator) {
            switch (operator) {
                case '+':
            return leftOperand + rightOperand;
            case '-':
                return leftOperand - rightOperand;
                case '*':
                    return leftOperand * rightOperand;
                    case '/':
                        return leftOperand / rightOperand;
                        default:
                            return "Terdapat operator tidak valid!";
                        }
                    }
                    
                    function validateCalculation(input) {
                        const operators = ['+', '-', '*', '/'];
                        const digits = '0123456789';
                        const parentheses = ['(', ')'];
                        const stack = [];
                        
                        for (let i = 0; i < input.length; i++) {
                            const char = input[i];
                            
                            if (digits.includes(char)) {
                                continue;
                            } else if (operators.includes(char)) {
                                const prevChar = input[i - 1];
                                const nextChar = input[i + 1];
                                if (!digits.includes(prevChar) || !digits.includes(nextChar)) {
                                    return false;
                                }
                            } else if (parentheses.includes(char)) {
                                if (char === '(') {
                                    stack.push(char);
                                } else {
                                    const lastOpenParenthesis = stack.pop();
                                    if (lastOpenParenthesis !== '(') {
                    return false;
                }
            }
        } else {
            return false;
        }
    }
    if (stack.length > 0) {
        return false;
    }
    return true;
}

module.exports = { calculate }