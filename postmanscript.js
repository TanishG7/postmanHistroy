const isObjectEmpty = (objectName) => {
  return (
    objectName &&
    Object.keys(objectName).length === 0 &&
    objectName.constructor === Object
  );
};



const url = pm.request.url.protocol + "://" + pm.request.url.host.join(".") + "/" + pm.request.url.path.join("/")

const type = pm.request.method

const paramsObj = {}
let query = pm.request.url.query
if (isObjectEmpty(query) === false) {
    paramsObj.query = query
}
let body = pm.request.body
if (isObjectEmpty(body) === false) {
console.log(body, isObjectEmpty(body))
    paramsObj.body = body
}
const params = JSON.stringify(paramsObj)
const status = pm.response.code
const response = pm.response.text()

const preRequest = {
    method: 'POST',
    url: 'http://localhost:8081/writeFile',
    header: 'Content-Type: application/json',
    body: {
        mode: 'formdata',
        formdata: [
            { key: "url", value: url, type: "text" },
            { key: "params", value: params, type: "text" },
            { key: "response", value: response, type: "text" },
            { key: "status", value: status, type: "text" },
            { key: "type", value: type, type: "text" },
        ]
    }
};

// Send the pre-request POST request.
pm.sendRequest(preRequest, (err, response) => {
    if (err) {
        console.error('Error:', err);
    } else {
        console.log('Pre-request script response:', response.json());
    }
});


console.log(pm.response)