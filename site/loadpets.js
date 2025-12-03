const xhttp = new XMLHttpRequest();

xhttp.onload = () => {
    const pets = JSON.parse(xhttp.response)

    let rowIndex = 0
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

xhttp.open("GET", "http://localhost:8080/pets");
xhttp.send();




