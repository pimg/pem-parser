enableCopyOnClick()

document.body.addEventListener('htmx:afterSettle', function (evt) {
    enableCopyOnClick()
});

function enableCopyOnClick() {
    document
        .querySelectorAll('[data-copy-on-click]')
        .forEach((el) => {
            // first remove the listener in case it was added in a previous HTML request
            el.removeEventListener('click', onClickHandler)
            el.addEventListener('click', onClickHandler)
        });
}

function onClickHandler(event) {
    copyTextFromElementToClipboard(event.currentTarget.dataset.copyOnClick)
    event.preventDefault()
}

function copyTextFromElementToClipboard(elementID) {
    const text = document.getElementById(elementID).value.trim();

    if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard.writeText(text);
    } else {
        const textArea = document.createElement("textarea");
        textArea.value = text;
        textArea.style.position = "absolute";
        textArea.style.left = "-999999px";

        document.body.prepend(textArea);
        textArea.select();

        try {
            document.execCommand("copy");
        } catch (error) {
            console.error(error);
        } finally {
            textArea.remove();
        }
    }
}
