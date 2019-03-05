// Assignment #2 - CMPT315 - Jamie Rajewski - 3020090

// TODO
// - Fix it so that you CANT review yourself
// - Fix the out of range error in server.go line 498

let loginID = "";

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

// Used the Promise.all() feature to process multiple requests
// Source:
// https://stackoverflow.com/questions/31710768/how-can-i-fetch-an-array-of-urls-with-promise-all
function clickPresenter(presenterID:number){
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

    }).catch((err) => {
        console.log(err);
    });
}

let clickLoginSubmit = (evt: Event): void => {

    let submitText = <HTMLInputElement>document.querySelector("#login-textbox");
    let token = submitText.value;

    let xhr = new XMLHttpRequest();

    xhr.open("GET", "/api/v1/people/presenters");
    xhr.setRequestHeader('Authorization', 'Bearer ' + token);
    xhr.send();

    xhr.onload = function(){
        if (xhr.status == 200){
            // Success; load template for next page and render it with the content 
            //received from this request
            loginID = token;
            loadPresentationList(xhr.responseText);
        }
        else {
            alert("Invalid login ID, please try again");
        }
    }
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