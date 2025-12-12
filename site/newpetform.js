// made with close reference to https://jordanfinners.dev/blogs/how-to-easily-convert-html-form-to-json/
// referenced the following to reload pets at the proper time:
// https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequestEventTarget/load_event
import { refreshPetData } from "./modules/utils.js"

const newPetFormElement = document.querySelector('form#newpetform')

const postNewPetRequest = new XMLHttpRequest();

const getFormJSON = (form) => {
    const data = new FormData(form);
    return Array.from(data.keys()).reduce((result, key) => {
        if (result[key]){
            result[key] = data.getAll(key)
            return result
        }
        result[key] = data.get(key);
        return result;
    }, {});
};

const handler = (event) => {
    event.preventDefault();
    const valid = newPetFormElement.reportValidity();
    if (valid){
        const result = getFormJSON(newPetFormElement);
        console.log(result);
        postNewPetRequest.open("POST", "http://localhost:8080/pets");
        postNewPetRequest.setRequestHeader("Content-Type", "application/json");
        postNewPetRequest.onload = () => console.log(postNewPetRequest.response);
        postNewPetRequest.send(JSON.stringify(result));
    }
}

newPetFormElement.addEventListener("submit", handler)
postNewPetRequest.addEventListener("loadend", refreshPetData)