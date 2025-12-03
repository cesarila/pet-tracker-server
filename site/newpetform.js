// made with close reference to https://jordanfinners.dev/blogs/how-to-easily-convert-html-form-to-json/

const newPetFormElement = document.querySelector('form#newpetform')

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
        console.log(result)
    }
}

newPetFormElement.addEventListener("submit", handler)