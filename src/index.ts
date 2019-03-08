// Assignment #2 - CMPT315 - Jamie Rajewski - 3020090

// TODO
// - Fix the alert() and replace with an HTML element being rendered
// - Fix the out of range error in server.go line 498
// - Fix case inconsistency (in JSON responses etc)
//      - Also look at the XML to verify it (Nick said there was an issue)
// - Update the rest-client file and test it in Emacs rest-client mode
//      - Strip out the HTTP/1.1 at the end of the requests since it doesnt work
// - Finally, VALIDATE the HTML

// Using a class would have allowed me to avoid globals
let loginID = "";
let presID = 0;
let priorResponseMap = new Map();

function sendResponse(method: string, qID: number, newAnswer: string){

    let xhr = new XMLHttpRequest();

    xhr.onload = function(){
        if (xhr.status == 200){
            // Update priorResponseMap (since the DB has changed)
            priorResponseMap.set(qID, newAnswer);
        } else {
            console.log("ERROR - " + xhr.responseText);
        }
    }
    xhr.open(method, "/api/v1/answers/");
    xhr.setRequestHeader('Authorization', 'Bearer ' + loginID);
    xhr.setRequestHeader('Content-Type', 'application/json');
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

function loadPresentationList(respText:string){

    // Clear the login page if it's visible; no need to clear possible form page
    // since we will be overwriting its element in the DOM
    let loginPage = <HTMLElement>document.querySelectorAll(".login-page")[0];
    loginPage.setAttribute("style", "display: none");

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
        
        // HACKY WORKAROUND; Since I don't have question type in my answer object,
        // I can't tell if its long or short answer unless I make more requests to
        // the DB. For now, check if length > 1 since this will cover 99% of cases
        if (answer.length > 1){
            let target = '[id="q' +questionID+ '"]';
            console.log("TARGET: " + target);

            // Find the appropriate element and fill in the answer
            let element = <HTMLInputElement>document.querySelector(target);
            element.innerHTML = answer;

        } else {
            let target = '[name="q' +questionID+ '"][value="'+answer+'"]';
            console.log("TARGET: " + target);

            // Find the appropriate element and select the answer
            let element = <HTMLInputElement>document.querySelector(target);
            element.checked = true;
        }
        // Save this selection for later use
        priorResponseMap.set(questionID, answer);
    }
}

// Process all required requests, then store them as an array of responses and
// utilize them as necessary
// Source:
// https://stackoverflow.com/questions/31710768/how-can-i-fetch-an-array-of-urls-with-promise-all
function clickPresenter(presenterID:number){

    presID = presenterID;
    let urls = ["http://localhost:8080/api/v1/people/presenters/"+presenterID,
                "http://localhost:8080/api/v1/questions/",
                "http://localhost:8080/api/v1/answers/"+presenterID]
    Promise.all(urls.map(u => fetch(u, {
        headers: {"Authorization": "Bearer " + loginID,}
    })))
    .then(responses => 
        Promise.all(responses.map(res => res.text())))
    .then(resp => {
        
        // Get the template stored in the HTML
        let templ = <HTMLElement>document.querySelector("#presenter-form-template");

        // Get the target element from the DOM
        let target = <HTMLElement>document.querySelector(".main-page");

        // Render the template and copy the result into the DOM
        let templFunc = doT.template(templ.innerHTML);
        target.innerHTML = templFunc(JSON.parse(resp[1]));

        // Now render the header with presenter info
        templ = <HTMLElement>document.querySelector("#presenter-header-template");
        target = <HTMLElement>document.querySelector(".presenter-header");

        templFunc = doT.template(templ.innerHTML);
        target.innerHTML = templFunc(JSON.parse(resp[0])); 

        // Then, populate the HTML with the results
        let priorResponses = JSON.parse(resp[2]);
        // Access the array of answers
        priorResponses = priorResponses.answers;

        parsePriorResponses(priorResponses);

        // **DO BETTER ERROR HANDLING HERE?**
    }).catch((err) => {
        console.log(err);
    });
}

let clickLoginSubmit = (evt: Event): void => {

    let submitText = <HTMLInputElement>document.querySelector("#login-textbox");
    let token = submitText.value;

    let xhr = new XMLHttpRequest();

    xhr.onload = function(){
        if (xhr.status == 200){
            // Success; load template for next page and render it with the content 
            //received from this request
            loginID = token;
            loadPresentationList(xhr.responseText);
        }
        else {
            // ***CHANGE THIS TO RENDER TEXT RATHER THAN ALERT***
            alert("Invalid login ID, please try again");
        }
    }

    // ***ADD ERROR HANDLING HERE (THESE ARE ONLY CONNECTION ERRORS)
    //xhr.onerror

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