document.querySelector("#research form").addEventListener("submit", (ev) => {
    ev.preventDefault()
    let target = ev.target
    let cats = document.querySelectorAll("#categories form input")
    for (let i=0; i<cats.length; i++){
        let cat = cats[i]
        if (cat.checked) target.appendChild(cat)
    }
    target.submit()
})