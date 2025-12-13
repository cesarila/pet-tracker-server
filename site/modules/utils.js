// Resources referenced in writing this file:
// https://reqbin.com/code/javascript/car0zbiq/javascript-xmlhttprequest-example
// https://www.w3schools.com/js/js_ajax_http.asp
// https://www.w3schools.com/jsref/met_table_insertrow.asp
// https://developer.mozilla.org/en-US/docs/Web/API/Element/replaceChildren
const api_url = "http://localhost:8080"

function refreshPetData() {
    const getPetsRequest = new XMLHttpRequest();

    getPetsRequest.onload = () => {
        const pets = JSON.parse(getPetsRequest.response);

        let rowIndex = 0;

        //Display the pet data
        const bodyElement = document.querySelector("body");
        const displayAreaElement = document.querySelector("table#petDisplay");
        displayAreaElement.replaceChildren(); //delete existing rows

        Object.keys(pets).forEach((name) => {
            let row = displayAreaElement.insertRow(rowIndex);
            let nameCell = row.insertCell(0);
            let statusCell = row.insertCell(1);
            let insideButton = row.insertCell(2);
            let outsideButton = row.insertCell(3);
            let deleteButton = row.insertCell(4);
            nameCell.innerHTML = name;
            statusCell.innerHTML = pets[name];
            insideButton.appendChild(createUpdateButton(name, "inside"));
            outsideButton.appendChild(createUpdateButton(name, "outside"));
            deleteButton.appendChild(createDeleteButton(name));
            rowIndex++;
        });
        //For some reason, the display table was being shoved within the form tag, causing all buttons to post a pet with the text box text,
        //not only the submit button. This resolves the issue by forcing the display to be in the body element.
        bodyElement.appendChild(displayAreaElement);
    };

    getPetsRequest.open("GET", `${api_url}/pets`);
    getPetsRequest.send();
}

function patchPetData(petId, newStatus) {
    const patchPetRequest = new XMLHttpRequest();
    patchPetRequest.open("PATCH", `${api_url}/pets/${petId}`);
    let body = { updated_status: `${newStatus}` };
    patchPetRequest.setRequestHeader("Content-Type", "application/json");
    patchPetRequest.onload = () => console.log(patchPetRequest.response);
    patchPetRequest.addEventListener("loadend", refreshPetData);
    patchPetRequest.send(JSON.stringify(body));
}

function postNewPet(bodyObject){
    const postNewPetRequest = new XMLHttpRequest();
    postNewPetRequest.open("POST", `${api_url}/pets`);
    postNewPetRequest.setRequestHeader("Content-Type", "application/json");
    postNewPetRequest.onload = () => console.log(postNewPetRequest.response);
    postNewPetRequest.addEventListener("loadend", refreshPetData)
    postNewPetRequest.send(JSON.stringify(bodyObject));
}

function deletePetData(petId) {
    const deletePetRequest = new XMLHttpRequest();
    deletePetRequest.open("DELETE", `${api_url}/pets/${petId}`);
    deletePetRequest.onload = () => console.log(deletePetRequest.response);
    deletePetRequest.addEventListener("loadend", refreshPetData);
    deletePetRequest.send();
}

function createUpdateButton(petId, newStatus) {
    let newButton = document.createElement("BUTTON");
    let handler = () => patchPetData(petId, newStatus);
    newButton.addEventListener("click", handler);
    newButton.innerHTML = `Set to ${newStatus}`;
    return newButton;
}

function createDeleteButton(petId) {
    let newButton = document.createElement("BUTTON");
    let handler = () => deletePetData(petId);
    newButton.addEventListener("click", handler);
    newButton.innerHTML = "Delete this pet";
    return newButton;
}

export { refreshPetData };
