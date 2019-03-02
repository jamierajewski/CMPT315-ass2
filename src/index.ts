// Assignment #2 - CMPT315 - Jamie Rajewski - 3020090

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
    attachLoginListener();
}