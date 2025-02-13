
let searchWord = "";
let interval = ""
let count = 0;


function syntaxHighlight(json) {
    return json.replace(
        /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
        match => {
            let cls = "number";
            if (/^"/.test(match)) {
                if (/:$/.test(match)) {
                    cls = "key";
                } else {
                    cls = "string";
                }
            } else if (/true|false/.test(match)) {
                cls = "boolean";
            } else if (/null/.test(match)) {
                cls = "null";
            }
            return `<span class="${cls}">${match}</span>`;
        }
    );
}

const FetchResult = () => {
    let url = "http://localhost:8081/getData?search=" + searchWord;
    axios.get(url).then(resp => {
        console.log(resp.data)
        if (resp.data.length == count) {
            clearInterval(interval)
        }
        if (resp.status == 200) {
            const resultCard = document.getElementById("resultCard")
            resultCard.innerHTML = ""
            let requestDataArr = resp.data
            requestDataArr.forEach(requestData => {
                console.log(resultCard)
                if (requestData.goResult == null) {
                    requestData.goResult = {}
                }
                if (requestData.phpResult == null) [
                    requestData.phpResult = {}
                ]
                const diff = window.getDiff(requestData.goResult, requestData.phpResult)
                // const diff = "No difference"
                console.log(requestData)
                const Div = document.createElement("div")
                Div.innerHTML = `<div class="card">
                 <div class="result" id="result">
            <div class="requestParams" id="requestParams">
                <h3>PARAMS</h3>
                <pre>${syntaxHighlight(JSON.stringify(requestData.params, null, 4))}</pre>
            </div>
            <div class="difference" id="difference">
                <h3>DIFFERENCE</h3>
                <pre>${syntaxHighlight(JSON.stringify(diff, null, 4))}</pre>
            </div>
            </div>
            <div class="result" id="result">
                <div class="goResult" id="goResult">
                    <h3>Go Result</h3>
                    <p>Status: ${requestData.goStatus}</p>
                    <h4>Response</h4>
                    <pre>${syntaxHighlight(JSON.stringify(requestData.goResult, null, 4))}</pre>
                    </div>
                    <div class="phpResult" id="phpResult">
                    <h3>PHP Result</h3>
                    <p>Status: ${requestData.phpStatus}</p>
                    <h4>Response</h4>
                    <pre>${syntaxHighlight(JSON.stringify(requestData.phpResult, null, 4))}</pre>
                </div>
            </div>
        </div>`
                resultCard.appendChild(Div)
            });
        } else {
            console.log(resp)
        }

    }).catch(err => console.log(err))
}



const updateResult = () => {
    interval = setInterval(FetchResult, 3000)
    return interval
}


function addRow() {
    const row = document.createElement("div");
    row.classList.add("row");
    row.innerHTML = `
        <input type="text" placeholder="Key" class="key">
        <textarea type="text" placeholder="Value" class="value"></textarea>
        <button onclick="removeRow(this)">Remove</button>
    `;
    document.getElementById("form").appendChild(row);
}

function removeRow(button) {
    button.parentElement.remove();
}


function submitData() {
    const keys = document.querySelectorAll(".key");
    const values = document.querySelectorAll(".value");
    let data = {};

    keys.forEach((key, index) => {
        const keyValue = key.value.trim();
        const value = values[index].value.trim();
        if (keyValue) {
            data[keyValue] = value;
        }
    });

    let url = "http://localhost:8081/testCases"
    axios.post(url, data, { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }).then(response => {
        document.getElementById("output").innerHTML = syntaxHighlight(JSON.stringify(response.data, null, 2));
        searchWord = response.data.id
        count = response.data.Counts
        console.log(searchWord)
        updateResult()
    }).catch(err => console.log(err))

}



