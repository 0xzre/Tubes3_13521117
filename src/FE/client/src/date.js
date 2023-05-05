function isLeapYear(year) {
    return (( (0 === year % 4) && (0 !== year % 100) ) || (0 === year % 400))
}

function isValidDate(string) {
    const date = string.slice(0, 2);
    const month = string.slice(3, 5);
    const year = string.slice(6, 10);
    if (month === '04' || month === '06' || month === '09' || month === '11') {
        if (date === '31') {
            return false;
        }
    } else if (month === '02') {
        let bottomBound = 28;
        if (isLeapYear(Number(year))) {
            bottomBound++;
        }
        if (Number(date) > bottomBound) {
            return false;
        }
    }
    return true;
}

function getDay(input) {
    const date = input.slice(0, 2);
    const month = input.slice(3, 5);
    const year = input.slice(6, 10);
    const date1 = new Date(month + '/' + date + '/' + year);
    const weekday = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
    return weekday[date1.getDay()];
}

module.exports = { getDay, isValidDate }