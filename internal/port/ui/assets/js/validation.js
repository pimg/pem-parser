
let form = document.getElementById("form");
let handler =()=> {
    console.log("validation triggered");
    event.preventDefault(); // Prevent default form submission

    validateForm()
}
['submit', 'blur', 'focusout'].forEach(event => form.addEventListener(event, handler));

function validateForm() {
    let isValid = true;
    let errorMessage = ''

    const pem = document.getElementById('pem');
    if (pem.value.trim() === '') {
        isValid = false;
        errorMessage = "Must submit a PEM file"
    }

    if (pem.value.includes("PRIVATE")) {
        isValid = false;
        errorMessage = "It seems a private key is entered, do not submit a private key! Even though we do not store anything you should never submit a private key!!"
    }

    if (!isLikelyPEM(pem.value.trim())) {
        isValid = false;
        errorMessage = "The file you provided does not seem to be a PEM file"
    }

    pem.setAttribute('aria-invalid', !isValid)
    const helperText = document.getElementById('pem-helper');
    helperText.textContent = isValid ? "detected PEM file" : errorMessage;
}

function isLikelyPEM(input) {
    const pemRegex = /-----BEGIN (?!PRIVATE KEY)([A-Z0-9 ]+)-----[\s\S]+?-----END \1-----/g;
    return pemRegex.test(input);
}