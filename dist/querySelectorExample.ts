var current = document.querySelector('#country');

var newElement = document.createElement('newCountry');
newElement.innerHTML = '<input type="text" id="country" name="country" placeholder="Your country...">'

current.parentNode.replaceChild(newElement, current)