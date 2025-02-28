const pem = document.getElementById("pem")
const form = document.getElementById("form");
['submit', 'input', 'focusout'].forEach(event => form.addEventListener(event, function (event) {
    event.preventDefault();
    validateForm();
}));

function validateForm() {
    console.log("validation triggered")
    let isValid = true;
    let errorMessage = ''

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

// required to make drop work on Firefox
pem.ondragover = function(e) {
    e.preventDefault();

    if (!e.target.classList.contains('drop-indicator')) {
        e.target.classList.add('drop-indicator')
    }
}

pem.ondragend = function (e) {
    e.target.classList.remove('drop-indicator')
}

pem.ondragleave = function (e) {
    e.target.classList.remove('drop-indicator')
}

pem.ondrop = function (e) {
    e.preventDefault();
    const file = e.dataTransfer.files[0];
    dropfile(file);
};

function dropfile(file) {
    const reader = new FileReader();
    reader.onload = function (e) {
        pem.value = e.target.result;
        validateForm();
    };
    reader.readAsText(file, "UTF-8");
}
