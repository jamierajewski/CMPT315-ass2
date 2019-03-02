// Assignment #2 - CMPT315 - Jamie Rajewski - 3020090

let loginID = "";

let clickSubmit = (evt: Event): void => {

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
            alert(xhr.responseText);
        }
        else {
            alert("Invalid login ID, please try again");
        }
    }
}

let attachListener = (): void => {
    let submitBtn = <HTMLElement>document.querySelector("#login-submit");
    submitBtn.addEventListener("click", clickSubmit);
}

// This always starts as the login page
window.onload = (): void => {
    attachListener();
}