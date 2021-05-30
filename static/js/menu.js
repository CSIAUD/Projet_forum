let self = document.querySelector("#menu *[href*='#']")
let others = document.querySelectorAll("#menu *:not([href*='#'])")
let main = document.getElementById("main")

window.addEventListener('load',()=>{
    main.children[1].classList.add('visible')
})

for (other of others){
    other.addEventListener('click', (ev) => {
        target = ev.originalTarget
        ev.preventDefault()
        self.classList.remove("active")
        target.classList.add("active")
        main.children[1].classList.remove("visible")
        main.style.width = target.attributes.href.value == "tickets.html" ? "70%" : "55%"

        setTimeout(function() {
            window.location.href = cut(window.location.href) + target.attributes.href.value
        },300);
    })
}