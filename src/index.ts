// Assignment #2 - CMPT315 - Jamie Rajewski - 3020090

// With more time, I realize after reviewing my code that I could have refactored it to
// condense some of the functions and make them more generic/re-usable
// (for example, the saveAndQuit() and clickLoginSubmit())
//      - And even some of those could have been nullified by utilizing the history API
//        to avoid repeating the requests
// I also would have liked to parse out the date/time better but it was low on the priorities

// Using a class would have allowed me to avoid globals
let loginID = "";
let presID = 0;
let priorResponseMap = new Map();

// Send the new/updated response
function sendResponse(method: string, qID: number, newAnswer: string){

    let xhr = new XMLHttpRequest();

    xhr.onload = function(){
        if (xhr.status === 200){
            // Update priorResponseMap (since the DB has changed)
            priorResponseMap.set(qID, newAnswer);
        } else {
            // Show as an alert AND extra details in the console for dev debug purposes
            alert(xhr.status + " - " + xhr.statusText);
            console.log("ERROR - " + xhr.responseText);
        }
    }
    xhr.onerror = function(){
        alert("There was an issue connecting to the server");
    }
    xhr.open(method, "/api/v1/answers/");
    xhr.setRequestHeader('Authorization', 'Bearer ' + loginID);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.setRequestHeader('Accept', 'application/json');
    xhr.send(JSON.stringify({"presenterid": presID, "questionid": qID, "answertext": newAnswer}));
}

// Take in the response and verify it against the prior response. 
function verifyResponse(qID: number, newAnswer: string){

    // If it's new, send as a POST
    if (!priorResponseMap.has(qID)){
        sendResponse("POST", qID, newAnswer);
    } else {
        let oldAnswer = priorResponseMap.get(qID);
        // Otherwise, send as a PUT
        if (newAnswer != oldAnswer){
            sendResponse("PUT", qID, newAnswer);
        } else {
            // Dont send since it isnt new
            return
        }
    }
}

function saveAndQuit(){
    let xhr = new XMLHttpRequest();

    xhr.onload = function(){
        if (xhr.status === 200){
            priorResponseMap.clear();
            loadPresentationList(xhr.responseText);
        }
        else {
            // Show as an alert AND extra details in the console for dev debug purposes
            let resp = JSON.parse(xhr.responseText);
            alert(xhr.status + " - " + xhr.statusText);
            console.log("ERROR - " + resp.Details);
        }
    }
    xhr.onerror = function(){
        alert("There was an issue connecting to the server");
    }
    xhr.open("GET", "/api/v1/people/presenters");
    xhr.setRequestHeader('Authorization', 'Bearer ' + loginID);
    xhr.send();
}

function loadPresentationList(respText:string){

    // Clear the login page if it's visible
    let loginPage = <HTMLElement>document.querySelector(".login-page");
    loginPage.setAttribute("style", "display: none");

    // Clear the header if it's visible; no need to clear possible form page
    // since we will be overwriting its element in the DOM
    let header = <HTMLElement>document.querySelector(".presenter-header");
    header.setAttribute("style", "display: none");

    // Get the template stored in the HTML
    let templ = <HTMLElement>document.querySelector("#presenter-list-template");

    // Get the target element from the DOM
    let target = <HTMLElement>document.querySelector(".main-page");

    // Render the template and copy the result into the DOM
    let templFunc = doT.template(templ.innerHTML);
    target.innerHTML = templFunc(JSON.parse(respText));
}

// Take the array of JSON answers and populate the survey with them,
// while saving the selections for later use
function parsePriorResponses(priorResponses:any){

    for (let i=0; i<priorResponses.length; i++){
        let questionID = priorResponses[i].questionId;
        let answer = priorResponses[i].answerText;
        
        if (answer.length > 1){
            let target = '[id="q' +questionID+ '"]';

            // Find the appropriate element and fill in the answer
            let element = <HTMLInputElement>document.querySelector(target);
            element.innerHTML = answer;

        } else {
            let target = '[name="q' +questionID+ '"][value="'+answer+'"]';

            // Find the appropriate element and select the answer
            let element = <HTMLInputElement>document.querySelector(target);
            element.checked = true;
        }
        // Save this selection for later use
        priorResponseMap.set(questionID, answer);
    }
}

function renderSurvey(questions: string){
    // Get the template stored in the HTML
    let templ = <HTMLElement>document.querySelector("#presenter-form-template");

    // Get the target element from the DOM
    let target = <HTMLElement>document.querySelector(".main-page");

    // Render the template and copy the result into the DOM
    let templFunc = doT.template(templ.innerHTML);
    target.innerHTML = templFunc(JSON.parse(questions));
}

function renderHeader(presenterInfo: string, otherPresenters: string){
    // Show the presenter header if it was previously hidden
    let header = <HTMLElement>document.querySelector(".presenter-header");
    header.setAttribute("style", "display: inline");

    // Now render the header with presenter info
    let templ = <HTMLElement>document.querySelector("#presenter-header-template");
    let target = <HTMLElement>document.querySelector(".presenter-header");

    let templFunc = doT.template(templ.innerHTML);
    target.innerHTML = templFunc(JSON.parse(presenterInfo)); 

    // After doing that, render the drop-downs using the template
    templ = <HTMLElement>document.querySelector("#dropdowns-template");
    target = <HTMLElement>document.querySelector(".other-dropdowns");

    templFunc = doT.template(templ.innerHTML);
    target.innerHTML = templFunc(JSON.parse(otherPresenters)); 
}

// Process all required requests, then store them as an array of responses and
// utilize them as necessary
// Source:
// https://stackoverflow.com/questions/31710768/how-can-i-fetch-an-array-of-urls-with-promise-all
function clickPresenter(presenterID:number){
    priorResponseMap.clear();

    presID = presenterID;
    let urls = ["http://localhost:8080/api/v1/people/presenters/"+presenterID,
                "http://localhost:8080/api/v1/questions/",
                "http://localhost:8080/api/v1/answers/"+presenterID,
                "http://localhost:8080/api/v1/people/presenters"];
    Promise.all(urls.map(u => fetch(u, {
        headers: {"Authorization": "Bearer " + loginID,}
    })))
    .then(responses => 
        Promise.all(responses.map(res => res.text())))
    .then(resp => {
        
        // Render the survey form and the header using the loaded resources
        renderSurvey(resp[1]);
        renderHeader(resp[0], resp[3]);

        // Then, populate the HTML with the results
        let priorResponses = JSON.parse(resp[2]);

        // Access the array of answers
        priorResponses = priorResponses.answers;

        parsePriorResponses(priorResponses);

    }).catch((err) => {
        if (err == null){
            return;
        } else {
            console.log(err.message);
        }
    });
}

let clickLoginSubmit = (evt: Event): void => {

    let errorText = <HTMLElement>document.querySelector("#login-error-text");
    errorText.setAttribute("style", "display: none");

    let submitText = <HTMLInputElement>document.querySelector("#login-textbox");
    let token = submitText.value;

    let xhr = new XMLHttpRequest();

    xhr.onload = function(){
        if (xhr.status === 200){
            loginID = token;
            loadPresentationList(xhr.responseText);
        }
        else if (xhr.status === 401){
            let errorText = <HTMLElement>document.querySelector("#login-error-text");
            errorText.innerHTML = "Invalid login ID; please try again";
            errorText.setAttribute("style", "display: block");
        }
        // This won't actually ever occur in this assignment, but it's here as extra
        else if (xhr.status === 403){
            let errorText = <HTMLElement>document.querySelector("#login-error-text");
            errorText.innerHTML = "Insufficient privileges; cannot login";
            errorText.setAttribute("style", "display: block");
        }
        else {
            let errorText = <HTMLElement>document.querySelector("#login-error-text");
            errorText.innerHTML = "Cannot log in - Server Error";
            errorText.setAttribute("style", "display: block");
            console.log("ERROR - " + xhr.responseText);
        }
    }

    xhr.onerror = function(){
        alert("There was an issue connecting to the server");
    }

    xhr.open("GET", "/api/v1/people/presenters");
    xhr.setRequestHeader('Authorization', 'Bearer ' + token);
    xhr.send();
}

let attachLoginListener = (): void => {
    let submitBtn = <HTMLElement>document.querySelector("#login-submit");
    submitBtn.addEventListener("click", clickLoginSubmit);
}

// This always starts as the login page
window.onload = (): void => {
    // Clear login field
    let submitText = <HTMLInputElement>document.querySelector("#login-textbox");
    submitText.value = "";

    attachLoginListener();
}