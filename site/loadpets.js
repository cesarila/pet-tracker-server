// Resources referenced in writing this file: 
// https://reqbin.com/code/javascript/car0zbiq/javascript-xmlhttprequest-example
// https://www.w3schools.com/js/js_ajax_http.asp
// https://www.w3schools.com/jsref/met_table_insertrow.asp

const getPetsRequest = new XMLHttpRequest();

getPetsRequest.onload = () => {
    const pets = JSON.parse(getPetsRequest.response)

    let rowIndex = 0
    
    //Display the pet data
    const displayAreaElement = document.querySelector("table#petDisplay")
    Object.keys(pets).forEach(name => {
        let row = displayAreaElement.insertRow(rowIndex);
        let nameCell = row.insertCell(0);
        let statusCell = row.insertCell(1);
        nameCell.innerHTML = name;
        statusCell.innerHTML = pets[name];
        rowIndex++;
    })


}

getPetsRequest.open("GET", "http://localhost:8080/pets");
getPetsRequest.send();




