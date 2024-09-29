document.getElementById("create-order-btn").addEventListener("click", function(){
    document.getElementById("order-form").classList.add("open")
})

document.getElementById("close-order-form-btn").addEventListener("click", function(){
    document.getElementById("order-form").classList.remove("open")
})

window.addEventListener('keydown', (e) => {
    if (e.key === "Escape") {
        document.getElementById("order-form").classList.remove("open")
    }
})

document.querySelector("#order-form .modal__box").addEventListener('click', event => {
    event._isClickWithInModal = true;
})

document.getElementById("order-form").addEventListener('click', event => {
    if (event._isClickWithInModal) return;
    event.currentTarget.classList.remove('open');
})